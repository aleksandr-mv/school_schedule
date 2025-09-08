package builder

import (
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/kafka"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/kafka/consumer"
)

// KafkaBuilder основной билдер для всех Kafka компонентов
type KafkaBuilder struct {
	consumerBuilder *ConsumerBuilder
	producerBuilder *ProducerBuilder
}

// NewKafkaBuilder создает главный билдер для Kafka
func NewKafkaBuilder(kafkaConfig contracts.KafkaConfig) *KafkaBuilder {
	return &KafkaBuilder{
		consumerBuilder: NewConsumerBuilder(kafkaConfig),
		producerBuilder: NewProducerBuilder(kafkaConfig),
	}
}

// WithLogger устанавливает кастомный логгер для консюмеров
func (b *KafkaBuilder) WithLogger(logger consumer.Logger) *KafkaBuilder {
	b.consumerBuilder.WithLogger(logger)
	return b
}

// BuildConsumer создает консюмера по имени
func (b *KafkaBuilder) BuildConsumer(consumerName string) (kafka.Consumer, error) {
	return b.consumerBuilder.BuildConsumer(consumerName)
}

// BuildProducer создает продюсера по имени
func (b *KafkaBuilder) BuildProducer(producerName string) (kafka.Producer, error) {
	return b.producerBuilder.BuildProducer(producerName)
}

// Consumers возвращает билдер консюмеров для дополнительной настройки
func (b *KafkaBuilder) Consumers() *ConsumerBuilder {
	return b.consumerBuilder
}

// Producers возвращает билдер продюсеров для дополнительной настройки
func (b *KafkaBuilder) Producers() *ProducerBuilder {
	return b.producerBuilder
}
