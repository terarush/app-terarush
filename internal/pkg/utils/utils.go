package utils

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain text password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CompareHashAndPassword verifies if the provided password matches the stored hashed password
func CompareHashAndPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GetUserIDFromContext(ctx context.Context) (uint, error) {
	userId, ok := ctx.Value("userId").(uint)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return userId, nil
}
