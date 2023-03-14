// Package handlers
// The package uses the Gin web framework for handling HTTP requests.
// The package also relies on other internal packages and models defined in the project.
package handlers

import (
	"log"
	"net/http"

	"yudinsv/gophkeeper/internal/gophkeeperserver/constans"
	"yudinsv/gophkeeper/internal/gophkeeperserver/container"

	"github.com/gin-gonic/gin"
)

// syncDataHandler handles requests for synchronizing user data.
// Handler: GET /api/v1/sync.
//
// The handler retrieves the UserID from the cookie and uses it to sync data for the user.
//
// Possible response codes:
//
// 200 - data successfully synced;
// 204 - no data to sync;
// 500 - internal server error.
func syncDataHandler(c *gin.Context) {
	storage := container.GetKeeperStorage()

	// Retrieve the UserID from the cookie
	UserID := c.Param(constans.CookeUserIDName)

	// Sync data for the user
	liteSecrets, err := storage.SyncSecret(c.Request.Context(), UserID)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, constans.ErrorWorkDataBase)
		return
	}
	if len(liteSecrets) == 0 {
		c.String(http.StatusNoContent, "")
		return
	}
	c.JSON(http.StatusOK, liteSecrets)
}
