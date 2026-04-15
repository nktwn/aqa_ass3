package repository

import (
	"context"
	"diploma/modules/notification/model"
	"diploma/pkg/client/db"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	db db.Client
}

func NewRepository(dbClient db.Client) *Repository {
	return &Repository{db: dbClient}
}

func (r *Repository) ListByUser(ctx context.Context, userID int64) ([]model.Notification, error) {
	builder := sq.
		Select("id", "user_id", "title", "message", "kind", "is_read", "order_id", "created_at").
		From("notifications").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
		Limit(50).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list notifications query: %w", err)
	}

	q := db.Query{Name: "notification_repository.ListByUser", QueryRaw: query}
	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var notification model.Notification
		if err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Title,
			&notification.Message,
			&notification.Kind,
			&notification.IsRead,
			&notification.OrderID,
			&notification.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	return notifications, rows.Err()
}

func (r *Repository) UnreadCount(ctx context.Context, userID int64) (int, error) {
	builder := sq.
		Select("COUNT(*)").
		From("notifications").
		Where(sq.Eq{"user_id": userID, "is_read": false}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build unread count query: %w", err)
	}

	q := db.Query{Name: "notification_repository.UnreadCount", QueryRaw: query}
	var count int
	if err := r.db.DB().QueryRowContext(ctx, q, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}
	return count, nil
}

func (r *Repository) MarkRead(ctx context.Context, userID, notificationID int64) error {
	builder := sq.
		Update("notifications").
		Set("is_read", true).
		Where(sq.Eq{"id": notificationID, "user_id": userID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build mark read query: %w", err)
	}

	q := db.Query{Name: "notification_repository.MarkRead", QueryRaw: query}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}

func (r *Repository) MarkAllRead(ctx context.Context, userID int64) error {
	builder := sq.
		Update("notifications").
		Set("is_read", true).
		Where(sq.Eq{"user_id": userID, "is_read": false}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build mark all read query: %w", err)
	}

	q := db.Query{Name: "notification_repository.MarkAllRead", QueryRaw: query}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}
	return nil
}

func (r *Repository) Create(ctx context.Context, input *model.CreateNotificationInput) error {
	builder := sq.
		Insert("notifications").
		Columns("user_id", "title", "message", "kind", "order_id").
		Values(input.UserID, input.Title, input.Message, input.Kind, input.OrderID).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build create notification query: %w", err)
	}

	q := db.Query{Name: "notification_repository.Create", QueryRaw: query}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}
