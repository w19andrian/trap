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

	inputField      lipgloss.Style
	FieldNameFormat lipgloss.Style

	sideBar        lipgloss.Style
	sidebarFocused lipgloss.Style

	tasksMenu   lipgloss.Style
	task        lipgloss.Style
	taskFocused lipgloss.Style

	detailsMenu   lipgloss.Style
	detailsHeader lipgloss.Style
	detailsBody   lipgloss.Style
	detailsFooter lipgloss.Style
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

	m.tasksViewPort.SetWidth(m.style.tasksMenu.GetWidth())
	m.tasksViewPort.SetHeight(m.style.tasksMenu.GetHeight() - m.style.tasksMenu.GetVerticalFrameSize())
}

func (s *Styles) refresh() {
	s.colorBorder = lipgloss.Color("#5D92D4")

	s.colorbgFocus = lipgloss.Color("#FFFFFF")
	s.colorFgFocus = lipgloss.Color("#000000")

	s.inputField = lipgloss.NewStyle().BorderForeground(s.colorBorder).Width(70).BorderStyle(lipgloss.NormalBorder())
	s.FieldNameFormat = lipgloss.NewStyle().Bold(true).Underline(true)

	s.sideBar = lipgloss.NewStyle().BorderForeground(s.colorBorder).Width(calculateSize(s.width, 15)).Height(s.height).BorderStyle(lipgloss.NormalBorder()).BorderRight(true).Align(lipgloss.Center)
	s.sidebarFocused = lipgloss.NewStyle().Foreground(s.colorFgFocus).Background(s.colorbgFocus).Width(s.sideBar.GetWidth() - s.sideBar.GetHorizontalFrameSize()).Align(lipgloss.Center)

	sideBarTotal := s.sideBar.GetWidth() + s.sideBar.GetHorizontalFrameSize()

	s.tasksMenu = lipgloss.NewStyle().Width(calculateSize(s.width, 60)).MarginLeft(1).Height(s.height)
	s.task = lipgloss.NewStyle().Width(s.tasksMenu.GetWidth()).Height(1)
	s.taskFocused = lipgloss.NewStyle().Width(s.tasksMenu.GetWidth() - s.tasksMenu.GetHorizontalFrameSize()).Background(s.colorbgFocus).Foreground(s.colorFgFocus).AlignVertical(0.5).Height(1)

	tasksMenuTotal := s.tasksMenu.GetWidth() + s.tasksMenu.GetHorizontalFrameSize()

	s.detailsMenu = lipgloss.NewStyle().Width(s.width - sideBarTotal - tasksMenuTotal).Height(s.height).BorderStyle(lipgloss.NormalBorder()).BorderForeground(s.colorBorder).BorderLeft(true)
	s.detailsHeader = lipgloss.NewStyle().Height(2)
	s.detailsFooter = lipgloss.NewStyle().Height(3)
	s.detailsBody = lipgloss.NewStyle().Height(s.detailsMenu.GetHeight() - s.detailsHeader.GetHeight() - s.detailsFooter.GetHeight()).Align(lipgloss.Left)
}

func calculateSize(a int, r int) int { return int((float32(r) / 100.0) * float32(a)) }
