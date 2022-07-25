//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
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
	// NOTE: mean:128, sigma:128/3, all of 99.7% are in [0, 255]
	const (
		mean  float64 = 128
		sigma float64 = 128 / 3
	)
	return gaussianDistributedUint8VectorGenerator(n, dim, mean, sigma)
}

// gaussianDistributedUint8VectorGenerator returns n uint8 vectors with dim dimension and their values under Gaussian distribution with user-specified mean and sigma
func gaussianDistributedUint8VectorGenerator(n, dim int, mean, sigma float64) [][]uint8 {
	// NOTE: The boundary test is the main purpose for refactoring. Now, passing this function is dependent on the seed of the random generator. We should fix the randomness of the passing test.
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

// GenF32Vec returns multiple float32 vectors.
func GenF32Vec(dist Distribution, num int, dim int) ([][]float32, error) {
	generator, err := Float32VectorGenerator(dist)
	if err != nil {
		return nil, err
	}
	return generator(num, dim), nil
}

// GenUint8Vec returns multiple uint8 vectors.
func GenUint8Vec(dist Distribution, num int, dim int) ([][]float32, error) {
	generator, err := Uint8VectorGenerator(dist)
	if err != nil {
		return nil, err
	}
	return ConvertVectorsUint8ToFloat32(generator(num, dim)), nil
}

// GenSameValueVec returns a float32 vector filled with value.
func GenSameValueVec(size int, val float32) []float32 {
	v := make([]float32, size)
	for i := 0; i < size; i++ {
		v[i] = val
	}
	return v
}

// ConvertVectorsUint8ToFloat32 converts uint8 vectors and return float32 vectors
func ConvertVectorsUint8ToFloat32(vectors [][]uint8) (ret [][]float32) {
	ret = make([][]float32, 0, len(vectors))
	for _, vec := range vectors {
		fvec := make([]float32, len(vec))
		for i, v := range vec {
			fvec[i] = float32(v)
		}
		ret = append(ret, fvec)
	}
	return
}
