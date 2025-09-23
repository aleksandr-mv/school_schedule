package producer

import (
	"fmt"
	"strings"

	"github.com/IBM/sarama"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/contracts"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/kafka"
)

// Builder создает Kafka producer'ы из конфигурации
type Builder struct {
	kafkaConfig contracts.KafkaConfig
	logger      kafka.Logger
	metrics     kafka.Metrics
}

// NewBuilder создает новый builder
func NewBuilder(kafkaConfig contracts.KafkaConfig) *Builder {
	return &Builder{
		kafkaConfig: kafkaConfig,
	}
}

// WithLogger устанавливает логгер
func (b *Builder) WithLogger(logger kafka.Logger) *Builder {
	b.logger = logger
	return b
}

// WithMetrics устанавливает метрики
func (b *Builder) WithMetrics(metrics kafka.Metrics) *Builder {
	b.metrics = metrics
	return b
}

// BuildProducer создает синхронный producer
func (b *Builder) BuildProducer(producerName string) (*producer, error) {
	if !b.kafkaConfig.IsEnabled() {
		return nil, fmt.Errorf("kafka is not enabled")
	}

	producerConfig, exists := b.kafkaConfig.Producers().GetProducer(producerName)
	if !exists {
		return nil, fmt.Errorf("producer '%s' not found", producerName)
	}

	if !producerConfig.IsEnabled() {
		return nil, fmt.Errorf("producer '%s' is disabled", producerName)
	}

	if producerConfig.Topic() == "" {
		return nil, fmt.Errorf("producer '%s' topic cannot be empty", producerName)
	}

	brokers := strings.Split(b.kafkaConfig.Brokers(), ",")
	config := b.createSaramaConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Надежность для sync

	syncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama sync producer: %w", err)
	}

	return NewProducer(syncProducer, producerConfig.Topic(), b.logger), nil
}

// BuildAsyncProducer создает асинхронный producer
func (b *Builder) BuildAsyncProducer(producerName string) (*asyncProducer, error) {
	if !b.kafkaConfig.IsEnabled() {
		return nil, fmt.Errorf("kafka is not enabled")
	}

	producerConfig, exists := b.kafkaConfig.Producers().GetProducer(producerName)
	if !exists {
		return nil, fmt.Errorf("producer '%s' not found", producerName)
	}

	if !producerConfig.IsEnabled() {
		return nil, fmt.Errorf("producer '%s' is disabled", producerName)
	}

	if producerConfig.Topic() == "" {
		return nil, fmt.Errorf("producer '%s' topic cannot be empty", producerName)
	}

	brokers := strings.Split(b.kafkaConfig.Brokers(), ",")
	config := b.createSaramaConfig()

	asyncProducer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama async producer: %w", err)
	}

	return NewAsyncProducer(
		asyncProducer,
		producerConfig.Topic(),
		b.logger,
		b.metrics,
	), nil
}

// createSaramaConfig создает базовую конфигурацию Sarama
func (b *Builder) createSaramaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Retry.Max = 3
	config.Producer.Compression = sarama.CompressionSnappy
	return config
}
