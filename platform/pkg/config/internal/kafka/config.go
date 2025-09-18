package kafka

import (
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// Компиляционная проверка
var _ contracts.KafkaConfig = (*Config)(nil)

// rawConfig для загрузки данных из YAML/ENV
type rawConfig struct {
	Brokers   string                       `mapstructure:"brokers"   yaml:"brokers"   env:"KAFKA_BROKERS"`
	Consumers map[string]rawConsumerConfig `mapstructure:"consumers" yaml:"consumers"`
	Producers map[string]rawProducerConfig `mapstructure:"producers" yaml:"producers"`
}

// Config публичная структура Kafka конфигурации
type Config struct {
	raw       rawConfig
	consumers *ConsumersConfig
	producers *ProducersConfig
}

// defaultConfig возвращает rawConfig с дефолтными значениями
func defaultConfig() rawConfig {
	return rawConfig{
		Brokers:   "",
		Consumers: make(map[string]rawConsumerConfig),
		Producers: make(map[string]rawProducerConfig),
	}
}

// Методы интерфейса KafkaConfig
func (c *Config) IsEnabled() bool { return c.raw.Brokers != "" }
func (c *Config) Brokers() string { return c.raw.Brokers }

func (c *Config) Consumers() contracts.ConsumersConfig {
	if c.consumers == nil {
		c.consumers = &ConsumersConfig{consumers: make(map[string]*ConsumerConfig)}
		for name, raw := range c.raw.Consumers {
			c.consumers.consumers[name] = &ConsumerConfig{name: name, raw: raw}
		}
	}
	return c.consumers
}

func (c *Config) Producers() contracts.ProducersConfig {
	if c.producers == nil {
		c.producers = &ProducersConfig{producers: make(map[string]*ProducerConfig)}
		for name, raw := range c.raw.Producers {
			c.producers.producers[name] = &ProducerConfig{name: name, raw: raw}
		}
	}
	return c.producers
}
