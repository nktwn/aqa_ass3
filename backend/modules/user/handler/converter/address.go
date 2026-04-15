package converter

import (
	modelApi "diploma/modules/user/handler/model"
	"diploma/modules/user/model"
)

func ToServiceAddressFromApi(userID int64, input modelApi.Address) model.Address {
	return model.Address{
		Street:      input.Street,
		Description: input.Description,
		UserID:      userID,
	}
}

func ToApiAddressFromService(input model.Address) modelApi.Address {
	return modelApi.Address{
		ID:          input.ID,
		Street:      input.Street,
		Description: input.Description,
	}
}

func ToApiAddressListFromService(input []model.Address) []modelApi.Address {
	addressList := make([]modelApi.Address, 0, len(input))
	for _, address := range input {
		addressList = append(addressList, ToApiAddressFromService(address))
	}
	return addressList
}
