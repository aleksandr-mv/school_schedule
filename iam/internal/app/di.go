package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	v1 "github.com/aleksandr-mv/school_schedule/iam/internal/api/auth/v1"
	userAPI "github.com/aleksandr-mv/school_schedule/iam/internal/api/user/v1"
	grpcClient "github.com/aleksandr-mv/school_schedule/iam/internal/client/grpc"
	rbacV1 "github.com/aleksandr-mv/school_schedule/iam/internal/client/grpc/rbac"
	"github.com/aleksandr-mv/school_schedule/iam/internal/repository"
	"github.com/aleksandr-mv/school_schedule/iam/internal/repository/notification"
	sessionRepo "github.com/aleksandr-mv/school_schedule/iam/internal/repository/session"
	userRepo "github.com/aleksandr-mv/school_schedule/iam/internal/repository/user"
	"github.com/aleksandr-mv/school_schedule/iam/internal/service"
	authService "github.com/aleksandr-mv/school_schedule/iam/internal/service/auth"
	userService "github.com/aleksandr-mv/school_schedule/iam/internal/service/user"
	userProducerService "github.com/aleksandr-mv/school_schedule/iam/internal/service/user_producer"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/cache"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/cache/builder"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/cache/redis"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/closer"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	grpcclient "github.com/aleksandr-mv/school_schedule/platform/pkg/grpc/client"
	kafkaBuilder "github.com/aleksandr-mv/school_schedule/platform/pkg/kafka/builder"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/migrator"
	authv1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/auth/v1"
	userV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user/v1"
	generatedRbacV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user_role/v1"
)

type diContainer struct {
	cfg contracts.Provider

	authV1 authv1.AuthServiceServer
	userV1 userV1.UserServiceServer

	authService service.AuthService
	userService service.UserService

	rbacClient grpcClient.RBACClient

	userRepository         repository.UserRepository
	sessionRepository      repository.SessionRepository
	notificationRepository repository.NotificationRepository

	userProducerService service.UserProducerService

	postgresPool *pgxpool.Pool
	redisClient  cache.RedisClient
	migrator     *migrator.Migrator
}

func NewDiContainer(cfg contracts.Provider) *diContainer {
	return &diContainer{cfg: cfg}
}

func (d *diContainer) AuthV1API(ctx context.Context) (authv1.AuthServiceServer, error) {
	if d.authV1 == nil {
		authService, err := d.AuthService(ctx)
		if err != nil {
			return nil, err
		}

		d.authV1 = v1.NewAPI(authService)
	}

	return d.authV1, nil
}

func (d *diContainer) UserV1API(ctx context.Context) (userV1.UserServiceServer, error) {
	if d.userV1 == nil {
		userService, err := d.UserService(ctx)
		if err != nil {
			return nil, err
		}

		d.userV1 = userAPI.NewAPI(userService)
	}

	return d.userV1, nil
}

func (d *diContainer) AuthService(ctx context.Context) (service.AuthService, error) {
	if d.authService == nil {
		userRepo, err := d.UserRepository(ctx)
		if err != nil {
			return nil, err
		}

		notificationRepo, err := d.NotificationRepository(ctx)
		if err != nil {
			return nil, err
		}

		sessionRepo, err := d.SessionRepository(ctx)
		if err != nil {
			return nil, err
		}

		rbacClient, err := d.RBACClient(ctx)
		if err != nil {
			return nil, err
		}

		d.authService = authService.NewService(userRepo, notificationRepo, sessionRepo, rbacClient, d.cfg.Session().TTL())
	}

	return d.authService, nil
}

func (d *diContainer) UserService(ctx context.Context) (service.UserService, error) {
	if d.userService == nil {
		userRepo, err := d.UserRepository(ctx)
		if err != nil {
			return nil, err
		}

		notificationRepo, err := d.NotificationRepository(ctx)
		if err != nil {
			return nil, err
		}

		userProducerService, err := d.UserProducerService(ctx)
		if err != nil {
			return nil, err
		}

		d.userService = userService.NewService(userRepo, notificationRepo, userProducerService)
	}

	return d.userService, nil
}

func (d *diContainer) UserRepository(ctx context.Context) (repository.UserRepository, error) {
	if d.userRepository == nil {
		pool, err := d.PostgresPool(ctx)
		if err != nil {
			return nil, err
		}

		d.userRepository = userRepo.NewRepository(pool)
	}

	return d.userRepository, nil
}

func (d *diContainer) NotificationRepository(ctx context.Context) (repository.NotificationRepository, error) {
	if d.notificationRepository == nil {
		pool, err := d.PostgresPool(ctx)
		if err != nil {
			return nil, err
		}

		d.notificationRepository = notification.NewRepository(pool)
	}

	return d.notificationRepository, nil
}

func (d *diContainer) PostgresPool(ctx context.Context) (*pgxpool.Pool, error) {
	if d.postgresPool == nil {
		dsn := d.cfg.Postgres().PrimaryURI()
		if dsn == "" {
			return nil, fmt.Errorf("postgres dsn is empty")
		}

		pgCfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return nil, fmt.Errorf("invalid postgres dsn: %w", err)
		}

		pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
		if err != nil {
			return nil, fmt.Errorf("create postgres pool failed: %w", err)
		}

		closer.AddNamed("PostgreSQL pool", func(ctx context.Context) error {
			logger.Info(ctx, "🐘 [Shutdown] Закрытие PostgreSQL pool")
			pool.Close()
			return nil
		})

		logger.Info(ctx, "✅ [Database] Пул PostgreSQL создан")
		d.postgresPool = pool
	}

	return d.postgresPool, nil
}

