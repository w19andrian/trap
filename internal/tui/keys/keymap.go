package keys

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Esc     key.Binding
	Confirm key.Binding

	Quit key.Binding

	NextTab key.Binding
	PrevTab key.Binding

	NextElement key.Binding
	PrevElement key.Binding

	NewItem  key.Binding
	EditItem key.Binding
	SaveItem key.Binding
	MarkItem key.Binding
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("up | k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("down | j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("left | h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("right | l", "move right"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit program"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm selection"),
	),
	NextTab: key.NewBinding(
		key.WithKeys("ctrl+tab"),
		key.WithHelp("ctrl+tab", "move to next tab"),
	),
	PrevTab: key.NewBinding(
		key.WithKeys("ctrl+shift+tab"),
		key.WithHelp("ctrl+shift+tab", "move to previous tab"),
	),
	NextElement: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "select next available element"),
	),
	PrevElement: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "select previous available element"),
	),
}
