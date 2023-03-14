package keepermemstorage

import (
	"context"
	"log"
	"sync"

	"yudinsv/gophkeeper/internal/constants"
	"yudinsv/gophkeeper/internal/models"
	"yudinsv/gophkeeper/internal/utils"

	"github.com/google/uuid"
)

type MemoryStorage struct {
	mu      sync.RWMutex
	secrets map[uuid.UUID]models.Secret
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		secrets: make(map[uuid.UUID]models.Secret),
	}
}
func (s *MemoryStorage) Ping() error {
	return nil
}

func (s *MemoryStorage) Close() error {
	return nil
}

// PutSecret adds a new secret to the store.
func (s *MemoryStorage) PutSecret(_ context.Context, secret models.Secret) error {
	log.Println(secret)
	s.mu.Lock()
	defer s.mu.Unlock()

	// Add the secret to the map
	s.secrets[secret.ID] = secret

	return nil
}

// GetSecret retrieves the first secret found in the store for a given secret ID.
func (s *MemoryStorage) GetSecret(_ context.Context, secretID uuid.UUID) (models.Secret, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Find the secret with the specified owner ID
	for _, secret := range s.secrets {
		if secret.ID == secretID {
			return secret, nil
		}
	}

	return models.Secret{}, constants.ErrSecretNotFound
}

// DeleteSecret removes the first secret found in the store for a given owner ID.
func (s *MemoryStorage) DeleteSecret(_ context.Context, secretID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find the secret with the specified owner ID
	for id, secret := range s.secrets {
		if secret.ID == secretID && !secret.IsDeleted {
			// Delete the secret from the map
			secret.IsDeleted = true
			s.secrets[id] = secret
			return nil
		}
	}

	return constants.ErrSecretNotFound
}

func (s *MemoryStorage) SyncSecret(_ context.Context, userID string) ([]models.LiteSecret, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var liteSecrets []models.LiteSecret
	for _, secret := range s.secrets {
		if secret.OwnerID == userID {

			liteSecrets = append(liteSecrets, models.LiteSecret{
				ID:              secret.ID,
				ValueHash:       utils.GetMD5Hash(secret.Value),
				DescriptionHash: utils.GetMD5Hash([]byte(secret.Description)),
				IsDeleted:       secret.IsDeleted,
				Ver:             secret.Ver,
			})
		}
	}

	return liteSecrets, nil
}
