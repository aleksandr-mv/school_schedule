package migrator

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
)

type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

// NewMigrator создаёт мигратор, способный применять SQL-миграции из каталога migrationsDir
// к базе данных, представленной *sql.DB.
func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{db: db, migrationsDir: migrationsDir}
}

// Up применяет все доступные миграции в порядке возрастания версии.
// Логирует процесс и ошибки; в случае ошибки возвращает её вызывающему коду.
func (m *Migrator) Up(ctx context.Context) error {
	logger.Info(ctx,
		"🔄 [Migrations] Применяем миграции",
		zap.String("dir", m.migrationsDir),
	)

	if err := goose.Up(m.db, m.migrationsDir); err != nil {
		logger.Error(ctx,
			"❌ [Migrations] Не удалось применить миграции",
			zap.Error(err),
			zap.String("dir", m.migrationsDir),
		)
		return err
	}

	logger.Info(ctx,
		"✅ [Migrations] Все миграции успешно применены",
		zap.String("dir", m.migrationsDir),
	)
	return nil
}
