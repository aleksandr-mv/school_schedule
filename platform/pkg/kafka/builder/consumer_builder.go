package builder

import (
	"fmt"
	"strings"

	"github.com/IBM/sarama"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/kafka"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/kafka/consumer"
	middleware "github.com/aleksandr-mv/school_schedule/platform/pkg/kafka/consumer/middleware"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
)

// ConsumerBuilder упрощает создание Kafka консюмеров
type ConsumerBuilder struct {
	kafkaConfig contracts.KafkaConfig
	logger      consumer.Logger
}

// NewConsumerBuilder создает новый билдер для консюмеров
func NewConsumerBuilder(kafkaConfig contracts.KafkaConfig) *ConsumerBuilder {
	return &ConsumerBuilder{
		kafkaConfig: kafkaConfig,
		logger:      logger.Logger(),
	}
}

// WithLogger устанавливает кастомный логгер
func (b *ConsumerBuilder) WithLogger(logger consumer.Logger) *ConsumerBuilder {
	b.logger = logger
	return b
}

// BuildConsumer создает готового к использованию консюмера
func (b *ConsumerBuilder) BuildConsumer(consumerName string) (kafka.Consumer, error) {
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

	consumerGroup, err := b.createConsumerGroup(consumerConfig.GroupID())
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group for '%s': %w", consumerName, err)
	}

	return consumer.NewConsumer(
		consumerGroup,
		[]string{consumerConfig.Topic()},
		b.logger,
		middleware.Logging(b.logger),
	), nil
}

// createConsumerGroup создает Sarama ConsumerGroup с дефолтной конфигурацией
func (b *ConsumerBuilder) createConsumerGroup(groupID string) (sarama.ConsumerGroup, error) {
	brokers := strings.Split(b.kafkaConfig.Brokers(), ",")

	config := b.createDefaultSaramaConfig()

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group '%s': %w", groupID, err)
	}

	return consumerGroup, nil
}

// createDefaultSaramaConfig создает стандартную конфигурацию Sarama
func (b *ConsumerBuilder) createDefaultSaramaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
