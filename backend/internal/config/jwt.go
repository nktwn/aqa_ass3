package config

import (
	"os"

	"github.com/pkg/errors"
)

const (
	jwtSecretKey = "JWT_SECRET_KEY"
)

type JWTConfig interface {
	GetSecretKey() string
}

type jwtConfig struct {
	secretKey string
}

func NewJWTConfig() (JWTConfig, error) {
	secretKey := os.Getenv(jwtSecretKey)
	if len(secretKey) == 0 {
		return nil, errors.New("jwt secret key not found")
	}

	return &jwtConfig{
		secretKey: secretKey,
	}, nil
}

func (cfg *jwtConfig) GetSecretKey() string {
	return cfg.secretKey
}
