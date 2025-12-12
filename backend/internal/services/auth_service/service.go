package auth_service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"devplus-backend/internal/config"
	"devplus-backend/internal/db"
	"devplus-backend/internal/models"
)

type AuthService struct {
	db     *gorm.DB
	Config *config.Config
}

func NewAuthService() *AuthService {
	return &AuthService{
		db:     db.GetInstance(),
		Config: config.LoadConfig(),
	}
}

// InitiateLogin generates a state token and returns the GitHub OAuth URL
func (s *AuthService) InitiateLogin() (string, error) {
	state := uuid.New().String()

	authState := models.AuthState{
		State:     state,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := s.db.Create(&authState).Error; err != nil {
		return "", err
	}

	// Correct URL construction:
	u, _ := url.Parse("https://github.com/login/oauth/authorize")
	q := u.Query()
	q.Set("client_id", s.Config.GithubClientID)
	if s.Config.GithubRedirectURI != "" {
		q.Set("redirect_uri", s.Config.GithubRedirectURI)
	}
	q.Set("scope", "user:email read:user repo")
	q.Set("state", state)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

type GitHubTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type GitHubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

// HandleCallback processes the implementation
func (s *AuthService) HandleCallback(code, state string) (*models.Session, error) {
	// 1. Validate State
	var authState models.AuthState
	if err := s.db.Where("state = ? AND expires_at > ?", state, time.Now()).First(&authState).Error; err != nil {
		return nil, errors.New("invalid or expired state token")
	}

	// Delete used state
	s.db.Delete(&authState)

	// 2. Exchange Code for Token
	tokenResp, err := s.exchangeCodeForToken(code)
	if err != nil {
		return nil, err
	}

	// 3. Fetch GitHub User
	ghUser, err := s.fetchGitHubUser(tokenResp.AccessToken)
	if err != nil {
		return nil, err
	}

	// 4. Fetch Email if missing
	if ghUser.Email == "" {
		ghUser.Email, _ = s.fetchGitHubEmail(tokenResp.AccessToken)
	}

	// 5. Upsert User
	user := models.User{
		GithubID:     ghUser.ID,
		Username:     ghUser.Login,
		Email:        ghUser.Email,
		AvatarURL:    ghUser.AvatarURL,
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
	}

	// Check if user exists
	var existingUser models.User
	result := s.db.Where("github_id = ?", ghUser.ID).First(&existingUser)
	if result.Error == nil {
		// Update existing
		existingUser.Username = user.Username
		existingUser.Email = user.Email
		existingUser.AvatarURL = user.AvatarURL
		existingUser.AccessToken = user.AccessToken
		existingUser.RefreshToken = user.RefreshToken
		s.db.Save(&existingUser)
		user = existingUser
	} else {
		// Create new
		if err := s.db.Create(&user).Error; err != nil {
			return nil, err
		}
	}

	// 6. Create Session
	session := models.Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 1 day session
	}

	if err := s.db.Create(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *AuthService) exchangeCodeForToken(code string) (*GitHubTokenResponse, error) {
	values := url.Values{}
	values.Set("client_id", s.Config.GithubClientID)
	values.Set("client_secret", s.Config.GithubClientSecret)
	values.Set("code", code)

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBufferString(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github token exchange failed: %s", string(body))
	}

	var tokenResp GitHubTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

func (s *AuthService) fetchGitHubUser(token string) (*GitHubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user from github")
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) fetchGitHubEmail(token string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []GitHubEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	if len(emails) > 0 {
		return emails[0].Email, nil
	}
	return "", nil
}
