package cli

import (
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
)

type keyMap struct {
	Quit         key.Binding
	Instructions key.Binding
	SwitchFocus  key.Binding
	Up           key.Binding
	Down         key.Binding
	PageUp       key.Binding
	PageDown     key.Binding
	Refresh      key.Binding
	Collapse     key.Binding
	Expand       key.Binding
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
	Instructions: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "toggle keybindings"),
	),
	SwitchFocus: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch focus"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Collapse: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "collapse"),
	),
	Expand: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "expand"),
	),
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Instructions, k.SwitchFocus, k.Refresh, k.Up, k.Down, k.PageUp, k.PageDown, k.Collapse, k.Expand},
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Instructions}
}

// instructionsView renders the focused instruction box shown over the application.
func (m Model) instructionsView() string {
	instructionsHelp := m.helpModel
	instructionsHelp.ShowAll = true
	instructionsHelp.SetWidth(max(0, m.windowWidth-8))
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		instructionsHelp.View(keys),
	)

	boxWidth := lipgloss.Width(content)
	if boxWidth < 42 {
		boxWidth = 42
	}
	if m.windowWidth > 0 && boxWidth > m.windowWidth-8 {
		boxWidth = m.windowWidth - 8
	}
	if boxWidth < 1 {
		boxWidth = 1
	}

	box := lipgloss.NewStyle().
		Width(boxWidth).
		Padding(1, 2).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(activeNodeBorderColor).
		Foreground(tooltipTextColor).
		Render(content)

	lines := strings.SplitN(box, "\n", 2)
	boxWidth = lipgloss.Width(box)
	title := "╔═ Keybindings "
	if remaining := boxWidth - lipgloss.Width(title) - lipgloss.Width("╗"); remaining > 0 {
		title += strings.Repeat("═", remaining)
	}
	title += "╗"
	lines[0] = lipgloss.NewStyle().Foreground(activeNodeBorderColor).Render(title)
	return strings.Join(lines, "\n")
}
