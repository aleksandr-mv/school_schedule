package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/cache"
	cacheBuilder "github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/cache/builder"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/closer"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/contracts"
	consumerBuilder "github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/kafka/consumer"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/migrator"
	permissionAPI "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/api/permission/v1"
	roleAPI "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/api/role/v1"
	rolePermissionAPI "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/api/role_permission/v1"
	userRoleAPI "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/api/user_role/v1"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository"
	enrichedRoleRepo "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/enriched_role"
	permissionRepo "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/permission"
	roleRepo "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/role"
	rolePermissionRepo "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/role_permission"
	userRoleRepo "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/repository/user_role"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/service"
	permissionService "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/service/permission"
	roleService "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/service/role"
	rolePermissionService "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/service/role_permission"
	userConsumerService "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/service/user_consumer"
	userRoleService "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/service/user_role"
	permissionV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/permission/v1"
	roleV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/role/v1"
	rolePermissionV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/role_permission/v1"
	userRoleV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/user_role/v1"
)

type diContainer struct {
	cfg contracts.Provider

	roleV1           roleV1.RoleServiceServer
	permissionV1     permissionV1.PermissionServiceServer
	rolePermissionV1 rolePermissionV1.RolePermissionServiceServer
	userRoleV1       userRoleV1.UserRoleServiceServer

	roleService           service.RoleServiceInterface
	permissionService     service.PermissionServiceInterface
	rolePermissionService service.RolePermissionServiceInterface
	userRoleService       service.UserRoleServiceInterface
	userConsumerService   service.UserConsumerService

	roleRepository           repository.RoleRepository
	permissionRepository     repository.PermissionRepository
	userRoleRepository       repository.UserRoleRepository
	rolePermissionRepository repository.RolePermissionRepository
	enrichedRoleRepository   repository.EnrichedRoleRepository

	postgresWritePool *pgxpool.Pool
	postgresReadPool  *pgxpool.Pool
	redisClient       cache.RedisClient
	migrator          *migrator.Migrator
}

func NewDiContainer(cfg contracts.Provider) *diContainer {
	return &diContainer{cfg: cfg}
}

func (d *diContainer) RoleV1API(ctx context.Context) (roleV1.RoleServiceServer, error) {
	if d.roleV1 == nil {
		roleService, err := d.RoleService(ctx)
		if err != nil {
			return nil, err
		}

		d.roleV1 = roleAPI.NewAPI(roleService)
	}

	return d.roleV1, nil
}

func (d *diContainer) PermissionV1API(ctx context.Context) (permissionV1.PermissionServiceServer, error) {
	if d.permissionV1 == nil {
		permissionService, err := d.PermissionService(ctx)
		if err != nil {
			return nil, err
		}

		d.permissionV1 = permissionAPI.NewAPI(permissionService)
	}

	return d.permissionV1, nil
}

func (d *diContainer) RolePermissionV1API(ctx context.Context) (rolePermissionV1.RolePermissionServiceServer, error) {
	if d.rolePermissionV1 == nil {
		rolePermissionService, err := d.RolePermissionService(ctx)
		if err != nil {
			return nil, err
		}

		d.rolePermissionV1 = rolePermissionAPI.NewAPI(rolePermissionService)
	}

	return d.rolePermissionV1, nil
}

func (d *diContainer) UserRoleV1API(ctx context.Context) (userRoleV1.UserRoleServiceServer, error) {
	if d.userRoleV1 == nil {
		userRoleService, err := d.UserRoleService(ctx)
		if err != nil {
			return nil, err
		}

		d.userRoleV1 = userRoleAPI.NewAPI(userRoleService)
	}

	return d.userRoleV1, nil
}

