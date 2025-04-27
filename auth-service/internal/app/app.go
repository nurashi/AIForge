package app 

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nurashi/AIForge/auth-service/oauth"
	"github.com/nurashi/AIForge/auth-service/repository"
	"github.com/nurashi/AIForge/auth-service/session"
	"github.com/redis/go-redis/v9"
)


// holds parts of Auth
type App struct {
	DB             *pgxpool.Pool
	Redis          *redis.Client
	OAuthConfig    *oauth.GoogleOAuthConfig
	UserRepo       repository.UserRepository
	SessionManager *session.SessionManager
}
