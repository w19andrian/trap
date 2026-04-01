package todolist

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

type Styles struct {
	colorBorder  color.Color
	colorFgFocus color.Color
	colorbgFocus color.Color

	height int
	width  int

	inputField lipgloss.Style

	sideBar        lipgloss.Style
	sidebarFocused lipgloss.Style

	tasksMenu   lipgloss.Style
	task        lipgloss.Style
	taskFocused lipgloss.Style
}

func DefaultStyle() *Styles {
	s := new(Styles)

	s.refresh()

	return s
}

func (m *Model) UpdateStyle(h, w int) {
	m.style.height = h
	m.style.width = w

	m.style.refresh()
}

func (s *Styles) refresh() {
	s.colorBorder = lipgloss.Color("36")

	s.colorbgFocus = lipgloss.Color("36")
	s.colorFgFocus = lipgloss.Complementary(s.colorbgFocus)

	s.inputField = lipgloss.NewStyle().BorderForeground(s.colorBorder).Width(70).BorderStyle(lipgloss.NormalBorder())

	s.sideBar = lipgloss.NewStyle().BorderForeground(s.colorBorder).Width(calculateSize(s.width, 15)).Height(s.height).BorderStyle(lipgloss.NormalBorder()).BorderRight(true).Align(lipgloss.Center)
	s.sidebarFocused = lipgloss.NewStyle().Foreground(s.colorFgFocus).Background(s.colorbgFocus).Width(s.sideBar.GetWidth() - 1).Align(lipgloss.Center)

	s.tasksMenu = lipgloss.NewStyle().BorderForeground(s.colorBorder).Width(calculateSize(s.width, 60)).Height(s.height).BorderStyle(lipgloss.NormalBorder()).BorderRight(true).MarginLeft(1).MarginTop(1)
	s.task = lipgloss.NewStyle().Width(s.tasksMenu.GetWidth() - 1).MarginBottom(0).PaddingBottom(1)
	s.taskFocused = lipgloss.NewStyle().Width(s.tasksMenu.GetWidth() - 1).Border(lipgloss.RoundedBorder()).BorderForeground(s.colorBorder).AlignVertical(0.5)
}

func calculateSize(a int, r int) int {
	return int(float32(r) / float32(100) * float32(a))
}
