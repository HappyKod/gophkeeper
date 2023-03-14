package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGeneratorStringUUID(t *testing.T) {
	uuid := GeneratorStringUUID()
	assert.NotEqual(t, "", uuid)
}

func TestValidContentType(t *testing.T) {
	// Create a test Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Test case 1: valid content type
	c.Request, _ = http.NewRequest("POST", "/", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	contentType := "application/json"
	result := ValidContentType(c, contentType)
	assert.True(t, result)

	// Test case 2: invalid content type
	c.Request.Header.Set("Content-Type", "text/plain")
	contentType = "application/json"
	result = ValidContentType(c, contentType)
	assert.False(t, result)
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
	assert.Contains(t, w.Body.String(), "invalid header, expected application/json")
}
