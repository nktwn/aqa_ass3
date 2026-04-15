package service

import (
	"context"
	"diploma/modules/notification/model"
	"diploma/pkg/service"
)

type Repository interface {
	ListByUser(ctx context.Context, userID int64) ([]model.Notification, error)
	UnreadCount(ctx context.Context, userID int64) (int, error)
	MarkRead(ctx context.Context, userID, notificationID int64) error
	MarkAllRead(ctx context.Context, userID int64) error
	Create(ctx context.Context, input *model.CreateNotificationInput) error
}

type Service struct {
	service.BaseService
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		BaseService: service.NewBaseService("notification"),
		repo:        repo,
	}
}

func (s *Service) ListByUser(ctx context.Context, userID int64) ([]model.Notification, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *Service) UnreadCount(ctx context.Context, userID int64) (int, error) {
	return s.repo.UnreadCount(ctx, userID)
}

func (s *Service) MarkRead(ctx context.Context, userID, notificationID int64) error {
	return s.repo.MarkRead(ctx, userID, notificationID)
}

func (s *Service) MarkAllRead(ctx context.Context, userID int64) error {
	return s.repo.MarkAllRead(ctx, userID)
}

func (s *Service) Notify(ctx context.Context, userID int64, title, message, kind string, orderID *int64) error {
	return s.repo.Create(ctx, &model.CreateNotificationInput{
		UserID:  userID,
		Title:   title,
		Message: message,
		Kind:    kind,
		OrderID: orderID,
	})
}
