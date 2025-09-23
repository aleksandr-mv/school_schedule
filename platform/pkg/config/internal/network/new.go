package network

import (
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/config/contracts"
)

type module struct {
	cors contracts.CORSConfig
}

func New() (contracts.NetworkConfig, error) {
	cors, err := NewCORSConfig()
	if err != nil {
		return nil, err
	}

	return &module{cors: cors}, nil
}

func (m *module) CORS() contracts.CORSConfig {
	return m.cors
}
