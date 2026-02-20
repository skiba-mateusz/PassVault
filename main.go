package main

import (
	"log"

	"github.com/skiba-mateusz/PassVault/db"
	"github.com/skiba-mateusz/PassVault/vault"
)

func main() {
	db, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	vault := vault.NewVault()
	vault.Setup("password")
}