func (d *diContainer) SessionRepository(ctx context.Context) (repository.SessionRepository, error) {
	if d.sessionRepository == nil {
		redis, err := d.RedisClient(ctx)
		if err != nil {
			return nil, err
		}

		d.sessionRepository = sessionRepo.NewRepository(redis)
	}

	return d.sessionRepository, nil
}

func (d *diContainer) RedisClient(ctx context.Context) (cache.RedisClient, error) {
	if d.redisClient == nil {
		redisBuilder := builder.NewRedisBuilder(d.cfg.Redis())
		redigoPool, err := redisBuilder.BuildPool()
		if err != nil {
			return nil, err
		}

		pool := d.cfg.Redis().Pool()
		client := redis.NewClient(redigoPool, nil, pool.PoolTimeout())

		if err = client.Ping(ctx); err != nil {
			return nil, fmt.Errorf("redis ping failed: %w", err)
		}

		closer.AddNamed("Redis client", func(ctx context.Context) error {
			logger.Info(ctx, "🔴 [Shutdown] Закрытие Redis клиента")
			return redigoPool.Close()
		})

		logger.Info(ctx, "✅ [Cache] Redis клиент создан")
		d.redisClient = client
	}

	return d.redisClient, nil
}

func (d *diContainer) RunMigrations(ctx context.Context) error {
	if d.migrator != nil {
		return nil
	}

	pool, err := d.PostgresPool(ctx)
	if err != nil {
		return fmt.Errorf("postgres pool failed: %w", err)
	}
	if pool == nil {
		return fmt.Errorf("cannot run migrations: DB pool is nil")
	}

	mig := migrator.NewMigrator(stdlib.OpenDB(*pool.Config().ConnConfig), d.cfg.App().MigrationsDir())

	if err = mig.Up(ctx); err != nil {
		return fmt.Errorf("migrations failed: %w", err)
	}

	d.migrator = mig

	return nil
}

func (d *diContainer) RBACClient(ctx context.Context) (grpcClient.RBACClient, error) {
	if d.rbacClient == nil {
		conn, addr, err := d.dialServiceConn(ctx, "rbac")
		if err != nil {
			return nil, fmt.Errorf("failed to dial rbac service: %w", err)
		}

		generatedRbacClient := generatedRbacV1.NewUserRoleServiceClient(conn)
		d.rbacClient = rbacV1.NewClient(generatedRbacClient)

		closer.AddNamed("gRPC RBAC conn", func(ctx context.Context) error {
			logger.Info(ctx, "🔐 [Shutdown] Закрытие gRPC RBAC соединения")
			return conn.Close()
		})

		logger.Info(ctx, "✅ [gRPC] Подключение к RBAC установлено", zap.String("address", addr))
	}

	return d.rbacClient, nil
}

func (d *diContainer) dialServiceConn(ctx context.Context, serviceName string) (*grpc.ClientConn, string, error) {
	svc, ok := d.cfg.Services().Get(serviceName)
	if !ok {
		logger.Error(ctx, "❌ [Config] Сервис не настроен", zap.String("service", serviceName))
		return nil, "", fmt.Errorf("%s service not configured", serviceName)
	}

	addr := svc.Address()
	limits := d.cfg.Transport().GRPCClientLimits()

	conn, err := grpcclient.NewClient(addr, limits.MaxRecvMsgSize(), limits.MaxSendMsgSize(), limits.Timeout())
	if err != nil {
		logger.Error(ctx, "❌ [gRPC] Не удалось подключиться к сервису", zap.String("service", serviceName), zap.String("address", addr), zap.Error(err))
		return nil, "", fmt.Errorf("connect to %s failed: %w", serviceName, err)
	}

	return conn, addr, nil
}

func (d *diContainer) UserProducerService(ctx context.Context) (service.UserProducerService, error) {
	if d.userProducerService == nil {
		if !d.cfg.Kafka().IsEnabled() {
			logger.Info(ctx, "⚠️ [Kafka] Kafka отключен, создаем no-op producer")
			d.userProducerService = userProducerService.NewNoOpService()
			return d.userProducerService, nil
		}

		kafkaBuilder := kafkaBuilder.NewKafkaBuilder(d.cfg.Kafka())
		userCreatedProducer, err := kafkaBuilder.BuildProducer("user_created")
		if err != nil {
			return nil, fmt.Errorf("failed to build user_created producer: %w", err)
		}

		d.userProducerService = userProducerService.NewService(userCreatedProducer)

		closer.AddNamed("Kafka user_created producer", func(ctx context.Context) error {
			logger.Info(ctx, "📤 [Shutdown] Закрытие Kafka user_created producer")
			return nil // Producer закрывается автоматически
		})

		logger.Info(ctx, "✅ [Kafka] UserCreated producer создан")
	}

	return d.userProducerService, nil
}
