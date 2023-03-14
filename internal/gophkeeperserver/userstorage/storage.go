package userstorage

import (
	"context"

	servermodels "yudinsv/gophkeeper/internal/gophkeeperserver/models"
	"yudinsv/gophkeeper/internal/gophkeeperserver/userstorage/memstorage"
	"yudinsv/gophkeeper/internal/gophkeeperserver/userstorage/pgstorage"
	"yudinsv/gophkeeper/internal/models"
)

type UserStorage interface {
	Ping() error
	Close() error
	AddUser(ctx context.Context, user models.User) error
	AuthenticationUser(ctx context.Context, user models.User) (bool, error)
}

func NewUserStorage(cfg servermodels.Config) (UserStorage, error) {
	var goferStorage UserStorage
	var err error
	if cfg.DataBaseURI != "" {
		goferStorage, err = pgstorage.New(cfg.DataBaseURI)
		if err != nil {
			return nil, err
		}
	} else {
		goferStorage, err = memstorage.New()
		if err != nil {
			return nil, err
		}
	}
	return goferStorage, nil
}
