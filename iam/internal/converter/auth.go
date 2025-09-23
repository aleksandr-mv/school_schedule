package converter

import (
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	authV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/auth/v1"
)

func LoginFromProto(req *authV1.LoginRequest) *model.LoginCredentials {
	return &model.LoginCredentials{
		Login:    req.Login,
		Password: req.Password,
	}
}
