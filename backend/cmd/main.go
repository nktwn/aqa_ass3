package main

import (
	"context"
	"diploma/internal/app"
	"diploma/pkg/logger"
)

// @title Go JWT Swagger Example API
// @description This is a sample server with JWT authorization.
// @version 1.0

// @SecurityDefinitions.bearer

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		logger.Fatal("Failed to initialize app", logger.Field("error", err))
	}

	err = a.Run()
	if err != nil {
		logger.Fatal("Failed to run app", logger.Field("error", err))
	}

	// Ensure all logs are flushed before exiting
	_ = logger.Sync()
}
