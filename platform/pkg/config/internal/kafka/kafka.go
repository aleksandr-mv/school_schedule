package kafka

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/helpers"
)

// rawKafkaConfig соответствует полям секции kafka в YAML и env-переменным.
type rawKafkaConfig struct {
	Brokers   string                       `mapstructure:"brokers"   yaml:"brokers"   env:"KAFKA_BROKERS"`
	Consumers map[string]rawConsumerConfig `mapstructure:"consumers" yaml:"consumers"`
	Producers map[string]rawProducerConfig `mapstructure:"producers" yaml:"producers"`
}

// kafkaConfig хранит данные секции kafka и реализует KafkaConfig.
type kafkaConfig struct {
	Raw rawKafkaConfig `yaml:"kafka"`
}

// defaultKafkaConfig возвращает конфигурацию Kafka с дефолтными значениями
func defaultKafkaConfig() *rawKafkaConfig {
	return &rawKafkaConfig{
		Brokers:   "",
		Consumers: make(map[string]rawConsumerConfig),
		Producers: make(map[string]rawProducerConfig),
	}
}

// DefaultKafkaConfig читает конфигурацию Kafka из ENV.
func DefaultKafkaConfig() (*kafkaConfig, error) {
	raw := defaultKafkaConfig()
	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse kafka env: %w", err)
	}
	return &kafkaConfig{Raw: *raw}, nil
}

// NewKafkaConfig создает конфигурацию Kafka, пытаясь сначала загрузить из YAML, затем из ENV.
func NewKafkaConfig() (*kafkaConfig, error) {
	if section := helpers.GetSection("kafka"); section != nil {
		raw := defaultKafkaConfig()
		if err := section.Unmarshal(raw); err == nil {
			return &kafkaConfig{Raw: *raw}, nil
		}
	}
	return DefaultKafkaConfig()
}

// Методы интерфейса KafkaConfig
func (k *kafkaConfig) IsEnabled() bool { return k.Raw.Brokers != "" }
func (k *kafkaConfig) Brokers() string { return k.Raw.Brokers }
