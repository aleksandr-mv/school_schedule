package converter

import (
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	repoModel "github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository/model"
)

func ToRepoUser(user *model.User) *repoModel.User {
	repoUser := &repoModel.User{
		ID:           user.ID,
		Login:        user.Login,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
	}

	if user.UpdatedAt != nil {
		repoUser.UpdatedAt = user.UpdatedAt
	}

	return repoUser
}

func ToDomainUser(user *repoModel.User) *model.User {
	domainUser := &model.User{
		ID:           user.ID,
		Login:        user.Login,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
	}

	if user.UpdatedAt != nil {
		domainUser.UpdatedAt = user.UpdatedAt
	}

	return domainUser
}
