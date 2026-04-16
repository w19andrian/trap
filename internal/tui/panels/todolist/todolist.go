// Package todolist is the to-do-list panel for trap
package todolist

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"repo.home.wmpandrian.dev/wmp/trap/internal/store"
	"repo.home.wmpandrian.dev/wmp/trap/internal/tui/keys"
)

type focusMode int

const (
	sideBar focusMode = iota
	tasksView
	editTask
	newTask
)

type sideBarMenu string

const (
	filterToday   sideBarMenu = "Today"
	filterAll     sideBarMenu = "All"
	filterDone    sideBarMenu = "Done"
	filterPending sideBarMenu = "Pending"
)

func getMenus() []sideBarMenu {
	return []sideBarMenu{
		filterToday,
		filterAll,
		filterDone,
		filterPending,
	}
}

func (s sideBarMenu) String() string { return string(s) }

type inputField struct {
	title       textinput.Model
	description textarea.Model
}

type Model struct {
	mode    focusMode
	menus   []sideBarMenu
	tasks   map[sideBarMenu][]store.Task
	current store.Task
	focus   map[focusMode]int
	err     error
	style   *Style

	tasksViewPort viewport.Model

	inputField
}

func InitToDoList(db *store.DB) (*Model, error) {
	m := new(Model)

	m.focus = make(map[focusMode]int)

	m.mode = sideBar

	m.menus = getMenus()

	m.tasks = make(map[sideBarMenu][]store.Task)

	m.style = DefaultStyle()

	m.title = textinput.New()
	m.title.CharLimit = 75
	m.title.SetWidth(m.style.inputField.GetWidth() - m.style.inputField.GetHorizontalFrameSize())
	m.title.Placeholder = "New title for your to-do-list"
	m.title.Prompt = ""

	m.description = textarea.New()
	m.description.ShowLineNumbers = false
	m.description.CharLimit = 300
	m.description.SetWidth(m.style.inputField.GetWidth() - m.style.inputField.GetHorizontalFrameSize())

	m.tasksViewPort = viewport.New()

	if err := m.loadTasks(db); err != nil {
		return m, err
	}

	return m, nil
}

func (m *Model) Update(msg tea.Msg, db *store.DB) tea.Cmd {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	m.title, cmd = m.title.Update(msg)
	cmds = append(cmds, cmd)

	m.description, cmd = m.description.Update(msg)
	cmds = append(cmds, cmd)

	if _, ok := msg.(tea.KeyPressMsg); !ok {
		m.tasksViewPort, cmd = m.tasksViewPort.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.err = nil

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch m.mode {
		case sideBar:
			m.sideBarHandler(msg, db)
		case tasksView:
			m.tasksViewHandler(msg, db)
		case newTask:
			m.newTaskHandler(msg, db)
		case editTask:
			m.editTaskHandler(msg, db)
		}
	}

	return tea.Batch(cmds...)
}

func (m *Model) sideBarHandler(msg tea.KeyPressMsg, db *store.DB) {
	switch {
	case key.Matches(msg, tdlKeyMap.newTask):
		m.current = store.Task{}

		m.prepareInput()

		m.mode = newTask
	case key.Matches(msg, keys.DefaultKeyMap.Up):
		if m.focus[sideBar] > 0 {
			m.focus[sideBar]--
			m.refreshTasks(db)
		}
	case key.Matches(msg, keys.DefaultKeyMap.Down):
		if m.focus[sideBar] < len(m.menus)-1 {
			m.focus[sideBar]++
			m.refreshTasks(db)
		}
	case key.Matches(msg, keys.DefaultKeyMap.Right, keys.DefaultKeyMap.Confirm):
		currMenu := m.menus[m.focus[sideBar]]

		if len(m.tasks[currMenu]) != 0 {
			m.current = m.tasks[currMenu][m.focus[tasksView]]
			m.mode = tasksView
		}
	}
}

