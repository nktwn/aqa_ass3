package model

import "time"

type Notification struct {
	ID        int64
	UserID    int64
	Title     string
	Message   string
	Kind      string
	IsRead    bool
	OrderID   *int64
	CreatedAt time.Time
}

type CreateNotificationInput struct {
	UserID  int64
	Title   string
	Message string
	Kind    string
	OrderID *int64
}
