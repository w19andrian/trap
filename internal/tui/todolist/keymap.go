package todolist

import (
	"charm.land/bubbles/v2/key"
	"repo.home.wmpandrian.dev/wmp/trap/internal/tui/keys"
)

var tdlKeyMap = keys.KeyMap{
	NewItem: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "create a new task for today"),
	),
	EditItem: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit currently highlighted task"),
	),
	SaveItem: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "save task"),
	),
	MarkItem: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("space", "mark task (un)done"),
	),
}
