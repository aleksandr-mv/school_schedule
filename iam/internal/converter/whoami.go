package converter

import (
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	authV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/auth/v1"
	commonV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/common/v1"
)

func WhoAMIToProto(i *model.WhoAMI) *authV1.WhoamiResponse {
	return &authV1.WhoamiResponse{
		Info: &commonV1.WhoamiInfo{
			User:                 UserToProto(&i.User),
			RolesWithPermissions: RoleWithPermissionsSliceToProto(i.RolesWithPermissions),
		},
	}
}

func WhoamiToProto(i *model.WhoAMI) *commonV1.WhoamiInfo {
	return &commonV1.WhoamiInfo{
		User:                 UserToProto(&i.User),
		RolesWithPermissions: RoleWithPermissionsSliceToProto(i.RolesWithPermissions),
	}
}
