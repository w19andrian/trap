// Package todolist is the to-do-list panel for trap
package todolist

import (
	"errors"
	"fmt"
	"maps"
	"slices"
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

type focusMode string

const (
	sideBar   focusMode = "list-dates"
	tasksView focusMode = "view-task"
	editTask  focusMode = "edit-task"
	newTask   focusMode = "new-task"
)

type inputField struct {
	title       textinput.Model
	description textarea.Model
}

type Model struct {
	mode    focusMode
	dates   []time.Time
	tasks   map[time.Time][]store.Task
	current store.Task
	focus   map[focusMode]int
	err     error
	style   *Style

	tasksViewPort viewport.Model

	inputField
}

func InitToDoList(db *store.DB) (*Model, error) {
	m := new(Model)

	if err := m.refresh(db); err != nil {
		return m, err
	}

	m.focus = make(map[focusMode]int)
	m.mode = sideBar

	m.style = DefaultStyle()

	m.title = textinput.New()
	m.title.CharLimit = 100
	m.title.SetWidth(m.style.inputField.GetWidth())
	m.title.Placeholder = "New title for your to-do-list"
	m.title.Prompt = ""

	m.description = textarea.New()
	m.description.ShowLineNumbers = false
	m.description.CharLimit = 300
	m.description.SetWidth(m.style.inputField.GetWidth() - m.style.inputField.GetHorizontalFrameSize())

	m.tasksViewPort = viewport.New()

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
			m.sideBarHandler(msg)
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

func (m *Model) sideBarHandler(msg tea.KeyPressMsg) {
	switch {
	case key.Matches(msg, tdlKeyMap.newTask):
		m.current = store.Task{}

		m.prepareInput()

		m.mode = newTask
	case key.Matches(msg, keys.DefaultKeyMap.Up):
		if m.focus[sideBar] > 0 {
			m.focus[sideBar]--
		}
	case key.Matches(msg, keys.DefaultKeyMap.Down):
		if m.focus[sideBar] < len(m.dates)-1 {
			m.focus[sideBar]++
		}
	case key.Matches(msg, keys.DefaultKeyMap.Right, keys.DefaultKeyMap.Confirm):
		currDate := m.dates[m.focus[sideBar]]

		if len(m.tasks[currDate]) != 0 {
			m.current = m.tasks[currDate][0]
			m.focus[tasksView] = 0
			m.mode = tasksView
		}
	}
}

func (m *Model) tasksViewHandler(msg tea.KeyPressMsg, db *store.DB) {
	currDate := m.dates[m.focus[sideBar]]
	switch {
	case key.Matches(msg, keys.DefaultKeyMap.Up):
		if m.focus[tasksView] > 0 {
			m.focus[tasksView]--
			m.current = m.tasks[currDate][m.focus[tasksView]]
			m.syncViewport()
		}
	case key.Matches(msg, keys.DefaultKeyMap.Down):
		if m.focus[tasksView] < len(m.tasks[currDate])-1 {
			m.focus[tasksView]++
			m.current = m.tasks[currDate][m.focus[tasksView]]
			m.syncViewport()
		}
	case key.Matches(msg, tdlKeyMap.markTask):
		m.current.IsDone = !m.current.IsDone

		if err := m.save(db); err != nil {
			m.err = err
		}

		if err := m.refresh(db); err != nil {
			m.err = err
		}

		currDate = m.dates[m.focus[sideBar]]
		m.current = m.tasks[currDate][m.focus[tasksView]]

	case key.Matches(msg, tdlKeyMap.editTask):
		m.current = m.tasks[currDate][m.focus[tasksView]]

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

		if err := m.save(db); err != nil {
			m.err = err
		}

		if err := m.refresh(db); err != nil {
			m.err = err
		}

		m.clearValues()

		m.focus[sideBar] = 0
		m.focus[tasksView] = 0

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

		if err := m.save(db); err != nil {
			m.err = err
		}

		m.clearValues()
		if err := m.refresh(db); err != nil {
			m.err = err
		}

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

	for i, date := range m.dates {
		if m.focus[sideBar] == i {
			builder.WriteString(stringNewLine(m.style.sidebarFocused.Render(getDateString(date))))
			continue
		}

		builder.WriteString(stringNewLine(getDateString(date)))
	}
	return m.style.sidebar.Render(builder.String())
}

func (m *Model) renderTasksMenu() string {
	date := m.dates[m.focus[sideBar]]

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

	if len(m.tasks[date]) == 0 {
		return m.style.tasksMenu.
			Width(m.style.tasksMenu.GetWidth() + m.style.detailsMenu.GetWidth()).
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Render("No tasks for today. Press 'n' to create new task")
	}

	builder := new(strings.Builder)
	for i, task := range m.tasks[date] {
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

	description := m.current.Description

	if description == "" {
		description = "(empty description)"
	}

	body := m.style.detailsBody.Render(
		lipgloss.Wrap(
			description,
			m.style.detailsBody.GetWidth(),
			" ",
		),
	)

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

func (m *Model) refresh(db *store.DB) error {
	data, err := db.GetTasks()
	if err != nil {
		return err
	}

	if m.tasks == nil {
		m.tasks = make(map[time.Time][]store.Task)
	}

	clear(m.tasks)

	for _, task := range data {
		m.tasks[truncDate(task.CreatedAt)] = append(m.tasks[truncDate(task.CreatedAt)], task)
	}

	today := truncDate(time.Now())
	if _, ok := m.tasks[today]; !ok {
		m.tasks[today] = make([]store.Task, 0)
	}

	m.dates = slices.Collect(maps.Keys(m.tasks))
	sortDatesDesc(m.dates)

	return nil
}

func (m *Model) prepareInput() {
	m.title.SetValue(m.current.Title)
	m.title.CursorEnd()
	m.title.Focus()

	m.description.SetValue(m.current.Description)
	m.description.CursorEnd()
	m.description.Blur()
}

func (m *Model) save(db *store.DB) error {
	if err := db.SaveTask(m.current); err != nil {
		return err
	}

	return nil
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

func sortDatesDesc(t []time.Time) {
	slices.SortFunc(t, func(a, b time.Time) int { return a.Compare(b) })
	slices.Reverse(t)
}

func getDateString(t time.Time) string {
	today := truncDate(time.Now())

	switch {
	case t.Equal(today):
		return "today"
	}
	return t.Format("02 Jan 2006")
}

func getDateTimeString(t time.Time) string {
	return t.Format(time.RFC822)
}
func stringNewLine(s string) string { return fmt.Sprintf("%s\n", s) }
