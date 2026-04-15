package middleware

import (
	"context"
	"diploma/modules/auth/jwt"
	contextkeys "diploma/pkg/context-keys"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is the middleware that checks if a user is authenticated.
func (m *AuthMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwt.ExtractTokenFromHeader(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort() // Stop further processing of the request
			return
		}
		user, err := m.jwt.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort() // Stop further processing of the request
			return
		}
		// Save the user info in the context
		ctx := context.WithValue(c.Request.Context(), contextkeys.UserKey, user)
		c.Request = c.Request.WithContext(ctx)

		c.Next() // Proceed to the next handler
	}
}
