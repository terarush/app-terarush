// internal/modules/user/interfaces/dto/response/user_response.go

package response

import (
	"go-modular/modules/users/domain/entity"
	"time"
)

// UserResponse represents a user response
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Role      string    `json:"role"`
	Bio       string    `json:"bio,omitempty"`
	Banner    string    `json:"banner,omitempty"`
	IsBanned  bool      `json:"is_banned"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token,omitempty"`
	TokenType    string        `json:"token_type"`
	ExpiresIn    int64         `json:"expires_in"`
	User         *UserResponse `json:"user"`
}

// FromEntity converts a user entity to a user response
func FromEntity(user *entity.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Role:      user.Role,
		Bio:       user.Bio,
		Banner:    user.Banner,
		IsBanned:  user.IsBanned,
		Provider:  user.Provider,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// FromEntities converts a slice of user entities to a slice of user responses
func FromEntities(users []*entity.User) []*UserResponse {
	userResponses := make([]*UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = FromEntity(user)
	}
	return userResponses
}
