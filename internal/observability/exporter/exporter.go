package exporter

import "context"

type Exporter interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

var (
	DefaultMillisecondsHistogramDistribution = []float64{
		0.01,
		0.05,
		0.1,
		0.3,
		0.6,
		0.8,
		1,
		2,
		3,
		4,
		5,
		6,
		8,
		10,
		13,
		16,
		20,
		25,
		30,
		40,
		50,
		65,
		80,
		100,
		130,
		160,
		200,
		250,
		300,
		400,
		500,
		650,
		800,
		1000,
		2000,
		5000,
		10000,
		20000,
		50000,
		100000,
	}
)
