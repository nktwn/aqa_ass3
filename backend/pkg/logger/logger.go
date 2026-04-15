package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Global logger instance
	log  *zap.Logger
	once sync.Once
)

// Logger returns the global logger instance
func Logger() *zap.Logger {
	once.Do(func() {
		var err error
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		log, err = config.Build(
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
		if err != nil {
			// If we can't create a logger, create a basic one
			log = zap.NewExample()
		}
	})
	return log
}

// Initialize sets up the logger with the given environment
func Initialize(environment string) error {
	var config zap.Config

	if environment == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var err error
	newLogger, err := config.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return err
	}

	// Replace the global logger
	log = newLogger
	zap.ReplaceGlobals(log)
	return nil
}

// Field creates a field for structured logging
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// Debug logs a debug message with fields
func Debug(msg string, fields ...zap.Field) {
	Logger().Debug(msg, fields...)
}

// Info logs an info message with fields
func Info(msg string, fields ...zap.Field) {
	Logger().Info(msg, fields...)
}

// Warn logs a warning message with fields
func Warn(msg string, fields ...zap.Field) {
	Logger().Warn(msg, fields...)
}

// Error logs an error message with fields
func Error(msg string, fields ...zap.Field) {
	Logger().Error(msg, fields...)
}

// Fatal logs a fatal message with fields and exits
func Fatal(msg string, fields ...zap.Field) {
	Logger().Fatal(msg, fields...)
	os.Exit(1)
}

// Sync flushes any buffered log entries
func Sync() error {
	return Logger().Sync()
}
