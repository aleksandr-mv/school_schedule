package mongo

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/spf13/viper"

	contracts "github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"
)

type rawTLSConfig struct {
	Enabled                  bool   `mapstructure:"enabled"                    yaml:"enabled"                       env:"MONGO_TLS_ENABLED"`
	CertificateKeyFile       string `mapstructure:"certificate_key_file"       yaml:"certificate_key_file"          env:"MONGO_TLS_CERTIFICATE_KEY_FILE"`
	CAFile                   string `mapstructure:"ca_file"                    yaml:"ca_file"                       env:"MONGO_TLS_CA_FILE"`
	AllowInvalidCertificates bool   `mapstructure:"allow_invalid_certificates" yaml:"allow_invalid_certificates"    env:"MONGO_TLS_ALLOW_INVALID_CERTIFICATES"`
	AllowInvalidHostnames    bool   `mapstructure:"allow_invalid_hostnames"    yaml:"allow_invalid_hostnames"       env:"MONGO_TLS_ALLOW_INVALID_HOSTNAMES"`
}

type tlsConfig struct {
	Raw rawTLSConfig `yaml:"tls"`
}

var _ contracts.TLSConfig = (*tlsConfig)(nil)

// defaultTLSConfig возвращает конфигурацию TLS MongoDB с дефолтными значениями
func defaultTLSConfig() *rawTLSConfig {
	return &rawTLSConfig{
		Enabled:                  false,
		CertificateKeyFile:       "",
		CAFile:                   "",
		AllowInvalidCertificates: false,
		AllowInvalidHostnames:    false,
	}
}

// DefaultTLSConfig парсит env и возвращает tlsConfig.
func DefaultTLSConfig() (*tlsConfig, error) {
	raw := defaultTLSConfig()
	if err := env.Parse(raw); err != nil {
		return nil, fmt.Errorf("failed to parse TLS config: %w", err)
	}
	return &tlsConfig{Raw: *raw}, nil
}

// NewTLSConfig читает TLS-конфиг Mongo из поддерева Viper (ожидается database.mongo.tls).
// При отсутствии секции — fallback на ENV.
func NewTLSConfig(sub *viper.Viper) (*tlsConfig, error) {
	if sub == nil {
		return DefaultTLSConfig()
	}

	tlsV := sub.Sub("tls")
	if tlsV == nil {
		return DefaultTLSConfig()
	}

	raw := defaultTLSConfig()
	if err := tlsV.Unmarshal(raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Mongo TLS config: %w", err)
	}
	return &tlsConfig{Raw: *raw}, nil
}

// BuildQueryParams возвращает параметры TLS для URI.
func (t *tlsConfig) BuildQueryParams() []string {
	if !t.Raw.Enabled {
		return nil
	}

	params := []string{"tls=true"}
	if t.Raw.CertificateKeyFile != "" {
		params = append(params, "tlsCertificateKeyFile="+t.Raw.CertificateKeyFile)
	}

	if t.Raw.CAFile != "" {
		params = append(params, "tlsCAFile="+t.Raw.CAFile)
	}

	if t.Raw.AllowInvalidCertificates {
		params = append(params, "tlsAllowInvalidCertificates=true")
	}

	if t.Raw.AllowInvalidHostnames {
		params = append(params, "tlsAllowInvalidHostnames=true")
	}

	return params
}

func (t *tlsConfig) IsEnabled() bool {
	return t.Raw.Enabled
}

func (t *tlsConfig) GetCertFile() string {
	return t.Raw.CertificateKeyFile
}

func (t *tlsConfig) GetKeyFile() string {
	// MongoDB использует единый файл для сертификата и ключа
	return t.Raw.CertificateKeyFile
}

func (t *tlsConfig) GetCAFile() string {
	return t.Raw.CAFile
}
