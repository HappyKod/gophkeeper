package main

import (
	"log"

	"yudinsv/gophkeeper/internal/gophkeeperclient/service"
	"yudinsv/gophkeeper/internal/gophkeeperclient/window"
	"yudinsv/gophkeeper/internal/gophkeeperserver/models"
	"yudinsv/gophkeeper/internal/keeperstorage"

	"github.com/caarlos0/env/v6"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var client service.Clienter
	var cfg models.Config
	client = &service.MyClient{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("error config read", err)
	}
	keeperStorage, err := keeperstorage.NewKeeperStorage(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	authorizationer := service.NewAuthorizationer(client, cfg.Address)
	registrationer := service.NewRegistrationer(client, cfg.Address)
	syncer := service.NewSyncer(keeperStorage, client, cfg.Address)
	serviceClient := service.ClientService{
		AuthService:     authorizationer,
		RegistryService: registrationer,
		SyncService:     syncer,
	}
	err = serviceClient.AuthService.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	err = serviceClient.RegistryService.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	window.RunWindow(serviceClient, keeperStorage)
}
