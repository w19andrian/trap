package todolist

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"repo.home.wmpandrian.dev/wmp/trap/internal/inventory"
	"repo.home.wmpandrian.dev/wmp/trap/internal/tui/keys"
)

type focusMode string

const (
	listDates focusMode = "list-dates"
	viewTasks focusMode = "view-task"
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
	tasks   map[time.Time][]inventory.Task
	current inventory.Task
	focus   map[focusMode]int
	err     error
	style   *Styles

	inputField
}

func InitToDoList(db *inventory.Inventory) (*Model, error) {
	m := new(Model)

	if err := m.refresh(db); err != nil {
		return m, err
	}

	m.focus = make(map[focusMode]int)
	m.mode = listDates

	m.title = textinput.New()
	m.title.CharLimit = 100
	m.title.SetWidth(70)
	m.title.Placeholder = "New title for your to-do-list"
	m.title.Prompt = ""

	m.description = textarea.New()
	m.description.ShowLineNumbers = false
	m.description.CharLimit = 300

	m.style = DefaultStyle()
	return m, nil
}

func (m *Model) Update(msg tea.Msg, db *inventory.Inventory) tea.Cmd {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	m.title, cmd = m.title.Update(msg)
	cmds = append(cmds, cmd)

	m.description, cmd = m.description.Update(msg)
	cmds = append(cmds, cmd)

	m.err = nil

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch m.mode {
		case listDates:
			switch {
			case key.Matches(msg, tdlKeyMap.NewItem):
				m.title.SetValue("")
				m.title.Focus()

				m.current = inventory.Task{}

				m.mode = newTask
			case key.Matches(msg, keys.DefaultKeyMap.Up):
				if m.focus[listDates] > 0 {
					m.focus[listDates]--
				}
			case key.Matches(msg, keys.DefaultKeyMap.Down):
				if m.focus[listDates] < len(m.dates)-1 {
					m.focus[listDates]++
				}
			case key.Matches(msg, keys.DefaultKeyMap.Right, keys.DefaultKeyMap.Confirm):
				m.mode = viewTasks
			}
		case viewTasks:
			currDate := m.dates[m.focus[listDates]]

			switch {
			case key.Matches(msg, keys.DefaultKeyMap.Up):
				if m.focus[viewTasks] > 0 {
					m.focus[viewTasks]--
				}
			case key.Matches(msg, keys.DefaultKeyMap.Down):
				if m.focus[viewTasks] < len(m.tasks[currDate])-1 {
					m.focus[viewTasks]++
				}
			case key.Matches(msg, tdlKeyMap.MarkItem):
				m.current = m.tasks[currDate][m.focus[viewTasks]]
				m.current.IsDone = !m.current.IsDone

				if err := m.save(db); err != nil {
					m.err = err
				}

				if err := m.refresh(db); err != nil {
					m.err = err
				}

				m.current = inventory.Task{}
			case key.Matches(msg, tdlKeyMap.EditItem):
				m.current = m.tasks[currDate][m.focus[viewTasks]]

				m.prepareInput()

				m.mode = editTask
			case key.Matches(msg, keys.DefaultKeyMap.Esc):
				m.current = inventory.Task{}

				m.mode = listDates
			}
		case newTask:
			switch {
			case key.Matches(msg, keys.DefaultKeyMap.Confirm):
				m.title.Blur()

				if m.title.Value() == "" {
					m.title.Focus()
					return tea.Batch(cmds...)
				}

				now := time.Now()

				m.current.ID = time.Now().UnixMicro()
				m.current.Title = m.title.Value()
				m.current.CreatedAt = now
				m.current.LastModified = now

				if err := m.save(db); err != nil {
					m.err = err
				}

				m.clearValues()
				m.refresh(db)

				m.focus[listDates] = 0
				m.focus[viewTasks] = 0

				m.mode = listDates
			case key.Matches(msg, keys.DefaultKeyMap.Esc):
				m.mode = viewTasks
			}
		case editTask:
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
			case key.Matches(msg, tdlKeyMap.SaveItem):
				m.title.Blur()
				m.description.Blur()

				clock := time.Now()

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
				m.current = inventory.Task{}

				m.mode = viewTasks
			case key.Matches(msg, keys.DefaultKeyMap.Esc):
				m.mode = viewTasks
			}
		}
	}

	return tea.Batch(cmds...)
}

func (m *Model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.renderSideBar(), m.renderTasksMenu())
}

func (m *Model) renderSideBar() string {
	builder := new(strings.Builder)

	for i, date := range m.dates {
		if m.focus[listDates] == i {
			builder.WriteString(stringNewLine(m.style.sidebarFocused.Render(getDateString(date))))
			continue
		}

		builder.WriteString(stringNewLine(getDateString(date)))
	}
	return m.style.sideBar.Render(builder.String())
}

func (m *Model) renderTasksMenu() string {
	builder := new(strings.Builder)

	date := m.dates[m.focus[listDates]]

	notDone := "🟥"
	done := "✅"

	switch m.mode {
	case newTask:
		builder.WriteString(m.style.inputField.Render(m.title.View()))
	case editTask:
		builder.WriteString(lipgloss.JoinVertical(
			lipgloss.Center,
			m.style.inputField.Render(m.title.View()),
			m.style.inputField.Render(m.description.View())),
		)
	case viewTasks:
		for i, task := range m.tasks[date] {
			marker := notDone

			if task.IsDone {
				marker = done
			}

			if m.focus[viewTasks] == i {
				builder.WriteString(stringNewLine(m.style.taskFocused.Render(fmt.Sprintf("%s %s", marker, task.Title))))
				continue
			}

			builder.WriteString(stringNewLine(m.style.task.Render(fmt.Sprintf("%s %s", marker, task.Title))))
		}
	default:
		if len(m.tasks[date]) == 0 {
			builder.WriteString(stringNewLine("No tasks for today. Press 'n' to create new task"))
		}
		for _, task := range m.tasks[date] {
			marker := notDone
			if task.IsDone {
				marker = done
			}
			builder.WriteString(stringNewLine(m.style.task.Render(fmt.Sprintf("%s %s", marker, task.Title))))
		}
	}
	return m.style.tasksMenu.Render(builder.String())
}

func (m *Model) refresh(db *inventory.Inventory) error {
	data, err := db.GetTasks()
	if err != nil {
		return err
	}

	if m.tasks == nil {
		m.tasks = make(map[time.Time][]inventory.Task)
	}

	clear(m.tasks)

	for _, task := range data {
		m.tasks[truncDate(task.CreatedAt)] = append(m.tasks[truncDate(task.CreatedAt)], task)
	}

	today := truncDate(time.Now())
	if _, ok := m.tasks[today]; !ok {
		m.tasks[today] = make([]inventory.Task, 0)
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

func (m *Model) save(db *inventory.Inventory) error {
	if err := db.SaveTask(m.current); err != nil {
		return err
	}

	return nil
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

func stringNewLine(s string) string { return fmt.Sprintf("%s\n", s) }
