package converter

import (
	modelApi "diploma/modules/user/handler/model"
	"diploma/modules/user/model"
)

func ToApiUserFromService(input model.User) modelApi.User {
	return modelApi.User{
		ID:          input.ID,
		Name:        input.Name,
		PhoneNumber: input.PhoneNumber,
		Address:     ToApiAddressFromService(input.Address),
	}
}

func ToUserFromApi(userID int64, input modelApi.UpdateUserProfileRequest) model.User {
	return model.User{
		ID:          userID,
		Name:        input.Name,
		PhoneNumber: input.PhoneNumber,
	}
}
