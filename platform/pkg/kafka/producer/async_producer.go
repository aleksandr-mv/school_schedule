package producer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/kafka"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/kafka/model"
)

// asyncProducer асинхронная реализация с простой retry логикой
type asyncProducer struct {
	saramaProducer sarama.AsyncProducer
	topic          string
	logger         kafka.Logger
	metrics        kafka.Metrics

	// Retry конфигурация
	maxRetries int
	maxDelay   time.Duration

	// Для graceful shutdown
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Retry queue
	retryQueue chan model.FailedMessage
}

// NewAsyncProducer создает async producer с retry логикой
func NewAsyncProducer(
	saramaAsyncProducer sarama.AsyncProducer,
	topic string,
	logger kafka.Logger,
	metrics kafka.Metrics,
) *asyncProducer {
	ctx, cancel := context.WithCancel(context.Background())

	p := &asyncProducer{
		saramaProducer: saramaAsyncProducer,
		topic:          topic,
		logger:         logger,
		metrics:        metrics,
		maxRetries:     5,                // 5 попыток
		maxDelay:       16 * time.Second, // до 16 секунд
		ctx:            ctx,
		cancel:         cancel,
		retryQueue:     make(chan model.FailedMessage, 1000),
	}

	// Запускаем обработчики
	p.startErrorTracking()
	p.startRetryProcessor()

	return p
}

// Send отправляет сообщение асинхронно
func (p *asyncProducer) Send(ctx context.Context, key, value []byte) error {
	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
		Metadata: model.FailedMessage{
			Topic:     p.topic,
			Key:       key,
			Value:     value,
			Attempt:   0,
			Timestamp: time.Now(),
		},
	}

	select {
	case p.saramaProducer.Input() <- message:
		p.recordQueued()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-p.ctx.Done():
		return fmt.Errorf("producer is shutting down")
	}
}

// Close корректно закрывает producer
func (p *asyncProducer) Close() error {
	// Отменяем контекст для остановки горутин
	p.cancel()

	// Ждем завершения горутин
	p.wg.Wait()

	// Обрабатываем оставшиеся сообщения в retry очереди
	p.drainRetryQueue()

	// Закрываем async producer
	return p.saramaProducer.Close()
}

// startErrorTracking запускает отслеживание ошибок доставки
func (p *asyncProducer) startErrorTracking() {
	p.wg.Add(1)

	go func() {
		defer p.wg.Done()
		for {
			select {
			case err := <-p.saramaProducer.Errors():
				if err != nil {
					p.handleDeliveryError(err)
				}
			case <-p.ctx.Done():
				return
			}
		}
	}()
}

// startRetryProcessor запускает обработчик retry очереди
func (p *asyncProducer) startRetryProcessor() {
	p.wg.Add(1)

	go func() {
		defer p.wg.Done()
		for {
			select {
			case failedMsg := <-p.retryQueue:
				p.processRetry(failedMsg)
			case <-p.ctx.Done():
				return
			}
		}
	}()
}

// handleDeliveryError обрабатывает ошибку доставки сообщения
func (p *asyncProducer) handleDeliveryError(err *sarama.ProducerError) {
	// Извлекаем метаданные из сообщения
	var failedMsg model.FailedMessage
	if metadata, ok := err.Msg.Metadata.(model.FailedMessage); ok {
		failedMsg = metadata
	} else {
		// Fallback если метаданные не найдены
		failedMsg = model.FailedMessage{
			Topic:     err.Msg.Topic,
			Key:       err.Msg.Key.(sarama.ByteEncoder),
			Value:     err.Msg.Value.(sarama.ByteEncoder),
			Attempt:   0,
			Timestamp: time.Now(),
		}
	}

	failedMsg.Error = err.Err.Error()
	failedMsg.Attempt++

	p.logger.Error(p.ctx, "Failed to deliver message",
		zap.Error(err.Err),
		zap.String("topic", err.Msg.Topic),
		zap.Int("attempt", failedMsg.Attempt),
	)

	p.recordError()

	// Решаем, что делать с сообщением
	if failedMsg.Attempt <= p.maxRetries {
		// Добавляем в retry очередь
		select {
		case p.retryQueue <- failedMsg:
			p.logger.Info(p.ctx, "Message queued for retry",
				zap.String("topic", failedMsg.Topic),
				zap.Int("attempt", failedMsg.Attempt),
			)
		default:
			// Очередь переполнена, логируем как критическую ошибку
			p.logger.Error(p.ctx, "Retry queue full, message lost",
				zap.String("topic", failedMsg.Topic),
				zap.Int("attempt", failedMsg.Attempt),
			)
		}
	} else {
		// Превышено максимальное количество попыток
		p.logger.Error(p.ctx, "Max retries reached, message lost",
			zap.String("topic", failedMsg.Topic),
			zap.Int("attempt", failedMsg.Attempt),
			zap.String("error", failedMsg.Error),
		)
	}
}

// processRetry обрабатывает retry сообщения
func (p *asyncProducer) processRetry(failedMsg model.FailedMessage) {
	// Вычисляем задержку с экспоненциальным backoff
	delay := p.calculateRetryDelay(failedMsg.Attempt)
	time.Sleep(delay)

	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(failedMsg.Key),
		Value: sarama.ByteEncoder(failedMsg.Value),
		Metadata: model.FailedMessage{
			Topic:     failedMsg.Topic,
			Key:       failedMsg.Key,
			Value:     failedMsg.Value,
			Attempt:   failedMsg.Attempt,
			Timestamp: time.Now(),
		},
	}

	select {
	case p.saramaProducer.Input() <- message:
		p.logger.Info(p.ctx, "Message retry sent",
			zap.String("topic", failedMsg.Topic),
			zap.Int("attempt", failedMsg.Attempt),
		)
	case <-p.ctx.Done():
		// Producer закрывается, логируем как критическую ошибку
		p.logger.Error(p.ctx, "Producer closing, retry message lost",
			zap.String("topic", failedMsg.Topic),
			zap.Int("attempt", failedMsg.Attempt),
		)
	}
}

// calculateRetryDelay вычисляет задержку для retry с экспоненциальным backoff
func (p *asyncProducer) calculateRetryDelay(attempt int) time.Duration {
	// Базовая задержка 1 секунда
	baseDelay := time.Second

	// Экспоненциальный backoff: 1s, 2s, 4s, 8s, 16s
	delay := baseDelay * time.Duration(1<<(attempt-1))

	// Ограничиваем максимальной задержкой
	if delay > p.maxDelay {
		delay = p.maxDelay
	}

	return delay
}

// drainRetryQueue обрабатывает оставшиеся сообщения при shutdown
func (p *asyncProducer) drainRetryQueue() {
	for {
		select {
		case failedMsg := <-p.retryQueue:
			p.logger.Error(p.ctx, "Retry message lost during shutdown",
				zap.String("topic", failedMsg.Topic),
				zap.Int("attempt", failedMsg.Attempt),
			)
		default:
			return
		}
	}
}

// Метрики
func (p *asyncProducer) recordQueued() {
	if p.metrics != nil {
		p.metrics.IncrementCounter("kafka_producer_messages_queued_total", map[string]string{
			"topic": p.topic,
		})
	}
}

func (p *asyncProducer) recordError() {
	if p.metrics != nil {
		p.metrics.IncrementCounter("kafka_producer_messages_failed_total", map[string]string{
			"topic": p.topic,
		})
	}
}
