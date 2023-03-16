// Package service implementation of an interface called Authorizationer.
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

// Authorizationer interface defines two methods: "Authorization" and "Ping".
type Authorizationer interface {
	Authorization(user models.User) error
	Ping() error
}

// NewAuthorizationer creates a new Authorizationer instance with the specified client, and address.
func NewAuthorizationer(client Clienter, address string) Authorizationer {
	return NewServiceAuthorization(client, address)
}

// Authorization type implements the Authorizationer interface and has fields for client and address.
type Authorization struct {
	client  Clienter
	address string
}

// NewServiceAuthorization creates a new Authorizationer instance.
func NewServiceAuthorization(client Clienter, address string) *Authorization {
	return &Authorization{client: client, address: address}
}

// Ping method sends a GET request to the server to check connectivity.
// If the response is not 200 OK, an error is returned.
func (s *Authorization) Ping() error {
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

// Authorization sends an HTTP POST request to the /api/v1/login endpoint with a JSON-encoded user object.
// It returns an error if the request fails or the response status code is not 200 OK.
func (s *Authorization) Authorization(user models.User) error {
	marshal, err := json.Marshal(user)
	if err != nil {
		return err
	}
	post, err := s.client.Post(s.address+"/api/v1/login", "application/json", bytes.NewReader(marshal))
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
