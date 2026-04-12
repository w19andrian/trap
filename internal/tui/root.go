// Package tui implements terminal user interface for 'trap'
package tui

import (
	"fmt"
	"image/color"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"repo.home.wmpandrian.dev/wmp/trap/internal/store"
	"repo.home.wmpandrian.dev/wmp/trap/internal/tui/keys"
	"repo.home.wmpandrian.dev/wmp/trap/internal/tui/panels/todolist"
)

const logo = `⢀⣠⣴⣶⣶⣾⣿⣿⣶⣶⣦⣤⡀
⢺⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿
⠈⢻⣿⠿⣿⣿⠏⢿⣿⡿⢿⡿⠁
⠀⠀⠙⠷⣬⣿⣶⣾⣯⡴⠋⠀⠀
⠀⠀⠀⢰⣿⠛⡟⡿⢻⡆⠀⠀⠀
⠀⠀⠀⠀⠻⡀⠀⠀⡸⠃   `

const textLogo = `▗       
▜▘▛▘▀▌▛▌
▐▖▌ █▌▙▌
      ▌ `

type Model struct {
	tdl *todolist.Model
	db  *store.DB

	style *Style
}

type Style struct {
	borderColor color.Color

	header      lipgloss.Style
	logoBlock   lipgloss.Style
	keyMapBlock lipgloss.Style
	infoBlock   lipgloss.Style

	body   lipgloss.Style
	status lipgloss.Style
}

func defaultStyle(h, w int) *Style {
	s := new(Style)

	s.borderColor = lipgloss.Color("#5D92D4")

	s.header = lipgloss.NewStyle().Height(7).Width(w)

	s.logoBlock = lipgloss.NewStyle().
		Width(min(26, w-26)).
		Height(s.header.GetHeight()).
		AlignHorizontal(lipgloss.Right).
		Padding(0, 1).
		Foreground(s.borderColor).
		AlignVertical(lipgloss.Bottom)
	s.infoBlock = lipgloss.NewStyle().
		Width(min(50, w-50)).
		Height(s.header.GetHeight()).
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Left)
	s.keyMapBlock = lipgloss.NewStyle().
		Height(s.header.GetHeight()).
		Width(w - s.infoBlock.GetWidth() - s.logoBlock.GetWidth()).
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center)

	s.status = lipgloss.NewStyle().Height(1).Width(w).Foreground(s.borderColor)

	headerVerticalFrame := max(
		0,
		s.header.GetVerticalFrameSize(),
		s.logoBlock.GetVerticalFrameSize(),
		s.infoBlock.GetVerticalFrameSize(),
		s.keyMapBlock.GetVerticalFrameSize(),
	)

	headerTotal := s.header.GetHeight() + headerVerticalFrame
	statusTotal := s.status.GetHeight() + s.status.GetVerticalFrameSize()

	s.body = lipgloss.NewStyle().
		Height(h - headerTotal - statusTotal).
		Width(w).
		Border(lipgloss.NormalBorder()).
		BorderForeground(s.borderColor)

	return s
}

func InitModel(db *store.DB) (Model, error) {
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
		style: defaultStyle(0, 0),
	}, nil
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.style = defaultStyle(msg.Height, msg.Width)
		m.tdl.UpdateStyle(
			m.style.body.GetHeight()-m.style.body.GetVerticalFrameSize(),
			m.style.body.GetWidth()-m.style.body.GetHorizontalFrameSize(),
		)

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
	header := m.style.header.Render(m.renderHeader())
	status := m.style.status.Render("STATUS")
	body := m.style.body.Render(m.tdl.View(), status)

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Center, header, body))
	v.AltScreen = true

	return v
}

func (m Model) renderHeader() string {
	logos := lipgloss.JoinHorizontal(
		lipgloss.Top,
		textLogo,
		logo,
	)
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.style.infoBlock.Render("info placeholder"),
		m.style.keyMapBlock.Render("keymap placeholder"),
		m.style.logoBlock.Render(logos),
	)
}
