package contracts

// KafkaConfig описывает конфигурацию для работы с Apache Kafka.
// Включает настройки брокеров, топиков и consumer groups.
// Конфигурация является опциональной - если брокеры не настроены,
// методы возвращают пустые строки и IsEnabled() возвращает false.
type KafkaConfig interface {
	// IsEnabled возвращает true, если Kafka настроена (есть брокеры)
	IsEnabled() bool

	// Brokers возвращает список адресов Kafka брокеров через запятую
	// Пример: "localhost:9092,localhost:9093"
	Brokers() string

	Consumers() ConsumersConfig
	Producers() ProducersConfig
}

type ConsumersConfig interface {
	GetConsumer(name string) (KafkaConsumerConfig, bool)
	GetAllConsumers() map[string]KafkaConsumerConfig
	GetEnabledConsumers() map[string]KafkaConsumerConfig
}

type ProducersConfig interface {
	GetProducer(name string) (KafkaProducerConfig, bool)
	GetAllProducers() map[string]KafkaProducerConfig
	GetEnabledProducers() map[string]KafkaProducerConfig
}

// KafkaConsumerConfig конфигурация отдельного Kafka консюмера
// Содержит ТОЛЬКО данные конфигурации, без привязки к конкретной реализации
type KafkaConsumerConfig interface {
	// Name возвращает имя консюмера
	Name() string

	// Topic возвращает топик для консюмера
	Topic() string

	// GroupID возвращает consumer group ID
	GroupID() string

	// IsEnabled возвращает true, если консюмер включен
	IsEnabled() bool
}

// KafkaProducerConfig конфигурация отдельного Kafka продюсера
// Содержит ТОЛЬКО данные конфигурации, без привязки к конкретной реализации
type KafkaProducerConfig interface {
	// Name возвращает имя продюсера
	Name() string

	// Topic возвращает топик для продюсера
	Topic() string

	// IsEnabled возвращает true, если продюсер включен
	IsEnabled() bool
}
