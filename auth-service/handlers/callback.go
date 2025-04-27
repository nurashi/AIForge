package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/nurashi/AIForge/auth-service/models"
	app "github.com/nurashi/AIForge/auth-service/internal/app" 
)

func OAuthCallbackHandler(appInstance *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		receivedState := r.FormValue("state")
		originalStateCookie, err := r.Cookie("oauthstate")
		if err != nil {
			log.Println("ERROR-state-cookie")
			http.Error(w, "ERROR-state-cookie", http.StatusBadRequest)
			return
		}
		if receivedState != originalStateCookie.Value {
			log.Println("ERROR-state-mismatch")
			http.Error(w, "ERROR-state-mismatch", http.StatusBadRequest)
			return
		}
		deleteCookie := &http.Cookie{Name: "oauthstate", MaxAge: -1, Path: "/"}
		http.SetCookie(w, deleteCookie)

		code := r.FormValue("code")
		if code == "" {
			googleError := r.FormValue("error")
			if googleError != "" {
				log.Printf("ERROR-google-auth: %s", googleError)
				http.Error(w, "ERROR-google-auth", http.StatusUnauthorized)
			} else {
				log.Println("ERROR-code-missing")
				http.Error(w, "ERROR-code-missing", http.StatusBadRequest)
			}
			return
		}

		userDataBytes, token, err := appInstance.OAuthConfig.GetUserDataFromGoogle(ctx, code)
		if err != nil {
			log.Printf("ERROR-get-user-data: %v", err)
			http.Error(w, "ERROR-get-user-data", http.StatusInternalServerError)
			return
		}

		var userInfo models.GoogleUserInfo
		err = json.Unmarshal(userDataBytes, &userInfo)
		if err != nil {
			log.Printf("ERROR-parse-user-json: %v", err)
			http.Error(w, "ERROR-parse-user-json", http.StatusInternalServerError)
			return
		}

		var appUser *models.User
		foundUser, err := appInstance.UserRepo.FindUserByGoogleID(ctx, userInfo.Sub)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				createdUser, createErr := appInstance.UserRepo.CreateUser(ctx, &userInfo)
				if createErr != nil {
					log.Printf("ERROR-create-user: %v", createErr)
					http.Error(w, "ERROR-create-user", http.StatusInternalServerError)
					return
				}
				appUser = createdUser
			} else {
				log.Printf("ERROR-find-user: %v", err)
				http.Error(w, "ERROR-find-user", http.StatusInternalServerError)
				return
			}
		} else {
			appUser = foundUser
		}

		_, err = appInstance.SessionManager.CreateSession(ctx, w, appUser.ID)
		if err != nil {
			log.Printf("ERROR-create-session: %v", err)
		}


		if token.RefreshToken != "" {
			log.Printf("refresh token")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK - Session Created"))
	}
}
