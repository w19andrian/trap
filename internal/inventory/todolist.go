package inventory

import "time"

type Task struct {
	ID           int64
	Title        string
	Description  string
	IsDone       bool
	CreatedAt    time.Time
	LastModified time.Time
}

func (c *Inventory) GetTasks() ([]Task, error) {
	rows, err := c.db.Query("SELECT id,title,description,done,created_at FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.IsDone,
			&task.CreatedAt,
			&task.LastModified,
		)

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (c *Inventory) SaveTask(t Task) error {
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
		t.CreatedAt,
		t.LastModified,
	); err != nil {
		return err
	}

	return nil
}
