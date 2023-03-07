package main

import (
	"flag"
	"log"

	"yudinsv/gophkeeper/internal/gophkeeperserver/api/handlers"
	"yudinsv/gophkeeper/internal/gophkeeperserver/api/server"
	"yudinsv/gophkeeper/internal/gophkeeperserver/container"
	"yudinsv/gophkeeper/internal/gophkeeperserver/models"
	"yudinsv/gophkeeper/internal/gophkeeperserver/userstorage"
	"yudinsv/gophkeeper/internal/keeperstorage"

	"github.com/caarlos0/env/v6"
)

func main() {
	var cfg models.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("error config read", err)
	}
	flag.Parse()
	userStorage, err := userstorage.NewUserStorage(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	keeperStorage, err := keeperstorage.NewKeeperStorage(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	if err := container.BuildContainer(cfg, userStorage, keeperStorage); err != nil {
		log.Fatalln("error starting container", err)
	}
	r := handlers.Router()
	server.NewServer(r, cfg.Address)
}
