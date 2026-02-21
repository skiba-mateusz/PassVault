package vault

import (
	"testing"

	"github.com/skiba-mateusz/PassVault/store"
)

func TestSetupVault(t *testing.T) {
	store := store.NewMockStore()
	vault := NewVault(store)

	ctx := t.Context()
	password := "123456"

	if err := vault.Setup(ctx, password); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	if vault.dek == nil {
		t.Fatal("DEK should not be nil after setup")
	}
}

func TestUnlockVault(t *testing.T) {
	store := store.NewMockStore()
	vault := NewVault(store)
	
	ctx := t.Context()
	password := "123456"

	if err := vault.Setup(ctx, password); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	dek := vault.dek

	if err := vault.Unlock(ctx, password); err != nil {
		t.Fatalf("Unlock failed: %v", err)
	}

	if string(dek) != string(vault.dek) {
		t.Fatalf("DEK mismatch after unlock")
	}
}