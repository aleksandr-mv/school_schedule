package v1

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
)

func (api *API) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	sessionID, err := api.extractSessionID(req)
	if err != nil {
		logger.Error(ctx, "❌ [External Auth] Не удалось извлечь session ID", zap.Error(err))
		return api.denyRequest("Missing or invalid session", 401), nil
	}

	whoami, err := api.whoAMIService.Whoami(ctx, sessionID)
	if err != nil {
		logger.Error(ctx, "❌ [External Auth] Невалидная сессия", zap.Error(err))
		return api.denyRequest("Invalid session", 401), nil
	}

	return api.allowRequest(whoami), nil
}
