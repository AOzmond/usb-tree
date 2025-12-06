package main

import (
	"fmt"
	"os"
	"time"

	"github.com/AOzmond/usb-tree/lib"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focusIndex int

type model struct {
	windowWidth    int
	windowHeight   int
	updates        chan []lib.Device
	roots          []*lib.TreeNode
	collapsed      map[*lib.TreeNode]bool // tracks which nodes are collapsed
	treeViewport   viewport.Model
	treeContent    string
	treeCursor     int
	nodeCount      int
	tooltip        string
	tooltipContent string
	log            []lib.Log
	logContent     string
	logViewport    viewport.Model
	help           help.Model
	lastUpdated    time.Time
	focusedView    focusIndex
}

const (
	gray          = "#888888"
	white         = "#ffffff"
	hotPink       = "#ff028d"
	red           = "#FF0000"
	green         = "#00FF00"
	splitRatio    = 0.7 // Ratio of tree view to log view
	borderSpacing = 2   // the space taken up by the border

	// Tooltip colors
	skyBlue   = "#00BFFF"
	gold      = "#FFD700"
	coralRed  = "#FF6B6B"
	paleGreen = "#98FB98"
	plum      = "#DDA0DD"
)

const (
	treeView focusIndex = 0
	logView  focusIndex = 1
)

var (
	activeStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(hotPink))

	inactiveStyle = lipgloss.NewStyle().
			BorderForeground(lipgloss.Color(gray)).
			Border(lipgloss.RoundedBorder())

	tooltipStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(white)).
			Border(lipgloss.RoundedBorder())
)

func initialModel() model {
	updates := make(chan []lib.Device, 1)

	treeViewport := viewport.New(0, 0)

	logViewport := viewport.New(0, 0)
	logViewport.SetContent("")

	m := model{
		treeViewport: treeViewport,
		logViewport:  logViewport,
		help:         help.New(),
		focusedView:  treeView,
		lastUpdated:  time.Now(),
		updates:      updates,
		collapsed:    make(map[*lib.TreeNode]bool),
	}
	return m
}

func formatLastUpdated(lastUpdated time.Time) string {
	return lastUpdated.Format("15:04:05")
}

func (m model) Init() tea.Cmd {
	lib.Init(func(devices []lib.Device) {
		m.updates <- devices
	})
	return waitForUpdate(m.updates)
}

func (m model) View() string {
	var treeStyle, logStyle lipgloss.Style

	fullWidthStyle := lipgloss.NewStyle().Width(m.windowWidth - borderSpacing)

	if m.focusedView == treeView {
		treeStyle = activeStyle
		logStyle = inactiveStyle
	} else {
		treeStyle = inactiveStyle
		logStyle = activeStyle
	}

	m.tooltipContent = fullWidthStyle.Render(m.getSelectedDeviceInfo())
	m.tooltip = tooltipStyle.Render(m.tooltipContent)

	lastUpdatedString := "Last Updated: " + formatLastUpdated(m.lastUpdated)
	lastUpdatedWidth := lipgloss.Width(lastUpdatedString)

	helpView := m.help.FullHelpView(keys.FullHelp())
	helpViewStyle := lipgloss.Style{}.Width(m.windowWidth - lastUpdatedWidth).Align(lipgloss.Center)
	helpView = helpViewStyle.Render(helpView)

	statusLine := lipgloss.JoinHorizontal(lipgloss.Left, lastUpdatedString, helpView)

	return lipgloss.JoinVertical(lipgloss.Center, treeStyle.Render(m.treeViewport.View()), m.tooltip, logStyle.Render(m.logViewport.View()), statusLine)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case []lib.Device:
		m.roots = lib.BuildDeviceTree(msg)
		m.refreshTreeContent()

		m.log = lib.GetLog()
		m.logContent = m.formatLogContent()
		m.logViewport.SetContent(m.logContent)

		m.treeViewport.SetContent(m.treeContent)
		return m, waitForUpdate(m.updates)

	case tea.WindowSizeMsg:
		m.windowWidth, m.windowHeight = msg.Width, msg.Height

		helpHeight := 3
		tooltipHeight := 5
		remainingHeight := m.windowHeight - helpHeight - tooltipHeight

		m.treeViewport.Height = int(float64(remainingHeight)*splitRatio) - borderSpacing
		m.treeViewport.Width = m.windowWidth - borderSpacing

		m.tooltip = tooltipStyle.Width(m.windowWidth - borderSpacing).Render(m.tooltipContent)

		m.logViewport.Height = remainingHeight - m.treeViewport.Height - (2 * borderSpacing)
		m.logViewport.Width = m.windowWidth - borderSpacing
		m.logContent = m.formatLogContent()
		m.logViewport.SetContent(m.logContent)

		m.help.Width = m.windowWidth

		m.refreshTreeContent()

		m.treeViewport.SetContent(m.treeContent)

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
			return m, nil

		case key.Matches(msg, keys.Up):
			if m.focusedView == treeView && m.treeCursor > 0 {
				m.treeCursor--
				m.refreshTreeContent()
				m.updateViewportForCursor()
				m.treeViewport.SetContent(m.treeContent)
			}

		case key.Matches(msg, keys.Down):
			if m.focusedView == treeView && m.treeCursor < m.nodeCount-1 {
				m.treeCursor++
				m.refreshTreeContent()
				m.updateViewportForCursor()
				m.treeViewport.SetContent(m.treeContent)
			}

		case key.Matches(msg, keys.Collapse):
			if m.focusedView == treeView {
				if node := m.getNodeAtCursor(); node != nil && len(node.Children) > 0 {
					m.collapsed[node] = true
					m.refreshTreeContent()
					m.treeViewport.SetContent(m.treeContent)
				}
			}

		case key.Matches(msg, keys.Expand):
			if m.focusedView == treeView {
				if node := m.getNodeAtCursor(); node != nil && len(node.Children) > 0 {
					delete(m.collapsed, node)
					m.refreshTreeContent()
					m.treeViewport.SetContent(m.treeContent)
				}
			}

		case key.Matches(msg, keys.Refresh):
			lastUpdate, newDevices := lib.Refresh()
			m.updates <- newDevices
			m.lastUpdated = lastUpdate
		}
	}

	if m.focusedView == treeView {
		m.treeViewport, cmd = m.treeViewport.Update(msg)
	} else if m.focusedView == logView {
		m.logViewport, cmd = m.logViewport.Update(msg)
	}

	return m, cmd
}

func main() {
	teaProgram := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := teaProgram.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
