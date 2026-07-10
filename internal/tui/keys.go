package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit       key.Binding
	Tab        key.Binding
	ShiftTab   key.Binding
	Tab1       key.Binding
	Tab2       key.Binding
	Tab3       key.Binding
	Up         key.Binding
	Down       key.Binding
	Sort       key.Binding
	Filter     key.Binding
	FilterDone key.Binding
	Export     key.Binding
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next tab"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev tab"),
	),
	Tab1: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "dashboard"),
	),
	Tab2: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "results"),
	),
	Tab3: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "logs"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("up/k", "scroll up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("down/j", "scroll down"),
	),
	Sort: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "sort"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	FilterDone: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "clear filter"),
	),
	Export: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "export"),
	),
}
