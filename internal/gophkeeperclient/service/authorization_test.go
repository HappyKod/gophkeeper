package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"yudinsv/gophkeeper/internal/models"
	mock "yudinsv/gophkeeper/mocks"

	"github.com/golang/mock/gomock"
)

func TestNewAuthorizationer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockClienter(ctrl)

	// Call the method being tested
	authorizationer := NewAuthorizationer(mockClient, "http://localhost:8080")

	// Check if the Authorizationer was created successfully
	if authorizationer == nil {
		t.Error("NewAuthorizationer failed to create a Authorizationer instance")
	}
}

func TestAuthorization_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockClienter(ctrl)
	response := http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))}
	mockClient.EXPECT().Get("http://localhost:8080/ping").Return(
		&response,
		nil,
	)
	// Call the method being tested
	authorizationer := NewAuthorizationer(mockClient, "http://localhost:8080")
	err := authorizationer.Ping()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthorization_Authorization(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	user := models.User{
		Login:    "1",
		Password: "test",
	}
	marshal, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}
	mockClient := mock.NewMockClienter(ctrl)
	response := http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))}
	mockClient.EXPECT().Post("http://localhost:8080/api/v1/login", "application/json", bytes.NewReader(marshal)).Return(
		&response,
		nil,
	)
	// Call the method being tested
	authorizationer := NewAuthorizationer(mockClient, "http://localhost:8080")
	err = authorizationer.Authorization(user)
	if err != nil {
		t.Fatal(err)
	}
}
