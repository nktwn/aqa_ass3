package handler

type Notification struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Kind      string `json:"kind"`
	IsRead    bool   `json:"is_read"`
	OrderID   *int64 `json:"order_id,omitempty"`
	CreatedAt string `json:"created_at"`
}

type ListResponse struct {
	Notifications []Notification `json:"notifications"`
}

type UnreadCountResponse struct {
	Count int `json:"count"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