func (d *diContainer) RoleService(ctx context.Context) (service.RoleServiceInterface, error) {
	if d.roleService == nil {
		roleRepo, err := d.RoleRepository(ctx)
		if err != nil {
			return nil, err
		}

		rolePermissionRepo, err := d.RolePermissionRepository(ctx)
		if err != nil {
			return nil, err
		}

		enrichedRoleRepo, err := d.EnrichedRoleRepository(ctx)
		if err != nil {
			return nil, err
		}

		enrichedRoleTTL := d.cfg.Session().TTL()

		d.roleService = roleService.NewService(roleRepo, rolePermissionRepo, enrichedRoleRepo, enrichedRoleTTL)
	}

	return d.roleService, nil
}

func (d *diContainer) PermissionService(ctx context.Context) (service.PermissionServiceInterface, error) {
	if d.permissionService == nil {
		permissionRepo, err := d.PermissionRepository(ctx)
		if err != nil {
			return nil, err
		}

		d.permissionService = permissionService.NewService(permissionRepo)
	}

	return d.permissionService, nil
}

func (d *diContainer) RolePermissionService(ctx context.Context) (service.RolePermissionServiceInterface, error) {
	if d.rolePermissionService == nil {
		rolePermissionRepo, err := d.RolePermissionRepository(ctx)
		if err != nil {
			return nil, err
		}

		d.rolePermissionService = rolePermissionService.NewService(rolePermissionRepo)
	}

	return d.rolePermissionService, nil
}

func (d *diContainer) UserRoleService(ctx context.Context) (service.UserRoleServiceInterface, error) {
	if d.userRoleService == nil {
		userRoleRepo, err := d.UserRoleRepository(ctx)
		if err != nil {
			return nil, err
		}

		roleService, err := d.RoleService(ctx)
		if err != nil {
			return nil, err
		}

		d.userRoleService = userRoleService.NewService(userRoleRepo, roleService)
	}

	return d.userRoleService, nil
}

func (d *diContainer) RoleRepository(ctx context.Context) (repository.RoleRepository, error) {
	if d.roleRepository == nil {
		writePool, err := d.PostgresWritePool(ctx)
		if err != nil {
			return nil, err
		}

		readPool, err := d.PostgresReadPool(ctx)
		if err != nil {
			return nil, err
		}

		d.roleRepository = roleRepo.NewRepository(writePool, readPool)
	}

	return d.roleRepository, nil
}

func (d *diContainer) PermissionRepository(ctx context.Context) (repository.PermissionRepository, error) {
	if d.permissionRepository == nil {
		readPool, err := d.PostgresReadPool(ctx)
		if err != nil {
			return nil, err
		}

		d.permissionRepository = permissionRepo.NewRepository(readPool)
	}

	return d.permissionRepository, nil
}

func (d *diContainer) UserRoleRepository(ctx context.Context) (repository.UserRoleRepository, error) {
	if d.userRoleRepository == nil {
		writePool, err := d.PostgresWritePool(ctx)
		if err != nil {
			return nil, err
		}

		readPool, err := d.PostgresReadPool(ctx)
		if err != nil {
			return nil, err
		}

		d.userRoleRepository = userRoleRepo.NewRepository(writePool, readPool)
	}

	return d.userRoleRepository, nil
}

func (d *diContainer) RolePermissionRepository(ctx context.Context) (repository.RolePermissionRepository, error) {
	if d.rolePermissionRepository == nil {
		writePool, err := d.PostgresWritePool(ctx)
		if err != nil {
			return nil, err
		}

		readPool, err := d.PostgresReadPool(ctx)
		if err != nil {
			return nil, err
		}

		d.rolePermissionRepository = rolePermissionRepo.NewRepository(writePool, readPool)
	}

	return d.rolePermissionRepository, nil
}

func (d *diContainer) EnrichedRoleRepository(ctx context.Context) (repository.EnrichedRoleRepository, error) {
	if d.enrichedRoleRepository == nil {
		redisClient, err := d.RedisClient(ctx)
		if err != nil {
			return nil, err
		}

		d.enrichedRoleRepository = enrichedRoleRepo.NewRepository(redisClient)
	}

	return d.enrichedRoleRepository, nil
}

