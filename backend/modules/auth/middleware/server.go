package middleware

import (
	"diploma/modules/auth/jwt"
)

// JWTVerifier defines the interface for JWT verification
type JWTVerifier interface {
	VerifyToken(token string) (*jwt.Claims, error)
}

type AuthMiddleware struct {
	jwt JWTVerifier
}

func NewAuthMiddleware(jwt JWTVerifier) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: jwt,
	}
}
