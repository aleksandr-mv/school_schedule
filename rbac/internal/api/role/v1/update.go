package v1

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/converter"
	roleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/role/v1"
)

func (api *API) Update(ctx context.Context, req *roleV1.UpdateRequest) (*emptypb.Empty, error) {
	role, err := converter.UpdateRoleToDomain(req)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка парсинга параметров запроса", zap.Error(err))
		return nil, mapError(err)
	}

	err = api.roleService.Update(ctx, role)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка обновления роли", zap.Error(err))
		return nil, mapError(err)
	}

	return &emptypb.Empty{}, nil
}
