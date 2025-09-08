package app

import "github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"

type module struct {
	appConfig *appConfig
}

func New() (contracts.AppModule, error) {
	appCfg, err := NewAppConfig()
	if err != nil {
		return nil, err
	}

	return &module{
		appConfig: appCfg,
	}, nil
}

// Методы сервиса
func (m *module) Name() string        { return m.appConfig.Name() }
func (m *module) Environment() string { return m.appConfig.Environment() }
func (m *module) Version() string     { return m.appConfig.Version() }

// Методы приложения
func (m *module) MigrationsDir() string { return m.appConfig.MigrationsDir() }
func (m *module) SwaggerPath() string   { return m.appConfig.SwaggerPath() }
func (m *module) SwaggerUIPath() string { return m.appConfig.SwaggerUIPath() }
