package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nurashi/AIForge/user-service/middleware"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong from user-service"})
	})

	router.Static("/uploads", "./uploads")

	protected := router.Group("/user")
	protected.Use(middleware.AuthRequired(rdb))
	{
		protected.GET("/profile", GetMyProfile(db))
		protected.PUT("/profile", UpdateMyProfile(db))
		protected.POST("/profile/avatar", UploadAvatar(db))
		protected.DELETE("/profile/avatar", DeleteAvatar(db))
	}
}
