package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"time"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"repo.home.wmpandrian.dev/wmp/trap/internal/inventory"
)

type taskMode uint

const (
	listDates taskMode = iota
	listTasks
	viewMode
	createMode
	editMode
)

type tabType string

const (
	newTab   tabType = "new-tab"
	toDoList tabType = "to-do-list"
	radio    tabType = "radio"
)

type tab struct {
	tabType tabType
	name    string
}

type model struct {
	inventory *inventory.Inventory
	task      taskModel
	tabs      []tab
	tabIdx    int
	txtInput  textinput.Model
	txtArea   textarea.Model
}

type taskModel struct {
	mode     taskMode
	dates    []string
	tasks    map[string][]inventory.Task
	currTask inventory.Task
	dateIdx  int
	taskIdx  int
}

func getDateString(t time.Time) string { return t.Format(time.DateOnly) }

func initTaskModel(inv *inventory.Inventory) taskModel {
	if err := inv.Migrate(); err != nil {
		return taskModel{}
	}

	tasks, err := getTasks(inv)
	if err != nil {
		return taskModel{}
	}

	return tasks
}

func initModel(inv *inventory.Inventory) model {
	return model{
		inventory: inv,
		tabs:      []tab{{tabType: toDoList, name: "default"}},
		task:      initTaskModel(inv),
		txtInput:  textinput.New(),
		txtArea:   textarea.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	ti, cmd := m.txtInput.Update(msg)
	cmds = append(cmds, cmd)

	ta, cmd := m.txtArea.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c", "alt+Q":
			return m, tea.Quit
		case "ctrl+tab":
			if m.tabIdx < len(m.tabs)-1 {
				m.tabIdx++
			}
		case "ctrl+shift+tab":
			if m.tabIdx > 0 {
				m.tabIdx--
			}
		}
		switch m.tabs[m.tabIdx].tabType {
		case toDoList:
			app := m.task

			tasks := app.tasks[app.dates[app.dateIdx]]
			switch app.mode {
			case listDates:
				switch key {
				case "up", "k":
					if app.dateIdx > 0 {
						app.dateIdx--
					}
				case "down", "j":
					if app.dateIdx < len(app.dates)-1 {
						app.dateIdx++
					}
				case "n":
					app.mode = createMode
				}
			case viewMode:
				switch key {
				case "up", "k":
					if app.taskIdx > 0 {
						app.taskIdx--
					}
				case "down", "j":
					if app.taskIdx < len(tasks)-1 {
						app.taskIdx++
					}
				case "left", "h", "esc":
					app.mode = listDates
				case " ":
					m.task.currTask = tasks[app.taskIdx]
					m.task.currTask.IsDone = !m.task.currTask.IsDone
					m.task.currTask.LastModified = time.Now()
				case "n":
					ti.Placeholder = "Create new task"
					ti.Prompt = ">"
					ti.CharLimit = 100
					ti.Focus()

					ta.CharLimit = 255
					ta.Placeholder = "Write short description"
					ta.Blur()

					app.mode = createMode
				case "e":
					app.mode = editMode
				case "ctrl+s":
					m.inventory.SaveTask(m.task.currTask)
				}
			case createMode, editMode:
				if app.mode == createMode {
					m.task.currTask = inventory.Task{}
				}

				if app.mode == editMode {
					m.task.currTask = tasks[app.taskIdx]
				}
				switch key {
				case "tab":
					if ti.Focused() {
						ti.Blur()

						ta.Focus()
						ta.CursorEnd()
					} else {
						ta.Blur()

						ti.Focus()
						ti.CursorEnd()
					}
				case "ctrl+s":
					ti.Blur()
					ta.Blur()

					clock := time.Now()

					m.task.currTask.Title = ti.Value()
					m.task.currTask.Description = ta.Value()
					m.task.currTask.LastModified = clock

					if app.mode == createMode {
						m.task.currTask.ID = clock.UTC().UnixMicro()
						m.task.currTask.CreatedAt = clock
					}

					err := m.inventory.SaveTask(m.task.currTask)
					if err != nil {
						// TODO handle error properly, don't quit
						return m, tea.Quit
					}

					tm, err := getTasks(m.inventory)
					if err != nil {
						// TODO handle error properly, don't quit
						return m, tea.Quit
					}

					m.task.tasks = tm.tasks
					m.task.dates = tm.dates
					m.task.currTask = inventory.Task{}
				}
			}
		}
	}
	return m, tea.Batch(cmds...)
}

func getTasks(inv *inventory.Inventory) (taskModel, error) {
	tasks, err := inv.GetTasks()
	if err != nil {
		return taskModel{}, err
	}

	dateTasks := make(map[string][]inventory.Task)
	for _, v := range tasks {
		dateTasks[getDateString(v.CreatedAt)] = append(dateTasks[getDateString(v.CreatedAt)], v)
	}

	today := getDateString(time.Now())
	if _, ok := dateTasks[today]; !ok {
		dateTasks[today] = make([]inventory.Task, 0)
	}

	return taskModel{
		mode:  listDates,
		dates: slices.Collect(maps.Keys(dateTasks)),
		tasks: dateTasks,
	}, nil
}

func (m model) View() tea.View {
	s := "to-do-list\n"

	for i, dates := range m.task.dates {
		cursor := " "
		if m.task.dateIdx == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, dates)
	}

	tasks := m.task.tasks[m.task.dates[m.task.dateIdx]]
	if len(tasks) == 0 {
		s += fmt.Sprint("No tasks for today(yet). Press 'n' to create a new task\n")
	}

	for i, task := range tasks {
		cursor := " "
		done := " "
		if m.task.taskIdx == i {
			cursor = "~>"
		}

		if task.IsDone {
			done = "x"
		}
		s += fmt.Sprintf("\t%s [%s] - %s\n", cursor, done, task.Title)
	}

	return tea.NewView(s)
}

func main() {
	db, err := inventory.Init("./app.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if _, err := tea.NewProgram(initModel(db)).Run(); err != nil {
		os.Exit(1)
	}
}
