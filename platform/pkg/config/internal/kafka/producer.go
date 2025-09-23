package kafka

import (
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/contracts"
)

// rawProducerConfig представляет конфигурацию одного продюсера
type rawProducerConfig struct {
	Topic   string `mapstructure:"topic"   yaml:"topic"`
	Enabled bool   `mapstructure:"enabled" yaml:"enabled"`
}

// ProducerConfig реализует KafkaProducerConfig
type ProducerConfig struct {
	name string
	raw  rawProducerConfig
}

// Компиляционная проверка
var _ contracts.KafkaProducerConfig = (*ProducerConfig)(nil)

// Методы интерфейса KafkaProducerConfig
func (p *ProducerConfig) Name() string    { return p.name }
func (p *ProducerConfig) Topic() string   { return p.raw.Topic }
func (p *ProducerConfig) IsEnabled() bool { return p.raw.Enabled }

// ProducersConfig реализует ProducersConfig интерфейс
type ProducersConfig struct {
	producers map[string]*ProducerConfig
}

// Компиляционная проверка
var _ contracts.ProducersConfig = (*ProducersConfig)(nil)

func (pc *ProducersConfig) GetProducer(name string) (contracts.KafkaProducerConfig, bool) {
	producer, exists := pc.producers[name]
	return producer, exists
}

func (pc *ProducersConfig) GetAllProducers() map[string]contracts.KafkaProducerConfig {
	result := make(map[string]contracts.KafkaProducerConfig, len(pc.producers))
	for name, producer := range pc.producers {
		result[name] = producer
	}
	return result
}

func (pc *ProducersConfig) GetEnabledProducers() map[string]contracts.KafkaProducerConfig {
	result := make(map[string]contracts.KafkaProducerConfig)
	for name, producer := range pc.producers {
		if producer.IsEnabled() {
			result[name] = producer
		}
	}
	return result
}
