package cli

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit        key.Binding
	SwitchFocus key.Binding
	Up          key.Binding
	Down        key.Binding
	Refresh     key.Binding
	Collapse    key.Binding
	Expand      key.Binding
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
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
		{k.Quit},
		{k.SwitchFocus, k.Refresh},
		{k.Up, k.Down},
		{k.Collapse, k.Expand},
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.SwitchFocus, k.Up, k.Down, k.Refresh, k.Collapse, k.Expand}
}
