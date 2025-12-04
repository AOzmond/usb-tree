package main

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
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Collapse: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "collapse"),
	),
	Expand: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "expand"),
	),
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Quit, k.SwitchFocus, k.Refresh, k.Collapse, k.Expand},
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.SwitchFocus, k.Up, k.Down, k.Refresh, k.Collapse, k.Expand}
}
