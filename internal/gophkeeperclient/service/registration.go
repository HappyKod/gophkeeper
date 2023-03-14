// Package service that contains an interface named Registrationer and its implementation Register.
package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"yudinsv/gophkeeper/internal/models"
)

// Registrationer interface defines two methods: Register and Ping.
// The Register method takes an argument of type models.
// User and returns an error, while the Ping method returns an error.
type Registrationer interface {
	Register(user models.User) error
	Ping() error
}

// NewRegistrationer creates a new Registrationer instance with the specified client, and address.
func NewRegistrationer(client Clienter, address string) Registrationer {
	return NewServiceRegister(client, address)
}

// Register implementation has two fields: client of type Clienter and address of type string.
// The NewServiceRegister function returns a pointer to a Register instance that takes client and address as arguments.
type Register struct {
	client  Clienter
	address string
}

// NewServiceRegister creates a new Register instance with the specified client, and address.
func NewServiceRegister(client Clienter, address string) *Register {
	return &Register{client: client, address: address}
}

// Ping method sends a GET request to the server to check connectivity.
// If the response is not 200 OK, an error is returned.
func (s *Register) Ping() error {
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

// Register sends an HTTP POST request with a JSON-encoded models.User object to the endpoint /api/v1/register.
// If the response status code is not 200, it returns an error with the response body as the error message.
// Otherwise, it returns nil.
func (s *Register) Register(user models.User) error {
	marshal, err := json.Marshal(user)
	if err != nil {
		return err
	}
	post, err := s.client.Post(s.address+"/api/v1/register", "application/json", bytes.NewReader(marshal))
	if err != nil {
		return err
	}
	all, err := io.ReadAll(post.Body)
	if err != nil {
		return err
	}
	defer func() {
		err = post.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	if post.StatusCode != http.StatusOK {
		return fmt.Errorf(string(all))
	}
	return nil
}
