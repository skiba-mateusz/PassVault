package vault

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/skiba-mateusz/PassVault/store"
	"golang.org/x/crypto/argon2"
)

const (
	time 		uint32 = 1
	memory 		uint32 = 32*64
	threads 	uint8 = 4
	keyLength 	uint32 = 32
)

type Vault struct {
	store *store.Store
	username string
	dek []byte
}

func NewVault(store *store.Store) *Vault {
	return &Vault{
		store: store,
		dek: nil,
	}
}

func (v *Vault) IsSetup(ctx context.Context) bool {
	_, _, _, _, err := v.store.Config.Get(ctx)
	
	return err == nil
}

func (v *Vault) Setup(ctx context.Context, username, password string) error {
	salt := make([]byte, keyLength)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	masterKey := v.deriveKey([]byte(password), salt)

	dek := make([]byte, keyLength)
	if _, err := rand.Read(dek); err != nil {
		return err
	}

	encryptedDek, nonce, err := v.encrypt(dek, masterKey)
	if err != nil {
		return err
	}

	if err = v.store.Config.Save(ctx, username, salt, encryptedDek, nonce); err != nil {
		return err
	}

	v.dek = []byte(dek)

	return nil
}

func (v *Vault) Unlock(ctx context.Context, password string) error {
	username, salt, dek, nonce, err := v.store.Config.Get(ctx)
	if err != nil {
		return err
	}

	masterKey := v.deriveKey([]byte(password), salt)

	decryptedDek, err := v.decrypt(dek, nonce, masterKey)
	if err != nil {
		return err
	}

	v.dek = []byte(decryptedDek)
	v.username = username

	return nil
}

func (v *Vault) List(ctx context.Context) ([]store.Password, error) {
	passwords, err := v.store.Password.List(ctx)
	if err != nil {
		return nil, err
	}

	for i := range passwords {
		pass := &passwords[i]
		decrypted, _ := v.decrypt(pass.EncryptedPassword, pass.Nonce, v.dek)
		pass.Password = decrypted
	}

	return passwords, nil
}

func (v *Vault) AddService(ctx context.Context, service string) error {
	password, err := v.generateRadnomPassword(16)
	if err != nil {
		return err
	}

	encrypted, nonce, err := v.encrypt([]byte(password), v.dek)
	if err != nil {
		return err
	}

	if err := v.store.Password.Add(ctx, service, encrypted, nonce); err != nil {
		return err
	}

	return nil
}

func (v *Vault) DeleteService(ctx context.Context, id int64) error {
	return v.store.Password.Delete(ctx, id)
}

func (v *Vault) Username() string {
	return v.username
}

func (v *Vault) encrypt(plaintext, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	return ciphertext, nonce, nil
}

func (v *Vault) decrypt(ciphertext, nonce, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func (v *Vault) deriveKey(password, salt []byte) []byte {
	return argon2.IDKey(password, salt, time, memory, threads, keyLength)
}

func (v *Vault) generateRadnomPassword(n int) (string , error) {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	
	password := make([]byte, n)
	if _, err := rand.Read(password); err != nil {
		return "", err
	}

	for i := range password {
		password[i] = letters[int(password[i])%len(letters)]
	}

	return string(password), nil
}