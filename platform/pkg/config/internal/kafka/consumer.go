package kafka

import (
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// rawConsumerConfig представляет конфигурацию одного консюмера
type rawConsumerConfig struct {
	Topic   string `mapstructure:"topic"    yaml:"topic"`
	GroupID string `mapstructure:"group_id" yaml:"group_id"`
	Enabled bool   `mapstructure:"enabled"  yaml:"enabled"`
}

// kafkaConsumerConfig реализует KafkaConsumerConfig
type kafkaConsumerConfig struct {
	name    string
	topic   string
	groupID string
	enabled bool
}

// Методы интерфейса KafkaConsumerConfig
func (c *kafkaConsumerConfig) Name() string    { return c.name }
func (c *kafkaConsumerConfig) Topic() string   { return c.topic }
func (c *kafkaConsumerConfig) GroupID() string { return c.groupID }
func (c *kafkaConsumerConfig) IsEnabled() bool { return c.enabled }

// consumersConfig реализует ConsumersConfig
type consumersConfig struct {
	kafkaConfig *kafkaConfig
}

// Методы для работы с консюмерами в kafkaConfig
func (k *kafkaConfig) Consumers() contracts.ConsumersConfig {
	return &consumersConfig{kafkaConfig: k}
}

func (c *consumersConfig) GetConsumer(name string) (contracts.KafkaConsumerConfig, bool) {
	raw, exists := c.kafkaConfig.Raw.Consumers[name]
	if !exists {
		return nil, false
	}

	consumer := &kafkaConsumerConfig{
		name:    name,
		topic:   raw.Topic,
		groupID: raw.GroupID,
		enabled: raw.Enabled,
	}

	return consumer, true
}

func (c *consumersConfig) GetAllConsumers() map[string]contracts.KafkaConsumerConfig {
	result := make(map[string]contracts.KafkaConsumerConfig)
	for name := range c.kafkaConfig.Raw.Consumers {
		if consumer, exists := c.GetConsumer(name); exists {
			result[name] = consumer
		}
	}

	return result
}

func (c *consumersConfig) GetEnabledConsumers() map[string]contracts.KafkaConsumerConfig {
	result := make(map[string]contracts.KafkaConsumerConfig)
	for name, consumer := range c.GetAllConsumers() {
		if consumer.IsEnabled() {
			result[name] = consumer
		}
	}

	return result
}
