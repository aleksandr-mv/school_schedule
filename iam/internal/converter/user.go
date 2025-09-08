package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	commonV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/common/v1"
)

func UserToProto(user *model.User) *commonV1.User {
	notificationMethods := NotificationMethodsToProto(user.NotificationMethods)

	protoUser := &commonV1.User{
		Id: user.ID.String(),
		Info: &commonV1.UserInfo{
			Login:               user.Login,
			Email:               user.Email,
			NotificationMethods: notificationMethods,
		},
		CreatedAt: timestamppb.New(user.CreatedAt),
	}

	if user.UpdatedAt != nil {
		protoUser.UpdatedAt = timestamppb.New(*user.UpdatedAt)
	}

	return protoUser
}
