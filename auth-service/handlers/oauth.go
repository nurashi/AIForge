package handlers

import (
	"net/http"

	"github.com/nurashi/AIForge/auth-service/oauth"
	"golang.org/x/oauth2"
)

func OAuthLoginHandler(oauthConf *oauth.GoogleOAuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := oauth.GenerateStateOauthCookie(w)

		if state == "" {
			http.Error(w, "fail with state in google.go", http.StatusInternalServerError)
			return
		}

		url := oauthConf.Config.AuthCodeURL(state, oauth2.AccessTypeOffline)

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}