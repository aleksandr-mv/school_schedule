package helpers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInitViper_EnvOnlyMode проверяет ENV-only режим
func TestInitViper_EnvOnlyMode(t *testing.T) {
	// Очищаем состояние перед тестом
	Reset()

	// Тест ENV-only режима
	err := InitViper("")
	if err != nil {
		t.Errorf("InitViper(\"\") should not return error, got: %v", err)
	}

	// Проверяем, что GetViper возвращает nil в ENV-only режиме
	if viper := GetViper(); viper != nil {
		t.Error("GetViper() should return nil in ENV-only mode")
	}

	// Проверяем, что GetSection возвращает nil в ENV-only режиме
	if section := GetSection("app"); section != nil {
		t.Error("GetSection() should return nil in ENV-only mode")
	}
}

// TestInitViper_InvalidPath проверяет обработку неверных путей
func TestInitViper_InvalidPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "non-yaml file",
			path:    "config.json",
			wantErr: true,
		},
		{
			name:    "not allowed config file",
			path:    "config/custom.yaml",
			wantErr: true,
		},
		{
			name:    "non-existent file",
			path:    "config/development.yaml",
			wantErr: true, // файл не существует в тестовом окружении
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Очищаем состояние перед каждым тестом
			Reset()

			err := InitViper(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitViper() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestInitViper_ValidYAML проверяет успешную загрузку YAML
func TestInitViper_ValidYAML(t *testing.T) {
	// Создаем временный YAML файл
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "development.yaml")

	// Подменяем allowedConfigFiles для теста
	originalAllowed := allowedConfigFiles
	defer func() { allowedConfigFiles = originalAllowed }()

	// Временно разрешаем наш тестовый файл
	allowedConfigFiles = map[string]bool{
		configPath: true,
	}

	yamlContent := `
app:
  name: "test-service"
  environment: "test"
  version: "1.0.0"

grpc_server:
  host: "localhost"
  port: 50051
`

	if err := os.WriteFile(configPath, []byte(yamlContent), 0o644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Очищаем состояние перед тестом
	Reset()

	// Тестируем загрузку
	err := InitViper(configPath)
	if err != nil {
		t.Errorf("InitViper() should not return error for valid YAML, got: %v", err)
	}

	// Проверяем, что конфигурация загружена
	viper := GetViper()
	if viper == nil {
		t.Error("GetViper() should not return nil after successful InitViper()")
	}

	// Проверяем доступ к секциям
	appSection := GetSection("app")
	if appSection == nil {
		t.Error("GetSection(\"app\") should not return nil")
	}

	// Проверяем значения
	if appSection != nil {
		if name := appSection.GetString("name"); name != "test-service" {
			t.Errorf("Expected app.name = 'test-service', got: '%s'", name)
		}
	}
}

// TestReset проверяет функцию Reset
func TestReset(t *testing.T) {
	// Инициализируем в ENV-only режиме
	err := InitViper("")
	if err != nil {
		t.Errorf("InitViper(\"\") should not return error, got: %v", err)
	}

	// Вызываем Reset
	Reset()

	// Проверяем, что состояние сброшено
	if viper := GetViper(); viper != nil {
		t.Error("GetViper() should return nil after Reset()")
	}

	if section := GetSection("app"); section != nil {
		t.Error("GetSection() should return nil after Reset()")
	}
}

// BenchmarkGetSection проверяет производительность GetSection
func BenchmarkGetSection(b *testing.B) {
	// Инициализируем в ENV-only режиме для теста
	Reset()
	_ = InitViper("")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetSection("app")
	}
}

// BenchmarkGetViper проверяет производительность GetViper
func BenchmarkGetViper(b *testing.B) {
	// Инициализируем в ENV-only режиме для теста
	Reset()
	_ = InitViper("")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetViper()
	}
}

// ============================================================================
// ТЕСТЫ ДЛЯ ПОЛНОГО ПОКРЫТИЯ
// ============================================================================

// TestResolveConfigPath проверяет функцию ResolveConfigPath
func TestResolveConfigPath(t *testing.T) {
	// Этот тест пропускается из-за глобального состояния resolver
	// Используйте прямое тестирование ConfigPathResolver для полного покрытия
	t.Skip("Skipping due to global resolver state. Use ConfigPathResolver tests instead.")
}

