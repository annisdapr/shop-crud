// File: middleware/auth.go
package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Standardized errors for consistent responses.
var (
	// Returned if the Authorization header is missing or malformed.
	ErrMissingAuthHeader = echo.NewHTTPError(http.StatusUnauthorized, "Missing or malformed JWT")
	// Returned if the token is invalid or expired.
	ErrInvalidJWT = echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired JWT")
)

// JWTAuthMiddleware returns an Echo middleware that validates JWT tokens.
func JWTAuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the Authorization header.
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return ErrMissingAuthHeader
			}

			// Check header format: must be "Bearer <token>".
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return ErrMissingAuthHeader
			}
			tokenString := parts[1]

			// Parse and validate the token.
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Ensure HMAC signing method to prevent downgrade attacks.
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return ErrInvalidJWT
			}

			// Extract claims from the token.
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return ErrInvalidJWT
			}

			// Store claims in the context for later use.
			c.Set("user", claims)

			return next(c)
		}
	}
}

// GetUserFromContext retrieves JWT claims from the Echo context.
func GetUserFromContext(c echo.Context) (jwt.MapClaims, bool) {
	user := c.Get("user")
	if user == nil {
		return nil, false
	}
	claims, ok := user.(jwt.MapClaims)
	return claims, ok
}
