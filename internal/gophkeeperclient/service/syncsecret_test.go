package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"yudinsv/gophkeeper/internal/models"
	mock "yudinsv/gophkeeper/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestNewSyncer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockClienter(ctrl)
	mockStorage := mock.NewMockKeeperStorage(ctrl)

	// Call the method being tested
	syncer := NewSyncer(mockStorage, mockClient, "http://localhost:8080")

	// Check if the syncer was created successfully
	if syncer == nil {
		t.Error("NewSyncer failed to create a Syncer instance")
	}
}

func TestSync_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockClienter(ctrl)
	response := http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))}
	mockClient.EXPECT().Get("http://localhost:8080/ping").Return(
		&response,
		nil,
	)
	mockStorage := mock.NewMockKeeperStorage(ctrl)

	// Call the method being tested
	syncer := NewSyncer(mockStorage, mockClient, "http://localhost:8080")
	err := syncer.Ping()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSync_PutService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockClienter(ctrl)
	response := http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))}
	mockClient.EXPECT().Put("http://localhost:8080/api/v1/", "application/json", gomock.Any()).Return(
		&response,
		nil,
	)
	key := uuid.New()
	mockStorage := mock.NewMockKeeperStorage(ctrl)
	mockStorage.EXPECT().GetSecret(gomock.Any(), key).Return(models.Secret{}, nil)
	// Call the method being tested
	syncer := NewSyncer(mockStorage, mockClient, "http://localhost:8080")
	err := syncer.PutService(context.Background(), key)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSync_GetService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mock.NewMockClienter(ctrl)
	key := uuid.New()
	s := models.Secret{ID: key}
	marshal, err := json.Marshal(s)
	if err != nil {
		t.Error(err)
	}
	response := http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(marshal))}
	mockClient.EXPECT().Post("http://localhost:8080/api/v1/", "application/json", gomock.Any()).Return(
		&response,
		nil,
	)

	mockStorage := mock.NewMockKeeperStorage(ctrl)
	mockStorage.EXPECT().PutSecret(gomock.Any(), s).Return(nil)
	// Call the method being tested
	syncer := NewSyncer(mockStorage, mockClient, "http://localhost:8080")
	err = syncer.GetService(context.Background(), key)
	if err != nil {
		t.Fatal(err)
	}
}
