package store

type Tab struct {
	ID       int64
	Name     string
	IsActive bool
}

func (c *DB) GetAllTabs() ([]Tab, error) {
	rows, err := c.db.Query("SELECT id,name,active FROM tabs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tabs := make([]Tab, 0)

	for rows.Next() {
		var tab Tab

		rows.Scan(
			&tab.ID,
			&tab.Name,
			&tab.IsActive,
		)

		tabs = append(tabs, tab)
	}

	return tabs, nil
}

func (c *DB) SaveTab(t Tab) error {
	upsertQuery := `
		INSERT INTO tasks(id,name,active)
		VALUES(?,?,?) ON CONFLICT (id) DO UPDATE
		SET name=excluded.name, active=excluded.active
	`
	if _, err := c.db.Exec(upsertQuery, t.ID, t.Name, t.IsActive); err != nil {
		return err
	}

	return nil
}
