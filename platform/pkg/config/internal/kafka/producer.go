package kafka

import (
	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// rawProducerConfig представляет конфигурацию одного продюсера
type rawProducerConfig struct {
	Topic   string `mapstructure:"topic"   yaml:"topic"`
	Enabled bool   `mapstructure:"enabled" yaml:"enabled"`
}

// kafkaProducerConfig реализует KafkaProducerConfig
type kafkaProducerConfig struct {
	name    string
	topic   string
	enabled bool
}

// Методы интерфейса KafkaProducerConfig
func (p *kafkaProducerConfig) Name() string    { return p.name }
func (p *kafkaProducerConfig) Topic() string   { return p.topic }
func (p *kafkaProducerConfig) IsEnabled() bool { return p.enabled }

// producersConfig реализует ProducersConfig
type producersConfig struct {
	kafkaConfig *kafkaConfig
}

// Методы для работы с продюсерами в kafkaConfig
func (k *kafkaConfig) Producers() contracts.ProducersConfig {
	return &producersConfig{kafkaConfig: k}
}

func (p *producersConfig) GetProducer(name string) (contracts.KafkaProducerConfig, bool) {
	raw, exists := p.kafkaConfig.Raw.Producers[name]
	if !exists {
		return nil, false
	}

	producer := &kafkaProducerConfig{
		name:    name,
		topic:   raw.Topic,
		enabled: raw.Enabled,
	}

	return producer, true
}

func (p *producersConfig) GetAllProducers() map[string]contracts.KafkaProducerConfig {
	result := make(map[string]contracts.KafkaProducerConfig)
	for name := range p.kafkaConfig.Raw.Producers {
		if producer, exists := p.GetProducer(name); exists {
			result[name] = producer
		}
	}

	return result
}

func (p *producersConfig) GetEnabledProducers() map[string]contracts.KafkaProducerConfig {
	result := make(map[string]contracts.KafkaProducerConfig)
	for name, producer := range p.GetAllProducers() {
		if producer.IsEnabled() {
			result[name] = producer
		}
	}

	return result
}
