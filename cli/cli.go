package cli

import (
	"time"

	"github.com/AOzmond/usb-tree/lib"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focusIndex int

// Model represents the primary structure containing application state and views.
type Model struct {
	windowWidth    int
	windowHeight   int
	statusHeight   int
	statusLine     string
	updates        chan []lib.Device
	roots          []*lib.TreeNode
	collapsed      map[string]bool // tracks which nodes are collapsed by their unique key
	treeViewport   viewport.Model
	logViewport    viewport.Model
	treeCursor     int
	nodeCount      int
	selectedDevice lib.Device
	helpModel      help.Model
	focusedView    focusIndex
	lastUpdated    time.Time
}

const (
	gray          = "#888888"
	white         = "#ffffff"
	hotPink       = "#ff028d"
	orange        = "#FF5c00"
	red           = "#FF0000"
	green         = "#00FF00"
	background    = "#003a3a"
	splitRatio    = 0.7 // Ratio of tree view to log view
	borderSpacing = 2   // the space taken up by the border
	tooltipHeight = 5
)

const (
	treeView focusIndex = iota
	logView
)

var (
	activeStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color(hotPink)).
			BorderBackground(lipgloss.Color(background)).
			Background(lipgloss.Color(background))

	inactiveStyle = lipgloss.NewStyle().
			BorderForeground(lipgloss.Color(gray)).
			Border(lipgloss.DoubleBorder()).
			BorderBackground(lipgloss.Color(background)).
			Background(lipgloss.Color(background))

	tooltipStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(white)).
			Background(lipgloss.Color(background)).
			Border(lipgloss.RoundedBorder()).
			BorderBackground(lipgloss.Color(background))

	windowStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(background))
)

// ***** Placeholder content *****
// TODO: replace with real data

var placeHolderDevice = "Bus 001 \nGaming Mouse \nhttps://www.google.com"

var placeholderLogContent = `00:00:00 Device xyz 100000 Gbps
00:00:01 Device abc 100000 Gbps
00:00:02 Device pqr 100000 Gbps
00:00:03 Device xyz 100000 Gbps`

// ***** End of placeholder content *****

// InitialModel initializes and returns a new Model instance with values for state and views.
func InitialModel() Model {
	updates := make(chan []lib.Device, 1)
	helpModel := help.New()
	helpModel.Styles.FullDesc = windowStyle
	helpModel.Styles.FullKey = windowStyle
	helpModel.Styles.FullSeparator = windowStyle
	helpModel.Styles.Ellipsis = windowStyle
	helpModel.Styles.ShortDesc = windowStyle
	helpModel.Styles.ShortKey = windowStyle
	helpModel.Styles.ShortSeparator = windowStyle
	m := Model{
		selectedDevice: lib.Device{},
		helpModel:      helpModel,
		focusedView:    treeView,
		lastUpdated:    time.Now(),
		treeCursor:     0,
		updates:        updates,
		collapsed:      make(map[string]bool),
	}
	return m
}

// Init initializes the Model, preparing it to handle updates and rendering. It returns an optional initial command.
func (m Model) Init() tea.Cmd {
	lib.Init(func(devices []lib.Device) {
		m.updates <- devices
	})
	return waitForUpdate(m.updates)
}

// View renders the current state of the Model, combining styled views for tree, log, tooltip, and status line.
func (m Model) View() string {
	if m.windowWidth == 0 || m.windowHeight == 0 {
		return ""
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
			topBorderColor = lipgloss.Color(orange)
		}
		if below {
			bottomBorderColor = lipgloss.Color(orange)
		}

		treeStyle = treeStyle.Border(borderStyle, true, true, true, true).
			BorderTopForeground(topBorderColor).
			BorderBottomForeground(bottomBorderColor)
	}

	tooltip := tooltipStyle.Width(m.windowWidth - borderSpacing).Render(placeHolderDevice)

	return lipgloss.JoinVertical(lipgloss.Center, treeStyle.Render(m.treeViewport.View()), tooltip, logStyle.Render(m.logViewport.View()), m.statusLine)
}

// Update processes incoming messages, updates the model state, and returns the updated model and an optional command.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case []lib.Device:
		m.roots = lib.BuildDeviceTree(msg)
		m.updateNodeCount()
		m.refreshContent()
		return m, waitForUpdate(m.updates)

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
				m.updateNodeCount()
				m.refreshContent()
				m.scrollUpToCursor()
			}
			return m, nil

		case key.Matches(msg, keys.Down):
			if m.focusedView == treeView && m.treeCursor < (m.nodeCount-1) {
				m.treeCursor++
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
					m.refreshContent()
					m.scrollDownToCursor()
				}
			}
			return m, nil

		case key.Matches(msg, keys.Refresh):
			lastUpdate, newDevices := lib.Refresh()
			m.updates <- newDevices
			m.lastUpdated = lastUpdate
		}
	}

	return m, cmd
}

func (m *Model) refreshContent() {
	lastUpdatedString := "Last Updated: " + m.lastUpdated.Format("15:04:05")
	lastUpdatedWidth := lipgloss.Width(lastUpdatedString)

	helpView := m.helpModel.View(keys)

	helpViewStyle := lipgloss.NewStyle().
		Width(m.windowWidth - lastUpdatedWidth).
		Align(lipgloss.Center)

	renderedHelp := helpViewStyle.Render(helpView)

	m.statusLine = lipgloss.NewStyle().
		Background(lipgloss.Color(background)).
		Width(m.windowWidth).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, lastUpdatedString, renderedHelp))

	m.recalculateDimensions(m.statusLine)
	m.treeViewport.SetContent(m.renderTree())
	m.logViewport.SetContent(placeholderLogContent)
}

func (m *Model) recalculateDimensions(statusLine string) {
	m.statusHeight = lipgloss.Height(statusLine)
	remainingHeight := m.windowHeight - m.statusHeight - tooltipHeight

	m.treeViewport.Height = int(float64(remainingHeight)*splitRatio) - borderSpacing
	m.treeViewport.Width = m.windowWidth - borderSpacing

	m.logViewport.Height = remainingHeight - m.treeViewport.Height - (2 * borderSpacing)
	m.logViewport.Width = m.windowWidth - borderSpacing
}
