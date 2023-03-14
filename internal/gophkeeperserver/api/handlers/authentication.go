// Package handlers
// The package uses the Gin web framework for handling HTTP requests,
// and the jwt-go library for generating and verifying JSON Web Tokens (JWTs) used for user authentication.
// The package also relies on other internal packages and models defined in the project.
package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"yudinsv/gophkeeper/internal/gophkeeperserver/constans"
	"yudinsv/gophkeeper/internal/gophkeeperserver/container"
	serverModels "yudinsv/gophkeeper/internal/gophkeeperserver/models"
	"yudinsv/gophkeeper/internal/gophkeeperserver/utils"
	"yudinsv/gophkeeper/internal/models"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

// authenticationHandler authenticates the user.
// Handler: POST /api/user/login.
//
// Authentication is performed using a login/password pair.
//
// Request format:
//
// POST /api/v1/login HTTP/1.1
// Content-Type: application/json
// ...
//
// {
// "login": "<login>",
// "password": "<password>"
// }
//
// Possible response codes:
//
// 200 - user successfully authenticated;
// 400 - invalid request format;
// 401 - invalid login/password pair;
// 500 - internal server error.
func authenticationHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), constans.TimeOutRequest)
	defer cancel()
	if !utils.ValidContentType(c, "application/json") {
		return
	}
	storage := container.GetUserStorage()
	var user models.User
	if err := c.Bind(&user); err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "password or username is not correct")
		return
	}
	authenticationUser, err := storage.AuthenticationUser(ctx, user)
	if err != nil {
		c.String(http.StatusInternalServerError, constans.ErrorWorkDataBase)
		return
	}
	if !authenticationUser {
		c.String(http.StatusUnauthorized, "password or username is not correct")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &serverModels.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * 100)),
			IssuedAt:  jwt.At(time.Now())},
		Login: user.Login,
	})
	accessToken, err := token.SignedString([]byte(container.GetConfig().SecretKey))
	if err != nil {
		c.String(http.StatusInternalServerError, "error token generation")
		return
	}
	c.Header("Authorization", "Bearer "+accessToken)
}
