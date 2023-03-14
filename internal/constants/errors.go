package constants

import "errors"

// ErrSecretNotFound secret not found in storage.
var ErrSecretNotFound = errors.New("secret not found")
