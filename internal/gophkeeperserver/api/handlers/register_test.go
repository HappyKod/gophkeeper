package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"yudinsv/gophkeeper/internal/gophkeeperserver/container"
	serverModels "yudinsv/gophkeeper/internal/gophkeeperserver/models"
	"yudinsv/gophkeeper/internal/gophkeeperserver/userstorage"
	"yudinsv/gophkeeper/internal/keeperstorage"
	"yudinsv/gophkeeper/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {
	// Setup test data
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", registerHandler)

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
			name:         "Successful registration",
			contentType:  "application/json",
			user:         models.User{Login: "newuser", Password: "newpass"},
			expectedCode: http.StatusPermanentRedirect,
			expectedBody: "",
		},
		{
			name:         "Invalid content type",
			contentType:  "text/plain",
			user:         models.User{Login: "newuser", Password: "newpass"},
			expectedCode: http.StatusUnsupportedMediaType,
			expectedBody: "invalid header, expected application/json",
		},
		{
			name:         "Missing login",
			contentType:  "application/json",
			user:         models.User{Password: "newpass"},
			expectedCode: http.StatusBadRequest,
			expectedBody: "error unmarshaling request body",
		},
		{
			name:         "Missing password",
			contentType:  "application/json",
			user:         models.User{Login: "newuser"},
			expectedCode: http.StatusBadRequest,
			expectedBody: "error unmarshaling request body",
		},
		{
			name:         "User already exists",
			contentType:  "application/json",
			user:         models.User{Login: "newuser", Password: "newpass"},
			expectedCode: http.StatusConflict,
			expectedBody: "there is already a user with this login",
		},
		{
			name:         "User already exists",
			contentType:  "application/json",
			user:         models.User{Login: "testuser", Password: "testpass"},
			expectedCode: http.StatusConflict,
			expectedBody: "there is already a user with this login",
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
			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonData))
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

		})
	}
}
