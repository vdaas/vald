package tcp

import (
	"crypto/tls"
	"reflect"
	"testing"
	"time"

	"github.com/kpango/gache"
	"github.com/vdaas/vald/internal/errors"
)

func TestWithCache(t *testing.T) {
	type test struct {
		name      string
		c         gache.Gache
		checkFunc func(DialerOption) error
	}

	tests := []test{
		func() test {
			c := gache.New()

			return test{
				name: "set success",
				c:    c,
				checkFunc: func(opt DialerOption) error {
					got := new(dialer)
					opt(got)

					if !reflect.DeepEqual(got.cache, c) {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithCache(tt.c)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithDNSRefreshDuration(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(DialerOption) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "10s",
			checkFunc: func(opt DialerOption) error {
				got := new(dialer)
				opt(got)

				if got.dnsRefreshDuration != 10*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "set default",
			dur:  "vald",
			checkFunc: func(opt DialerOption) error {
				got := new(dialer)
				opt(got)

				if got.dnsRefreshDuration != 30*time.Minute {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing",
			checkFunc: func(opt DialerOption) error {
				got := new(dialer)
				opt(got)

				if got.dnsRefreshDuration != 0 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithDNSRefreshDuration(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithDNSCacheExpiration(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(DialerOption) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "10s",
			checkFunc: func(opt DialerOption) error {
				got := new(dialer)
				opt(got)

				if got.dnsCacheExpiration != 10*time.Second {
					return errors.New("invalid param (dnsCacheExpiration) was set")
				}

				if !got.dnsCache {
					return errors.New("invalid param (dnsCache) was set")
				}

				return nil
			},
		},

		{
			name: "set default",
			dur:  "vald",
			checkFunc: func(opt DialerOption) error {
				got := new(dialer)
				opt(got)

				if got.dnsCacheExpiration != time.Hour {
					return errors.New("invalid param (dnsCacheExpiration) was set")
				}

				if !got.dnsCache {
					return errors.New("invalid param (dnsCache) was set")
				}

				return nil
			},
		},

		{
			name: "do nothing",
			checkFunc: func(opt DialerOption) error {
				got := new(dialer)
				opt(got)

				if got.dnsCacheExpiration != 0 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithDNSCacheExpiration(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithDialerTimeout(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(DialerOption) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "10s",
			checkFunc: func(opt DialerOption) error {
				got := new(dialer)
				opt(got)

				if got.dialerTimeout != 10*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "set default",
			dur:  "vald",
			checkFunc: func(opt DialerOption) error {
				got := new(dialer)
				opt(got)

				if got.dialerTimeout != 30*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing",
			checkFunc: func(opt DialerOption) error {
				got := new(dialer)
				opt(got)

				if got.dialerTimeout != 0 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithDialerTimeout(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithDialerKeepAlive(t *testing.T) {
	type test struct {
		name      string
		dur       string
		checkFunc func(DialerOption) error
	}

	tests := []test{
		{
			name: "set success",
			dur:  "10s",
			checkFunc: func(opt DialerOption) error {
				got := new(dialer)
				opt(got)

				if got.dialerKeepAlive != 10*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "set default",
			dur:  "vald",
			checkFunc: func(opt DialerOption) error {
				got := new(dialer)
				opt(got)

				if got.dialerKeepAlive != 30*time.Second {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "do nothing",
			checkFunc: func(opt DialerOption) error {
				got := new(dialer)
				opt(got)

				if got.dialerKeepAlive != 0 {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithDialerKeepAlive(tt.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithTLS(t *testing.T) {
	type test struct {
		name      string
		cfg       *tls.Config
		checkFunc func(DialerOption) error
	}

	tests := []test{
		func() test {
			cfg := new(tls.Config)
			return test{
				name: "set success",
				cfg:  cfg,
				checkFunc: func(opt DialerOption) error {
					got := new(dialer)
					opt(got)

					if !reflect.DeepEqual(got.tlsConfig, cfg) {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithTLS(tt.cfg)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithEnableDNSCache(t *testing.T) {
	type test struct {
		name      string
		checkFunc func(DialerOption) error
	}

	tests := []test{
		{
			name: "set success",
			checkFunc: func(opt DialerOption) error {
				got := new(dialer)
				opt(got)

				if !got.dnsCache {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithEnableDNSCache()
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithDisableDNSCache(t *testing.T) {
	type test struct {
		name      string
		checkFunc func(DialerOption) error
	}

	tests := []test{
		{
			name: "set success",
			checkFunc: func(opt DialerOption) error {
				got := &dialer{
					dnsCache: true,
				}
				opt(got)

				if got.dnsCache {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithDisableDNSCache()
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithEnableDialerDualStack(t *testing.T) {
	type test struct {
		name      string
		checkFunc func(DialerOption) error
	}

	tests := []test{
		{
			name: "set success",
			checkFunc: func(opt DialerOption) error {
				got := new(dialer)
				opt(got)

				if !got.dialerDualStack {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithEnableDialerDualStack()
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithDisableDialerDualStack(t *testing.T) {
	type test struct {
		name      string
		checkFunc func(DialerOption) error
	}

	tests := []test{
		{
			name: "set success",
			checkFunc: func(opt DialerOption) error {
				got := &dialer{
					dialerDualStack: true,
				}
				opt(got)

				if got.dialerDualStack {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithDisableDialerDualStack()
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
