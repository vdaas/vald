package vector

import (
	"math"
	"math/rand"

	"github.com/vdaas/vald/internal/errors"
)

type Distribution int

const (
	Gaussian Distribution = iota
	Uniform
)

var ErrUnsupportedDistribution = errors.Errorf("Unsupported distribution generator type")

func Float32VectorGenerator(d Distribution) (func(int, int) [][]float32, error) {
	switch d {
	case Gaussian:
		return GaussianDistributedFloat32VectorGenerator, nil
	case Uniform:
		return UniformDistributedFloat32VectorGenerator, nil
	default:
		return nil, ErrUnsupportedDistribution
	}
}

func Uint8VectorGenerator(d Distribution) (func(int, int) [][]uint8, error) {
	switch d {
	case Gaussian:
		return GaussianDistributedUint8VectorGenerator, nil
	case Uniform:
		return UniformDistributedUint8VectorGenerator, nil
	default:
		return nil, ErrUnsupportedDistribution
	}
}

func float32Generator(n, dim int, gen func() float32) (ret [][]float32) {
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

func UniformDistributedFloat32VectorGenerator(n, dim int) [][]float32 {
	return float32Generator(n, dim, rand.Float32)
}

func GaussianDistributedFloat32VectorGenerator(n, dim int) [][]float32 {
	return float32Generator(n, dim, func() float32 {
		return float32(rand.NormFloat64())
	})
}

func uint8Generator(n, dim int, gen func() uint8) (ret [][]uint8) {
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

func UniformDistributedUint8VectorGenerator(n, dim int) [][]uint8 {
	return uint8Generator(n, dim, func() uint8 {
		return uint8(rand.Intn(int(math.MaxUint8) + 1))
	})
}

func GaussianDistributedUint8VectorGenerator(n, dim int) [][]uint8 {
	return uint8Generator(n, dim, func() uint8 {
		if val := rand.NormFloat64() + (math.MaxUint8 / 2); val < 0 {
			return 0
		} else if val > math.MaxUint8 {
			return math.MaxUint8
		} else {
			return uint8(val)
		}
	})
}
