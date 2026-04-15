package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	AuthRedisDB = iota
)
const (
	redisPortEnvName     = "REDIS_PORT"
	redisHostEnvName     = "REDIS_HOST"
	redisPasswordEnvName = "REDIS_PASSWORD"
	redisAuthDBEnvName   = "REDIS_AUTH_DB"
)

type RedisConfig interface {
	Addr() string
	Password() string
	RedisAuthDB() int
}

type redisConfig struct {
	authRedisDb int
	host        string
	port        string
	password    string
}

// NewRedisConfig создает Redis-конфиг из переменных окружения.
func NewRedisConfig() (RedisConfig, error) {
	host := os.Getenv(redisHostEnvName)
	if host == "" {
		return nil, errors.New("REDIS_HOST not found")
	}

	port := os.Getenv(redisPortEnvName)
	if port == "" {
		return nil, errors.New("REDIS_PORT not found")
	}

	authRedisDb := os.Getenv(redisAuthDBEnvName)

	if authRedisDb == "" {
		return nil, errors.New("REDIS_AUTH_DB not found")
	}
	// Convert authRedisDb to int
	authRedisDbInt, err := strconv.Atoi(authRedisDb)
	if err != nil {
		return nil, fmt.Errorf("REDIS_AUTH_DB is not a valid integer: %s", err.Error())
	}

	// Password may be empty (no password)
	password := os.Getenv(redisPasswordEnvName)

	return &redisConfig{
		host:        host,
		port:        port,
		password:    password,	
		authRedisDb: authRedisDbInt,
	}, nil
}

// Addr returns the Redis address.
func (cfg *redisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", cfg.host, cfg.port)
}

// Password returns the Redis password.
func (cfg *redisConfig) Password() string {
	return cfg.password
}
func (cfg *redisConfig) RedisAuthDB() int {
	return cfg.authRedisDb
}