// TestResolveConfigPath_WithFlags проверяет поведение с флагами в отдельном процессе
// Примечание: тестирование флагов сложно из-за глобального состояния flag.CommandLine
// В реальном приложении флаги парсятся один раз при запуске
func TestResolveConfigPath_WithFlags(t *testing.T) {
	// Этот тест пропускается, так как глобальный resolver может быть изменен
	// другими тестами. Используйте прямое тестирование ConfigPathResolver.
	t.Skip("Skipping due to global state. Use direct ConfigPathResolver tests instead.")
}

// TestInitViper_ParseError проверяет обработку ошибок парсинга YAML
func TestInitViper_ParseError(t *testing.T) {
	// Создаем временный файл с неверным YAML
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "development.yaml")

	// Подменяем allowedConfigFiles для теста
	originalAllowed := allowedConfigFiles
	defer func() { allowedConfigFiles = originalAllowed }()

	allowedConfigFiles = map[string]bool{
		configPath: true,
	}

	// Неверный YAML контент
	invalidYaml := `
app:
  name: "test-service"
  invalid_yaml: [unclosed array
`

	if err := os.WriteFile(configPath, []byte(invalidYaml), 0o644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	Reset()

	err := InitViper(configPath)
	if err == nil {
		t.Error("InitViper() should return error for invalid YAML")
	}

	// Проверяем, что состояние не изменилось при ошибке
	if viper := GetViper(); viper != nil {
		t.Error("GetViper() should return nil after failed InitViper()")
	}
}

// TestValidateRequired проверяет функцию validateRequired полностью
func TestValidateRequired_FullCoverage(t *testing.T) {
	tests := []struct {
		name     string
		yamlData string
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid grpc_server config",
			yamlData: `
grpc_server:
  host: "localhost"
  port: 50051
`,
			wantErr: false,
		},
		{
			name: "empty grpc_server host",
			yamlData: `
grpc_server:
  host: ""
  port: 50051
`,
			wantErr: true,
			errMsg:  "grpc_server.host is empty",
		},
		{
			name: "invalid grpc_server port",
			yamlData: `
grpc_server:
  host: "localhost"
  port: 0
`,
			wantErr: true,
			errMsg:  "grpc_server.port is empty or invalid",
		},
		{
			name: "valid mongo config",
			yamlData: `
database:
  mongo:
    connection:
      host: "localhost"
      port: 27017
      database: "testdb"
      user: "testuser"
      password: "testpass"
`,
			wantErr: false,
		},
		{
			name: "empty mongo host",
			yamlData: `
database:
  mongo:
    connection:
      host: ""
      port: 27017
      database: "testdb"
`,
			wantErr: true,
			errMsg:  "database.mongo.connection.host is empty",
		},
		{
			name: "invalid mongo port",
			yamlData: `
database:
  mongo:
    connection:
      host: "localhost"
      port: 0
      database: "testdb"
`,
			wantErr: true,
			errMsg:  "database.mongo.connection.port is empty or invalid",
		},
		{
			name: "empty mongo database",
			yamlData: `
database:
  mongo:
    connection:
      host: "localhost"
      port: 27017
      database: ""
`,
			wantErr: true,
			errMsg:  "database.mongo.connection.database is empty",
		},
		{
			name: "mismatched mongo credentials - user without password",
			yamlData: `
database:
  mongo:
    connection:
      host: "localhost"
      port: 27017
      database: "testdb"
      user: "testuser"
      password: ""
`,
			wantErr: true,
			errMsg:  "both database.mongo.connection.user and password must be set together",
		},
		{
			name: "mismatched mongo credentials - password without user",
			yamlData: `
database:
  mongo:
    connection:
      host: "localhost"
      port: 27017
      database: "testdb"
      user: ""
      password: "testpass"
`,
			wantErr: true,
			errMsg:  "both database.mongo.connection.user and password must be set together",
		},
		{
			name: "valid config without grpc_server and mongo",
			yamlData: `
app:
  name: "test-app"
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем временный файл
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, "development.yaml")

			// Подменяем allowedConfigFiles для теста
			originalAllowed := allowedConfigFiles
			defer func() { allowedConfigFiles = originalAllowed }()

			allowedConfigFiles = map[string]bool{
				configPath: true,
			}

			if err := os.WriteFile(configPath, []byte(tt.yamlData), 0o644); err != nil {
				t.Fatalf("Failed to create test config file: %v", err)
			}

			Reset()

			err := InitViper(configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitViper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("InitViper() error = %v, expected to contain %v", err, tt.errMsg)
				}
			}
		})
	}
}
