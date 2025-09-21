package kafka

import (
	"context"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/kafka/model"
)

// MessageHandler — обработчик сообщений.
type MessageHandler func(ctx context.Context, msg model.Message) error

type Consumer interface {
	Consume(ctx context.Context, handler MessageHandler) error
	Close() error
}

type Producer interface {
	Send(ctx context.Context, key, value []byte) error
	Close() error
}
