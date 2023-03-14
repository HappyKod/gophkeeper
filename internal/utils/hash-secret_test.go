package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	// Generate a random 32-byte key and 16-byte IV
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	iv := make([]byte, 16)
	if _, err := rand.Read(iv); err != nil {
		t.Fatal(err)
	}

	// Test encryption and decryption of empty plaintext
	plaintext := []byte{}
	ciphertext, err := Encrypt(plaintext, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err := Decrypt(ciphertext, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Encryption and decryption of empty plaintext failed")
	}

	// Test encryption and decryption of short plaintext
	plaintext = []byte("hello world")
	ciphertext, err = Encrypt(plaintext, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err = Decrypt(ciphertext, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Encryption and decryption of short plaintext failed")
	}

	// Test encryption and decryption of long plaintext
	plaintext = make([]byte, 100000)
	if _, err := rand.Read(plaintext); err != nil {
		t.Fatal(err)
	}
	ciphertext, err = Encrypt(plaintext, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err = Decrypt(ciphertext, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Encryption and decryption of long plaintext failed")
	}
}

func TestPadUnpad(t *testing.T) {
	// Test padding and unpadding of empty data
	data := []byte{}
	padded := pad(data, aes.BlockSize)
	unpadded := unpad(padded)
	if !bytes.Equal(data, unpadded) {
		t.Errorf("Padding and unpadding of empty data failed")
	}

	// Test padding and unpadding of data with length equal to block size
	data = []byte("hello world")
	padded = pad(data, aes.BlockSize)
	unpadded = unpad(padded)
	if !bytes.Equal(data, unpadded) {
		t.Errorf("Padding and unpadding of data with length equal to block size failed")
	}

	// Test padding and unpadding of data with length less than block size
	data = []byte("hello")
	padded = pad(data, aes.BlockSize)
	unpadded = unpad(padded)
	if !bytes.Equal(data, unpadded) {
		t.Errorf("Padding and unpadding of data with length less than block size failed")
	}

	// Test padding and unpadding of data with length greater than block size
	data = make([]byte, 100)
	if _, err := rand.Read(data); err != nil {
		t.Fatal(err)
	}
	padded = pad(data, aes.BlockSize)
	unpadded = unpad(padded)
	if !bytes.Equal(data, unpadded) {
		t.Errorf("Padding and unpadding of data with length greater than block size failed")
	}
}
