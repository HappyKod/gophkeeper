package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"yudinsv/gophkeeper/internal/gophkeeperserver/container"
	serverModels "yudinsv/gophkeeper/internal/gophkeeperserver/models"
	"yudinsv/gophkeeper/internal/gophkeeperserver/userstorage"
	"yudinsv/gophkeeper/internal/keeperstorage"
	"yudinsv/gophkeeper/internal/models"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticationHandler(t *testing.T) {
	// Setup test data
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/auth", authenticationHandler)

	cfg := serverModels.Config{}
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

	// Register a test user
	user := models.User{Login: "testuser", Password: "testpass"}
	ctx := context.Background()
	err = container.GetUserStorage().AddUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	// Setup test cases
	testCases := []struct {
		name         string
		contentType  string
		user         models.User
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Successful authentication",
			contentType:  "application/json",
			user:         models.User{Login: "testuser", Password: "testpass"},
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:         "Invalid content type",
			contentType:  "text/plain",
			user:         models.User{Login: "testuser", Password: "testpass"},
			expectedCode: http.StatusUnsupportedMediaType,
			expectedBody: "invalid header, expected application/json",
		},
		{
			name:         "Missing login",
			contentType:  "application/json",
			user:         models.User{Password: "testpass", Login: ""},
			expectedCode: http.StatusBadRequest,
			expectedBody: "password or username is not correct",
		},
		{
			name:         "Missing password",
			contentType:  "application/json",
			user:         models.User{Login: "testuser", Password: ""},
			expectedCode: http.StatusBadRequest,
			expectedBody: "password or username is not correct",
		},
		{
			name:         "Wrong password",
			contentType:  "application/json",
			user:         models.User{Login: "testuser", Password: "wrongpass"},
			expectedCode: http.StatusUnauthorized,
			expectedBody: "password or username is not correct",
		},
		{
			name:         "Non-existent user",
			contentType:  "application/json",
			user:         models.User{Login: "nonexistent", Password: "password"},
			expectedCode: http.StatusUnauthorized,
			expectedBody: "password or username is not correct",
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encode user data as JSON
			jsonData, err := json.Marshal(tc.user)
			if err != nil {
				t.Fatal(err)
			}

			// Create request
			req, err := http.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", tc.contentType)

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check response status code
			assert.Equal(t, tc.expectedCode, w.Code)

			// Check response body
			assert.Equal(t, tc.expectedBody, w.Body.String())
			// Check Authorization header if authentication was successful
			if tc.expectedCode == http.StatusOK {
				authHeader := w.Header().Get("Authorization")
				assert.NotEmpty(t, authHeader, "authorization header should not be empty")
				token := strings.TrimPrefix(authHeader, "Bearer ")
				claims := &serverModels.Claims{}
				_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(container.GetConfig().SecretKey), nil
				})
				assert.NoError(t, err, "error parsing token")
				assert.Equal(t, tc.user.Login, claims.Login, "login in token should be equal to the login provided")
			}
		})
	}
}
