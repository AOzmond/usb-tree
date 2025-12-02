package main

import (
	"fmt"
	"os"

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
	treeViewport   viewport.Model
	treeContent    string
	treeCursor     int
	nodeCount      int
	logViewport    viewport.Model
	tooltip        string
	tooltipContent string
	help           help.Model
	focusedView    focusIndex
	lastUpdated    string
	updates        chan []lib.Device
	roots          []*lib.TreeNode
}

const (
	gray          = "#888888"
	white         = "#ffffff"
	hotPink       = "#ff028d"
	red           = "#FF0000"
	green         = "#00FF00"
	splitRatio    = 0.7 // Ratio of tree view to log view
	borderSpacing = 2   // the space taken up by the border
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

// ***** Placeholder content *****
// TODO: replace with real data

var placeHolderContent = "Bus 001 \nGaming Mouse \nhttps://www.google.com"

var placeholderLogContent = `00:00:00 Device xyz 100000 Gbps
00:00:01 Device abc 100000 Gbps
00:00:02 Device pqr 100000 Gbps
00:00:03 Device xyz 100000 Gbps`

var placeHolderUpdated = " Last updated: 00:00:00"

// ***** End of placeholder content *****

func initialModel() model {
	updates := make(chan []lib.Device, 1)

	treeViewport := viewport.New(0, 0)

	logViewport := viewport.New(0, 0)
	logViewport.SetContent(placeholderLogContent)

	m := model{
		treeViewport:   treeViewport,
		logViewport:    logViewport,
		tooltipContent: placeHolderContent,
		tooltip:        tooltipStyle.Render(placeHolderContent),
		help:           help.New(),
		focusedView:    treeView,
		lastUpdated:    placeHolderUpdated,
		updates:        updates,
	}
	return m
}

func (m model) Init() tea.Cmd {
	lib.Init(func(devices []lib.Device) {
		m.updates <- devices
	})
	return waitForUpdate(m.updates)
}

func (m model) View() string {
	var treeStyle, logStyle lipgloss.Style

	if m.focusedView == treeView {
		treeStyle = activeStyle
		logStyle = inactiveStyle
	} else {
		treeStyle = inactiveStyle
		logStyle = activeStyle
	}

	lastUpdatedWidth := lipgloss.Width(m.lastUpdated)

	helpView := m.help.View(keys)
	helpViewStyle := lipgloss.Style{}.Width(m.windowWidth - lastUpdatedWidth).Align(lipgloss.Center)
	helpView = helpViewStyle.Render(helpView)

	statusLine := lipgloss.JoinHorizontal(lipgloss.Left, m.lastUpdated, helpView)

	return lipgloss.JoinVertical(lipgloss.Center, treeStyle.Render(m.treeViewport.View()), m.tooltip, logStyle.Render(m.logViewport.View()), statusLine)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case []lib.Device:
		m.roots = lib.BuildDeviceTree(msg)
		m.refreshTreeContent()

		m.treeViewport.SetContent(m.treeContent)
		return m, waitForUpdate(m.updates)

	case tea.WindowSizeMsg:
		m.windowWidth, m.windowHeight = msg.Width, msg.Height

		helpHeight := lipgloss.Height(m.help.View(keys))
		tooltipHeight := lipgloss.Height(m.tooltip)
		remainingHeight := m.windowHeight - helpHeight - tooltipHeight

		m.treeViewport.Height = int(float64(remainingHeight)*splitRatio) - borderSpacing
		m.treeViewport.Width = m.windowWidth - borderSpacing

		m.tooltip = tooltipStyle.Width(m.windowWidth - borderSpacing).Render(m.tooltipContent)

		m.logViewport.Height = remainingHeight - m.treeViewport.Height - (2 * borderSpacing)
		m.logViewport.Width = m.windowWidth - borderSpacing

		m.help.Width = m.windowWidth

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
