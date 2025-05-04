package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
	"log"

	"github.com/redis/go-redis/v9"
)

const (
	SessionCookieName = "aiforge_session_id"
	SessionDuration   = 24 * time.Hour
)

type SessionManager struct {
	rdb *redis.Client
}

func NewSessionManager(rdb *redis.Client) *SessionManager {
	return &SessionManager{rdb: rdb}
}

func (sm *SessionManager) CreateSession(ctx context.Context, w http.ResponseWriter, userID int64) (string, error) {
	
	sessionIDBytes := make([]byte, 32)
	
	_, err := rand.Read(sessionIDBytes)

	if err != nil { return "", fmt.Errorf("failed to generate session id bytes: %w", err) }

	sessionID := base64.URLEncoding.EncodeToString(sessionIDBytes)

	redisKey := fmt.Sprintf("session:%s", sessionID)

	err = sm.rdb.Set(ctx, redisKey, userID, SessionDuration).Err()

	if err != nil { return "", fmt.Errorf("failed to set session in redis: %w", err) }
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Expires:  time.Now().Add(SessionDuration),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	return sessionID, nil
}

// DeleteSession removes the session data from Redis based on the session ID
func (sm *SessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}
	redisKey := fmt.Sprintf("session:%s", sessionID)
	_, err := sm.rdb.Del(ctx, redisKey).Result()
	if err != nil {
		log.Printf("Error deleting session key %s from Redis: %v", redisKey, err)
		return fmt.Errorf("failed to delete session from redis: %w", err)
	}
	log.Printf("Deleted session key %s from Redis", redisKey)
	return nil
}

