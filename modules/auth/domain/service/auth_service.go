package service

import (
	"context"
	"errors"
	"go-modular/internal/pkg/jwt"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/users/domain/entity"
	"go-modular/modules/users/domain/repository"
)

// Errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyUsed   = errors.New("email already in use")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidOldPassword = errors.New("old password is incorrect")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrForbidden          = errors.New("forbidden: insufficient permissions")
)

// AuthService handles user authentication
type AuthService struct {
	userRepo repository.UserRepository
	jwt      jwt.JWT
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo repository.UserRepository, jwtService jwt.JWT) *AuthService {
	if userRepo == nil {
		panic("userRepo cannot be nil")
	}
	if jwtService == nil {
		panic("jwtService cannot be nil")
	}
	return &AuthService{
		userRepo: userRepo,
		jwt:      jwtService,
	}
}

// CreateUser creates a new user
func (s *AuthService) CreateUser(ctx context.Context, user *entity.User) error {
	if user.Email == "" || user.Password == "" {
		return errors.New("email and password cannot be empty")
	}

	existingUser, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil && err != repository.ERR_RECORD_NOT_FOUND {
		return err
	}
	if existingUser != nil {
		return ErrEmailAlreadyUsed
	}

	// Hash the password before saving the user
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Set default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	return s.userRepo.Create(ctx, user)
}

// ProcessLogin handles user login and password verification
func (s *AuthService) ProcessLogin(ctx context.Context, email, password string) (*entity.User, error) {
	// Validate input
	if email == "" || password == "" {
		return nil, errors.New("email and password cannot be empty")
	}

	// Find user by email
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == repository.ERR_RECORD_NOT_FOUND {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Compare the provided password with the hashed password in the database
	if !utils.CompareHashAndPassword(existingUser.Password, password) {
		return nil, ErrInvalidPassword
	}

	// Return the authenticated user
	return existingUser, nil
}

// GenerateTokens generates access and refresh tokens for a user
func (s *AuthService) GenerateTokens(user *entity.User) (accessToken string, refreshToken string, expiresIn int64, err error) {
	// Generate access token
	tokenData := map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"role":    user.Role,
	}

	accessToken, err = s.jwt.GenerateToken(tokenData)
	if err != nil {
		return "", "", 0, err
	}

	// Generate refresh token (longer expiration)
	refreshTokenData := map[string]interface{}{
		"user_id": user.ID,
		"type":    "refresh",
	}
	refreshToken, err = s.jwt.GenerateToken(refreshTokenData)
	if err != nil {
		return "", "", 0, err
	}

	// Expiration time in seconds (e.g., 24 hours)
	expiresIn = 86400

	return accessToken, refreshToken, expiresIn, nil
}

// ChangePassword changes user password with old password verification
func (s *AuthService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	if oldPassword == "" || newPassword == "" {
		return errors.New("old password and new password cannot be empty")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err == repository.ERR_RECORD_NOT_FOUND {
			return ErrUserNotFound
		}
		return err
	}

	// Verify old password
	if !utils.CompareHashAndPassword(user.Password, oldPassword) {
		return ErrInvalidOldPassword
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = hashedPassword

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

// GetProfile gets user profile by ID
func (s *AuthService) GetProfile(ctx context.Context, userID uint) (*entity.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err == repository.ERR_RECORD_NOT_FOUND {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// UpdateProfile updates user profile
func (s *AuthService) UpdateProfile(ctx context.Context, userID uint, name, email string) (*entity.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err == repository.ERR_RECORD_NOT_FOUND {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Check if email is already used by another user
	if email != user.Email {
		existingUser, err := s.userRepo.FindByEmail(ctx, email)
		if err != nil && err != repository.ERR_RECORD_NOT_FOUND {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != userID {
			return nil, ErrEmailAlreadyUsed
		}
	}

	user.Name = name
	user.Email = email

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, errors.New("failed to update profile")
	}

	return user, nil
}

// ValidateToken validates JWT token and returns user info
func (s *AuthService) ValidateToken(token string) (map[string]interface{}, error) {
	claims, err := s.jwt.ParseToken(token)
	if err != nil {
		return nil, ErrUnauthorized
	}
	return claims, nil
}

// RefreshAccessToken refreshes access token using refresh token
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, int64, error) {
	// Parse and validate refresh token
	claims, err := s.jwt.ParseToken(refreshToken)
	if err != nil {
		return "", 0, ErrUnauthorized
	}

	// Check if it's a refresh token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", 0, ErrUnauthorized
	}

	// Get user ID from claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return "", 0, ErrUnauthorized
	}
	userID := uint(userIDFloat)

	// Get user from database
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err == repository.ERR_RECORD_NOT_FOUND {
			return "", 0, ErrUserNotFound
		}
		return "", 0, err
	}

	// Generate new access token
	tokenData := map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"role":    user.Role,
	}

	accessToken, err := s.jwt.GenerateToken(tokenData)
	if err != nil {
		return "", 0, err
	}

	expiresIn := int64(86400) // 24 hours

	return accessToken, expiresIn, nil
}
