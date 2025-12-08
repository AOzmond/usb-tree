package cli

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

type focusIndex int

// Model represents the primary structure containing application state and views.
type Model struct {
	windowWidth    int
	windowHeight   int
	treeViewport   viewport.Model
	logViewport    viewport.Model
	tooltip        string
	tooltipContent string
	help           help.Model
	focusedView    focusIndex
	lastUpdated    string
}

const (
	gray          = "#888888"
	white         = "#ffffff"
	hotPink       = "#ff028d"
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

var placeHolderUpdated = "Last updated: 00:00:00"

// ***** End of placeholder content *****

// InitialModel initializes and returns a new Model instance with values for state and views.
func InitialModel() Model {
	m := Model{
		tooltipContent: placeHolderContent,
		tooltip:        tooltipStyle.Render(placeHolderContent),
		help:           help.New(),
		focusedView:    treeView,
		lastUpdated:    placeHolderUpdated,
	}
	return m
}

// Init initializes the Model, preparing it to handle updates and rendering. It returns an optional initial command.
func (m Model) Init() tea.Cmd {
	return nil
}

// View renders the current state of the Model, combining styled views for tree, log, tooltip, and status line.
func (m Model) View() string {
	var treeStyle, logStyle lipgloss.Style

	m.recalculateDimensions()
	m.treeViewport.SetContent(deviceTreePlaceHolder.String())
	m.logViewport.SetContent(placeholderLogContent)

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

// Update processes incoming messages, updates the model state, and returns the updated model and an optional command.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.windowWidth, m.windowHeight = msg.Width, msg.Height

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

func (m *Model) recalculateDimensions() {
	helpHeight := lipgloss.Height(m.help.View(keys))
	tooltipHeight := lipgloss.Height(m.tooltip)
	remainingHeight := m.windowHeight - helpHeight - tooltipHeight

	m.treeViewport.Height = int(float64(remainingHeight)*splitRatio) - borderSpacing
	m.treeViewport.Width = m.windowWidth - borderSpacing

	m.tooltip = tooltipStyle.Width(m.windowWidth - borderSpacing).Render(m.tooltipContent)

	m.logViewport.Height = remainingHeight - m.treeViewport.Height - (2 * borderSpacing)
	m.logViewport.Width = m.windowWidth - borderSpacing

	m.help.Width = m.windowWidth
}
