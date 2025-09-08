package transport

import "github.com/aleksandr-mv/school_schedule/platform/pkg/config/contracts"

type module struct {
	http       contracts.HTTPServer
	grpc       contracts.GRPCServer
	grpcLimits contracts.GRPCClientLimits
	tls        contracts.TLSConfig
}

func New() (contracts.TransportModule, error) {
	http, err := NewHTTPServerConfig()
	if err != nil {
		return nil, err
	}

	grpc, err := NewGRPCServerConfig()
	if err != nil {
		return nil, err
	}

	grpcLimits, err := NewGRPCClientLimitsConfig()
	if err != nil {
		return nil, err
	}

	tls, err := NewTLSConfig()
	if err != nil {
		return nil, err
	}

	return &module{
		http:       http,
		grpc:       grpc,
		grpcLimits: grpcLimits,
		tls:        tls,
	}, nil
}

func (m *module) HTTP() contracts.HTTPServer                   { return m.http }
func (m *module) GRPC() contracts.GRPCServer                   { return m.grpc }
func (m *module) GRPCClientLimits() contracts.GRPCClientLimits { return m.grpcLimits }
func (m *module) TLS() contracts.TLSConfig                     { return m.tls }
