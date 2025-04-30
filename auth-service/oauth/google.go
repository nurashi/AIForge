package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc" 
	"github.com/nurashi/AIForge/auth-service/config"
	"github.com/nurashi/AIForge/auth-service/models" 
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleOAuthConfig holds the OAuth2 config and OIDC provider/verifier.
type GoogleOAuthConfig struct {
	Config   *oauth2.Config
	Verifier *oidc.IDTokenVerifier 
	Provider *oidc.Provider        
}

func NewGoogleOAuthConfig(ctx context.Context, cfg *config.GoogleConfig) (*GoogleOAuthConfig, error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	oauthConfig := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes: []string{
			oidc.ScopeOpenID, 
			"profile",        
			"email",          
		},
		Endpoint: google.Endpoint, 
	}

	return &GoogleOAuthConfig{
		Config:   oauthConfig,
		Verifier: verifier,
		Provider: provider,
	}, nil
}

func GenerateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(10 * time.Minute)
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error generating random bytes for state:", err)
		return ""
	}
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  expiration,
		HttpOnly: true,
		Path:     "/",
		Secure:   false, 
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	return state
}

func (g *GoogleOAuthConfig) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := g.Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	return token, nil
}

func (g *GoogleOAuthConfig) VerifyIDTokenAndExtractUserInfo(ctx context.Context, token *oauth2.Token) (*models.GoogleUserInfo, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("id_token field missing or not a string")
	}

	idToken, err := g.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	var userInfo models.GoogleUserInfo
	if err := idToken.Claims(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to extract claims from ID token: %w", err)
	}

	return &userInfo, nil
}


