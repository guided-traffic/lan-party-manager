package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/guided-traffic/rate-your-mate/backend/auth"
)

const (
	// ContextKeyClaims is the key used to store JWT claims in the Gin context
	ContextKeyClaims = "claims"
)

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			return
		}

		// Check Bearer token format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format. Use: Bearer <token>",
			})
			return
		}

		tokenString := parts[1]

		// Validate the token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		// Store claims in context for handlers to use
		c.Set(ContextKeyClaims, claims)
		c.Next()
	}
}

// GetClaims retrieves the JWT claims from the Gin context
func GetClaims(c *gin.Context) (*auth.Claims, bool) {
	claims, exists := c.Get(ContextKeyClaims)
	if !exists {
		return nil, false
	}

	jwtClaims, ok := claims.(*auth.Claims)
	return jwtClaims, ok
}

// GetUserID retrieves the user ID from the Gin context
func GetUserID(c *gin.Context) (uint64, bool) {
	claims, ok := GetClaims(c)
	if !ok {
		return 0, false
	}
	return claims.UserID, true
}

// GetSteamID retrieves the Steam ID from the Gin context
func GetSteamID(c *gin.Context) (string, bool) {
	claims, ok := GetClaims(c)
	if !ok {
		return "", false
	}
	return claims.SteamID, true
}
