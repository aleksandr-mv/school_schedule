package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/closer"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/app"
)

func main() {
	if err := logger.InitDefault(); err != nil {
		panic(fmt.Errorf("failed to init default logger: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	cfg, err := config.Load(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось загрузить конфигурацию", zap.Error(err))
		return
	}

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx, cfg)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось инициализировать приложение", zap.Error(err))
		return
	}

	if err = a.Start(appCtx); err != nil {
		logger.Error(appCtx, "❌ Не удалось запустить приложение", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
