package config

// CircuitBreaker represents the configuration for the internal circuitbreaker package.
type CircuitBreaker struct {
	ClosedErrorRate float32 `yaml:"closed_error_rate" json:"closed_error_rate,omitempty"`
	OpenTimeout     string  `yaml:"open_timeout" json:"open_timeout,omitempty"`
}

func (cb *CircuitBreaker) Bind() *CircuitBreaker {
	cb.OpenTimeout = GetActualValue(cb.OpenTimeout)
	return cb
}
