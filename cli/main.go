package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

type model struct {
	windowWidth    int
	windowHeight   int
	treeViewport   viewport.Model
	logViewport    viewport.Model
	tooltip        string
	tooltipContent string
	help           help.Model
	focusedView    int
	lastUpdated    string
}

const (
	gray       = "#888888"
	white      = "#ffffff"
	hotPink    = "#ff028d"
	splitRatio = 0.7 // Ratio of tree view to log view
)

const (
	treeView = iota
	logView
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
var deviceTreePlaceHolder = tree.Root(".").
	Child("Hub 1").
	Child(
		tree.New().
			Root("Hub 2").
			Child("Device 1      300Gbps").
			Child("Device 2      300Gbps").
			Child("Device 3      300Gbps"),
	).
	Child(
		tree.New().
			Root("Hub 3").
			Child("Device 4      300Gbps").
			Child("Device 5      300Gbps"),
	)

var placeHolderContent = "Bus 001 \nGaming Mouse \nhttps://www.google.com"

var placeholderLogContent = `00:00:00 Device xyz 100000 Gbps
00:00:01 Device abc 100000 Gbps
00:00:02 Device pqr 100000 Gbps
00:00:03 Device xyz 100000 Gbps`

var placeHolderUpdated = " Last updated: 00:00:00"

// ***** End of placeholder content *****

func initialModel() model {
	treeViewport := viewport.New(0, 0)
	treeViewport.SetContent(deviceTreePlaceHolder.String())

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
	}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
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

	case tea.WindowSizeMsg:
		m.windowWidth, m.windowHeight = msg.Width, msg.Height

		helpHeight := lipgloss.Height(m.help.View(keys))
		tooltipHeight := lipgloss.Height(m.tooltip)
		remainingHeight := m.windowHeight - helpHeight - tooltipHeight

		m.treeViewport.Height = int(float64(remainingHeight)*splitRatio) - 2
		m.treeViewport.Width = m.windowWidth - 2

		m.tooltip = tooltipStyle.Width(m.windowWidth - 2).Render(m.tooltipContent)

		m.logViewport.Height = remainingHeight - m.treeViewport.Height - 4
		m.logViewport.Width = m.windowWidth - 2

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
