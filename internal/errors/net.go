// Package errors provides error types and function
package errors

import "time"

var (
	// tcp

	ErrFailedInitDialer = New("failed to init dialer")
	ErrInvalidDNSConfig = func(dnsRefreshDur, dnsCacheExp time.Duration) error {
		return Errorf("dnsRefreshDuration  > dnsCacheExp, %s, %s", dnsRefreshDur, dnsCacheExp)
	}
)
