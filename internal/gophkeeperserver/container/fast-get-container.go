package container

import (
	"yudinsv/gophkeeper/internal/gophkeeperserver/models"
	"yudinsv/gophkeeper/internal/gophkeeperserver/userstorage"
	"yudinsv/gophkeeper/internal/keeperstorage"
)

func GetUserStorage() userstorage.UserStorage {
	return DiContainer.Get("userstorage").(userstorage.UserStorage)
}

func GetKeeperStorage() keeperstorage.KeeperStorage {
	return DiContainer.Get("keeperstorage").(keeperstorage.KeeperStorage)
}

func GetConfig() models.Config {
	return DiContainer.Get("server-config").(models.Config)
}
