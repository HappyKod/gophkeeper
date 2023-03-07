// Package constans provides constants used throughout the application.
package constans

import (
	"time"
)

const CookeUserIDName = "UserID" // Name of the user ID.

const (
	TimeOutRequest = time.Duration(5 * time.Second) //Duration to wait for a request to complete before timing out.
)
