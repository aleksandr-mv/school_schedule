package contracts

// AppModule содержит основную конфигурацию приложения и сервиса.
// Включает базовые настройки сервиса, которые используются во всех компонентах.
type AppModule interface {
	// Основные настройки сервиса
	Name() string
	Environment() string
	Version() string

	// Конфигурация приложения
	MigrationsDir() string
	SwaggerPath() string
	SwaggerUIPath() string
}
