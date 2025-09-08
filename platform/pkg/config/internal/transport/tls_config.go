package transport

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

type rawTLSConfig struct {
	Enabled  bool   `mapstructure:"enabled" yaml:"enabled" env:"TLS_ENABLED"`
	CertFile string `mapstructure:"cert_file" yaml:"cert_file" env:"TLS_CERT_FILE"`
	KeyFile  string `mapstructure:"key_file" yaml:"key_file" env:"TLS_KEY_FILE"`
	CAFile   string `mapstructure:"ca_file" yaml:"ca_file" env:"TLS_CA_FILE"`
}

type tlsConfig struct {
	Raw rawTLSConfig `yaml:"tls"`
}

var _ contracts.TLSConfig = (*tlsConfig)(nil)

// defaultTLSConfig возвращает конфигурацию TLS с дефолтными значениями
func defaultTLSConfig() *rawTLSConfig {
	return &rawTLSConfig{
		Enabled:  false,
		CertFile: "",
		KeyFile:  "",
		CAFile:   "",
	}
}

// DefaultTLSConfig читает конфиг TLS из ENV
func DefaultTLSConfig() (*tlsConfig, error) {
	raw := defaultTLSConfig()

	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse TLS config: %w", err)
	}

	return &tlsConfig{Raw: *raw}, nil
}

// NewTLSConfig читает конфиг TLS из поддерева Viper (ожидается transport.tls)
// Если sub == nil, используется fallback на переменные окружения через DefaultTLSConfig().
func NewTLSConfig() (*tlsConfig, error) {
	return DefaultTLSConfig()
}

func (t *tlsConfig) BuildQueryParams() []string {
	return nil
}

func (t *tlsConfig) IsEnabled() bool {
	return t.Raw.Enabled
}

func (t *tlsConfig) GetCertFile() string {
	return t.Raw.CertFile
}

func (t *tlsConfig) GetKeyFile() string {
	return t.Raw.KeyFile
}

func (t *tlsConfig) GetCAFile() string {
	return t.Raw.CAFile
}

func (t *tlsConfig) String() string {
	return fmt.Sprintf(
		"TLSConfig{Enabled:%t, CertFile:%s, KeyFile:%s, CAFile:%s}",
		t.Raw.Enabled, t.Raw.CertFile, t.Raw.KeyFile, t.Raw.CAFile,
	)
}
