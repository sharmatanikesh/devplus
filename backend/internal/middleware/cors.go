package middleware

import (
	"net/http"
	"net/url"
	"os"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Allowed Origin from Env or default to localhost:3000
		allowedOrigin := os.Getenv("FRONTEND_URL")
		if allowedOrigin == "" {
			allowedOrigin = "http://localhost:3000"
		}

		// Parse URL to get just the origin (scheme + host)
		if parsedURL, err := url.Parse(allowedOrigin); err == nil {
			allowedOrigin = parsedURL.Scheme + "://" + parsedURL.Host
		}

		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Cookie")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle Preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
