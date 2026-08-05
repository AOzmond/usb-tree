package cli

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/AOzmond/usb-tree/lib"
)

type focusIndex int

// Model represents the primary structure containing application state and views.
type Model struct {
	windowWidth    int
	windowHeight   int
	statusHeight   int
	statusLine     string
	updateChan     chan []lib.Device
	roots          []*lib.TreeNode
	collapsed      map[string]bool // tracks which nodes are collapsed by their unique key
	treeViewport   viewport.Model
	treeCursor     int
	nodeCount      int
	selectedDevice *lib.Device
	logViewport    viewport.Model
	log            []lib.Log
	logContent     string
	helpModel      help.Model
	focusedView    focusIndex
	lastUpdated    time.Time
}

const (
	splitRatio    = 0.7 // Ratio of tree view to log view
	borderSpacing = 2   // the space taken up by the border
	tooltipHeight = 5
)

const (
	treeView focusIndex = iota
	logView
)

// InitialModel initializes and returns a new Model instance with values for state and views.
func InitialModel() Model {
	updates := make(chan []lib.Device, 1)

	helpModel := help.New()
	helpModel.Styles.ShortDesc = windowStyle
	helpModel.Styles.ShortKey = windowStyle
	helpModel.Styles.ShortSeparator = windowStyle

	m := Model{
		helpModel:   helpModel,
		focusedView: treeView,
		lastUpdated: time.Now(),
		treeCursor:  0,
		updateChan:  updates,
		collapsed:   make(map[string]bool),
	}
	return m
}

// Init initializes the Model, preparing it to handle updateChan and rendering. It returns an optional initial command.
func (m Model) Init() tea.Cmd {
	lib.Init(func(devices []lib.Device) {
		m.updateChan <- devices
	})
	return waitForUpdate(m.updateChan)
}

// View renders the current state of the Model, combining styled views for tree, log, tooltip, and status line.
func (m Model) View() tea.View {
	if m.windowWidth == 0 || m.windowHeight == 0 {
		return tea.NewView("")
	}
	var treeStyle, logStyle lipgloss.Style

	if m.focusedView == treeView {
		treeStyle = activeStyle
		logStyle = inactiveStyle
	} else {
		treeStyle = inactiveStyle
		logStyle = activeStyle
	}

	// Check for offscreen changes to highlight borders
	above, below := m.checkOffscreenChanges()
	if above || below {
		borderStyle := treeStyle.GetBorderStyle()
		topBorderColor := treeStyle.GetBorderTopForeground()
		bottomBorderColor := treeStyle.GetBorderBottomForeground()

		if above {
			topBorderColor = edgeHighlightChangeColor
		}
		if below {
			bottomBorderColor = edgeHighlightChangeColor
		}

		treeStyle = treeStyle.Border(borderStyle, true, true, true, true).
			BorderTopForeground(topBorderColor).
			BorderBottomForeground(bottomBorderColor)
	}

	tooltip := tooltipStyle.
		Width(m.windowWidth).
		Render(m.getSelectedDeviceInfo())

	appContent := lipgloss.JoinVertical(
		lipgloss.Center,
		treeStyle.Render(m.treeViewport.View()),
		tooltip,
		logStyle.Render(m.logViewport.View()),
		m.statusLine,
	)

	view := tea.NewView(appContent)
	view.AltScreen = true
	view.BackgroundColor = backgroundColor
	return view
}

