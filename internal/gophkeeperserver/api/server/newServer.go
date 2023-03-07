package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const timeOutShutdownService = time.Duration(5) * time.Second

// NewServer is a function that creates and starts a new HTTP server instance using the provided Gin engine and address string.
// The function listens on the specified address for incoming requests and routes them using the Gin engine.
//
// Syntax:
// func NewServer(r *gin.Engine, addressService string)
//
// Parameters:
//
// r: A pointer to a Gin engine instance that will handle the incoming requests.
// addressService: A string representing the address on which the server will listen for incoming requests. This can be an IP address and port number combination, e.g. "127.0.0.1:8080".
// Return Values: This function does not return any values.
//
// Errors:
//
// If an error occurs during the server's ListenAndServe method, the function logs the error and exits.
// If an error occurs during the server's Shutdown method, the function logs the error.
//
// Usage:
// To use NewServer, pass in a pointer to a Gin engine and the desired address for the server to listen on. The function will start listening for incoming requests on that address, and will route them to the provided Gin engine.
func NewServer(r *gin.Engine, addressService string) {
	srv := &http.Server{
		Addr:    addressService,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error: %s\n\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), timeOutShutdownService)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln("emergency shutdown of the service", err)
	}
	log.Println("service shutdown")
}
