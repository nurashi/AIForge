package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nurashi/AIForge/api-gateway/config"
	"github.com/nurashi/AIForge/api-gateway/handlers" 
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("FATAL: Failed to load api-gateway config: %v", err)
	}

	router := gin.Default()

	handlers.SetupRoutes(router, cfg) 
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	log.Printf("API Gateway starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("FATAL: Failed to start api-gateway server: %v", err)
	}
}
