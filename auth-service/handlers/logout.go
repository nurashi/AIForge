package handlers

import (
	"log"
	"net/http"
	"time"

	app "github.com/nurashi/AIForge/auth-service/internal/app"
	"github.com/nurashi/AIForge/auth-service/session"          
)

func LogoutHandler(appInstance *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sessionCookie, err := r.Cookie(session.SessionCookieName)
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/", http.StatusSeeOther) 
				return
			}
			log.Printf("ERROR-logout-read-cookie: %v", err)
			http.Error(w, "ERROR-reading-cookie", http.StatusBadRequest)
			return
		}

		sessionID := sessionCookie.Value

		err = appInstance.SessionManager.DeleteSession(ctx, sessionID)
		if err != nil {
			log.Printf("ERROR-logout-delete-session: %v (SessionID: %s)", err, sessionID)
		} else {
			log.Printf("Logout successfull: Deleted session  %s from Redis", sessionID)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     session.SessionCookieName,
			Value:    "",            
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,           
			Path:     "/",
			HttpOnly: true,
			Secure:   false, 
			SameSite: http.SameSiteLaxMode,
		})

		log.Println("User logged out, redirecting.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

