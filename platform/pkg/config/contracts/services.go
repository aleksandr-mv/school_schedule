package contracts

import "time"

// ServicesConfig описывает интерфейс набора служб.
type ServicesConfig interface {
	Get(name string) (ServiceConfig, bool)
	All() map[string]ServiceConfig
}

// ServiceConfig описывает интерфейс одной службы.
type ServiceConfig interface {
	Host() string
	Port() int
	Timeout() time.Duration
	Address() string
}
