package model

const (
	CustomerRole = iota
	SupplierRole
	AdminRole
)

type User struct {
	ID          int64
	Name        string
	PhoneNumber string
	Role        int
	Address     Address
}

type Address struct {
	ID          int64
	Street      string
	Description string
	UserID      int64
}
