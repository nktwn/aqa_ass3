package model

type RegisterInput struct {
	Name            string `json:"name" example:"John Doe"`
	PhoneNumber     string `json:"phone_number" example:"+123456789"`
	Password        string `json:"password" example:"secure123"`
	ConfirmPassword string `json:"confirm_password" example:"secure123"`
}

type RegisterResponse struct {
	ID int64 `json:"id" example:"1"`
}

type LoginInput struct {
	PhoneNumber string `json:"phone_number" example:"+123456789"`
	Password    string `json:"password" example:"secure123"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Err string `json:"error"`
}
