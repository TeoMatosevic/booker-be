package server // Or server.middleware

import (
	"booker-be/internal/session" // Adjust import path if needed
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "Bearer"
	authorizationPayloadKey = "userID" // Key to store user ID in Gin context
)

// AuthMiddleware creates a Gin middleware for authentication.
// It expects a SessionValidator (like the session.Store) to validate tokens.
func AuthMiddleware(sessionValidator session.SessionValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeaderKey)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 || !strings.EqualFold(fields[0], authorizationTypeBearer) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format. Expected 'Bearer <token>'"})
			return
		}

		accessToken := fields[1]
		userID, err := sessionValidator.ValidateToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired access token: " + err.Error()})
			return
		}

		// Set the userID in the context for subsequent handlers to use
		c.Set(authorizationPayloadKey, userID)
		c.Next()
	}
}
