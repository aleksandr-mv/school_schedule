package postgres

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/spf13/viper"

	contracts "github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

var validSSLModes = []string{"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}

type rawTLSConfig struct {
	Enabled        bool   `mapstructure:"enabled"         yaml:"enabled"         env:"DB_TLS_ENABLED"`
	SSLMode        string `mapstructure:"ssl_mode"        yaml:"ssl_mode"        env:"DB_SSL_MODE"`
	SSLCert        string `mapstructure:"ssl_cert"        yaml:"ssl_cert"        env:"DB_SSL_CERT"`
	SSLKey         string `mapstructure:"ssl_key"         yaml:"ssl_key"         env:"DB_SSL_KEY"`
	SSLRootCert    string `mapstructure:"ssl_root_cert"   yaml:"ssl_root_cert"   env:"DB_SSL_ROOT_CERT"`
	SSLPassword    string `mapstructure:"ssl_password"    yaml:"ssl_password"    env:"DB_SSL_PASSWORD"`
	SSLCRL         string `mapstructure:"ssl_crl"         yaml:"ssl_crl"         env:"DB_SSL_CRL"`
	SSLCompression bool   `mapstructure:"ssl_compression" yaml:"ssl_compression" env:"DB_SSL_COMPRESSION"`
}

type tlsConfig struct {
	Raw rawTLSConfig `yaml:"tls"`
}

var _ contracts.TLSConfig = (*tlsConfig)(nil)

func defaultTLSConfig() *rawTLSConfig {
	return &rawTLSConfig{
		Enabled:        false,
		SSLMode:        "disable",
		SSLCert:        "",
		SSLKey:         "",
		SSLRootCert:    "",
		SSLPassword:    "",
		SSLCRL:         "",
		SSLCompression: false,
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

// NewTLSConfig читает конфиг TLS из поддерева Viper (ожидается database.postgres.tls)
// Если sub == nil, используется fallback на переменные окружения через DefaultTLSConfig().
func NewTLSConfig(sub *viper.Viper) (*tlsConfig, error) {
	if sub == nil {
		return DefaultTLSConfig()
	}

	tls := sub.Sub("tls")
	if tls == nil {
		return DefaultTLSConfig()
	}

	raw := defaultTLSConfig()
	if err := tls.Unmarshal(raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TLS config: %w", err)
	}

	return &tlsConfig{Raw: *raw}, nil
}

func (t *tlsConfig) BuildQueryParams() []string {
	if !t.Raw.Enabled {
		return []string{"sslmode=disable"}
	}

	params := []string{fmt.Sprintf("sslmode=%s", t.Raw.SSLMode)}
	if t.Raw.SSLCert != "" {
		params = append(params, "sslcert="+t.Raw.SSLCert)
	}
	if t.Raw.SSLKey != "" {
		params = append(params, "sslkey="+t.Raw.SSLKey)
	}
	if t.Raw.SSLRootCert != "" {
		params = append(params, "sslrootcert="+t.Raw.SSLRootCert)
	}
	if t.Raw.SSLPassword != "" {
		params = append(params, "sslpassword="+t.Raw.SSLPassword)
	}
	if t.Raw.SSLCRL != "" {
		params = append(params, "sslcrl="+t.Raw.SSLCRL)
	}
	params = append(params, fmt.Sprintf("sslcompression=%t", t.Raw.SSLCompression))

	return params
}

func (t *tlsConfig) IsEnabled() bool {
	return t.Raw.Enabled
}

func (t *tlsConfig) GetCertFile() string {
	return t.Raw.SSLCert
}

func (t *tlsConfig) GetKeyFile() string {
	return t.Raw.SSLKey
}

func (t *tlsConfig) GetCAFile() string {
	return t.Raw.SSLRootCert
}
