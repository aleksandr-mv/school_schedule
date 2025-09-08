package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/closer"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/grpc/health"
	platformgrpc "github.com/aleksandr-mv/school_schedule/platform/pkg/grpc/server"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/metric"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
	authV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/auth/v1"
	userV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user/v1"
)

type App struct {
	cfg         contracts.Provider
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

func New(ctx context.Context, cfg contracts.Provider) (*App, error) {
	app := &App{cfg: cfg}

	if err := app.initDeps(ctx); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Start(ctx context.Context) error {
	return app.runGRPCServer(ctx)
}

func (app *App) initDeps(ctx context.Context) error {
	steps := []func(context.Context) error{
		app.initDI,
		app.initLogger,
		app.initTracer,
		app.initMetrics,
		app.initCloser,
		app.initDatabase,
		app.initMigrations,
		app.initListener,
		app.initGRPCServer,
	}
	for _, step := range steps {
		if err := step(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (app *App) initDI(_ context.Context) error {
	app.diContainer = NewDiContainer(app.cfg)
	return nil
}

func (app *App) initLogger(ctx context.Context) error {
	return logger.Reinit(ctx,
		logger.WithLevel(app.cfg.Logger().Level()),
		logger.WithJSON(app.cfg.Logger().AsJSON()),
		logger.WithName(app.cfg.App().Name()),
		logger.WithEnvironment(app.cfg.App().Environment()),
		logger.WithOTLPEnable(app.cfg.Logger().OTLP().Enable()),
		logger.WithOTLPEndpoint(app.cfg.Logger().OTLP().Endpoint()),
		logger.WithOTLPTimeout(time.Duration(app.cfg.Logger().OTLP().ShutdownTimeout())*time.Second),
	)
}

func (app *App) initTracer(ctx context.Context) error {
	if err := tracing.Init(ctx,
		tracing.WithName(app.cfg.App().Name()),
		tracing.WithEnvironment(app.cfg.App().Environment()),
		tracing.WithVersion(app.cfg.App().Version()),
		tracing.WithEnable(app.cfg.Tracing().Enable()),
		tracing.WithEndpoint(app.cfg.Tracing().Endpoint()),
		tracing.WithTimeout(app.cfg.Tracing().Timeout()),
		tracing.WithSampleRatio(app.cfg.Tracing().SampleRatio()),
		tracing.WithRetryEnabled(app.cfg.Tracing().RetryEnabled()),
		tracing.WithRetryInitialInterval(app.cfg.Tracing().RetryInitialInterval()),
		tracing.WithRetryMaxInterval(app.cfg.Tracing().RetryMaxInterval()),
		tracing.WithRetryMaxElapsedTime(app.cfg.Tracing().RetryMaxElapsedTime()),
		tracing.WithTraceContext(app.cfg.Tracing().EnableTraceContext()),
		tracing.WithBaggage(app.cfg.Tracing().EnableBaggage()),
	); err != nil {
		return fmt.Errorf("failed to init tracer: %w", err)
	}

	tracing.SetLogger(logger.Logger())

	logger.Info(ctx, "✅ [Tracing] Трейсинг успешно инициализирован",
		zap.String("serviceName", app.cfg.App().Name()))

	closer.AddNamed("Tracer", func(ctx context.Context) error {
		logger.Info(ctx, "🔍 [Shutdown] Закрытие Tracer provider")
		return tracing.Shutdown(ctx, app.cfg.Tracing().ShutdownTimeout())
	})

	return nil
}

func (app *App) initMetrics(ctx context.Context) error {
	if err := metric.Init(ctx,
		metric.WithName(app.cfg.App().Name()),
		metric.WithEnvironment(app.cfg.App().Environment()),
		metric.WithVersion(app.cfg.App().Version()),
		metric.WithEnable(app.cfg.Metric().Enable()),
		metric.WithEndpoint(app.cfg.Metric().Endpoint()),
		metric.WithTimeout(app.cfg.Metric().Timeout()),
		metric.WithNamespace(app.cfg.Metric().Namespace()),
		metric.WithAppName(app.cfg.Metric().AppName()),
		metric.WithExportInterval(app.cfg.Metric().ExportInterval()),
		metric.WithShutdownTimeout(app.cfg.Metric().ShutdownTimeout()),
	); err != nil {
		return fmt.Errorf("failed to init metrics: %w", err)
	}

	metric.SetLogger(logger.Logger())

	logger.Info(ctx, "✅ [Metrics] Метрики успешно инициализированы",
		zap.String("serviceName", app.cfg.App().Name()))

	closer.AddNamed("Metrics", func(ctx context.Context) error {
		logger.Info(ctx, "📊 [Shutdown] Закрытие Metrics provider")
		return metric.Shutdown(ctx, app.cfg.Metric().ShutdownTimeout())
	})

	return nil
}

func (app *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (app *App) initDatabase(ctx context.Context) error {
	if _, err := app.diContainer.PostgresPool(ctx); err != nil {
		return fmt.Errorf("postgres pool failed: %w", err)
	}

	if _, err := app.diContainer.RedisClient(ctx); err != nil {
		return fmt.Errorf("redis client failed: %w", err)
	}

	return nil
}

func (app *App) initMigrations(ctx context.Context) error {
	return app.diContainer.RunMigrations(ctx)
}

func (app *App) initListener(_ context.Context) error {
	addr := app.cfg.Transport().GRPC().Address()
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		logger.Info(ctx, "🔌 [Shutdown] Закрытие TCP listener")
		if err = listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			return err
		}
		return nil
	})

	app.listener = listener

	return nil
}

func (app *App) initGRPCServer(ctx context.Context) error {
	app.grpcServer = platformgrpc.New(ctx,
		app.cfg.Transport().GRPC().Timeout(),
		app.cfg.Transport().GRPCClientLimits().MaxRecvMsgSize(),
		app.cfg.Transport().GRPCClientLimits().MaxSendMsgSize(),
		tracing.UnaryServerInterceptor(app.cfg.App().Name()),
		metric.UnaryServerInterceptor(ctx, app.cfg.Metric().BucketBoundaries()))

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		logger.Info(ctx, "⚡ [Shutdown] Остановка gRPC сервера")
		app.grpcServer.GracefulStop()
		return nil
	})

	authAPI, err := app.diContainer.AuthV1API(ctx)
	if err != nil {
		return fmt.Errorf("create auth v1 api: %w", err)
	}

	userAPI, err := app.diContainer.UserV1API(ctx)
	if err != nil {
		return fmt.Errorf("create user v1 api: %w", err)
	}

	reflection.Register(app.grpcServer)
	health.RegisterService(app.grpcServer)
	authV1.RegisterAuthServiceServer(app.grpcServer, authAPI)
	userV1.RegisterUserServiceServer(app.grpcServer, userAPI)

	logger.Info(ctx, "✅ [App] Auth API инициализирован")
	logger.Info(ctx, "✅ [App] User API инициализирован")
	logger.Info(ctx, "✅ [gRPC] Сервер успешно инициализирован")

	return nil
}

func (app *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx,
		"🚀 [gRPC] IAM сервис слушает адрес",
		zap.String("address", app.cfg.Transport().GRPC().Address()),
	)

	err := app.grpcServer.Serve(app.listener)
	if err != nil {
		return err
	}

	return nil
}
