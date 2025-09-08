package builder

import (
	"fmt"
	"strings"

	"github.com/IBM/sarama"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/kafka"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/kafka/producer"
)

// ProducerBuilder упрощает создание Kafka продюсеров
type ProducerBuilder struct {
	kafkaConfig contracts.KafkaConfig
}

// NewProducerBuilder создает новый билдер для продюсеров
func NewProducerBuilder(kafkaConfig contracts.KafkaConfig) *ProducerBuilder {
	return &ProducerBuilder{
		kafkaConfig: kafkaConfig,
	}
}

// BuildProducer создает готового к использованию продюсера
func (b *ProducerBuilder) BuildProducer(producerName string) (kafka.Producer, error) {
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

	saramaSyncProducer, err := b.createSaramaSyncProducer()
	if err != nil {
		return nil, fmt.Errorf("failed to create producer for '%s': %w", producerName, err)
	}

	return producer.NewSyncProducer(saramaSyncProducer, producerConfig.Topic()), nil
}

// createSaramaSyncProducer создает Sarama SyncProducer с дефолтной конфигурацией
func (b *ProducerBuilder) createSaramaSyncProducer() (sarama.SyncProducer, error) {
	brokers := strings.Split(b.kafkaConfig.Brokers(), ",")

	config := b.createDefaultSaramaConfig()

	saramaSyncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama sync producer: %w", err)
	}

	return saramaSyncProducer, nil
}

// createDefaultSaramaConfig создает стандартную конфигурацию Sarama для продюсера
func (b *ProducerBuilder) createDefaultSaramaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Retry.Max = 3
	return config
}
