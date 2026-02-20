package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/skiba-mateusz/PassVault/db"
	"golang.org/x/crypto/argon2"
)

func main() {
	db, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("PassVault")

	password := "password"
	salt := "salt"
	key := derieveKey([]byte(password), []byte(salt))
	fmt.Println(key)

	plaintext := "secret"

	ciphertext, err := encrypt([]byte(plaintext), key)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(ciphertext)

	decrypted, err := decrypt(ciphertext, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(decrypted)
}

func derieveKey(password, salt []byte) []byte {
	return argon2.IDKey(password, salt, 1, 32*64, 4, 32)
}

func encrypt(plaintext, key []byte) ([]byte, error) {
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

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

func decrypt(ciphertext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}