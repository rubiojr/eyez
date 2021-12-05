package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rubiojr/eyez/internal/db"
	"github.com/rubiojr/eyez/internal/tui"
)

func main() {
	dbPath := flag.String("db", "records.db", "Database file path")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	db.InitRODB(*dbPath)
	m := tui.NewModel()
	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
