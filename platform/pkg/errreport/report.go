package errreport

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
)

// Report логирует ошибку и отправляет в трейс
// Используется в сервисном слое для обработки бизнес-ошибок
func Report(ctx context.Context, message string, err error) {
	// Логируем ошибку
	logger.Error(ctx, message, zap.Error(err))

	// Отправляем в трейс
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}
