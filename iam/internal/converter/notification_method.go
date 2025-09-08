package converter

import (
	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	commonV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/common/v1"
)

func NotificationMethodsToProto(methods []*model.NotificationMethod) []*commonV1.NotificationMethod {
	if methods == nil {
		return nil
	}

	protoMethods := make([]*commonV1.NotificationMethod, len(methods))
	for i, method := range methods {
		protoMethods[i] = &commonV1.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		}
	}

	return protoMethods
}
