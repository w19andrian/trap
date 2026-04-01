package inventory

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Radio struct {
	ID       int64
	Name     string
	URI      string
	IsActive bool
}

const schema = `
	CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER PRIMARY KEY,
	title TEXT  NOT NULL,
	description TEXT,
	done BOOLEAN NOT NULL,
	created_at INTEGER,
	last_modified INTEGER
	);
	CREATE TABLE IF NOT EXISTS radio_providers (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	uri TEXT NOT NULL,
	active BOOLEAN NOT NULL
	);
	CREATE TABLE IF NOT EXISTS tabs (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
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

func (c *Inventory) Close() error {
	return c.db.Close()
}
