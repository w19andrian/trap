package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"

	"repo.home.wmpandrian.dev/wmp/trap/internal/inventory"
)

type Model struct {
	tdl toDoList
	db  *inventory.Inventory
}

func InitModel(inv *inventory.Inventory) (Model, error) {
	db, err := inventory.Init("./app.db")
	if err != nil {
		return Model{}, fmt.Errorf("error initializing database: %w", err)
	}

	tdl, err := initToDoList(db)
	if err != nil {
		return Model{}, fmt.Errorf("error initializing to-do-list: %w", err)
	}

	return Model{
		tdl: tdl,
		db:  db,
	}, nil
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() tea.View {
	return tea.NewView("")
}