// Update processes incoming messages, updateChan the model state, and returns the updated model and an optional command.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case deviceMessage:
		devices := []lib.Device(msg)
		previousKey := ""
		if m.selectedDevice != nil {
			previousKey = m.selectedDevice.Key()
		}
		m.roots = lib.BuildDeviceTree(devices)
		m.updateNodeCount()
		if previousKey != "" {
			if cursor, found := m.visibleNodeIndexByKey(previousKey); found {
				m.treeCursor = cursor
			} else if m.treeCursor >= m.nodeCount {
				m.treeCursor = max(0, m.nodeCount-1)
			}
		} else {
			m.treeCursor = 0
		}
		m.updateSelectedDevice()
		m.refreshContent()
		m.scrollToCursor()

		m.log = lib.GetLog()
		m.logContent = m.formatLogContent()
		m.logViewport.SetContent(m.logContent)
		m.clampLogViewport()

		return m, waitForUpdate(m.updateChan)

	case tea.WindowSizeMsg:
		m.windowWidth, m.windowHeight = msg.Width, msg.Height
		m.refreshContent()
		m.scrollToCursor()
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.SwitchFocus):
			if m.focusedView == treeView {
				m.focusedView = logView
			} else {
				m.focusedView = treeView
			}
			m.refreshContent()
			return m, nil

		case key.Matches(msg, keys.Up):
			if m.focusedView == treeView && m.treeCursor > 0 {
				m.treeCursor--
				m.updateSelectedDevice()
				m.updateNodeCount()
				m.refreshContent()
				m.scrollUpToCursor()
			}
			return m, nil

		case key.Matches(msg, keys.Down):
			if m.focusedView == treeView && m.treeCursor < (m.nodeCount-1) {
				m.treeCursor++
				m.updateSelectedDevice()
				m.updateNodeCount()
				m.refreshContent()
				m.scrollDownToCursor()
			}
			return m, nil

		case key.Matches(msg, keys.Collapse):
			if m.focusedView == treeView {
				if node := m.getNodeAtCursor(); node != nil && len(node.Children) > 0 {
					m.collapsed[node.Key()] = true
					if m.treeCursor > 0 {
						m.treeCursor--
					}
					m.updateNodeCount()
					m.updateSelectedDevice()
					m.refreshContent()
					m.scrollUpToCursor()
				}
			}
			return m, nil

		case key.Matches(msg, keys.Expand):
			if m.focusedView == treeView {
				if node := m.getNodeAtCursor(); node != nil && len(node.Children) > 0 {
					delete(m.collapsed, node.Key())
					m.treeCursor++
					m.updateNodeCount()
					m.updateSelectedDevice()
					m.refreshContent()
					m.scrollDownToCursor()
				}
			}
			return m, nil

		case key.Matches(msg, keys.Refresh):
			lastUpdate, _ := lib.Refresh()
			m.lastUpdated = lastUpdate
		}
	}

	return m, cmd
}

// refreshContent updateChan the UI content, including status line, tree viewport, and log viewport, based on current state.
func (m *Model) refreshContent() {
	lastUpdatedString := "Last Updated: " + m.lastUpdated.Format("15:04:05")
	lastUpdatedWidth := lipgloss.Width(lastUpdatedString) + 1

	helpView := m.helpModel.View(keys)

	helpViewStyle := windowStyle.
		Width(m.windowWidth - lastUpdatedWidth).
		Align(lipgloss.Center)

	renderedHelp := helpViewStyle.Render(helpView)

	m.statusLine = statusStyle.
		Width(m.windowWidth).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, lastUpdatedString, " ", renderedHelp))

	m.recalculateDimensions(m.statusLine)
	m.treeViewport.SetContent(m.renderTree())
	m.logViewport.SetContent(m.logContent)
}

// recalculateDimensions adjusts the dimensions of the tree and log viewports based on window size and status line height.
func (m *Model) recalculateDimensions(statusLine string) {
	m.statusHeight = lipgloss.Height(statusLine)
	remainingHeight := m.windowHeight - m.statusHeight - tooltipHeight

	m.treeViewport.SetHeight(int(float64(remainingHeight)*splitRatio) - borderSpacing)
	m.treeViewport.SetWidth(m.windowWidth - borderSpacing)

	m.logViewport.SetHeight(remainingHeight - m.treeViewport.Height() - (2 * borderSpacing))
	m.logViewport.SetWidth(m.windowWidth - borderSpacing)
}

func (m *Model) formatLogContent() string {
	var sb strings.Builder
	for _, entry := range m.log {
		sb.WriteString(m.formatLogEntry(entry))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (m *Model) formatLogEntry(log lib.Log) string {
	stateString := " "
	if log.State == lib.StateRemoved {
		stateStyle = removedLogStyle
		stateString = "-"
	} else if log.State == lib.StateAdded {
		stateStyle = addedLogStyle
		stateString = "+"
	}
	lhsString := stateStyle.Render(log.Time.Format("15:04:05") + " " + stateString + " " + log.Text + " ")
	rhsString := formatSpeed(log.Speed)
	paddingSize := m.windowWidth - lipgloss.Width(rhsString) - lipgloss.Width(lhsString) - borderSpacing
	if paddingSize < 0 {
		paddingSize = 0
	}
	padding := strings.Repeat(" ", paddingSize)
	return lipgloss.JoinHorizontal(lipgloss.Left, lhsString, padding, rhsString)
}

// clampLogViewport prevents the log viewport from scrolling past the available content.
func (m *Model) clampLogViewport() {
	if m.logViewport.Height() <= 0 {
		return
	}

	contentHeight := lipgloss.Height(m.logContent)
	maxYOffset := contentHeight - m.logViewport.Height()
	if maxYOffset < 0 {
		maxYOffset = 0
	}
	if m.logViewport.YOffset() > maxYOffset {
		m.logViewport.SetYOffset(maxYOffset)
		return
	}
}
