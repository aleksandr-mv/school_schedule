package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	commonV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/common/v1"
)

func SessionToProto(session *model.Session) *commonV1.Session {
	return &commonV1.Session{
		Id:        session.ID.String(),
		CreatedAt: timestamppb.New(session.CreatedAt),
		UpdatedAt: timestamppb.New(session.UpdatedAt),
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}
}
