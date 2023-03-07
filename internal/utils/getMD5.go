// Package utils provides a collection of utility functions.

package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// GetMD5Hash computes the MD5 hash of a given byte slice and returns it as a hexadecimal string.
// This function is useful for generating unique identifiers or checksums for data.
//
// Parameters:
// - text: A byte slice containing the data to be hashed.
//
// Returns:
// - A string representing the hexadecimal value of the MD5 hash of the input data.
func GetMD5Hash(text []byte) string {
	hash := md5.Sum(text)
	return hex.EncodeToString(hash[:])
}
