package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/skiba-mateusz/PassVault/db"
	"github.com/skiba-mateusz/PassVault/store"
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
	
	if vault.IsSetup(ctx) {
		vault.Unlock(ctx, "password")
	}

	vault.Setup(ctx, "password")
}