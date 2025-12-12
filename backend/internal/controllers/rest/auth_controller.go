package rest

import (
	"net/http"
	"time"

	"devplus-backend/internal/services/auth_service"
)

type AuthController struct {
	service *auth_service.AuthService
}

func NewAuthController(service *auth_service.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	authURL, err := c.service.InitiateLogin()
	if err != nil {
		http.Error(w, "Failed to initiate login", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (c *AuthController) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" || state == "" {
		http.Error(w, "Missing code or state", http.StatusBadRequest)
		return
	}

	session, err := c.service.HandleCallback(code, state)
	if err != nil {
		http.Error(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Set Session Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   false, // Set to true in production (checking TLS)
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to Frontend Dashboard
	http.Redirect(w, r, c.service.Config.FrontendURL, http.StatusTemporaryRedirect)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out"))
}

// GetCurrentUser is a helper to return current user info based on session
// This requires the middleware to have populated the context
func (c *AuthController) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Implementation depends on Middleware populating User
	// For now, simple return
	w.WriteHeader(http.StatusOK)
}
