// Package container contains the code related to dependency injection.
package container

import (
	"yudinsv/gophkeeper/internal/gophkeeperserver/models"
	"yudinsv/gophkeeper/internal/gophkeeperserver/userstorage"
	"yudinsv/gophkeeper/internal/keeperstorage"

	"github.com/sarulabs/di"
)

// DiContainer contains the dependency injection container.
var DiContainer di.Container

// BuildContainer creates a new dependency injection container and initializes it
// with the necessary dependencies used throughout the code. It assigns the result
// to the DiContainer variable.
func BuildContainer(cfg models.Config, storage userstorage.UserStorage, keeperStorage keeperstorage.KeeperStorage) error {
	builder, err := di.NewBuilder()
	if err != nil {
		return err
	}
	if err = storage.Ping(); err != nil {
		return err
	}
	if err = keeperStorage.Ping(); err != nil {
		return err
	}
	if err = builder.Add(di.Def{
		Name:  "server-config",
		Build: func(ctn di.Container) (interface{}, error) { return cfg, nil }}); err != nil {
		return err
	}
	if err = builder.Add(di.Def{
		Name:  "userstorage",
		Build: func(ctn di.Container) (interface{}, error) { return storage, nil }}); err != nil {
		return err
	}
	if err = builder.Add(di.Def{
		Name:  "keeperstorage",
		Build: func(ctn di.Container) (interface{}, error) { return keeperStorage, nil }}); err != nil {
		return err
	}
	DiContainer = builder.Build()
	return nil
}
