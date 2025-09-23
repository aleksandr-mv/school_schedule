package converter

import (
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	authV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/auth/v1"
)

func WhoAMIToProto(i *model.WhoAMI) *authV1.WhoamiResponse {
	return &authV1.WhoamiResponse{
		Session:              SessionToProto(&i.Session),
		User:                 UserToProto(&i.User),
		RolesWithPermissions: RoleWithPermissionsSliceToProto(i.RolesWithPermissions),
	}
}
