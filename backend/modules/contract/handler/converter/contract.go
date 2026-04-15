package converter

import (
	"database/sql"
	apiModel "diploma/modules/contract/handler/model"
	serviceModel "diploma/modules/contract/model"
)

func ToAPI(c *serviceModel.Contract) *apiModel.ContractResponse {
	return &apiModel.ContractResponse{
		ID:          c.ID,
		Content:     c.Content,
		Status:      c.Status,
		SupplierSig: nullStringToString(c.SupplierSig),
		CustomerSig: nullStringToString(c.CustomerSig),
	}

}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
