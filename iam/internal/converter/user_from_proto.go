package converter

import (
	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	userV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user/v1"
)

func UserFromProto(req *userV1.RegisterRequest) *model.CreateUser {
	return &model.CreateUser{
		Login:    req.Info.Login,
		Email:    req.Info.Email,
		Password: req.Password,
	}
}
