// Package tui implements terminal user interface for 'trap'
package tui

import (
	"fmt"
	"image/color"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"repo.home.wmpandrian.dev/wmp/trap/internal/inventory"
	"repo.home.wmpandrian.dev/wmp/trap/internal/tui/keys"
	"repo.home.wmpandrian.dev/wmp/trap/internal/tui/todolist"
)

type Model struct {
	tdl *todolist.Model
	db  *inventory.Inventory

	style *Styles
}

type Styles struct {
	borderColor color.Color

	header lipgloss.Style
	body   lipgloss.Style
}

func defaultStyles(h, w int) *Styles {
	s := new(Styles)

	s.borderColor = lipgloss.Color("36")

	s.header = lipgloss.NewStyle().Height(getFinalHeight(h, 20)).Width(w).Padding(0).Margin(0)
	s.body = lipgloss.NewStyle().Height(h - s.header.GetHeight()).Width(w).BorderForeground(s.borderColor).Border(lipgloss.NormalBorder()).Padding(0).Margin(0)

	return s
}

func InitModel(db *inventory.Inventory) (Model, error) {
	if err := db.Migrate(); err != nil {
		return Model{}, fmt.Errorf("error migrating database: %w", err)
	}

	tdl, err := todolist.InitToDoList(db)
	if err != nil {
		return Model{}, fmt.Errorf("error initializing to-do-list: %w", err)
	}

	return Model{
		db:    db,
		tdl:   tdl,
		style: defaultStyles(0, 0),
	}, nil
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.style = defaultStyles(msg.Height, msg.Width)

		m.tdl.UpdateStyle(m.style.body.GetHeight(), m.style.body.GetWidth())
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.DefaultKeyMap.Quit):
			return m, tea.Quit
		}
	}
	cmd := m.tdl.Update(msg, m.db)
	return m, cmd
}

func (m Model) View() tea.View {
	header := m.style.header.Render("")
	body := m.style.body.Render(m.tdl.View())

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Top, header, body))
	v.AltScreen = true

	return v
}

func getFinalHeight(h int, r int) int { return int(float32(r) / float32(100) * float32(h)) }