func (m *Model) tasksViewHandler(msg tea.KeyPressMsg, db *store.DB) {
	currMenu := m.menus[m.focus[sideBar]]

	switch {
	case key.Matches(msg, keys.DefaultKeyMap.Up):
		if m.focus[tasksView] > 0 {
			m.focus[tasksView]--
			m.current = m.tasks[currMenu][m.focus[tasksView]]
			m.syncViewport()
		}

	case key.Matches(msg, keys.DefaultKeyMap.Down):
		if m.focus[tasksView] < len(m.tasks[currMenu])-1 {
			m.focus[tasksView]++
			m.current = m.tasks[currMenu][m.focus[tasksView]]
			m.syncViewport()
		}

	case key.Matches(msg, tdlKeyMap.markTask):
		m.current.IsDone = !m.current.IsDone

		m.save(db)

		m.refreshTasks(db)

	case key.Matches(msg, tdlKeyMap.editTask):
		m.current = m.tasks[currMenu][m.focus[tasksView]]

		m.prepareInput()

		m.mode = editTask

	case key.Matches(msg, keys.DefaultKeyMap.Left, keys.DefaultKeyMap.Esc):
		m.current = store.Task{}

		m.mode = sideBar
	}
}

func (m *Model) newTaskHandler(msg tea.KeyPressMsg, db *store.DB) {
	switch {
	case key.Matches(msg, keys.DefaultKeyMap.Confirm):
		m.title.Blur()

		if m.title.Value() == "" {
			m.err = errors.New("title cannot be empty")
			m.title.Focus()
			return
		}

		now := time.Now()

		m.current.ID = time.Now().UnixMicro()
		m.current.Title = m.title.Value()
		m.current.CreatedAt = now
		m.current.LastModified = now

		m.save(db)

		m.clearValues()

		m.focus[sideBar] = 0
		m.focus[tasksView] = 0

		m.refreshTasks(db)

		m.mode = sideBar

	case key.Matches(msg, keys.DefaultKeyMap.Esc):
		m.mode = sideBar
	}
}

func (m *Model) editTaskHandler(msg tea.KeyPressMsg, db *store.DB) {
	switch {
	case key.Matches(msg, keys.DefaultKeyMap.NextElement):
		switch {
		case m.title.Focused():
			m.title.Blur()
			m.description.Focus()

		case m.description.Focused():
			m.title.Focus()
			m.description.Blur()
		}

	case key.Matches(msg, tdlKeyMap.saveTask):
		m.title.Blur()
		m.description.Blur()

		clock := time.Now()

		if m.title.Value() == "" {
			m.err = errors.New("title cannot be empty")
			m.title.SetValue(m.current.Title)
			m.title.Focus()
			return
		}

		m.current.Title = m.title.Value()
		m.current.Description = m.description.Value()
		m.current.LastModified = clock

		m.save(db)

		m.clearValues()

		m.refreshTasks(db)

		m.mode = tasksView

	case key.Matches(msg, keys.DefaultKeyMap.Esc):
		m.mode = tasksView
	}
}

func (m *Model) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.renderSideBar(),
		m.renderTasksMenu(),
		m.renderDetailsMenu(),
	)
}

func (m *Model) renderSideBar() string {
	builder := new(strings.Builder)

	for i, menu := range m.menus {
		if m.focus[sideBar] == i {
			builder.WriteString(stringNewLine(m.style.sidebarFocused.Render(menu.String())))
			continue
		}

		builder.WriteString(stringNewLine(menu.String()))
	}
	return m.style.sidebar.Render(builder.String())
}

