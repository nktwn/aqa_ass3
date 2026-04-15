package service

import (
	"context"
	"diploma/pkg/logger"

	"go.uber.org/zap"
)

// BaseService provides common functionality for all services
type BaseService struct {
	log *zap.Logger
}

// NewBaseService creates a new base service with logging
func NewBaseService(serviceName string) BaseService {
	return BaseService{
		log: logger.Logger().With(zap.String("service", serviceName)),
	}
}

// Log returns the service's logger
func (s *BaseService) Log() *zap.Logger {
	if s.log == nil {
		s.log = logger.Logger()
	}
	return s.log
}

// LogInfo logs an info message with context and fields
func (s *BaseService) LogInfo(ctx context.Context, msg string, fields ...zap.Field) {
	s.Log().Info(msg, append(fields, extractContextFields(ctx)...)...)
}

// LogError logs an error message with context and fields
func (s *BaseService) LogError(ctx context.Context, msg string, err error, fields ...zap.Field) {
	s.Log().Error(msg, append(fields, append(extractContextFields(ctx), zap.Error(err))...)...)
}

// LogDebug logs a debug message with context and fields
func (s *BaseService) LogDebug(ctx context.Context, msg string, fields ...zap.Field) {
	s.Log().Debug(msg, append(fields, extractContextFields(ctx)...)...)
}

// LogWarn logs a warning message with context and fields
func (s *BaseService) LogWarn(ctx context.Context, msg string, fields ...zap.Field) {
	s.Log().Warn(msg, append(fields, extractContextFields(ctx)...)...)
}

// extractContextFields extracts relevant fields from context
func extractContextFields(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	// Add request ID if present
	if reqID := ctx.Value("request_id"); reqID != nil {
		fields = append(fields, zap.Any("request_id", reqID))
	}

	// Add user ID if present
	if userID := ctx.Value("user_id"); userID != nil {
		fields = append(fields, zap.Any("user_id", userID))
	}

	return fields
}
