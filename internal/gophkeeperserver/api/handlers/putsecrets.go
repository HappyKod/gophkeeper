// Package handlers
// The package uses the Gin web framework for handling HTTP requests.
// The package also relies on other internal packages and models defined in the project.
package handlers

import (
	"log"
	"net/http"

	"yudinsv/gophkeeper/internal/gophkeeperserver/constans"
	"yudinsv/gophkeeper/internal/gophkeeperserver/container"
	"yudinsv/gophkeeper/internal/models"

	"github.com/gin-gonic/gin"
)

// putDataHandler handles requests for putting user data.
// Handler: PUT /api/v1/secret.
//
// The handler retrieves the secret from the request body and stores it in the database.
//
// Possible response codes:
//
// 200 - data successfully stored;
// 400 - bad request;
// 500 - internal server error.
func putDataHandler(c *gin.Context) {
	var secret models.Secret
	if err := c.ShouldBindJSON(&secret); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	storage := container.GetKeeperStorage()
	err := storage.PutSecret(c.Request.Context(), secret)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, constans.ErrorWorkDataBase)
		return
	}
}
