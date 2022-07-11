package config

// CircuitBreaker represents the configuration for the internal circuitbreaker package.
type CircuitBreaker struct {
	ClosedErrorRate   float32 `yaml:"closed_error_rate" json:"closed_error_rate,omitempty"`
	HalfOpenErrorRate float32 `yaml:"half_open_error_rate" json:"half_open_error_rate,omitempty"`
	MinSamples        int64   `yaml:"min_samples" json:"min_samples,omitempty"`
	OpenTimeout       string  `yaml:"open_timeout" json:"open_timeout,omitempty"`
}

func (cb *CircuitBreaker) Bind() *CircuitBreaker {
	cb.OpenTimeout = GetActualValue(cb.OpenTimeout)
	return cb
}
