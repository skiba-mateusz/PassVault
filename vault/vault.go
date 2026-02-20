package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	time 		uint32 = 1
	memory 		uint32 = 32*64
	threads 	uint8 = 4
	keyLength 	uint32 = 32
)

type Vault struct {
	dek []byte
}

func NewVault() *Vault {
	return &Vault{
		dek: nil,
	}
}

func (v *Vault) Setup(password string) error {
	salt := make([]byte, keyLength)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	masterKey := v.derieveKey([]byte(password), salt)

	dek := make([]byte, keyLength)
	if _, err := rand.Read(dek); err != nil {
		return err
	}

	encryptedDek, err := v.encrypt(dek, masterKey)
	if err != nil {
		return err
	}

	fmt.Println(encryptedDek)

	return nil
}

func (v *Vault) Unlock(password string) error {
	_ = v.derieveKey([]byte(password), nil)

	// implement after fetching config from db

	return nil
}

func (v *Vault) encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	return ciphertext, nil
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
		return "", nil
	}

	return string(plaintext), nil
}

func (v *Vault) derieveKey(password, salt []byte) []byte {
	return argon2.IDKey(password, salt, time, memory, threads, keyLength)
}