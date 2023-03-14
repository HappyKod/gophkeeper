package handlers

import (
	"net/http"

	"yudinsv/gophkeeper/internal/gophkeeperserver/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Router указание маршрутов севера
func Router() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	api := r.Group("/api")
	api.Use(middleware.JwtValid())
	v1 := api.Group("/v1")

	{
		v1.POST("/register", registerHandler)
		v1.POST("/login", authenticationHandler)

		v1.GET("/sync", syncDataHandler)
		v1.PUT("/", putDataHandler)
		v1.POST("/", getDataHandler)
		v1.DELETE("/", deleteDataHandler)
	}
	r.GET("/ping", func(context *gin.Context) {
		context.String(http.StatusOK, "pong")
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	return r
}
