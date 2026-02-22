package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skiba-mateusz/PassVault/db"
	"github.com/skiba-mateusz/PassVault/store"
	"github.com/skiba-mateusz/PassVault/ui"
	"github.com/skiba-mateusz/PassVault/vault"
)

func main() {
	db, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	store := store.NewStore(db)
	vault := vault.NewVault(store)

	program := tea.NewProgram(ui.NewModel(ctx, vault))
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}