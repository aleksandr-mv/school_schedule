package redis

import (
	"strings"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

// rawCluster для загрузки данных из YAML/ENV
type rawCluster struct {
	Nodes            []string `mapstructure:"nodes"             yaml:"nodes"             env:"REDIS_CLUSTER_NODES" envSeparator:","`
	Password         string   `mapstructure:"password"          yaml:"password"          env:"REDIS_CLUSTER_PASSWORD"`
	MaxRedirects     int      `mapstructure:"max_redirects"     yaml:"max_redirects"     env:"REDIS_CLUSTER_MAX_REDIRECTS"`
	ReadOnlyCommands bool     `mapstructure:"readonly_commands" yaml:"readonly_commands" env:"REDIS_CLUSTER_READONLY_COMMANDS"`
	RouteByLatency   bool     `mapstructure:"route_by_latency"  yaml:"route_by_latency"  env:"REDIS_CLUSTER_ROUTE_BY_LATENCY"`
	RouteRandomly    bool     `mapstructure:"route_randomly"    yaml:"route_randomly"    env:"REDIS_CLUSTER_ROUTE_RANDOMLY"`
}

// Cluster публичная структура для использования
type Cluster struct {
	raw rawCluster
}

// defaultCluster возвращает rawCluster с дефолтными значениями
func defaultCluster() rawCluster {
	return rawCluster{
		Nodes:            []string{}, // Пустой список - Redis отключен по умолчанию
		Password:         "",
		MaxRedirects:     3,
		ReadOnlyCommands: true,
		RouteByLatency:   false,
		RouteRandomly:    false,
	}
}

// Реализуем contracts.RedisClusterConfig
var _ contracts.RedisClusterConfig = (*Cluster)(nil)

// Cluster методы
func (c *Cluster) IsEnabled() bool        { return len(c.raw.Nodes) > 0 }
func (c *Cluster) Nodes() []string        { return c.raw.Nodes }
func (c *Cluster) Password() string       { return c.raw.Password }
func (c *Cluster) MaxRedirects() int      { return c.raw.MaxRedirects }
func (c *Cluster) ReadOnlyCommands() bool { return c.raw.ReadOnlyCommands }
func (c *Cluster) RouteByLatency() bool   { return c.raw.RouteByLatency }
func (c *Cluster) RouteRandomly() bool    { return c.raw.RouteRandomly }

// NodesAddresses возвращает узлы кластера через запятую
func (c *Cluster) NodesAddresses() string {
	return strings.Join(c.raw.Nodes, ",")
}
