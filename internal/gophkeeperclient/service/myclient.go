// Package service The Clienter interface and MyClient struct are used for making HTTP requests.
// The MyClient struct implements the Clienter interface and provides the implementation for the Get, Post, and Put methods.
package service

import (
	"io"
	"net/http"
)

// Clienter interface defines methods: Get and Post, and Put.
type Clienter interface {
	Get(url string) (resp *http.Response, err error)
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
	Put(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

// MyClient struct implements the Clienter interface and provides the implementation for the Get, Post, and Put methods.
type MyClient struct {
	client http.Client
	Jwt    string
}

// Get method creates a new GET request with the specified URL and sends it using the http.Client client.
// The Authorization header is set to the JWT token stored in the MyClient struct.
// It returns the HTTP response and an error if any.
func (c *MyClient) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", c.Jwt)
	return c.client.Do(req)
}

// Post method creates a new POST request with the specified URL, content type, and request body,
// and sends it using the http.Client client.
// The Authorization header is set to the JWT token stored in the MyClient struct.
// If the response contains an Authorization header, the JWT token is updated accordingly.
// It returns the HTTP response and an error if any.
func (c *MyClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", c.Jwt)
	do, err := c.client.Do(req)
	if do.Header.Get("Authorization") != "" {
		c.Jwt = do.Header.Get("Authorization")
	}
	return do, err
}

// Put  method creates a new PUT request with the specified URL, content type, and request body,
// and sends it using the http.Client client.
// The Authorization header is set to the JWT token stored in the MyClient struct.
// If the response contains an Authorization header, the JWT token is updated accordingly.
// It returns the HTTP response and an error if any.
func (c *MyClient) Put(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", c.Jwt)
	do, err := c.client.Do(req)
	if do.Header.Get("Authorization") != "" {
		c.Jwt = do.Header.Get("Authorization")
	}
	return do, err
}
