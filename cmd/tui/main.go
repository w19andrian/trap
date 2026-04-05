package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"repo.home.wmpandrian.dev/wmp/trap/internal/inventory"
	"repo.home.wmpandrian.dev/wmp/trap/internal/tui"
)

func main() {
	db, err := inventory.Init("./app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	m, err := tui.InitModel(db)
	if err != nil {
		log.Fatal(err)
	}

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := tea.NewProgram(m)

	if _, err = p.Run(); err != nil {
		log.Fatal(err)
	}
}
