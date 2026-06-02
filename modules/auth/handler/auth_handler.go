package handler

import (
	"fmt"
	"strings"
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/jwt"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/auth/domain/service"
	"go-modular/modules/users/domain/entity"
	"go-modular/modules/users/dto/request"
	"go-modular/modules/users/dto/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

// AuthHandler struct handles HTTP request for auth.
type AuthHandler struct {
	authService   *service.AuthService
	githubService *service.GitHubOAuthService
	log           *logger.Logger
	event         *bus.EventBus
	jwt           jwt.JWT
	r             *utils.Response
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(log *logger.Logger, event *bus.EventBus, authService *service.AuthService, jwt jwt.JWT) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		githubService: service.NewGitHubOAuthService(),
		log:           log,
		event:         event,
		jwt:           jwt,
		r:             &utils.Response{},
	}
}

// Initialize Event Handle.
func (h *AuthHandler) Handle(event bus.Event) {
	fmt.Printf("User created: %v", event.Payload)
}

// Register handles user registration.
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.CreateUserRequest true "User registration request"
// @Success 200 {object} response.AuthResponse "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 409 {object} map[string]interface{} "Email already in use"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	h.log.Info("Handling register request")

	req := new(request.CreateUserRequest)
	if err := c.Bind(req); err != nil {
		h.log.Error("Failed to bind request:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		h.log.Error("Validation failed:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	h.log.Debug("Request validated successfully:", req)

	// Create user with role (default to "user" if not specified)
	user := entity.NewUser(req.Name, req.Email, req.Password)
	if req.Role != "" {
		user.Role = req.Role
	} else {
		user.Role = "user"
	}

	err := h.authService.CreateUser(c.Request().Context(), user)
	if err != nil {
		if err == service.ErrEmailAlreadyUsed {
			h.log.Warn("Email already in use:", req.Email)
			return h.r.ErrorResponse(c, http.StatusConflict, "Email already in use")
		}
		h.log.Error("Failed to create user:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	h.log.Debug("User created successfully:", user)

	h.event.Publish(bus.Event{Type: "user.created", Payload: user})
	h.log.Debug("Event 'user.created' published successfully")

	// Generate tokens
	accessToken, refreshToken, expiresIn, err := h.authService.GenerateTokens(user)
	if err != nil {
		h.log.Error("Failed to generate tokens:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate tokens")
	}

	authResponse := &response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         response.FromEntity(user),
	}

	return h.r.SuccessResponse(c, authResponse, "User registered successfully")
}

// Login handles user login.
// @Summary User login
// @Description Authenticate user with email and password, returns access and refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "User login request"
// @Success 200 {object} response.AuthResponse "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Invalid email or password"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	h.log.Info("Handling login request")

	req := new(request.LoginRequest)
	if err := c.Bind(req); err != nil {
		h.log.Error("Failed to bind request:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		h.log.Error("Validation failed:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	h.log.Debug("Request validated successfully:", req)

	user, err := h.authService.ProcessLogin(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		if err == service.ErrUserNotFound || err == service.ErrInvalidPassword {
			h.log.Warn("Invalid email or password for:", req.Email)
			return h.r.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		}
		h.log.Error("Failed to process login:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	h.log.Debug("User authenticated successfully:", user)

	// Generate tokens
	accessToken, refreshToken, expiresIn, err := h.authService.GenerateTokens(user)
	if err != nil {
		h.log.Error("Failed to generate tokens:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate tokens")
	}

	authResponse := &response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         response.FromEntity(user),
	}

	return h.r.SuccessResponse(c, authResponse, "Login successful")
}

// GetProfile retrieves the authenticated user's profile.
// @Summary Get user profile
// @Description Retrieve the authenticated user's profile information
// @Tags Auth
// @Produce json
// @Success 200 {object} response.UserResponse "User profile"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/profile [get]
// @Security Bearer
func (h *AuthHandler) GetProfile(c echo.Context) error {
	h.log.Info("Handling get profile request")

	// Get user ID from context (set by auth middleware)
	claims := c.Get("user").(map[string]interface{})
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		h.log.Error("Invalid user ID in token claims")
		return h.r.ErrorResponse(c, http.StatusUnauthorized, "Invalid token claims")
	}
	userID := uint(userIDFloat)

	user, err := h.authService.GetProfile(c.Request().Context(), userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			h.log.Warn("User not found:", userID)
			return h.r.ErrorResponse(c, http.StatusNotFound, "User not found")
		}
		h.log.Error("Failed to get profile:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntity(user), "Profile retrieved successfully")
}

// UpdateProfile updates the authenticated user's profile.
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.UpdateUserRequest true "Profile update request"
// @Success 200 {object} response.UserResponse "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 409 {object} map[string]interface{} "Email already in use"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/profile [put]
// @Security Bearer
func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	h.log.Info("Handling update profile request")

	// Get user ID from context
	claims := c.Get("user").(map[string]interface{})
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		h.log.Error("Invalid user ID in token claims")
		return h.r.ErrorResponse(c, http.StatusUnauthorized, "Invalid token claims")
	}
	userID := uint(userIDFloat)

	req := new(request.UpdateUserRequest)
	if err := c.Bind(req); err != nil {
		h.log.Error("Failed to bind request:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		h.log.Error("Validation failed:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := h.authService.UpdateProfile(c.Request().Context(), userID, req.Name, req.Email)
	if err != nil {
		if err == service.ErrEmailAlreadyUsed {
			h.log.Warn("Email already in use:", req.Email)
			return h.r.ErrorResponse(c, http.StatusConflict, "Email already in use")
		}
		if err == service.ErrUserNotFound {
			h.log.Warn("User not found:", userID)
			return h.r.ErrorResponse(c, http.StatusNotFound, "User not found")
		}
		h.log.Error("Failed to update profile:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntity(user), "Profile updated successfully")
}

// ChangePassword handles password change request.
// @Summary Change user password
// @Description Change the authenticated user's password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.ChangePasswordRequest true "Password change request"
// @Success 200 {object} map[string]interface{} "Password changed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or old password"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/change-password [post]
// @Security Bearer
func (h *AuthHandler) ChangePassword(c echo.Context) error {
	h.log.Info("Handling change password request")

	// Get user ID from context
	claims := c.Get("user").(map[string]interface{})
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		h.log.Error("Invalid user ID in token claims")
		return h.r.ErrorResponse(c, http.StatusUnauthorized, "Invalid token claims")
	}
	userID := uint(userIDFloat)

	req := new(request.ChangePasswordRequest)
	if err := c.Bind(req); err != nil {
		h.log.Error("Failed to bind request:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		h.log.Error("Validation failed:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err := h.authService.ChangePassword(c.Request().Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		if err == service.ErrInvalidOldPassword {
			h.log.Warn("Invalid old password for user:", userID)
			return h.r.ErrorResponse(c, http.StatusBadRequest, "Old password is incorrect")
		}
		if err == service.ErrUserNotFound {
			h.log.Warn("User not found:", userID)
			return h.r.ErrorResponse(c, http.StatusNotFound, "User not found")
		}
		h.log.Error("Failed to change password:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return h.r.SuccessResponse(c, nil, "Password changed successfully")
}

// RefreshToken handles refresh token request.
// @Summary Refresh access token
// @Description Get a new access token using a valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} map[string]interface{} "Token refreshed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Invalid refresh token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	h.log.Info("Handling refresh token request")

	req := new(request.RefreshTokenRequest)
	if err := c.Bind(req); err != nil {
		h.log.Error("Failed to bind request:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		h.log.Error("Validation failed:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	accessToken, expiresIn, err := h.authService.RefreshAccessToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		if err == service.ErrUnauthorized {
			h.log.Warn("Invalid refresh token")
			return h.r.ErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token")
		}
		h.log.Error("Failed to refresh token:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return h.r.SuccessResponse(c, map[string]interface{}{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   expiresIn,
	}, "Token refreshed successfully")
}

// Logout handles user logout (invalidate token on client side).
// @Summary User logout
// @Description Logout the authenticated user (client-side token removal)
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]interface{} "Logged out successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/logout [post]
// @Security Bearer
func (h *AuthHandler) Logout(c echo.Context) error {
	h.log.Info("Handling logout request")

	// In a stateless JWT system, logout is handled on the client side
	// by removing the token. Here we just return a success response.
	return h.r.SuccessResponse(c, nil, "Logged out successfully")
}

// GitHubCallback handles the GitHub OAuth callback
// @Summary GitHub OAuth callback
// @Description Handle GitHub OAuth callback and authenticate user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "GitHub authorization code"
// @Success 200 {object} response.AuthResponse "GitHub login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request or missing code"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/github/callback [post]
func (h *AuthHandler) GitHubCallback(c echo.Context) error {
	h.log.Info("Handling GitHub OAuth callback")

	// Parse request body
	var req struct {
		Code string `json:"code" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		h.log.Error("Failed to bind request:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if req.Code == "" {
		h.log.Error("No code provided in callback")
		return h.r.ErrorResponse(c, http.StatusBadRequest, "No authorization code provided")
	}

	code := req.Code

	// Exchange code for access token
	accessToken, err := h.githubService.ExchangeCode(c.Request().Context(), code)
	if err != nil {
		h.log.Error("Failed to exchange code:", err)
		// Return more descriptive error message (BadRequest for code issues, InternalServerError for other issues)
		statusCode := http.StatusBadRequest
		if !strings.Contains(err.Error(), "github oauth error") {
			statusCode = http.StatusInternalServerError
		}
		return h.r.ErrorResponse(c, statusCode, fmt.Sprintf("GitHub authentication failed: %v", err))
	}

	// Get user info from GitHub
	githubUser, err := h.githubService.GetUserInfo(c.Request().Context(), accessToken)
	if err != nil {
		h.log.Error("Failed to get GitHub user info:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch user information")
	}

	h.log.Debug("GitHub user info:", githubUser)

	// Check if user already exists by GitHub ID
	user, err := h.authService.FindOrCreateGitHubUser(
		c.Request().Context(),
		fmt.Sprintf("%d", githubUser.ID),
		githubUser.Name,
		githubUser.Email,
		githubUser.AvatarURL,
	)
	if err != nil {
		h.log.Error("Failed to find or create user:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, "Failed to process user")
	}

	// Generate tokens
	accessTokenJWT, refreshToken, expiresIn, err := h.authService.GenerateTokens(user)
	if err != nil {
		h.log.Error("Failed to generate tokens:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate tokens")
	}

	authResponse := &response.AuthResponse{
		AccessToken:  accessTokenJWT,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         response.FromEntity(user),
	}

	return h.r.SuccessResponse(c, authResponse, "GitHub login successful")
}

// RegisterRoutes sets up the auth routes.
func (h *AuthHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath + "/auth")

	// Public routes
	group.POST("/register", h.Register)
	group.POST("/login", h.Login)
	group.POST("/refresh", h.RefreshToken)

	// GitHub OAuth routes
	group.POST("/github/callback", h.GitHubCallback)

	// Protected routes (require authentication)
	protected := group.Group("", middleware.Auth)
	protected.GET("/profile", h.GetProfile)
	protected.PUT("/profile", h.UpdateProfile)
	protected.POST("/change-password", h.ChangePassword)
	protected.POST("/logout", h.Logout)
}
