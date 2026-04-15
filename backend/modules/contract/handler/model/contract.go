package model

type SignRequest struct {
	ContractID int64  `json:"contract_id" binding:"required"`
	Signature  string `json:"signature" binding:"required"`
}

type ContractResponse struct {
	ID          int64  `json:"id"`
	Content     string `json:"content"`
	Status      int    `json:"status"`
	SupplierSig string `json:"supplier_signature,omitempty"`
	CustomerSig string `json:"customer_signature,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
