package tui

import "charm.land/bubbles/v2/key"

type keyMap struct {
	up      key.Binding
	down    key.Binding
	left    key.Binding
	right   key.Binding
	esc     key.Binding
	confirm key.Binding

	nextTab key.Binding
	prevTab key.Binding

	newItem  key.Binding
	editItem key.Binding
	delItem  key.Binding

	nextElement key.Binding
	prevElement key.Binding
}

var defaultKeyMap = keyMap{
	up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("up | k", "move up"),
	),
	down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("down | j", "move down"),
	),
	left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("left | h", "move left"),
	),
	right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("right | l", "move right"),
	),
	esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm selection"),
	),
	nextTab: key.NewBinding(
		key.WithKeys("ctrl+tab"),
		key.WithHelp("ctrl+tab", "move to next tab"),
	),
	prevTab: key.NewBinding(
		key.WithKeys("ctrl+shift+tab"),
		key.WithHelp("ctrl+shift+tab", "move to previous tab"),
	),
	nextElement: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "select next available element"),
	),
	prevElement: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "select previous available element"),
	),
}
