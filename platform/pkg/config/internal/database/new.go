package database

import (
	"fmt"
	"strings"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/database/mongo"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/internal/database/postgres"
)

type dbConfig struct {
	Pg contracts.PostgresConfig
	Mg contracts.MongoConfig
}

// NewDatabaseConfig создает агрегированную конфигурацию баз данных.
// Пробует создать оба конфига (PostgreSQL и MongoDB).
func New() (contracts.DatabaseConfig, error) {
	var c dbConfig

	// Попытка создания PostgreSQL конфига
	if pg, err := postgres.NewPostgresConfig(); err == nil {
		c.Pg = pg
	}

	// Попытка создания MongoDB конфига
	if mg, err := mongo.NewMongoConfig(); err == nil {
		c.Mg = mg
	}

	return &c, nil
}

func (c *dbConfig) HasPostgres() bool { return c.Pg != nil }
func (c *dbConfig) HasMongo() bool    { return c.Mg != nil }

func (c *dbConfig) PostgresDSN() string {
	if c.Pg == nil {
		return ""
	}

	return c.Pg.URI()
}

func (c *dbConfig) MongoURI() string {
	if c.Mg == nil {
		return ""
	}
	return c.Mg.URI()
}

func (c *dbConfig) PostgresConnection() (contracts.DBConnection, bool) {
	if c.Pg == nil {
		return nil, false
	}
	return c.Pg.Database(), true
}

func (c *dbConfig) MongoConnection() (contracts.DBConnection, bool) {
	if c.Mg == nil {
		return nil, false
	}
	return c.Mg.Database(), true
}

func (c *dbConfig) PostgresPool() (contracts.PostgresPoolConfig, bool) {
	if c.Pg == nil {
		return nil, false
	}
	return c.Pg.Pool(), true
}

func (c *dbConfig) MongoPool() (contracts.MongoPoolConfig, bool) {
	if c.Mg == nil {
		return nil, false
	}
	return c.Mg.Pool(), true
}

func (c *dbConfig) String() string {
	parts := []string{}
	if c.HasPostgres() {
		parts = append(parts, "PostgresDSN="+c.PostgresDSN())
	}

	if c.HasMongo() {
		parts = append(parts, "MongoURI="+c.MongoURI())
	}

	return fmt.Sprintf("DatabaseConfig{%s}", strings.Join(parts, ", "))
}
