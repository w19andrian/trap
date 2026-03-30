package tui

import (
	"maps"
	"slices"
	"time"

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
)

type toDoList struct {
	mode     taskMode
	dates    []time.Time
	tasks    map[time.Time][]inventory.Task
	currTask inventory.Task
	index    map[menu]int
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
		index: make(map[menu]int),
	}, nil
}

func (tdl *toDoList) update(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	return tea.Batch(cmds...)
}

func truncDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func getDateString(t time.Time) string { return t.Format("02 Jan 2006") }
