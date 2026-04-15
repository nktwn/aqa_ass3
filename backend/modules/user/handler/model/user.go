package model

type User struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	PhoneNumber string  `json:"phone_number"`
	Address     Address `json:"address"`
}

type GetUserProfileResponse struct {
	User User `json:"user"`
}

type UpdateUserProfileRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type UserRoleResponse struct {
	Role int `json:"role"`
}

type StatusResponse struct {
	Status string `json:"status"`
}
