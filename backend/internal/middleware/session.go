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
		// Get token from Authorization header (format: "Bearer <token>")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: No authorization header", http.StatusUnauthorized)
			return
		}

		// Extract token
		var sessionID string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			sessionID = authHeader[7:]
		} else {
			http.Error(w, "Unauthorized: Invalid authorization format", http.StatusUnauthorized)
			return
		}

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
