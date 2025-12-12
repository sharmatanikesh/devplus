package middleware

import (
	"context"
	"net/http"
	"time"

	"devplus-backend/internal/db"
	"devplus-backend/internal/models"
)

// SessionMiddleware validates the session_token cookie
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized: No session cookie", http.StatusUnauthorized)
			return
		}

		sessionID := cookie.Value
		dbInstance := db.GetInstance()

		var session models.Session
		// Preload User to get AccessToken
		if err := dbInstance.Preload("User").Where("id = ? AND expires_at > ?", sessionID, time.Now()).First(&session).Error; err != nil {
			http.Error(w, "Unauthorized: Invalid or expired session", http.StatusUnauthorized)
			return
		}

		// Add User and Token to Context
		ctx := context.WithValue(r.Context(), UserContextKey, session.User)
		// Populate GithubTokenContextKey with the stored access token so existing controllers work!
		ctx = context.WithValue(ctx, GithubTokenContextKey, session.User.AccessToken)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
