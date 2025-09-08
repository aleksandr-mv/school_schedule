package tracing

import (
	"time"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

type module struct {
	config *rawTracingConfig
}

// New создает модуль конфигурации трассировки.
func New() (contracts.TracingModule, error) {
	tracingCfg, err := NewTracingConfig()
	if err != nil {
		return nil, err
	}

	return &module{config: tracingCfg}, nil
}

// Методы TracingModule
func (m *module) Enable() bool                        { return m.config.Enable }
func (m *module) Endpoint() string                    { return m.config.Endpoint }
func (m *module) Timeout() time.Duration              { return m.config.Timeout }
func (m *module) SampleRatio() int                    { return m.config.SampleRatio }
func (m *module) RetryEnabled() bool                  { return m.config.RetryEnabled }
func (m *module) RetryInitialInterval() time.Duration { return m.config.RetryInitialInterval }
func (m *module) RetryMaxInterval() time.Duration     { return m.config.RetryMaxInterval }
func (m *module) RetryMaxElapsedTime() time.Duration  { return m.config.RetryMaxElapsedTime }
func (m *module) EnableTraceContext() bool            { return m.config.EnableTraceContext }
func (m *module) EnableBaggage() bool                 { return m.config.EnableBaggage }
func (m *module) ShutdownTimeout() time.Duration      { return m.config.ShutdownTimeout }
