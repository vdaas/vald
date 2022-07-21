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
package circuitbreaker

// Tripper is a function type to determine if the CircuitBreaker should trip.
type Tripper interface {
	ShouldTrip(Counter) bool
}

type TripperFunc func(Counter) bool

func (f TripperFunc) ShouldTrip(c Counter) bool {
	return f(c)
}

func NewRateTripper(rate float32, min int64) Tripper {
	return TripperFunc(func(c Counter) bool {
		successes, fails := c.Successes(), c.Fails()

		if fails+successes <= min {
			return false
		}
		return float32(fails)/float32(successes+fails) >= rate
	})
}
