package service

type ClientService struct {
	AuthService     Authorizationer
	RegistryService Registrationer
	SyncService     Syncer
}
