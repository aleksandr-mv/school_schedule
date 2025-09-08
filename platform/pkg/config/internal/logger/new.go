package logger

import "github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"

type module struct {
	loggerConfig *loggerConfig
}

func New() (contracts.LoggerModule, error) {
	loggerCfg, err := NewLoggerConfig()
	if err != nil {
		return nil, err
	}

	return &module{
		loggerConfig: loggerCfg,
	}, nil
}

// Методы логирования
func (m *module) Level() string { return m.loggerConfig.Level() }
func (m *module) AsJSON() bool  { return m.loggerConfig.AsJSON() }

// OTLP возвращает модуль OTLP конфигурации
func (m *module) OTLP() contracts.OTLPModule { return m.loggerConfig.OTLP() }
