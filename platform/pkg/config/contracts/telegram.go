package contracts

// TelegramConfig описывает конфигурацию для работы с Telegram Bot API.
// Включает настройки токена, webhook URL и дополнительные параметры.
// Конфигурация является опциональной - если токен не настроен,
// методы возвращают пустые строки и IsEnabled() возвращает false.
type TelegramConfig interface {
	// IsEnabled возвращает true, если Telegram бот настроен (есть токен)
	IsEnabled() bool

	// Token возвращает токен бота от @BotFather
	// Пример: "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz"
	Token() string

	// WebhookURL возвращает URL для webhook (опционально)
	// Пример: "https://your-domain.com/webhook/telegram"
	WebhookURL() string

	// DebugMode возвращает true, если включен режим отладки
	DebugMode() bool

	// Timeout возвращает таймаут для HTTP запросов к Telegram API
	Timeout() int // в секундах

	// MaxRetries возвращает максимальное количество попыток при ошибках
	MaxRetries() int
}
