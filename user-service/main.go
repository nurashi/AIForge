package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nurashi/AIForge/user-service/config"
	"github.com/nurashi/AIForge/user-service/database"
	"github.com/nurashi/AIForge/user-service/handlers"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("FATAL: Failed to load configuration: %v", err)
	}

	dbPool, err := database.InitPostgres(&cfg.Postgres)
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to PostgreSQL: %v", err)
	}
	defer dbPool.Close()

	redisClient, err := database.InitRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	router := gin.Default()

	handlers.SetupRoutes(router, dbPool, redisClient)

	addr := fmt.Sprintf(":%d", cfg.App.Port)
	log.Printf("User service starting on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("FATAL: Failed to start server: %v", err)
	}
}
