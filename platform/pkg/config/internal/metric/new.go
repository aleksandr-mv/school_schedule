package metric

import (
	"time"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

type module struct {
	config *rawMetricConfig
}

// New создает модуль конфигурации метрик.
func New() (contracts.MetricModule, error) {
	metricCfg, err := NewMetricConfig()
	if err != nil {
		return nil, err
	}

	return &module{config: metricCfg}, nil
}

// Методы MetricModule
func (m *module) Enable() bool {
	return m.config.EnableFlag
}

func (m *module) Endpoint() string {
	return m.config.EndpointAddr
}

func (m *module) Timeout() time.Duration {
	return m.config.TimeoutDur
}

func (m *module) Namespace() string {
	return m.config.NamespaceStr
}

func (m *module) AppName() string {
	return m.config.AppNameStr
}

func (m *module) ExportInterval() time.Duration {
	return m.config.ExportIntervalDur
}

func (m *module) ShutdownTimeout() time.Duration {
	return m.config.ShutdownTimeoutDur
}

func (m *module) BucketBoundaries() []float64 {
	return m.config.BucketBoundariesSlice
}
