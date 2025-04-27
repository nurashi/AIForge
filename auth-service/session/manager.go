package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

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
	if err != nil {
		return "", fmt.Errorf("failed to generate session id bytes: %w", err)
	}
	sessionID := base64.URLEncoding.EncodeToString(sessionIDBytes)

	redisKey := fmt.Sprintf("session:%s", sessionID)
	err = sm.rdb.Set(ctx, redisKey, userID, SessionDuration).Err()
	if err != nil {
		return "", fmt.Errorf("failed to set session in redis: %w", err)
	}

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


