package consumer

import (
	"fmt"
	"strings"

	"github.com/IBM/sarama"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/kafka"
)

// Builder создает Kafka consumer'ы из конфигурации
type Builder struct {
	kafkaConfig contracts.KafkaConfig
	logger      kafka.Logger
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

// BuildConsumer создает consumer
func (b *Builder) BuildConsumer(consumerName string) (*consumer, error) {
	if !b.kafkaConfig.IsEnabled() {
		return nil, fmt.Errorf("kafka is not enabled")
	}

	consumerConfig, exists := b.kafkaConfig.Consumers().GetConsumer(consumerName)
	if !exists {
		return nil, fmt.Errorf("consumer '%s' not found", consumerName)
	}

	if !consumerConfig.IsEnabled() {
		return nil, fmt.Errorf("consumer '%s' is disabled", consumerName)
	}

	if consumerConfig.Topic() == "" {
		return nil, fmt.Errorf("consumer '%s' topic cannot be empty", consumerName)
	}

	brokers := strings.Split(b.kafkaConfig.Brokers(), ",")
	config := b.createSaramaConfig()

	consumerGroup, err := sarama.NewConsumerGroup(brokers, consumerConfig.GroupID(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama consumer group: %w", err)
	}

	return NewConsumer(consumerGroup, []string{consumerConfig.Topic()}, b.logger), nil
}

// createSaramaConfig создает базовую конфигурацию Sarama
func (b *Builder) createSaramaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	return config
}
