// Package helpers содержит общие helper функции для загрузки конфигураций.
//
// Этот пакет реализует принцип KISS (Keep It Simple, Stupid) и предоставляет
// минимальный API для работы с Viper конфигурацией. Все загрузчики используют
// единый экземпляр Viper через функции этого пакета.
//
// Особенности загрузки:
//   - Один глобальный Viper экземпляр для всех загрузчиков
//   - Простые функции без сложных абстракций
//   - Поддержка fallback на ENV-переменные (через теги env в структурах)
//   - Централизованная инициализация из YAML файла
//   - Перед парсингом YAML выполняется os.ExpandEnv для подстановки ${VAR} из окружения
//   - После загрузки выполняется базовая валидация обязательных полей
package helpers

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// ============================================================================
// ГЛОБАЛЬНОЕ СОСТОЯНИЕ И ПЕРЕМЕННЫЕ
// ============================================================================

// globalViper хранит общий экземпляр Viper для всех загрузчиков конфигурации.
// Может быть nil, если YAML файл недоступен (ENV-only режим).
var globalViper *viper.Viper

// allowedConfigFiles — белый список допустимых конфигурационных файлов.
var allowedConfigFiles = map[string]bool{
	"config/development.yaml": true,
	"config/production.yaml":  true,
	"config/staging.yaml":     true,
	"config/test.yaml":        true,
}

// ============================================================================
// ФУНКЦИИ ИНИЦИАЛИЗАЦИИ И УПРАВЛЕНИЯ
// ============================================================================

// InitViper создаёт и инициализирует глобальный экземпляр Viper из YAML файла.
//
// Параметры:
//   - path: путь к YAML файлу конфигурации
//
// Поведение:
//   - Если path пустой, устанавливает globalViper в nil (ENV-only режим)
//   - Читает файл как текст, выполняет os.ExpandEnv для подстановки ${VAR}
//   - Парсит результат через Viper (ReadConfig)
//   - Выполняет базовую валидацию обязательных полей (см. validateRequired)
//   - При ошибке возвращает её и не сохраняет глобальное состояние
//
// Возвращает:
//   - nil при успешной инициализации или ENV-only режиме
//   - error при ошибке чтения YAML файла
func InitViper(path string) error {
	if path == "" {
		globalViper = nil
		return nil
	}

	cleanPath := filepath.Clean(path)
	if !strings.HasSuffix(cleanPath, ".yaml") && !strings.HasSuffix(cleanPath, ".yml") {
		return fmt.Errorf("only YAML config files are allowed: %s", cleanPath)
	}

	if !allowedConfigFiles[cleanPath] {
		return fmt.Errorf("config file not in allowed list: %s", cleanPath)
	}

	v := viper.New()
	v.SetConfigType("yaml")

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		globalViper = nil
		return fmt.Errorf("failed to read config file %s: %w", cleanPath, err)
	}

	expanded := os.ExpandEnv(string(data))
	if err := v.ReadConfig(bytes.NewReader([]byte(expanded))); err != nil {
		globalViper = nil
		return fmt.Errorf("failed to parse config file %s: %w", cleanPath, err)
	}

	if err := validateRequired(v); err != nil {
		globalViper = nil
		return err
	}

	globalViper = v
	return nil
}

// validateRequired выполняет базовую проверку, что критичные поля не пустые
// и имеют валидные значения. Это страховка на случай, если os.ExpandEnv
// подставил пустые строки из незаданных переменных окружения.
func validateRequired(v *viper.Viper) error {
	if gs := v.Sub("grpc_server"); gs != nil {
		if host := gs.GetString("host"); host == "" {
			return fmt.Errorf("invalid config: grpc_server.host is empty")
		}

		if port := gs.GetInt("port"); port <= 0 {
			return fmt.Errorf("invalid config: grpc_server.port is empty or invalid")
		}
	}

	if db := v.Sub("database"); db != nil {
		if mongo := db.Sub("mongo"); mongo != nil {
			if conn := mongo.Sub("connection"); conn != nil {
				if host := conn.GetString("host"); host == "" {
					return fmt.Errorf("invalid config: database.mongo.connection.host is empty")
				}

				if port := conn.GetInt("port"); port <= 0 {
					return fmt.Errorf("invalid config: database.mongo.connection.port is empty or invalid")
				}

				if dbn := conn.GetString("database"); dbn == "" {
					return fmt.Errorf("invalid config: database.mongo.connection.database is empty")
				}

				user := conn.GetString("user")
				pass := conn.GetString("password")
				if (user == "" && pass != "") || (user != "" && pass == "") {
					return fmt.Errorf("invalid config: both database.mongo.connection.user and password must be set together or both empty")
				}
			}
		}
	}

	return nil
}

// ============================================================================
// ФУНКЦИИ ДОСТУПА К КОНФИГУРАЦИИ
// ============================================================================

// GetSection возвращает подсекцию конфигурации по имени.
//
// Параметры:
//   - sectionName: имя секции в YAML файле (например, "logger", "database")
//
// Возвращает:
//   - *viper.Viper: подсекцию конфигурации для unmarshaling
//   - nil: если globalViper не инициализирован или секция не найдена
//
// Использование:
//
//	section := GetSection("logger")
//	if section != nil {
//	    section.Unmarshal(&config)
//	}
func GetSection(sectionName string) *viper.Viper {
	if globalViper == nil {
		return nil
	}
	return globalViper.Sub(sectionName)
}

// GetViper возвращает полный экземпляр Viper для unmarshaling всей конфигурации.
//
// Используется для загрузки сложных структур, таких как массивы сервисов,
// которые требуют доступа ко всему дереву конфигурации.
//
// Возвращает:
//   - *viper.Viper: полный экземпляр конфигурации
//   - nil: если YAML файл не был загружен (ENV-only режим)
//
// Использование:
//
//	if v := GetViper(); v != nil {
//	    v.Unmarshal(&fullConfig)
//	}
func GetViper() *viper.Viper {
	return globalViper
}
