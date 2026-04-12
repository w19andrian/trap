package store

import (
	"time"
)

type Task struct {
	ID           int64
	Title        string
	Description  string
	IsDone       bool
	CreatedAt    time.Time
	LastModified time.Time
}

func (c *DB) GetTasks() ([]Task, error) {
	rows, err := c.db.Query("SELECT id,title,description,done,created_at,last_modified FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		var ca, lm int64

		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.IsDone,
			&ca,
			&lm,
		); err != nil {
			return tasks, err
		}

		task.CreatedAt = time.Unix(ca, 0)
		task.LastModified = time.Unix(lm, 0)

		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (c *DB) SaveTask(t Task) error {
	upsertQuery := `
		INSERT INTO tasks(id,title,description,done,created_at,last_modified)
		VALUES(?,?,?,?,?,?) ON CONFLICT (id) DO UPDATE
		SET title=excluded.title, description=excluded.description,
			done=excluded.done, last_modified=excluded.last_modified
	`
	if _, err := c.db.Exec(
		upsertQuery,
		t.ID,
		t.Title,
		t.Description,
		t.IsDone,
		t.CreatedAt.Unix(),
		t.LastModified.Unix(),
	); err != nil {
		return err
	}

	return nil
}
