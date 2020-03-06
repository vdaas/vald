package agent

import (
	iconfig "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/pkg/agent/ngt/config"
)

type Option func(*server)

var (
	defaultOptions = []Option{
		WithConfig(&config.Data{
			GlobalConfig: config.GlobalConfig{
				Version: "v0.0.0",
			},
			Server: &iconfig.Servers{
				Servers: []*iconfig.Server{
					{
						Name:          "agent-grpc",
						Host:          "127.0.0.1",
						Port:          8082,
						Mode:          "GRPC",
						ProbeWaitTime: "0s",
						HTTP: &iconfig.HTTP{
							ShutdownDuration: "0s",
						},
					},
					{
						Name:          "agent-rest",
						Host:          "127.0.0.1",
						Port:          8081,
						Mode:          "REST",
						ProbeWaitTime: "0s",
						HTTP: &iconfig.HTTP{
							ShutdownDuration:  "0s",
							HandlerTimeout:    "60s",
							IdleTimeout:       "60s",
							ReadHeaderTimeout: "60s",
							ReadTimeout:       "60s",
							WriteTimeout:      "60s",
						},
					},
				},
				StartUpStrategy: []string{
					"agent-grpc",
					"agent-rest",
				},
				ShutdownStrategy: []string{
					"agent-grpc",
					"agent-rest",
				},
				FullShutdownDuration: "600s",
				TLS: &iconfig.TLS{
					Enabled: false,
				},
			},
			NGT: &iconfig.NGT{
				Dimension:          0,
				DistanceType:       "unknown",
				ObjectType:         "unknown",
				CreationEdgeSize:   20,
				SearchEdgeSize:     10,
				EnableInMemoryMode: true,
			},
		}),
	}
)

func WithConfig(cfg *config.Data) Option {
	return func(s *server) {
		if cfg != nil {
			s.cfg = cfg
		}
	}
}

func WithDimentaion(d int) Option {
	return func(s *server) {
		if s.cfg != nil && s.cfg.NGT != nil {
			if d > 0 {
				s.cfg.NGT.Dimension = d
			}
		}
	}
}

func WithDistanceType(dtype string) Option {
	return func(s *server) {
		if s.cfg != nil && s.cfg.NGT != nil {
			if len(dtype) != 0 {
				s.cfg.NGT.DistanceType = dtype
			}
		}
	}
}

func WithObjectType(otype string) Option {
	return func(s *server) {
		if s.cfg != nil && s.cfg.NGT != nil {
			if len(otype) != 0 {
				s.cfg.NGT.ObjectType = otype
			}
		}
	}
}
