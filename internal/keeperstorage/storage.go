package keeperstorage

import (
	"context"

	servermodels "yudinsv/gophkeeper/internal/gophkeeperserver/models"
	"yudinsv/gophkeeper/internal/keeperstorage/keeperpgstorage"
	keepermemstorage "yudinsv/gophkeeper/internal/keeperstorage/memstorage"

	"yudinsv/gophkeeper/internal/models"
)

type KeeperStorage interface {
	Ping() error
	Close() error
	PutSecret(ctx context.Context, secret models.Secret) error
	GetSecret(ctx context.Context, userID string) (models.Secret, error)
	DeleteSecret(ctx context.Context, userID string) error
	SyncSecret(ctx context.Context, userID string) ([]models.LiteSecret, error)
}

func NewKeeperStorage(cfg servermodels.Config) (KeeperStorage, error) {
	var goferStorage KeeperStorage
	var err error
	if cfg.DataBaseURI != "" {
		goferStorage, err = keeperpgstorage.NewPostgresStorage(cfg.DataBaseURI)
		if err != nil {
			return nil, err
		}
	} else {
		goferStorage = keepermemstorage.NewMemoryStorage()
		if err != nil {
			return nil, err
		}
	}
	goferStorage = keepermemstorage.NewMemoryStorage()
	if err != nil {
		return nil, err
	}
	return goferStorage, nil
}
