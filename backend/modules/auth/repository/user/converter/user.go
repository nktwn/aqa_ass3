package converter

import (
	"diploma/modules/auth/model"
	modelRepo "diploma/modules/auth/repository/user/model"
)

func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Info:      ToUserInfoFromRepo(user.Info),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserInfoFromRepo(info *modelRepo.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:        info.Name,
		PhoneNumber: info.PhoneNumber,
		Role:        info.Role,
	}
}

func ToAuthUserFromRepo(user *modelRepo.AuthUser) *model.AuthUser {
	return &model.AuthUser{
		ID:             user.ID,
		Info:           ToUserInfoFromRepo(user.Info),
		HashedPassword: user.HashedPassword,
	}
}
