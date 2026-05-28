package middleware

import (
	"fmt"
	"go-modular/internal/pkg/jwt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var jwtService jwt.JWT

func InitializeAuth(service jwt.JWT) {
	jwtService = service
}

// Auth middleware validates JWT token
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   "Authorization header is missing",
				"message": "Unauthorized",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   "Invalid Authorization header format",
				"message": "Unauthorized",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse token to get claims
		claims, err := jwtService.ParseToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   fmt.Sprintf("Invalid token: %v", err),
				"message": "Unauthorized",
			})
		}

		// Store user claims in context for use in handlers
		c.Set("user", claims)

		return next(c)
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user claims from context (set by Auth middleware)
			claims, ok := c.Get("user").(map[string]interface{})
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error":   "User claims not found",
					"message": "Unauthorized",
				})
			}

			// Get user role from claims
			userRole, ok := claims["role"].(string)
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"error":   "Role not found in token",
					"message": "Forbidden",
				})
			}

			// Check if user has required role
			hasRole := false
			for _, role := range roles {
				if userRole == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"error":   fmt.Sprintf("Insufficient permissions. Required role: %v, User role: %s", roles, userRole),
					"message": "Forbidden",
				})
			}

			return next(c)
		}
	}
}

// AdminOnly middleware allows only admin users
func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return RequireRole("admin")(next)
}

// UserOrAdmin middleware allows both user and admin roles
func UserOrAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return RequireRole("user", "admin")(next)
}
