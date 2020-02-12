package tls

import (
	"crypto/tls"
	"reflect"
	"testing"

	"github.com/cockroachdb/errors"
)

func TestWithCert(t *testing.T) {
	type test struct {
		name      string
		cert      string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success when cert is not empty",
			cert: "cert",
			checkFunc: func(opt Option) error {
				got := new(credentials)
				if err := opt(got); err != nil {
					return errors.Errorf("err is not nil: %v", err)
				}

				if got.cert != "cert" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when cert is empty",
			checkFunc: func(opt Option) error {
				got := &credentials{
					cert: "cert",
				}
				if err := opt(got); err != nil {
					return errors.Errorf("err is not nil: %v", err)
				}

				if got.cert != "cert" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithCert(tt.cert)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithKey(t *testing.T) {
	type test struct {
		name      string
		key       string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success when key is not empty",
			key:  "key",
			checkFunc: func(opt Option) error {
				got := new(credentials)
				if err := opt(got); err != nil {
					return errors.Errorf("err is not nil: %v", err)
				}

				if got.key != "key" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when key is empty",
			checkFunc: func(opt Option) error {
				got := &credentials{
					key: "key",
				}
				if err := opt(got); err != nil {
					return errors.Errorf("err is not nil: %v", err)
				}

				if got.key != "key" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithKey(tt.key)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithCa(t *testing.T) {
	type test struct {
		name      string
		ca        string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success when ca is not empty",
			ca:   "ca",
			checkFunc: func(opt Option) error {
				got := new(credentials)
				if err := opt(got); err != nil {
					return errors.Errorf("err is not nil: %v", err)
				}

				if got.ca != "ca" {
					return errors.New("invalid param was set")
				}

				return nil
			},
		},

		{
			name: "set set when ca is empty",
			checkFunc: func(opt Option) error {
				got := &credentials{
					ca: "ca",
				}
				if err := opt(got); err != nil {
					return errors.Errorf("err is not nil: %v", err)
				}

				if got.ca != "ca" {
					return errors.New("invalid param was set")
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithCa(tt.ca)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithTLSConfig(t *testing.T) {
	type test struct {
		name      string
		cfg       *tls.Config
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			cfg := new(tls.Config)

			return test{
				name: "set success when cfg is not nil",
				cfg:  cfg,
				checkFunc: func(opt Option) error {
					got := new(credentials)
					if err := opt(got); err != nil {
						return errors.Errorf("err is not nil: %v", err)
					}

					if !reflect.DeepEqual(got.cfg, cfg) {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),

		{
			name: "not set when cfg is nil",
			checkFunc: func(opt Option) error {
				cfg := new(tls.Config)

				got := &credentials{
					cfg: cfg,
				}

				if err := opt(got); err != nil {
					return errors.Errorf("err is not nil: %v", err)
				}

				if !reflect.DeepEqual(got.cfg, cfg) {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithTLSConfig(tt.cfg)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
