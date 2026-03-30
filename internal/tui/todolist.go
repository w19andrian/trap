package tui

import (
	"maps"
	"slices"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"repo.home.wmpandrian.dev/wmp/trap/internal/inventory"
)

type menu int

const (
	dateMenu menu = iota
	taskMenu
)

type taskMode string

const (
	listDates taskMode = "list-dates"
	viewTasks taskMode = "view-task"
	editTask  taskMode = "edit-task"
	newTask   taskMode = "new-task"
)

type toDoList struct {
	mode     taskMode
	dates    []time.Time
	tasks    map[time.Time][]inventory.Task
	currTask inventory.Task
	index    map[taskMode]int

	inputField
}

type inputField struct {
	title       textinput.Model
	description textarea.Model
}

func initToDoList(inv *inventory.Inventory) (toDoList, error) {
	data, err := inv.GetTasks()
	if err != nil {
		return toDoList{}, err
	}

	mapTasks := make(map[time.Time][]inventory.Task)

	for _, task := range data {
		mapTasks[truncDate(task.CreatedAt)] = append(mapTasks[truncDate(task.CreatedAt)], task)
	}

	today := truncDate(time.Now())
	if _, ok := mapTasks[today]; !ok {
		mapTasks[today] = make([]inventory.Task, 0)
	}

	dates := slices.Collect(maps.Keys(mapTasks))

	slices.SortFunc(dates, func(a, b time.Time) int { return a.Compare(b) })

	slices.Reverse(dates)

	return toDoList{
		mode:  listDates,
		dates: dates,
		tasks: mapTasks,
		index: make(map[taskMode]int),
		inputField: inputField{
			title:       textinput.New(),
			description: textarea.New(),
		},
	}, nil
}

func (tdl *toDoList) update(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	inputTitle, cmd := tdl.title.Update(msg)
	cmds = append(cmds, cmd)

	inputDesc, cmd := tdl.description.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch tdl.mode {
		case listDates:
			switch {
			case key.Matches(msg, defaultKeyMap.up):
				if tdl.index[listDates] > 0 {
					tdl.index[tdl.mode]--
				}
			case key.Matches(msg, defaultKeyMap.down):
				if tdl.index[listDates] < len(tdl.dates)-1 {
					tdl.index[listDates]++
				}
			case key.Matches(msg, defaultKeyMap.right, defaultKeyMap.confirm):
				tdl.mode = viewTasks
			case key.Matches(msg, tdlKeyMap.newItem):
				tdl.mode = newTask
			}
		case viewTasks:
			currDate := tdl.dates[tdl.index[listDates]]

			switch {
			case key.Matches(msg, defaultKeyMap.up):
				if tdl.index[viewTasks] > 0 {
					tdl.index[viewTasks]--
				}
			case key.Matches(msg, defaultKeyMap.down):
				if tdl.index[tdl.mode] < len(tdl.tasks[currDate])-1 {
					tdl.index[viewTasks]++
				}
			case key.Matches(msg, tdlKeyMap.editItem):
				tdl.currTask = tdl.tasks[currDate][tdl.index[viewTasks]]

				inputTitle.CharLimit = 100
				inputTitle.SetValue(tdl.currTask.Title)
				inputTitle.CursorEnd()
				inputTitle.Focus()

				inputDesc.CharLimit = 255
				inputDesc.SetValue(tdl.currTask.Description)
				inputDesc.CursorEnd()
				inputDesc.Blur()

				tdl.mode = editTask
			case key.Matches(msg, defaultKeyMap.esc):
				tdl.currTask = inventory.Task{}
				tdl.mode = listDates
			}

		}
	}
	return tea.Batch(cmds...)
}

func truncDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func getDateString(t time.Time) string { return t.Format("02 Jan 2006") }
