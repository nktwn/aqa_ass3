package handler

type SupplierHandler struct {
	service ISupplierService
}

func NewHandler(service ISupplierService) *SupplierHandler {
	return &SupplierHandler{service: service}
}

type ISupplierService interface {
}
