package memstorage

import (
	"context"
	"sync"

	"yudinsv/gophkeeper/internal/gophkeeperserver/constans"
	"yudinsv/gophkeeper/internal/models"

	"github.com/google/uuid"
)

type MemStorage struct {
	userCash map[uuid.UUID]models.User
	mu       *sync.RWMutex
}

func New() (*MemStorage, error) {
	return &MemStorage{
		userCash: make(map[uuid.UUID]models.User),
		mu:       new(sync.RWMutex),
	}, nil
}

func (MS *MemStorage) Ping() error {
	return nil
}

func (MS *MemStorage) Close() error {
	return nil
}

func (MS *MemStorage) AddUser(ctx context.Context, user models.User) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	MS.mu.Lock()
	defer MS.mu.Unlock()
	for _, v := range MS.userCash {
		if v.Login == user.Login {
			return constans.ErrorNoUNIQUE
		}
	}
	MS.userCash[uuid.New()] = user
	select {
	default:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (MS *MemStorage) AuthenticationUser(_ context.Context, user models.User) (bool, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	for _, v := range MS.userCash {
		if v.Login == user.Login && v.Password == user.Password {
			return true, nil
		}
	}
	return false, nil
}
