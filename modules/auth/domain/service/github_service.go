package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go-modular/internal/pkg/config"
)

// GitHubUser represents the user data from GitHub API
type GitHubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// GitHubEmail represents email data from GitHub API
type GitHubEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

// GitHubOAuthService handles GitHub OAuth operations
type GitHubOAuthService struct {
	clientID     string
	clientSecret string
	redirectURI  string
}

// NewGitHubOAuthService creates a new GitHub OAuth service
func NewGitHubOAuthService() *GitHubOAuthService {
	// Get redirect URI from environment, fallback to localhost for development
	redirectURI := config.GetString("GITHUB_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "http://localhost:5173/auth/github/callback"
	}
	
	return &GitHubOAuthService{
		clientID:     config.GetString("GITHUB_CLIENT"),
		clientSecret: config.GetString("GITHUB_SECRET"),
		redirectURI:  redirectURI,
	}
}

// GetAuthURL returns the GitHub OAuth authorization URL
func (s *GitHubOAuthService) GetAuthURL(state string) string {
	baseURL := "https://github.com/login/oauth/authorize"
	params := url.Values{}
	params.Add("client_id", s.clientID)
	params.Add("redirect_uri", s.redirectURI)
	params.Add("scope", "user:email")
	params.Add("state", state)

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// ExchangeCode exchanges the authorization code for an access token
func (s *GitHubOAuthService) ExchangeCode(ctx context.Context, code string) (string, error) {
	tokenURL := "https://github.com/login/oauth/access_token"

	data := url.Values{}
	data.Set("client_id", s.clientID)
	data.Set("client_secret", s.clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", s.redirectURI)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	accessToken, ok := result["access_token"].(string)
	if !ok || accessToken == "" {
		return "", fmt.Errorf("no access token in response: %v", result)
	}

	return accessToken, nil
}

// GetUserInfo fetches user information from GitHub API
func (s *GitHubOAuthService) GetUserInfo(ctx context.Context, accessToken string) (*GitHubUser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var user GitHubUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// If email is not public, fetch from emails endpoint
	if user.Email == "" {
		email, err := s.GetPrimaryEmail(ctx, accessToken)
		if err == nil {
			user.Email = email
		}
	}

	return &user, nil
}

// GetPrimaryEmail fetches the primary verified email from GitHub API
func (s *GitHubOAuthService) GetPrimaryEmail(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch emails: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var emails []GitHubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", fmt.Errorf("failed to parse emails: %w", err)
	}

	// Find primary verified email
	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	// If no primary email, return first verified email
	for _, email := range emails {
		if email.Verified {
			return email.Email, nil
		}
	}

	return "", fmt.Errorf("no verified email found")
}
