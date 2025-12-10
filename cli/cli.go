package cli

import (
	"time"

	"github.com/AOzmond/usb-tree/lib"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/go-playground/locales"
)

type focusIndex int

// Model represents the primary structure containing application state and views.
type Model struct {
	windowWidth    int
	windowHeight   int
	treeViewport   viewport.Model
	logViewport    viewport.Model
	deviceTree     *tree.Tree
	selectedDevice lib.Device
	help           help.Model
	focusedView    focusIndex
	lastUpdated    time.Time
	translator     locales.Translator
}

const (
	gray          = "#888888"
	white         = "#ffffff"
	hotPink       = "#ff028d"
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
			BorderForeground(lipgloss.Color(hotPink))

	inactiveStyle = lipgloss.NewStyle().
			BorderForeground(lipgloss.Color(gray)).
			Border(lipgloss.DoubleBorder())

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

var placeHolderDevice = "Bus 001 \nGaming Mouse \nhttps://www.google.com"

var placeholderLogContent = `00:00:00 Device xyz 100000 Gbps
00:00:01 Device abc 100000 Gbps
00:00:02 Device pqr 100000 Gbps
00:00:03 Device xyz 100000 Gbps`

// ***** End of placeholder content *****

// InitialModel initializes and returns a new Model instance with values for state and views.
func InitialModel() Model {
	m := Model{
		deviceTree:     deviceTreePlaceHolder,
		selectedDevice: lib.Device{},
		help:           help.New(),
		focusedView:    treeView,
		lastUpdated:    time.Now(),
		translator:     getSystemLocale(),
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

	tooltip := tooltipStyle.Width(m.windowWidth - borderSpacing).Render(placeHolderDevice)

	lastUpdatedString := "Last Updated: " + m.translator.FmtTimeMedium(m.lastUpdated)
	lastUpdatedWidth := lipgloss.Width(lastUpdatedString)

	helpView := m.help.View(keys)
	helpViewStyle := lipgloss.Style{}.Width(m.windowWidth - lastUpdatedWidth).Align(lipgloss.Center)
	helpView = helpViewStyle.Render(helpView)

	statusLine := lipgloss.JoinHorizontal(lipgloss.Left, lastUpdatedString, helpView)

	return lipgloss.JoinVertical(lipgloss.Center, treeStyle.Render(m.treeViewport.View()), tooltip, logStyle.Render(m.logViewport.View()), statusLine)
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
	remainingHeight := m.windowHeight - helpHeight - tooltipHeight

	m.treeViewport.Height = int(float64(remainingHeight)*splitRatio) - borderSpacing
	m.treeViewport.Width = m.windowWidth - borderSpacing

	m.logViewport.Height = remainingHeight - m.treeViewport.Height - (2 * borderSpacing)
	m.logViewport.Width = m.windowWidth - borderSpacing

	m.help.Width = m.windowWidth
}
