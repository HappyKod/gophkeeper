// Package handlers
// The package uses the Gin web framework for handling HTTP requests.
// The package also relies on other internal packages and models defined in the project.
package handlers

import (
	"errors"
	"net/http"

	"yudinsv/gophkeeper/internal/constants"
	"yudinsv/gophkeeper/internal/gophkeeperserver/constans"
	"yudinsv/gophkeeper/internal/gophkeeperserver/container"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getDataHandler handles requests for retrieving a secret.
// Handler: POST /api/v1/
//
// The handler retrieves the secretID from the request body and uses it to get the secret from the storage.
// If the secret is not found, a 204 No Content response is returned.
//
// Possible response codes:
//
// 200 - secret successfully retrieved;
// 204 - secret not found;
// 400 - invalid request body;
// 500 - internal server error.
func getDataHandler(c *gin.Context) {
	var secretID uuid.UUID
	if err := c.ShouldBindJSON(&secretID); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	storage := container.GetKeeperStorage()
	secret, err := storage.GetSecret(c.Request.Context(), secretID)
	if err != nil {
		if errors.Is(constants.ErrSecretNotFound, err) {
			c.String(http.StatusNoContent, "")
			return
		}
		c.String(http.StatusInternalServerError, constans.ErrorWorkDataBase, err.Error())
		return
	}
	c.JSON(http.StatusOK, secret)
}
