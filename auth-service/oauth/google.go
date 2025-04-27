package oauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/nurashi/AIForge/auth-service/config"
	"net/http"
	"time"
	"fmt"
	"encoding/base64"
	"crypto/rand"
	"io"
	"context"
)


/* init of Google config from config.go */
type GoogleOAuthConfig struct {
	Config *oauth2.Config
}


// function that reads a config from to get access to google auth service 
func NewGoogleAuthConfig(cfg *config.GoogleConfig) *GoogleOAuthConfig {
	return &GoogleOAuthConfig{
		Config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
				"openid",
			},
			Endpoint: google.Endpoint,
		},
	}

}


// generating cookie, something like JWT
func GenerateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(10 * time.Minute) 

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error generating random bytes:", err)
		return ""
	}
	state := base64.URLEncoding.EncodeToString(b)

	cookie := &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  expiration,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	return state
}

// taking data from auth of google
func (g *GoogleOAuthConfig) GetUserDataFromGoogle(ctx context.Context, code string) ([]byte, *oauth2.Token, error) {
	token, err := g.Config.Exchange(ctx, code)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	response, err := g.Config.Client(ctx, token).Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, token, fmt.Errorf("failed getting user info: %w", err)
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, token, fmt.Errorf("failed reading response body: %w", err)
	}

	return contents, token, nil
}