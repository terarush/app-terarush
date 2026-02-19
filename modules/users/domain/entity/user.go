package entity

import (
	"time"
)

// User represents a user entity
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Role      string    `json:"role" gorm:"type:enum('admin', 'user');default:'user'"`
	Password  string    `json:"-"`
	Provider  string    `json:"provider" gorm:"type:enum('local', 'github');default:'local'"`
	GithubID  string    `json:"github_id,omitempty" gorm:"column:github_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for User
func (*User) TableName() string {
	return "users"
}

// NewUser creates a new user
func NewUser(name, email, password string) *User {
	now := time.Now()
	return &User{
		Name:      name,
		Email:     email,
		Password:  password,
		Provider:  "local",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewGithubUser creates a new user from GitHub OAuth
func NewGithubUser(name, email, avatar, githubID string) *User {
	now := time.Now()
	return &User{
		Name:      name,
		Email:     email,
		Avatar:    avatar,
		Provider:  "github",
		GithubID:  githubID,
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
}
