package vector

import (
	"math"
	"math/rand"

	"github.com/vdaas/vald/internal/errors"
)

type (
	Distribution               int
	Float32VectorGeneratorFunc func(int, int) [][]float32
	Uint8VectorGeneratorFunc   func(int, int) [][]uint8
)

const (
	Gaussian Distribution = iota
	Uniform
)

// ErrUnknownDistritbution represents an error which the distribution is unknown.
var ErrUnknownDistribution = errors.New("Unknown distribution generator type")

// Float32VectorGenerator returns float32 vector generator function which has selected distribution
func Float32VectorGenerator(d Distribution) (Float32VectorGeneratorFunc, error) {
	switch d {
	case Gaussian:
		return GaussianDistributedFloat32VectorGenerator, nil
	case Uniform:
		return UniformDistributedFloat32VectorGenerator, nil
	default:
		return nil, ErrUnknownDistribution
	}
}

// Uint8VectorGenerator returns uint8 vector generator function which has selected distribution
func Uint8VectorGenerator(d Distribution) (Uint8VectorGeneratorFunc, error) {
	switch d {
	case Gaussian:
		return GaussianDistributedUint8VectorGenerator, nil
	case Uniform:
		return UniformDistributedUint8VectorGenerator, nil
	default:
		return nil, ErrUnknownDistribution
	}
}

// float32VectorGenerator return n float32 vectors with dim dimension
func float32VectorGenerator(n, dim int, gen func() float32) (ret [][]float32) {
	ret = make([][]float32, 0, n)

	for i := 0; i < n; i++ {
		v := make([]float32, dim)
		for j := 0; j < dim; j++ {
			v[j] = gen()
		}
		ret = append(ret, v)
	}
	return
}

// UniformDistributedFloat32VectorGenerator returns n float32 vectors with dim dimension and their values under Uniform distribution
func UniformDistributedFloat32VectorGenerator(n, dim int) [][]float32 {
	return float32VectorGenerator(n, dim, rand.Float32)
}

// GaussianDistributedFloat32VectorGenerator returns n float32 vectors with dim dimension and their values under Gaussian distribution
func GaussianDistributedFloat32VectorGenerator(n, dim int) [][]float32 {
	return float32VectorGenerator(n, dim, func() float32 {
		return float32(rand.NormFloat64())
	})
}

// uint8VectorGenerator return n uint8 vectors with dim dimension
func uint8VectorGenerator(n, dim int, gen func() uint8) (ret [][]uint8) {
	ret = make([][]uint8, 0, n)

	for i := 0; i < n; i++ {
		v := make([]uint8, dim)
		for j := 0; j < dim; j++ {
			v[j] = gen()
		}
		ret = append(ret, v)
	}
	return
}

// UniformDistributedUint8VectorGenerator returns n uint8 vectors with dim dimension and their values under Uniform distribution
func UniformDistributedUint8VectorGenerator(n, dim int) [][]uint8 {
	return uint8VectorGenerator(n, dim, func() uint8 {
		return uint8(rand.Intn(int(math.MaxUint8) + 1))
	})
}

// GaussianDistributedUint8VectorGenerator returns n uint8 vectors with dim dimension and their values under Gaussian distribution
func GaussianDistributedUint8VectorGenerator(n, dim int) [][]uint8 {
	// mean:128, sigma:128/3, all of 99.7% are in [0, 255]
	const (
		mean  float64 = 128
		sigma float64 = 128 / 3
	)
	return uint8VectorGenerator(n, dim, func() uint8 {
		val := rand.NormFloat64()*sigma + mean
		if val < 0 {
			return 0
		} else if val > math.MaxUint8 {
			return math.MaxUint8
		} else {
			return uint8(val)
		}
	})
}
