package kafka

import (
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/contracts"
)

// rawConsumerConfig представляет конфигурацию одного консюмера
type rawConsumerConfig struct {
	Topic   string `mapstructure:"topic"    yaml:"topic"`
	GroupID string `mapstructure:"group_id" yaml:"group_id"`
	Enabled bool   `mapstructure:"enabled"  yaml:"enabled"`
}

// ConsumerConfig реализует KafkaConsumerConfig
type ConsumerConfig struct {
	name string
	raw  rawConsumerConfig
}

// Компиляционная проверка
var _ contracts.KafkaConsumerConfig = (*ConsumerConfig)(nil)

// Методы интерфейса KafkaConsumerConfig
func (c *ConsumerConfig) Name() string    { return c.name }
func (c *ConsumerConfig) Topic() string   { return c.raw.Topic }
func (c *ConsumerConfig) GroupID() string { return c.raw.GroupID }
func (c *ConsumerConfig) IsEnabled() bool { return c.raw.Enabled }

// ConsumersConfig реализует ConsumersConfig интерфейс
type ConsumersConfig struct {
	consumers map[string]*ConsumerConfig
}

// Компиляционная проверка
var _ contracts.ConsumersConfig = (*ConsumersConfig)(nil)

func (cc *ConsumersConfig) GetConsumer(name string) (contracts.KafkaConsumerConfig, bool) {
	consumer, exists := cc.consumers[name]
	return consumer, exists
}

func (cc *ConsumersConfig) GetAllConsumers() map[string]contracts.KafkaConsumerConfig {
	result := make(map[string]contracts.KafkaConsumerConfig, len(cc.consumers))
	for name, consumer := range cc.consumers {
		result[name] = consumer
	}
	return result
}

func (cc *ConsumersConfig) GetEnabledConsumers() map[string]contracts.KafkaConsumerConfig {
	result := make(map[string]contracts.KafkaConsumerConfig)
	for name, consumer := range cc.consumers {
		if consumer.IsEnabled() {
			result[name] = consumer
		}
	}
	return result
}
