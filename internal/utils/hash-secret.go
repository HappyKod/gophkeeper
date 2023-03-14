// Package utils provides functions for encrypting and decrypting data using the AES encryption algorithm with the CBC mode of operation.
// The package imports the "crypto/aes" and "crypto/cipher" packages.
package utils

import (
	"crypto/aes"
	"crypto/cipher"
)

// Encrypt takes a plaintext byte array, a key byte array, and an initialization vector (IV) byte array as inputs.
// It creates a new AES cipher using the key and pads the plaintext to a multiple of the block size.
// It then creates a new CBC cipher using the AES cipher and IV and encrypts the plaintext using the CBC cipher.
// The function returns the ciphertext byte array and any errors encountered.
func Encrypt(plaintext []byte, key []byte, iv []byte) ([]byte, error) {
	// Create a new AES cipher using the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Pad the plaintext to a multiple of the block size
	plaintext = pad(plaintext, aes.BlockSize)

	// Create a new CBC cipher using the AES cipher and IV
	mode := cipher.NewCBCEncrypter(block, iv)

	// Encrypt the plaintext and return the ciphertext
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

// Decrypt takes a ciphertext byte array, a key byte array, and an initialization vector (IV) byte array as inputs.
// It creates a new AES cipher using the key and creates a new CBC cipher using the AES cipher and IV.
// It then decrypts the ciphertext using the CBC cipher and removes any padding bytes added during encryption.
// The function returns the plaintext byte array and any errors encountered.
func Decrypt(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	// Create a new AES cipher using the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a new CBC cipher using the AES cipher and IV
	mode := cipher.NewCBCDecrypter(block, iv)

	// Decrypt the ciphertext and return the plaintext
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Unpad the plaintext to remove any padding bytes added during encryption
	plaintext = unpad(plaintext)

	return plaintext, nil
}

// pad takes a data byte array and a block size as inputs and pads the data to a multiple of the block size using
// PKCS7 padding. The function returns the padded data byte array.
func pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padBytes := make([]byte, padding)
	for i := range padBytes {
		padBytes[i] = byte(padding)
	}
	return append(data, padBytes...)
}

// unpad takes a data byte array as input and removes PKCS7 padding from the end of the array.
// The function returns the unpadded data byte array.
func unpad(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}
