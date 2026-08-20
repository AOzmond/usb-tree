package cli

import (
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
	windowWidth         int
	windowHeight        int
	statusHeight        int
	statusLine          string
	updateChan          chan []lib.Device
	roots               []*lib.TreeNode
	collapsed           map[string]bool // tracks which nodes are collapsed by their unique key
	treeViewport        viewport.Model
	treeCursor          int
	nodeCount           int
	selectedDevice      *lib.Device
	logViewport         viewport.Model
	log                 []lib.Log
	logContent          string
	helpModel           help.Model
	focusedView         focusIndex
	lastUpdated         time.Time
	logHasNew           bool
	instructionsVisible bool
}

const (
	splitRatio        = 0.7 // Ratio of tree view to log view
	borderSpacing     = 2   // the space taken up by the border
	horizontalPadding = 1
	tooltipHeight     = 4
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
	helpModel.Styles.FullKey = windowStyle
	helpModel.Styles.FullDesc = windowStyle

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
	topBorderColor := treeStyle.GetBorderTopForeground()
	middleBorderColor := treeStyle.GetBorderBottomForeground()
	if above || below {
		borderStyle := treeStyle.GetBorderStyle()

		if above {
			topBorderColor = edgeHighlightChangeColor
		}
		if below {
			middleBorderColor = edgeHighlightChangeColor
		}

		treeStyle = treeStyle.BorderStyle(borderStyle)
	}

	// The tree and tooltip form one panel. The tree has no bottom edge; the
	// tooltip's top edge uses tee junctions to provide their shared divider.
	treeStyle = treeStyle.
		Border(treeStyle.GetBorderStyle(), true, true, false, true).
		BorderTopForeground(topBorderColor)
	tooltipBorder := treeStyle.GetBorderStyle()
	tooltipBorder.TopLeft = tooltipBorder.MiddleLeft
	tooltipBorder.TopRight = tooltipBorder.MiddleRight
	tooltipStyle := treeStyle.
		Foreground(tooltipTextColor).
		Border(tooltipBorder, true, true, true, true).
		BorderTopForeground(middleBorderColor)
	if m.logHasNew {
		borderStyle := logStyle.GetBorderStyle()
		logStyle = logStyle.Border(borderStyle, true, true, true, true).
			BorderBottomForeground(edgeHighlightChangeColor)
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
	if m.instructionsVisible {
		appContent = lipgloss.Place(m.windowWidth, m.windowHeight, lipgloss.Center, lipgloss.Center, m.instructionsView())
	}

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
		wasAtBottom := m.logViewport.AtBottom()
		previousLogCount := len(m.log)
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
		if wasAtBottom {
			m.logViewport.GotoBottom()
			m.logHasNew = false
		} else if len(m.log) > previousLogCount {
			m.logHasNew = true
		}
		m.clampLogViewport()

		return m, waitForUpdate(m.updateChan)

	case tea.WindowSizeMsg:
		wasAtBottom := m.logViewport.AtBottom()
		m.windowWidth, m.windowHeight = msg.Width, msg.Height
		m.refreshContent()
		m.logContent = m.formatLogContent()
		m.logViewport.SetContent(m.logContent)
		if wasAtBottom {
			m.logViewport.GotoBottom()
		}
		m.clampLogViewport()
		m.scrollToCursor()
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, keys.Quit) {
			return m, tea.Quit
		}
		if key.Matches(msg, keys.Instructions) {
			m.instructionsVisible = !m.instructionsVisible
			return m, nil
		}
		if m.instructionsVisible {
			return m, nil
		}

		switch {
		case key.Matches(msg, keys.SwitchFocus):
			if m.focusedView == treeView {
				m.focusedView = logView
			} else {
				m.focusedView = treeView
			}
			m.refreshContent()
			return m, nil

		case key.Matches(msg, keys.Up):
			if m.focusedView == logView {
				m.scrollLogUp()
			} else if m.treeCursor > 0 {
				m.treeCursor--
				m.updateSelectedDevice()
				m.updateNodeCount()
				m.refreshContent()
				m.scrollUpToCursor()
			}
			return m, nil

		case key.Matches(msg, keys.Down):
			if m.focusedView == logView {
				m.scrollLogDown()
			} else if m.treeCursor < (m.nodeCount - 1) {
				m.treeCursor++
				m.updateSelectedDevice()
				m.updateNodeCount()
				m.refreshContent()
				m.scrollDownToCursor()
			}
			return m, nil

		case key.Matches(msg, keys.PageUp):
			if m.focusedView == logView {
				m.logViewport.PageUp()
				m.clampLogViewport()
			}
			return m, nil

		case key.Matches(msg, keys.PageDown):
			if m.focusedView == logView {
				m.logViewport.PageDown()
				m.updateLogScrollState()
				m.clampLogViewport()
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