func (m *Model) renderTasksMenu() string {
	menu := m.menus[m.focus[sideBar]]

	if m.mode == newTask {
		return m.style.tasksMenu.Render(m.style.inputField.Render(m.title.View()))
	}

	if m.mode == editTask {
		return m.style.tasksMenu.Render(lipgloss.JoinVertical(
			lipgloss.Center,
			m.style.inputField.Render(m.title.View()),
			m.style.inputField.Render(m.description.View()),
		))
	}

	if len(m.tasks[menu]) == 0 {
		return m.style.tasksMenu.
			Width(m.style.tasksMenu.GetWidth() + m.style.detailsMenu.GetWidth()).
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Render("No tasks for today. Chill day?")
	}

	builder := new(strings.Builder)
	for i, task := range m.tasks[menu] {
		done := " "
		if task.IsDone {
			done = "x"
		}
		lineFormat := fmt.Sprintf("[%s] %s", done, task.Title)
		if m.mode == tasksView && m.focus[tasksView] == i {
			builder.WriteString(stringNewLine(m.style.taskFocused.Render(lineFormat)))
			continue

		}
		builder.WriteString(stringNewLine(m.style.task.Render(lineFormat)))

	}

	m.tasksViewPort.SetContent(builder.String())

	return m.style.tasksMenu.Render(m.tasksViewPort.View())
}

func (m *Model) renderDetailsMenu() string {
	if m.mode != tasksView {
		return m.style.detailsMenu.BorderLeft(false).Render("")
	}

	header := m.style.detailsHeader.Render(
		lipgloss.Wrap(
			m.current.Title,
			m.style.detailsHeader.GetWidth(),
			" ",
		),
	)

	body := m.style.detailsBody.Render(
		lipgloss.Wrap(
			m.current.Description,
			m.style.detailsBody.GetWidth(),
			" ",
		),
	)

	if m.current.Description == "" {
		body = ""
	}

	builder := new(strings.Builder)

	fmt.Fprintf(
		builder,
		"%s: %s\n",
		m.style.FieldNameFormat.Render("Created At"),
		getDateTimeString(m.current.CreatedAt),
	)

	fmt.Fprintf(
		builder,
		"%s: %s",
		m.style.FieldNameFormat.Render("Last Modified"),
		getDateTimeString(m.current.LastModified),
	)

	footer := m.style.detailsFooter.Render(stringNewLine(builder.String()))

	return m.style.detailsMenu.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
		footer,
	))
}

func (m *Model) loadTasks(db *store.DB) error {
	data, err := db.GetTasks()
	if err != nil {
		return err
	}

	if m.tasks == nil {
		m.tasks = make(map[sideBarMenu][]store.Task)
	}

	clear(m.tasks)

	currMenu := m.menus[m.focus[sideBar]]

	for _, task := range data {
		switch {
		case currMenu == filterToday && truncDate(task.CreatedAt).Equal(truncDate(time.Now())):
			m.tasks[filterToday] = append(m.tasks[filterToday], task)
		case currMenu == filterDone && task.IsDone:
			m.tasks[currMenu] = append(m.tasks[currMenu], task)
		case currMenu == filterPending && !task.IsDone:
			m.tasks[currMenu] = append(m.tasks[currMenu], task)
		case currMenu == filterAll:
			m.tasks[currMenu] = append(m.tasks[currMenu], task)
		}
	}

	if _, ok := m.tasks[filterToday]; !ok {
		m.tasks[filterToday] = make([]store.Task, 0)
	}

	return nil
}

func (m *Model) refreshTasks(db *store.DB) {
	if err := m.loadTasks(db); err != nil {
		m.err = err
	}
}

func (m *Model) prepareInput() {
	m.title.SetValue(m.current.Title)
	m.title.CursorEnd()
	m.title.Focus()

	m.description.SetValue(m.current.Description)
	m.description.CursorEnd()
	m.description.Blur()
}

func (m *Model) save(db *store.DB) {
	if err := db.SaveTask(m.current); err != nil {
		m.err = err
	}
}

func (m *Model) syncViewport() {
	vpHeight := m.tasksViewPort.Height()
	offset := m.tasksViewPort.YOffset()
	focus := m.focus[tasksView]

	if focus < offset {
		m.tasksViewPort.SetYOffset(focus)
	} else if focus >= offset+vpHeight {
		m.tasksViewPort.SetYOffset(focus - vpHeight + 1)
	}
}

func (m *Model) clearValues() {
	m.title.SetValue("")
	m.description.SetValue("")
}

func truncDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func getDateTimeString(t time.Time) string {
	return t.Format(time.RFC822)
}

func stringNewLine(s string) string { return fmt.Sprintf("%s\n", s) }
