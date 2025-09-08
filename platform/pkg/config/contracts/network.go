package contracts

type NetworkConfig interface {
	CORS() CORSConfig
}

// CORSConfig публичный интерфейс конфигурации CORS.
type CORSConfig interface {
	AllowedOrigins() []string
	AllowedMethods() []string
	AllowedHeaders() []string
	ExposedHeaders() []string
	AllowCredentials() bool
	MaxAge() int
}
