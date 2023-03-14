// Package constans defines project errors.
package constans

import "errors"

// Errors string guide
const (
	ErrorWorkDataBase  = "error working with the database"
	ErrorUnmarshalBody = "error unmarshaling request body"
)

// ErrorNoUNIQUE occurs when a value is not unique.
var ErrorNoUNIQUE = errors.New("value is not unique")
