package metric

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
)

// Metrics управляет состоянием метрик и их инициализацией
type Metrics struct {
	initOnce      sync.Once
	exporter      *otlpmetricgrpc.Exporter
	meterProvider *sdkmetric.MeterProvider
	logger        Logger
	config        *Config
}

// Глобальный экземпляр для использования по всему приложению
var globalMetrics = NewWithLogger(&logger.NoopLogger{})

// SetLogger устанавливает логгер для глобального экземпляра метрик
func SetLogger(logger Logger) {
	globalMetrics.SetLogger(logger)
}

// getMetricName генерирует полное имя метрики с namespace и appName
func getMetricName(metricName string) string {
	if globalMetrics.config == nil {
		return metricName
	}
	return globalMetrics.config.namespace + "_" + globalMetrics.config.appName + "_" + metricName
}

// NewWithLogger создает новый экземпляр Metrics с указанным логгером
func NewWithLogger(logger Logger) *Metrics {
	return &Metrics{
		logger: logger,
	}
}

// SetLogger устанавливает логгер для Metrics
func (m *Metrics) SetLogger(logger Logger) {
	m.logger = logger
}

// Init инициализирует OpenTelemetry MeterProvider и все инструменты метрик
func Init(ctx context.Context, opts ...Option) error {
	return globalMetrics.Init(ctx, opts...)
}

// Init инициализирует OpenTelemetry MeterProvider и все инструменты метрик для конкретного экземпляра
func (m *Metrics) Init(ctx context.Context, opts ...Option) error {
	var initErr error

	m.initOnce.Do(func() {
		cfg := defaultConfig()
		for _, opt := range opts {
			opt(cfg)
		}

		// Сохраняем конфигурацию в структуре для доступа к namespace и appName
		m.config = cfg

		if !cfg.enable {
			return // Метрики отключены - meterProvider остается nil
		}

		initErr = m.initMetrics(ctx, cfg)
		if initErr != nil {
			return
		}
	})

	return initErr
}

// initMetrics выполняет фактическую инициализацию метрик для конкретного экземпляра
func (m *Metrics) initMetrics(ctx context.Context, cfg *Config) error {
	var err error

	m.exporter, err = otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(cfg.endpoint),
		otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()),
		otlpmetricgrpc.WithTimeout(cfg.timeout),
	)
	if err != nil {
		m.logger.Error(ctx, "❌ [Metrics] Ошибка создания OTLP экспортера", zap.Error(err))
		return errors.Wrap(err, "failed to create OTLP exporter")
	}

	// 2. Создаем ресурс с метаданными о сервисе
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", cfg.name),
			attribute.String("service.version", cfg.version),
			attribute.String("deployment.environment", cfg.environment),
		),
	)
	if err != nil {
		m.logger.Error(ctx, "❌ [Metrics] Ошибка создания ресурса", zap.Error(err))
		return errors.Wrap(err, "failed to create resource")
	}

	// 3. Создаем MeterProvider с периодическим reader
	m.meterProvider = sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				m.exporter,
				sdkmetric.WithInterval(cfg.exportInterval),
			),
		),
	)

	otel.SetMeterProvider(m.meterProvider)

	return nil
}

// GetMeterProvider возвращает текущий провайдер метрик
func GetMeterProvider() *sdkmetric.MeterProvider {
	return globalMetrics.GetMeterProvider()
}

// GetMeterProvider возвращает текущий провайдер метрик для конкретного экземпляра
func (m *Metrics) GetMeterProvider() *sdkmetric.MeterProvider {
	if m.meterProvider == nil {
		return sdkmetric.NewMeterProvider()
	}
	return m.meterProvider
}

// Shutdown закрывает провайдер метрик и экспортер в правильном порядке.
// MeterProvider должен закрываться первым, чтобы корректно завершить отправку данных в экспортер.
func Shutdown(ctx context.Context, timeout time.Duration) error {
	return globalMetrics.Shutdown(ctx, timeout)
}

// Shutdown закрывает провайдер метрик и экспортер в правильном порядке для конкретного экземпляра.
func (m *Metrics) Shutdown(ctx context.Context, timeout time.Duration) error {
	if m.meterProvider == nil && m.exporter == nil {
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var err error

	// 1. Сначала закрываем MeterProvider - он может продолжать отправлять данные в экспортер
	if m.meterProvider != nil {
		err = m.meterProvider.Shutdown(shutdownCtx)
		if err != nil {
			return errors.Wrap(err, "failed to shutdown meter provider")
		}
		m.meterProvider = nil
	}

	// 2. Затем закрываем экспортер - после того как MeterProvider завершил отправку данных
	if m.exporter != nil {
		err = m.exporter.Shutdown(shutdownCtx)
		if err != nil {
			return errors.Wrap(err, "failed to shutdown exporter")
		}
		m.exporter = nil
	}

	return nil
}
