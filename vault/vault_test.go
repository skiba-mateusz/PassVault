package vault

import (
	"testing"

	"github.com/skiba-mateusz/PassVault/store"
)

func TestSetupVault(t *testing.T) {
	store := store.NewMockStore()
	vault := NewVault(store)

	ctx := t.Context()
	username := "tester"
	password := "123456"

	if err := vault.Setup(ctx, username, password); err != nil {
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
	username := "tester"
	password := "123456"

	if err := vault.Setup(ctx, username, password); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	dek := vault.dek

	if err := vault.Unlock(ctx, password); err != nil {
		t.Fatalf("Unlock failed: %v", err)
	}

	if string(dek) != string(vault.dek) {
		t.Fatal("DEK mismatch after unlock")
	}
}

func TestAddService(t *testing.T) {
	store := store.NewMockStore()
	vault := NewVault(store)

	ctx := t.Context()
	username := "tester"
	password := "123456"

	if err := vault.Setup(ctx, username, password); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	services := []string{"facebook", "spotify"}

	for _, svc := range services {
		if err := vault.AddService(ctx, svc); err != nil {
			t.Fatalf("Add service faield: %v", err)
		}
	}

	passwords, err := vault.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(passwords) != len(services) {
		t.Fatalf("Expected %d passwords, got %d", len(services), len(passwords))
	}

	for i, svc := range services {
		if passwords[i].Service != svc {
			t.Fatalf("Expected %s, got %s", svc, passwords[i].Service)
		} 

		if len(passwords[i].EncryptedPassword) == 0 {
			t.Fatalf("Service %s has empty password", svc)
		}

		if len(passwords[i].Nonce) == 0 {
			t.Fatalf("Service %s has empty nonce", svc)
		}
	}
}