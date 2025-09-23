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

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/closer"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/contracts"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/grpc/health"
	platformgrpc "github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/grpc/server"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/metric"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/tracing"
	permissionV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/permission/v1"
	roleV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/role/v1"
	rolePermissionV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/role_permission/v1"
	userRoleV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/user_role/v1"
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
	return app.Run(ctx)
}

func (app *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := app.runKafkaConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("kafka consumer crashed: %w", err)
		}
	}()

	go func() {
		if err := app.runGRPCServer(ctx); err != nil {
			errCh <- fmt.Errorf("gRPC server crashed: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		cancel()
		<-ctx.Done()
		return err
	}

	return nil
}

func (app *App) runKafkaConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 [Kafka] Запуск Kafka consumer для UserCreated событий")

	consumerService, err := app.diContainer.UserConsumerService(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user consumer service: %w", err)
	}

	if err = consumerService.Run(ctx); err != nil {
		return fmt.Errorf("failed to run user consumer: %w", err)
	}

	return nil
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
		logger.WithOTLPEnable(app.cfg.Logger().Enable()),
		logger.WithOTLPEndpoint(app.cfg.Logger().Endpoint()),
		logger.WithOTLPTimeout(time.Duration(app.cfg.Logger().ShutdownTimeout())*time.Second),
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
	addr := app.cfg.GRPC().Address()
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
		app.cfg.GRPC().Timeout(),
		app.cfg.GRPC().MaxRecvMsgSize(),
		app.cfg.GRPC().MaxSendMsgSize(),
		tracing.UnaryServerInterceptor(app.cfg.App().Name()),
		metric.UnaryServerInterceptor(ctx, app.cfg.Metric().BucketBoundaries()),
	)

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		logger.Info(ctx, "⚡ [Shutdown] Остановка gRPC сервера")
		app.grpcServer.GracefulStop()
		return nil
	})

	roleAPI, err := app.diContainer.RoleV1API(ctx)
	if err != nil {
		return fmt.Errorf("create role v1 api: %w", err)
	}

	permissionAPI, err := app.diContainer.PermissionV1API(ctx)
	if err != nil {
		return fmt.Errorf("create permission v1 api: %w", err)
	}

	rolePermissionAPI, err := app.diContainer.RolePermissionV1API(ctx)
	if err != nil {
		return fmt.Errorf("create role_permission v1 api: %w", err)
	}

	userRoleAPI, err := app.diContainer.UserRoleV1API(ctx)
	if err != nil {
		return fmt.Errorf("create user_role v1 api: %w", err)
	}

	reflection.Register(app.grpcServer)
	health.RegisterService(app.grpcServer)
	roleV1.RegisterRoleServiceServer(app.grpcServer, roleAPI)
	permissionV1.RegisterPermissionServiceServer(app.grpcServer, permissionAPI)
	rolePermissionV1.RegisterRolePermissionServiceServer(app.grpcServer, rolePermissionAPI)
	userRoleV1.RegisterUserRoleServiceServer(app.grpcServer, userRoleAPI)

	logger.Info(ctx, "✅ [App] Role API инициализирован")
	logger.Info(ctx, "✅ [App] Permission API инициализирован")
	logger.Info(ctx, "✅ [App] RolePermission API инициализирован")
	logger.Info(ctx, "✅ [App] UserRole API инициализирован")
	logger.Info(ctx, "✅ [gRPC] Сервер успешно инициализирован")

	return nil
}

func (app *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx,
		"🚀 [gRPC] RBAC сервис слушает адрес",
		zap.String("address", app.cfg.GRPC().Address()),
	)

	err := app.grpcServer.Serve(app.listener)
	if err != nil && !errors.Is(err, net.ErrClosed) {
		return err
	}

	return nil
}
