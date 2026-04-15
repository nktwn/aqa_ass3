package converter

import (
	modelApi "diploma/modules/auth/handler/model"
	"diploma/modules/auth/model"
)

func ToServiceFromRegisterInput(user modelApi.RegisterInput) *model.AuthUser {
	return &model.AuthUser{
		Info: &model.UserInfo{
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
		},
		Password: user.Password,
	}
}
