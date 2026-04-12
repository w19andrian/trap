// Package keys defines the main key bindings for trap by using Bubbles' key package
package keys

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Left  key.Binding
	Down  key.Binding
	Up    key.Binding
	Right key.Binding

	Esc     key.Binding
	Confirm key.Binding

	Quit key.Binding

	NextTab key.Binding
	PrevTab key.Binding

	NextElement key.Binding
	PrevElement key.Binding
}

var DefaultKeyMap = KeyMap{
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
	),

	Esc: key.NewBinding(
		key.WithKeys("esc"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
	),

	Quit: key.NewBinding(
		key.WithKeys("ctrl+q"),
		key.WithHelp("^C+q", "quit"),
	),

	NextTab: key.NewBinding(
		key.WithKeys("ctrl+tab"),
		key.WithHelp("^CT", "next tab"),
	),
	PrevTab: key.NewBinding(
		key.WithKeys("ctrl+shift+tab"),
		key.WithHelp("^CST", "previous tab"),
	),

	NextElement: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("^T", "next element"),
	),
	PrevElement: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("^ST", "previous element"),
	),
}
