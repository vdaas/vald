package exporter

import "context"

type Exporter interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
