package handlers

import (
	"context"
	"errors"
	"net/http"

	"yudinsv/gophkeeper/internal/gophkeeperserver/constans"
	"yudinsv/gophkeeper/internal/gophkeeperserver/container"
	"yudinsv/gophkeeper/internal/gophkeeperserver/utils"
	"yudinsv/gophkeeper/internal/models"

	"github.com/gin-gonic/gin"
)

// registerHandler User registration
// Handler: POST /api/user/register.
//
// Registration is done by username/password pair. Each username must be unique.
// After successful registration, the user should be automatically authenticated.
//
// Request format:
//
// POST /api/v1/register HTTP/1.1
// Content-Type: application/json
// ...
//
//	{
//		"login": "<login>",
//		"password": "<password>"
//	}
//
// Possible response codes:
//
// 200 - user successfully registered and authenticated;
// 400 - wrong request format;
// 409 - username is already taken;
// 500 - internal server error.
func registerHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), constans.TimeOutRequest)
	defer cancel()
	if !utils.ValidContentType(c, "application/json") {
		return
	}
	storage := container.GetUserStorage()
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.String(http.StatusInternalServerError, constans.ErrorUnmarshalBody)
		return
	}
	if user.Login == "" || user.Password == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err := storage.AddUser(ctx, user)
	if err != nil {
		if errors.Is(err, constans.ErrorNoUNIQUE) {
			c.String(http.StatusConflict, "there is already a user with this login")
			return
		}
		c.String(http.StatusInternalServerError, constans.ErrorWorkDataBase)
		return
	}
	c.Redirect(http.StatusPermanentRedirect, "/api/v1/login")
}
