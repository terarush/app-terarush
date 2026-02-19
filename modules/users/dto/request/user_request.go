package request

// LoginRequest represents a request to login a user
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"omitempty,oneof=admin user"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"omitempty,oneof=admin user"`
}

// ChangePasswordRequest represents a request to change password
type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" validate:"required,min=6"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// ResetPasswordRequest represents a request to reset password
type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// VerifyResetTokenRequest represents a request to verify reset token and set new password
type VerifyResetTokenRequest struct {
	Token           string `json:"token" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// RefreshTokenRequest represents a request to refresh access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// UpdateProfileRequest represents a request to update user profile
type UpdateProfileRequest struct {
	Name   string `json:"name" validate:"required"`
	Avatar string `json:"avatar" validate:"omitempty"`
}
