// Package service implements a Syncer interface for synchronizing secrets between a client and a server.
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"yudinsv/gophkeeper/internal/gophkeeperclient/constatns"
	"yudinsv/gophkeeper/internal/keeperstorage"
	"yudinsv/gophkeeper/internal/models"

	"github.com/google/uuid"
	"github.com/pterm/pterm"
)

// Syncer interface has several methods, including Sync() for syncing secrets,
// Ping() for checking connectivity, StartSync() for starting the synchronization process,
// PutService() for sending a secret to the server, and GetService() for receiving a secret from the server.
type Syncer interface {
	Sync() error
	Ping() error
	StartSync(string)
	PutService(ctx context.Context, secretID uuid.UUID) error
	GetService(ctx context.Context, secretID uuid.UUID) error
}

// NewSyncer creates a new Syncer instance with the specified storage, client, and address.
func NewSyncer(storage keeperstorage.KeeperStorage, client Clienter, address string) Syncer {
	return NewSync(storage, client, address)
}

// Sync type implements the Syncer interface and has fields for storage,
// client, clientID, and address.
type Sync struct {
	storage  keeperstorage.KeeperStorage
	client   Clienter
	clientID string
	address  string
}

// NewSync creates a new Syncer instance with the specified storage, client, and address.
func NewSync(storage keeperstorage.KeeperStorage, client Clienter, address string) *Sync {
	return &Sync{storage: storage, client: client, address: address}
}

// Ping method sends a GET request to the server to check connectivity.
// If the response is not 200 OK, an error is returned.
func (s *Sync) Ping() error {
	get, err := s.client.Get(s.address + "/ping")
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(get.Body)
	if get.StatusCode != http.StatusOK {
		return fmt.Errorf("ping failed")
	} else {
		return nil
	}
}

// StartSync starts the synchronization process by calling Ping() and then Sync() in a loop.
// If either method returns an error, the loop continues.
func (s *Sync) StartSync(clientID string) {
	s.clientID = clientID
	for {
		time.Sleep(constatns.TimeSleepSync)
		err := s.Ping()
		if err != nil {
			continue
		}
		err = s.Sync()
		if err != nil {
			pterm.Error.Println(err)
		}
	}
}

// Sync sends a GET request to the server to get a list of secrets.
// If the response is 204 No Content, there are no secrets to sync, so the method returns.
// Otherwise, it unmarshals the response into an array of models.LiteSecret structs.
func (s *Sync) Sync() error {
	get, err := s.client.Get(s.address + "/api/v1/sync")
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(get.Body)
	all, err := io.ReadAll(get.Body)
	if err != nil {
		return err
	}
	var serviceSecrets []models.LiteSecret

	if get.StatusCode != http.StatusOK && get.StatusCode != http.StatusNoContent {
		return fmt.Errorf("sync failed %s", all)
	}

	if get.StatusCode != http.StatusNoContent {
		err = json.Unmarshal(all, &serviceSecrets)
		if err != nil {
			return err
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), constatns.TimeOutSync)
	defer cancel()
	secret, err := s.storage.SyncSecret(ctx, s.clientID)
	if err != nil {
		return err
	}

	if len(secret) == 0 && len(serviceSecrets) != 0 {
		for _, v := range serviceSecrets {
			err = s.GetService(ctx, v.ID)
			if err != nil {
				return err
			}
		}
		return nil
	}

	serviceSecretsMap := make(map[uuid.UUID]models.LiteSecret)
	for _, servicelite := range serviceSecrets {
		serviceSecretsMap[servicelite.ID] = servicelite
	}

	for _, locallite := range secret {
		tmps := serviceSecretsMap[locallite.ID]
		if locallite.IsDeleted != tmps.IsDeleted {
			if locallite.IsDeleted {
				//	load in service
				err = s.PutService(ctx, locallite.ID)
				if err != nil {
					return err
				}
			} else {
				//	load in client
				err = s.GetService(ctx, tmps.ID)
				if err != nil {
					return err
				}
			}
		}
		if locallite.Ver.Sub(tmps.Ver) > (1 * time.Second) {
			//	load in service
			err = s.PutService(ctx, locallite.ID)
			if err != nil {
				return err
			}
		} else {
			if locallite.ValueHash != tmps.ValueHash || locallite.DescriptionHash != tmps.DescriptionHash {
				//	load in client
				err = s.GetService(ctx, tmps.ID)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// PutService gets a secret from local storage with the specified ID, marshals it into JSON
// and sends it to the server with a PUT request.
func (s *Sync) PutService(ctx context.Context, secretID uuid.UUID) error {
	getSecret, err := s.storage.GetSecret(ctx, secretID)
	if err != nil {
		return err
	}
	marshal, err := json.Marshal(getSecret)
	if err != nil {
		return err
	}
	resp, err := s.client.Put(s.address+"/api/v1/", "application/json", bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

// GetService sends a POST request to the server with a JSON payload containing the ID of the secret to retrieve.
// It unmarshals the response into a models.Secret struct and stores it in local storage.
func (s *Sync) GetService(ctx context.Context, secretID uuid.UUID) error {
	marshal, err := json.Marshal(secretID)
	if err != nil {
		return err
	}
	post, err := s.client.Post(s.address+"/api/v1/", "application/json", bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}
	defer func() {
		err = post.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	all, err := io.ReadAll(post.Body)
	if err != nil {
		return err
	}
	var secret models.Secret
	err = json.Unmarshal(all, &secret)
	if err != nil {
		return err
	}
	err = s.storage.PutSecret(ctx, secret)
	if err != nil {
		return err
	}
	return nil
}
