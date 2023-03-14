// Package utils provides helper functions.
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GeneratorStringUUID generates a unique uuid.
func GeneratorStringUUID() string {
	return uuid.New().String()
}

// ValidContentType validates ContentType of request header.
func ValidContentType(c *gin.Context, ContentType string) bool {
	if c.GetHeader("Content-Type") != ContentType {
		c.String(http.StatusUnsupportedMediaType, "invalid header, expected %s", ContentType)
		return false
	}
	return true
}
