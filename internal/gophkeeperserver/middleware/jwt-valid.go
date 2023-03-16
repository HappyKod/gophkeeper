/*Package middleware

Function Name: JwtValid

Description: JwtValid is a middleware function that performs JWT token validation for each incoming request. If the token is not present in the Authorization header or is invalid, it aborts the request and returns the appropriate HTTP status code. If the token is valid, it adds the login user parameter to the context for further use in the request.

Function Signature:

func JwtValid() gin.HandlerFunc

Parameters: None

Return Values: gin.HandlerFunc - a middleware function that accepts the gin context and performs the JWT validation.

Functionality:

Gets the "Authorization" header from the incoming request.
If the "Authorization" header is empty, aborts the request with HTTP status code 401 (Unauthorized).
Splits the header value into two parts using whitespace as the separator.
If the header value does not contain two parts, aborts the request with HTTP status code 401 (Unauthorized).
If the first part of the header value is not "Bearer", aborts the request with HTTP status code 401 (Unauthorized).
Calls the "parseToken" function to validate the JWT token using the secret key.
If the token is invalid, aborts the request with HTTP status code 401 (Unauthorized) or 400 (Bad Request) based on the error.
If the token is valid, adds the login user parameter to the context and calls the next middleware function.
Function Name: parseToken

Description: parseToken is a helper function that validates the JWT token and returns the login user.

Function Signature:

func parseToken(accessToken string, signingKey []byte) (string, error)

Parameters:

accessToken string: the JWT token to be validated.
signingKey []byte: the secret key used for validating the token.
Return Values:

string: the login user extracted from the token.
error: an error if the token is invalid.
Functionality:

Parses the JWT token using jwt.ParseWithClaims function.
Validates the token using the signing key.
Returns an error if the token is invalid.
Returns the login user extracted from the token if the token is valid.
*/

package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"yudinsv/gophkeeper/internal/gophkeeperserver/constans"
	"yudinsv/gophkeeper/internal/gophkeeperserver/container"
	"yudinsv/gophkeeper/internal/gophkeeperserver/models"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/auth/pkg/auth"
)

// JwtValid Validate the JWT token.
func JwtValid() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, "register") || strings.Contains(c.Request.URL.Path, "login") {
			return
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if headerParts[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		login, err := parseToken(headerParts[1],
			[]byte(container.GetConfig().SecretKey),
		)
		if err != nil {
			status := http.StatusBadRequest
			if err == auth.ErrInvalidAccessToken {
				status = http.StatusUnauthorized
			}

			c.AbortWithStatus(status)
			return
		}
		c.AddParam(constans.CookeUserIDName, login)
	}
}

// parseToken parses the JWT token and returns the login user.
func parseToken(accessToken string, signingKey []byte) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims.Login, nil
	}

	return "", auth.ErrInvalidAccessToken
}
