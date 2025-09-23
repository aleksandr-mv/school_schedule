package postgres

import (
	"os"
	"strings"
	"testing"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/helpers"
)

func TestNew_DefaultsOnly(t *testing.T) {
	// Очищаем состояние
	helpers.Reset()

	// Очищаем ENV переменные
	clearPostgresEnv()

	// Инициализируем без YAML файла (ENV-only режим)
	err := helpers.InitViper("")
	if err != nil {
		t.Fatalf("Failed to init viper: %v", err)
	}

	cfg, err := New()
	if err != nil {
		t.Errorf("New() should not fail with defaults, got error: %v", err)
	}

	// Проверяем defaults
	if cfg.Primary().Address() != "localhost:5432" {
		t.Errorf("Expected localhost:5432, got %s", cfg.Primary().Address())
	}

	if cfg.Primary().Database() != "postgres" {
		t.Errorf("Expected postgres, got %s", cfg.Primary().Database())
	}

	// Проверяем DSN формат
	dsn := cfg.PrimaryURI()
	if !strings.HasPrefix(dsn, "postgresql://") {
		t.Errorf("DSN should start with postgresql://, got: %s", dsn)
	}
}

func TestNew_WithYAMLConfig(t *testing.T) {
	// Тест с YAML пропускается из-за белого списка файлов в helpers
	// В реальном приложении YAML файлы находятся в config/ директории
	t.Skip("Skipping YAML test due to allowed files whitelist. Test YAML loading in integration tests.")
}

func TestNew_WithENVOverride(t *testing.T) {
	// Очищаем состояние
	helpers.Reset()
	clearPostgresEnv()

	// Устанавливаем ENV переменные
	os.Setenv("POSTGRES_USER", "envuser")
	os.Setenv("POSTGRES_PASSWORD", "envpass")
	os.Setenv("POSTGRES_DB", "envdb")
	defer clearPostgresEnv()

	err := helpers.InitViper("")
	if err != nil {
		t.Fatalf("Failed to init viper: %v", err)
	}

	cfg, err := New()
	if err != nil {
		t.Errorf("New() should not fail with ENV vars, got error: %v", err)
	}

	// Проверяем, что ENV переменные применились
	dsn := cfg.PrimaryURI()
	if !strings.Contains(dsn, "envuser") {
		t.Errorf("DSN should contain envuser, got: %s", dsn)
	}

	if !strings.Contains(dsn, "envpass") {
		t.Errorf("DSN should contain envpass, got: %s", dsn)
	}

	if !strings.Contains(dsn, "envdb") {
		t.Errorf("DSN should contain envdb, got: %s", dsn)
	}
}

func TestConnection_DSN_Format(t *testing.T) {
	tests := []struct {
		name     string
		conn     rawConnection
		expected string
	}{
		{
			name: "with password",
			conn: rawConnection{
				Host:     "localhost",
				Port:     5432,
				User:     "testuser",
				Password: "testpass",
				DB:       "testdb",
				AppName:  "testapp",
			},
			expected: "postgresql://testuser:testpass@localhost:5432/testdb?sslmode=disable&application_name=testapp",
		},
		{
			name: "without password",
			conn: rawConnection{
				Host:     "localhost",
				Port:     5432,
				User:     "testuser",
				Password: "",
				DB:       "testdb",
				AppName:  "testapp",
			},
			expected: "postgresql://testuser@localhost:5432/testdb?sslmode=disable&application_name=testapp",
		},
		{
			name: "without app name",
			conn: rawConnection{
				Host:     "localhost",
				Port:     5432,
				User:     "testuser",
				Password: "testpass",
				DB:       "testdb",
				AppName:  "",
			},
			expected: "postgresql://testuser:testpass@localhost:5432/testdb?sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connection := &Connection{raw: tt.conn}
			dsn := connection.DSN()

			if dsn != tt.expected {
				t.Errorf("DSN() = %v, expected %v", dsn, tt.expected)
			}

			// Проверяем, что URI() возвращает то же самое
			if connection.URI() != dsn {
				t.Errorf("URI() should return same as DSN()")
			}
		})
	}
}

func TestConfig_Replicas(t *testing.T) {
	cfg := &Config{
		replicaConns: []*Connection{
			{raw: rawConnection{Host: "replica1", Port: 5432, DB: "db1"}},
			{raw: rawConnection{Host: "replica2", Port: 5433, DB: "db2"}},
		},
	}

	replicas := cfg.Replicas()
	if len(replicas) != 2 {
		t.Errorf("Expected 2 replicas, got %d", len(replicas))
	}

	if replicas[0].Address() != "replica1:5432" {
		t.Errorf("Expected replica1:5432, got %s", replicas[0].Address())
	}
}

func TestConfig_ReplicaURI(t *testing.T) {
	tests := []struct {
		name     string
		replicas []*Connection
		wantURI  string
	}{
		{
			name:     "no replicas - fallback to primary",
			replicas: []*Connection{},
			wantURI:  "postgresql://postgres@localhost:5432/postgres?sslmode=disable&application_name=app",
		},
		{
			name: "single replica",
			replicas: []*Connection{
				{raw: rawConnection{Host: "replica1", Port: 5432, User: "user", DB: "db", AppName: "app"}},
			},
			wantURI: "postgresql://user@replica1:5432/db?sslmode=disable&application_name=app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				primaryConn:  &Connection{raw: defaultConnection()},
				replicaConns: tt.replicas,
			}

			uri := cfg.ReplicaURI()
			if tt.name == "no replicas - fallback to primary" {
				if uri != tt.wantURI {
					t.Errorf("ReplicaURI() = %v, expected %v", uri, tt.wantURI)
				}
			} else {
				// Для единственной реплики должен вернуться её URI
				if uri != tt.wantURI {
					t.Errorf("ReplicaURI() = %v, expected %v", uri, tt.wantURI)
				}
			}
		})
	}
}

func TestConfig_Pool(t *testing.T) {
	cfg := &Config{
		poolSettings: &Pool{raw: defaultPool()},
	}

	pool := cfg.Pool()
	if pool.MaxSize() != 10 {
		t.Errorf("Expected MaxSize 10, got %d", pool.MaxSize())
	}

	if pool.MinSize() != 2 {
		t.Errorf("Expected MinSize 2, got %d", pool.MinSize())
	}
}

func TestConnection_Methods(t *testing.T) {
	conn := &Connection{
		raw: rawConnection{
			Host:    "testhost",
			Port:    5432,
			User:    "testuser",
			DB:      "testdb",
			AppName: "testapp",
		},
	}

	if conn.Address() != "testhost:5432" {
		t.Errorf("Address() = %v, expected testhost:5432", conn.Address())
	}

	if conn.Database() != "testdb" {
		t.Errorf("Database() = %v, expected testdb", conn.Database())
	}

	if conn.ApplicationName() != "testapp" {
		t.Errorf("ApplicationName() = %v, expected testapp", conn.ApplicationName())
	}
}

// clearPostgresEnv очищает все PostgreSQL ENV переменные
func clearPostgresEnv() {
	envVars := []string{
		"DB_HOST", "DB_PORT", "POSTGRES_USER", "POSTGRES_PASSWORD",
		"POSTGRES_DB", "DB_APPLICATION_NAME",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}
