package keeperstorage

import (
	"context"

	servermodels "yudinsv/gophkeeper/internal/gophkeeperserver/models"
	"yudinsv/gophkeeper/internal/keeperstorage/keeperpgstorage"
	"yudinsv/gophkeeper/internal/keeperstorage/keepsqlstorage"
	keepermemstorage "yudinsv/gophkeeper/internal/keeperstorage/memstorage"

	"yudinsv/gophkeeper/internal/models"

	"github.com/google/uuid"
)

type KeeperStorage interface {
	Ping() error
	Close() error
	PutSecret(ctx context.Context, secret models.Secret) error
	GetSecret(ctx context.Context, secretID uuid.UUID) (models.Secret, error)
	DeleteSecret(ctx context.Context, secretID uuid.UUID) error
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
	} else if cfg.DbPath != "" {
		goferStorage, err = keepsqlstorage.NewSqliteStorage(cfg.DbPath)
		if err != nil {
			return nil, err
		}
	} else {
		goferStorage = keepermemstorage.NewMemoryStorage()
	}
	//goferStorage = keepermemstorage.NewMemoryStorage()
	return goferStorage, nil
}
