package keepermemstorage

import (
	"context"
	"testing"
	"time"

	"yudinsv/gophkeeper/internal/models"
	"yudinsv/gophkeeper/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_PutSecret(t *testing.T) {
	s := NewMemoryStorage()

	err := s.PutSecret(context.Background(), models.Secret{
		ID:          uuid.New(),
		OwnerID:     "user1",
		Value:       []byte("secret1"),
		Type:        "type1",
		Description: "description1",
		IsDeleted:   false,
	})

	assert.NoError(t, err)
}

func TestMemoryStorage_GetSecret(t *testing.T) {
	s := NewMemoryStorage()
	key := uuid.New()
	err := s.PutSecret(context.Background(), models.Secret{
		ID:          key,
		OwnerID:     "user1",
		Value:       []byte("secret1"),
		Type:        "type1",
		Description: "description1",
		IsDeleted:   false,
	})

	assert.NoError(t, err)

	secret, err := s.GetSecret(context.Background(), key)

	assert.NoError(t, err)
	assert.Equal(t, "user1", secret.OwnerID)
	assert.Equal(t, []byte("secret1"), secret.Value)
	assert.Equal(t, "type1", secret.Type)
	assert.Equal(t, "description1", secret.Description)
	assert.Equal(t, false, secret.IsDeleted)
}

func TestMemoryStorage_DeleteSecret(t *testing.T) {
	s := NewMemoryStorage()
	key := uuid.New()
	err := s.PutSecret(context.Background(), models.Secret{
		ID:          key,
		OwnerID:     "user1",
		Value:       []byte("secret1"),
		Type:        "type1",
		Description: "description1",
		IsDeleted:   false,
	})

	assert.NoError(t, err)

	err = s.DeleteSecret(context.Background(), key)
	assert.NoError(t, err)
}

func TestSyncSecret(t *testing.T) {
	// Create a new MemoryStorage instance
	s := NewMemoryStorage()

	// Add some secrets to the storage
	secret1 := models.Secret{ID: uuid.New(), OwnerID: "user1", Value: []byte("secret1"), Type: "type1", Description: "desc1", IsDeleted: false, Ver: time.Now()}
	secret2 := models.Secret{ID: uuid.New(), OwnerID: "user1", Value: []byte("secret2"), Type: "type2", Description: "desc2", IsDeleted: true, Ver: time.Now()}
	secret3 := models.Secret{ID: uuid.New(), OwnerID: "user2", Value: []byte("secret3"), Type: "type1", Description: "desc3", IsDeleted: false, Ver: time.Now()}
	if err := s.PutSecret(context.Background(), secret1); err != nil {
		t.Fatal(err)
	}
	if err := s.PutSecret(context.Background(), secret2); err != nil {
		t.Fatal(err)
	}
	if err := s.PutSecret(context.Background(), secret3); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		userID string
		want   []models.LiteSecret
	}{
		{"user1", []models.LiteSecret{
			{ID: secret2.ID, ValueHash: utils.GetMD5Hash([]byte("secret2")), DescriptionHash: utils.GetMD5Hash([]byte("desc2")), IsDeleted: true, Ver: secret2.Ver},
			{ID: secret1.ID, ValueHash: utils.GetMD5Hash([]byte("secret1")), DescriptionHash: utils.GetMD5Hash([]byte("desc1")), IsDeleted: false, Ver: secret1.Ver},
		}},
		{"user2", []models.LiteSecret{
			{ID: secret3.ID, ValueHash: utils.GetMD5Hash([]byte("secret3")), DescriptionHash: utils.GetMD5Hash([]byte("desc3")), IsDeleted: false, Ver: secret3.Ver},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.userID, func(t *testing.T) {
			_, err := s.SyncSecret(context.Background(), tt.userID)
			if err != nil {
				t.Fatalf("s.SyncSecret() error = %v", err)
			}
		})
	}
}

func TestMemoryStorage_Ping(t *testing.T) {
	// Create a new instance of the MemoryStorage struct
	storage := &MemoryStorage{}

	// Call the Ping method and check that it returns nil
	err := storage.Ping()
	if err != nil {
		t.Errorf("expected nil, but got %v", err)
	}
}

func TestMemoryStorage_Close(t *testing.T) {
	// Create a new instance of the MemoryStorage struct
	storage := &MemoryStorage{}

	// Call the Close method and check that it returns nil
	err := storage.Close()
	if err != nil {
		t.Errorf("expected nil, but got %v", err)
	}
}
