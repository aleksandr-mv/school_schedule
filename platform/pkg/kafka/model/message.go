package model

import "time"

// Message — универсальная обёртка над сообщением Kafka.
type Message struct {
	Headers        map[string][]byte
	Timestamp      time.Time
	BlockTimestamp time.Time

	Key       []byte
	Value     []byte
	Topic     string
	Partition int32
	Offset    int64
}

// FailedMessage представляет сообщение, которое не удалось доставить
type FailedMessage struct {
	Topic     string    `json:"topic"`
	Key       []byte    `json:"key"`
	Value     []byte    `json:"value"`
	Error     string    `json:"error"`
	Attempt   int       `json:"attempt"`
	Timestamp time.Time `json:"timestamp"`
}
