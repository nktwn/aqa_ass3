package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken is returned when the token is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenExpired is returned when the token has expired
	ErrTokenExpired = errors.New("token has expired")
)

type JSONWebToken struct {
	jwtKey []byte
}

func NewJSONWebToken(jwtSecret string) *JSONWebToken {
	return &JSONWebToken{
		jwtKey: []byte(jwtSecret),
	}
}

type Claims struct {
	Role     int    `json:"role"`
	UserID   int64  `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (sa *JSONWebToken) GetJWTKey() []byte {
	return sa.jwtKey
}

func (sa *JSONWebToken) GenerateJSONWebTokens(id int64, username string, role int) (string, string, error) {
	accessToken, err := sa.generateShortLivedJSONWebToken(id, role, username)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := sa.generateLongLivedJSONWebToken(id, role, username)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (sa *JSONWebToken) generateShortLivedJSONWebToken(id int64, role int, username string) (string, error) {
	expiration := time.Now().Add(3000 * time.Hour)
	return sa.generateJSONWebToken(id, role, username, expiration)
}

func (sa *JSONWebToken) generateLongLivedJSONWebToken(id int64, role int, username string) (string, error) {
	expiration := time.Now().Add(24 * time.Hour)
	return sa.generateJSONWebToken(id, role, username, expiration)
}

func (sa *JSONWebToken) generateJSONWebToken(id int64, role int, username string, expirationTime time.Time) (string, error) {
	claims := &Claims{
		Username: username,
		UserID:   id,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(sa.jwtKey)
}

func (sa *JSONWebToken) RefreshAccessToken(refreshToken string) (string, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return sa.jwtKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	newAccessToken, err := sa.generateShortLivedJSONWebToken(claims.UserID, claims.Role, claims.Username)
	if err != nil {
		return "", err
	}
	return newAccessToken, nil
}

func (sa *JSONWebToken) VerifyToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method is HMAC (HS256).
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return sa.jwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func ExtractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}
	return parts[1], nil
}
