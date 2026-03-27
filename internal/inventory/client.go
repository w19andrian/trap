package inventory

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID           int64
	Title        string
	Description  string
	IsDone       bool
	CreatedAt    time.Time
	LastModified time.Time
}

type Radio struct {
	ID       int64
	Name     string
	URI      string
	IsActive bool
}

const schema = `
	CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT  NOT NULL,
	description TEXT,
	done BOOLEAN NOT NULL,
	created_at DATETIME,
	last_modified DATETIME
	);
	CREATE TABLE IF NOT EXISTS radio_providers (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	uri TEXT NOT NULL,
	active BOOLEAN NOT NULL
	)
	`

type Inventory struct {
	db *sql.DB
}

func Init(dbPath string) (*Inventory, error) {
	dsn := dbPath + "?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL&_cache_size=-64000&_foreign_keys=ON"

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return &Inventory{}, fmt.Errorf("error opening database %s: %w", dbPath, err)
	}

	if err := db.Ping(); err != nil {
		return &Inventory{}, fmt.Errorf("error connecting to database %s: %w", dbPath, err)
	}

	return &Inventory{
		db: db,
	}, nil
}

func (c *Inventory) Migrate() error {
	if _, err := c.db.Exec(schema); err != nil {
		return err
	}
	return nil
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

func (c *Inventory) Close() error {
	return c.db.Close()
}