func (d *diContainer) PostgresWritePool(ctx context.Context) (*pgxpool.Pool, error) {
	if d.postgresWritePool == nil {
		dsn := d.cfg.Postgres().PrimaryURI()
		if dsn == "" {
			return nil, fmt.Errorf("postgres primary dsn is empty")
		}

		pgCfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return nil, fmt.Errorf("invalid postgres primary dsn: %w", err)
		}

		pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
		if err != nil {
			return nil, fmt.Errorf("create postgres write pool failed: %w", err)
		}

		closer.AddNamed("PostgreSQL write pool", func(ctx context.Context) error {
			logger.Info(ctx, "🐘 [Shutdown] Закрытие PostgreSQL write pool")
			pool.Close()
			return nil
		})

		logger.Info(ctx, "✅ [Database] PostgreSQL write pool создан")
		d.postgresWritePool = pool
	}

	return d.postgresWritePool, nil
}

func (d *diContainer) PostgresReadPool(ctx context.Context) (*pgxpool.Pool, error) {
	if d.postgresReadPool == nil {
		dsn := d.cfg.Postgres().ReplicaURI()
		if dsn == "" {
			// Fallback на primary если реплика не настроена
			dsn = d.cfg.Postgres().PrimaryURI()
		}

		if dsn == "" {
			return nil, fmt.Errorf("postgres replica dsn is empty")
		}

		pgCfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return nil, fmt.Errorf("invalid postgres replica dsn: %w", err)
		}

		pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
		if err != nil {
			return nil, fmt.Errorf("create postgres read pool failed: %w", err)
		}

		closer.AddNamed("PostgreSQL read pool", func(ctx context.Context) error {
			logger.Info(ctx, "🐘 [Shutdown] Закрытие PostgreSQL read pool")
			pool.Close()
			return nil
		})

		logger.Info(ctx, "✅ [Database] PostgreSQL read pool создан")
		d.postgresReadPool = pool
	}

	return d.postgresReadPool, nil
}

func (d *diContainer) PostgresPool(ctx context.Context) (*pgxpool.Pool, error) {
	return d.PostgresWritePool(ctx)
}

func (d *diContainer) RedisClient(ctx context.Context) (cache.RedisClient, error) {
	if d.redisClient == nil {
		redisBuilder := cacheBuilder.NewRedisBuilder(d.cfg.Redis())
		client, err := redisBuilder.BuildClient()
		if err != nil {
			return nil, err
		}

		if err = client.Ping(ctx); err != nil {
			return nil, fmt.Errorf("redis ping failed: %w", err)
		}

		closer.AddNamed("Redis client", func(ctx context.Context) error {
			logger.Info(ctx, "🔴 [Shutdown] Закрытие Redis клиента")
			return nil // Redis client закрывается автоматически
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

func (d *diContainer) UserConsumerService(ctx context.Context) (service.UserConsumerService, error) {
	if d.userConsumerService == nil {
		if !d.cfg.Kafka().IsEnabled() {
			logger.Info(ctx, "⚠️ [Kafka] Kafka отключен, создаем no-op consumer")
			d.userConsumerService = userConsumerService.NewNoOpService()
			return d.userConsumerService, nil
		}

		builder := consumerBuilder.NewBuilder(d.cfg.Kafka())
		userCreatedConsumer, err := builder.BuildConsumer("user_created")
		if err != nil {
			return nil, fmt.Errorf("get user created consumer: %w", err)
		}

		userRoleService, err := d.UserRoleService(ctx)
		if err != nil {
			return nil, fmt.Errorf("get user role service: %w", err)
		}

		d.userConsumerService = userConsumerService.NewService(
			userCreatedConsumer,
			userRoleService,
		)

		closer.AddNamed("Kafka user_created consumer", func(ctx context.Context) error {
			logger.Info(ctx, "📥 [Shutdown] Закрытие Kafka user_created consumer")
			return nil // Consumer закрывается автоматически
		})

		logger.Info(ctx, "✅ [Kafka] UserCreated consumer создан")
	}

	return d.userConsumerService, nil
}
