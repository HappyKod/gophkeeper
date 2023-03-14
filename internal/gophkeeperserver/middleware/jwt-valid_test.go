package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"yudinsv/gophkeeper/internal/gophkeeperserver/container"
	"yudinsv/gophkeeper/internal/gophkeeperserver/models"
	"yudinsv/gophkeeper/internal/gophkeeperserver/userstorage"
	"yudinsv/gophkeeper/internal/keeperstorage"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJwtValidNoAuthHeader(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	JwtValid()(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJwtValidInvalidHeaderParts(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "invalid_header")

	JwtValid()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJwtValidInvalidBearer(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Basic token")

	JwtValid()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJwtValidInvalidToken(t *testing.T) {
	cfg := models.Config{}
	userStorage, err := userstorage.NewUserStorage(cfg)
	if err != nil {
		t.Fatal(err)
	}
	keeperStorage, err := keeperstorage.NewKeeperStorage(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err = container.BuildContainer(cfg, userStorage, keeperStorage); err != nil {
		t.Fatal("error starting container", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid_token")

	JwtValid()(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJwtValidValidToken(t *testing.T) {
	cfg := models.Config{}
	userStorage, err := userstorage.NewUserStorage(cfg)
	if err != nil {
		t.Fatal(err)
	}
	keeperStorage, err := keeperstorage.NewKeeperStorage(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err = container.BuildContainer(cfg, userStorage, keeperStorage); err != nil {
		t.Fatal("error starting container", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Replace with valid token

	c.Request, _ = http.NewRequest("GET", "/test", nil)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * 100)),
			IssuedAt:  jwt.At(time.Now())},
		Login: "test",
	})
	accessToken, err := token.SignedString([]byte(container.GetConfig().SecretKey))
	if err != nil {
		t.Fatal(err)
	}
	c.Request.Header.Set("Authorization", "Bearer "+accessToken)

	JwtValid()(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestParseTokenInvalidToken(t *testing.T) {
	token := "invalid_token"
	signingKey := []byte("secret-key")

	login, err := parseToken(token, signingKey)

	assert.EqualError(t, err, "token is malformed: token contains an invalid number of segments")
	assert.Empty(t, login)
}

func TestParseTokenInvalidSigningKey(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NSIsIm5hbWUiOiJKb2huIERvZSIsImlhdCI6MTUxNjIzOTAyMn0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	signingKey := []byte("invalid_secret_key")

	login, err := parseToken(token, signingKey)

	assert.EqualError(t, err, "token signature is invalid")
	assert.Empty(t, login)
}
