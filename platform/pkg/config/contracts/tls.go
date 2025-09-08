package contracts

// TLSConfig предоставляет интерфейс для настройки TLS/SSL соединений.
//
// Этот интерфейс используется как для баз данных (PostgreSQL, MongoDB),
// так и для транспортного слоя (HTTP/gRPC серверы и клиенты).
//
// Методы:
//   - BuildQueryParams(): возвращает параметры для строки подключения БД
//   - IsEnabled(): проверяет, включен ли TLS
//   - GetCertFile(): путь к сертификату (для серверов)
//   - GetKeyFile(): путь к приватному ключу (для серверов)
//   - GetCAFile(): путь к CA сертификату (для клиентов)
type TLSConfig interface {
	// BuildQueryParams возвращает параметры для строки подключения БД
	// Например: ["sslmode=require", "sslcert=/path/to/cert.pem"]
	BuildQueryParams() []string

	// IsEnabled проверяет, включен ли TLS
	IsEnabled() bool

	// GetCertFile возвращает путь к сертификату (для серверов)
	GetCertFile() string

	// GetKeyFile возвращает путь к приватному ключу (для серверов)
	GetKeyFile() string

	// GetCAFile возвращает путь к CA сертификату (для клиентов)
	GetCAFile() string
}
