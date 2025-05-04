package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/nurashi/AIForge/api-gateway/config"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) {

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong from gateway"})
	})

    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Welcome to AIForge API Gateway (c) Nurasyl Orazbek's first API_GATEWAY"})
    })

	authGroup := router.Group("/auth")
	{
		authProxy := NewReverseProxy(cfg.Services.Auth)
		authGroup.GET("/login", gin.WrapH(authProxy))
		authGroup.GET("/callback", gin.WrapH(authProxy))
		authGroup.GET("/logout", gin.WrapH(authProxy))
	}

}
