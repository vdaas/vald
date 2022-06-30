package circuitbreaker

// Tripper is a function type to determine if the CircuitBreaker should trip.
type Tripper interface {
	ShouldTrip(Counter) bool
}

type TripperFunc func(Counter) bool

func (f TripperFunc) ShouldTrip(c Counter) bool {
	return f(c)
}

func NewRateTripper(rate float32) Tripper {
	return TripperFunc(func(c Counter) bool {
		successes, fails := c.Successes(), c.Fails()
		return float32(fails/successes+fails) > rate
	})
}
