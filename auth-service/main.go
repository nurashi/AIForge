package main

import (
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

	dbPool, err := database.InitPostgres(&cfg.Postgres)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize PostgreSQL: %v", err)
	}
	defer dbPool.Close()

	redisClient, err := database.InitRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize Redis: %v", err)
	}

	oauthConf := oauth.NewGoogleAuthConfig(&cfg.Google)
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
		fmt.Fprintf(w, "General endpoint works. Use /auth/login to use Google OAuth")
	})
	
	// Assuming you've updated your handlers to be consistent
	mux.HandleFunc("/auth/login", handlers.OAuthLoginHandler(appInstance.OAuthConfig))
	mux.HandleFunc("/auth/callback", handlers.OAuthCallbackHandler(appInstance))

	addr := fmt.Sprintf(":%d", cfg.App.Port)
	log.Printf("Server starting on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("FATAL: Failed to start server: %v", err)
	}
}