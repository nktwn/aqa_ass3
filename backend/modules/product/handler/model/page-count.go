package model

type PageCountInput struct {
	PageSize int `json:"pageSize"`
}

type PageCountResponse struct {
	Pages int `json:"pages"`
}
