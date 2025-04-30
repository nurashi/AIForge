package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/nurashi/AIForge/auth-service/config"
	"github.com/nurashi/AIForge/auth-service/database"
	"github.com/nurashi/AIForge/auth-service/handlers"
	"github.com/nurashi/AIForge/auth-service/oauth" 
	"github.com/nurashi/AIForge/auth-service/repository"
	"github.com/nurashi/AIForge/auth-service/session"
	app "github.com/nurashi/AIForge/auth-service/internal/app"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("FATAL: Failed to load configuration: %v", err)
	}

	if cfg.Google.ClientID == "" || cfg.Google.ClientSecret == "" {
		log.Fatalf("FATAL: GOOGLE_CLIENT_ID and/or GOOGLE_CLIENT_SECRET are not set")
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

	ctx := context.Background() 
	oauthConf, err := oauth.NewGoogleOAuthConfig(ctx, &cfg.Google)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize OAuth/OIDC config: %v", err)
	}

	userRepo := repository.NewUserRepository(dbPool)
	sessionMgr := session.NewSessionManager(redisClient)

	appInstance := &app.App{
		DB:             dbPool,
		Redis:          redisClient,
		OAuthConfig:    oauthConf, 
		UserRepo:       userRepo,
		SessionManager: sessionMgr,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "general endpoint is works, Use -> /auth/login to use Google OAuth")
	})
	mux.HandleFunc("/auth/login", handlers.OAuthLoginHandler(appInstance.OAuthConfig))
	mux.HandleFunc("/auth/callback", handlers.OAuthCallbackHandler(appInstance))
	mux.HandleFunc("/auth/logout", handlers.LogoutHandler(appInstance))

	addr := fmt.Sprintf(":%d", cfg.App.Port)
	log.Printf("Server starting on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("FATAL: Failed to start server: %v", err)
	}
}
