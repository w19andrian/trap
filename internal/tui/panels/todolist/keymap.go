package todolist

import "charm.land/bubbles/v2/key"

type keyMap struct {
	newTask    key.Binding
	editTask   key.Binding
	saveTask   key.Binding
	markTask   key.Binding
	deleteTask key.Binding
}

var tdlKeyMap = keyMap{
	newTask: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "create a new task for today"),
	),
	editTask: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit currently highlighted task"),
	),
	saveTask: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "save task"),
	),
	markTask: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("space", "mark task (un)done"),
	),
	deleteTask: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "delete task"),
	),
}
