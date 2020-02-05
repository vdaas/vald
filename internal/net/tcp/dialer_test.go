package tcp

import (
	"context"
	"crypto/tls"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/kpango/gache"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

func TestNewDialer(t *testing.T) {
	type test struct {
		name      string
		opts      []DialerOption
		checkFunc func(got *dialer) error
	}

	tests := []test{
		{
			name: "returns dialer when option is empty",
			checkFunc: func(d *dialer) error {
				if d == nil {
					return errors.New("dialer is nil")
				}

				if d.der == nil {
					return errors.New("der is nil")
				}

				if d.der.KeepAlive != 30*time.Second {
					return errors.New("invalid param was set about der.KeepAlive")
				}

				if d.der.DualStack != true {
					return errors.New("invalid param was set about der.DualStack")
				}

				if d.dnsCache != false {
					return errors.New("invalid param was set about dnsCache")
				}

				return nil
			},
		},

		{
			name: "returns dialer when option is not empty",
			opts: []DialerOption{
				WithTLS(new(tls.Config)),
				WithCache(gache.New()),
				WithEnableDNSCache(),
				WithDNSRefreshDuration("10s"),
				WithDNSCacheExpiration("5s"),
			},
			checkFunc: func(d *dialer) error {
				if d == nil {
					return errors.New("dialer is nil")
				}

				if d.cache == nil {
					return errors.New("invalid param was set about cache")
				}

				if d.dnsRefreshDuration != 5*time.Second {
					return errors.New("invalid param was set about dnsRefreshDuration")
				}

				if d.dnsCacheExpiration != 10*time.Second {
					return errors.New("invalid param was set about dnsCacheExpiration")
				}

				return nil
			},
		},

		{
			name: "returns dialer when tls option is not empty",
			opts: []DialerOption{
				WithTLS(new(tls.Config)),
			},
			checkFunc: func(d *dialer) error {
				if d == nil {
					return errors.New("dialer is nil")
				}

				if d.tlsConfig == nil {
					return errors.New("invalid param was set about tlsConfig")
				}

				if d.dialer == nil {
					return errors.New("invalid param was set about dialer")
				}

				if d.cache != nil {
					return errors.New("invalid param was set about cache")
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDialer(tt.opts...)
			if err := tt.checkFunc(got.(*dialer)); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGetDialer(t *testing.T) {
	type test struct {
		name string
		fn   func(ctx context.Context, network, addr string) (net.Conn, error)
	}

	tests := []test{
		func() test {
			fn := func(ctx context.Context, network, addr string) (net.Conn, error) {
				return nil, nil
			}
			return test{
				name: "get success",
				fn:   fn,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dialer{
				dialer: tt.fn,
			}
			got := d.GetDialer()
			if reflect.ValueOf(tt.fn).Pointer() != reflect.ValueOf(got).Pointer() {
				t.Errorf("not equals.")
			}
		})
	}
}

func Test_lookup(t *testing.T) {
	type args struct {
		ctx  context.Context
		addr string
	}

	type field struct {
		cache gache.Gache
		der   *net.Dialer
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func(map[int]string, error) error
	}

	tests := []test{
		func() test {
			return test{
				name: "return ips and nil when lookupIpAddr returns ips",
				args: args{
					ctx:  context.Background(),
					addr: "google.com",
				},
				field: field{
					cache: gache.New(),
					der: &net.Dialer{
						Resolver: &net.Resolver{
							PreferGo: false,
						},
					},
				},
				checkFunc: func(ips map[int]string, err error) error {
					if err != nil {
						return errors.New("err is not nil")
					}

					if len(ips) == 0 {
						return errors.New("ips is empty")
					}
					return nil
				},
			}
		}(),

		func() test {
			addr := "addr"

			wantIPs := map[int]string{
				1: "0.0.0.0",
			}

			cache := gache.New()
			cache.SetWithExpire(addr, wantIPs, 1*time.Hour)

			return test{
				name: "return ips and nil when the cache hits",
				args: args{
					addr: addr,
				},
				field: field{
					cache: cache,
				},
				checkFunc: func(ips map[int]string, err error) error {
					if err != nil {
						return errors.New("err is not nil")
					}

					if !reflect.DeepEqual(ips, wantIPs) {
						return errors.Errorf("not equals. want: %v, but got: %v", wantIPs, ips)
					}
					return nil
				},
			}
		}(),

		{
			name: "return nil and error when addr is empty",
			args: args{
				ctx: context.Background(),
			},
			field: field{
				cache: gache.GetGache(),
				der: &net.Dialer{
					Resolver: &net.Resolver{
						PreferGo: false,
					},
				},
			},
			checkFunc: func(ips map[int]string, err error) error {
				if ips != nil {
					return errors.New("ips is not nil")
				}

				if err == nil {
					return errors.New("err is nil")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dialer{
				cache: tt.field.cache,
				der:   tt.field.der,
			}
			ips, err := d.lookup(tt.args.ctx, tt.args.addr)
			if err := tt.checkFunc(ips, err); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestStartDialerCache(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	type field struct {
		dnsCache           bool
		cache              gache.Gache
		dnsRefreshDuration time.Duration
		dnsCacheExpiration time.Duration
		der                *net.Dialer
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func(cache gache.Gache) error
	}

	tests := []test{
		func() test {
			addr := "google.com"
			ips := map[int]string{
				1: "0.0.0.0",
			}
			cache := gache.New()
			cache.SetWithExpire(addr, ips, 1*time.Nanosecond)

			return test{
				name: "hook is called when expired",
				args: args{
					ctx: context.Background(),
				},
				field: field{
					dnsCache:           true,
					cache:              cache,
					dnsRefreshDuration: 1 * time.Nanosecond,
					dnsCacheExpiration: 1 * time.Nanosecond,
					der: &net.Dialer{
						Resolver: &net.Resolver{
							PreferGo: false,
						}},
				},
				checkFunc: func(cache gache.Gache) error {
					time.Sleep(1 * time.Second)
					val, _ := cache.Get(addr)
					if reflect.DeepEqual(val, ips) {
						return errors.New("cache is not cleared")
					}
					return nil
				},
			}
		}(),
	}

	log.Init()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dialer{
				dnsCache:           tt.field.dnsCache,
				cache:              tt.field.cache,
				dnsRefreshDuration: tt.field.dnsRefreshDuration,
				dnsCacheExpiration: tt.field.dnsRefreshDuration,
				der:                tt.field.der,
			}
			d.StartDialerCache(tt.args.ctx)
			if err := tt.checkFunc(d.cache); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestDialContext(t *testing.T) {
	type args struct {
		ctx     context.Context
		network string
		address string
	}

	type field struct {
		dialer func(ctx context.Context, network, addr string) (net.Conn, error)
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func(net.Conn, error) error
	}

	tests := []test{
		func() test {
			wantConn := new(net.IPConn)
			want := func(ctx context.Context, network, address string) (net.Conn, error) {
				return wantConn, nil
			}

			return test{
				name: "returns conn and nil when dialer returns conn and nil",
				args: args{
					ctx:     context.Background(),
					network: "network",
					address: "address",
				},
				field: field{
					dialer: want,
				},
				checkFunc: func(gotConn net.Conn, gotErr error) error {
					if gotErr != nil {
						return errors.New("err is not nil")
					}

					if !reflect.DeepEqual(gotConn, wantConn) {
						return errors.Errorf("conn not equals. want: %v, got: %v", wantConn, gotConn)
					}
					return nil
				},
			}
		}(),

		func() test {
			var wantErr error = errors.New("fail")
			want := func(ctx context.Context, network, address string) (net.Conn, error) {
				return nil, wantErr
			}

			return test{
				name: "returns nil and error when dialer returns nil and error",
				args: args{
					ctx:     context.Background(),
					network: "network",
					address: "address",
				},
				field: field{
					dialer: want,
				},
				checkFunc: func(gotConn net.Conn, gotErr error) error {
					if gotErr == nil {
						return errors.New("err is nil")
					} else if !reflect.DeepEqual(wantErr, gotErr) {
						return errors.Errorf("err not equals. want: %v, got: %v", wantErr, gotErr)
					}

					if gotConn != nil {
						return errors.New("conn is not nil")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dialer{
				dialer: tt.field.dialer,
			}

			conn, err := d.DialContext(tt.args.ctx, tt.args.network, tt.args.address)
			if err := tt.checkFunc(conn, err); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_cachedDialer(t *testing.T) {
	type args struct {
		dctx    context.Context
		network string
		address string
	}

	type field struct {
		der   *net.Dialer
		cache gache.Gache
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func(net.Conn, error) error
	}

	tests := []test{
		func() test {
			return test{
				name: "returns conn and nil when dialer returns conn and nil",
				args: args{
					dctx:    context.Background(),
					network: "tcp",
					address: "google.com",
				},
				field: field{
					der: &net.Dialer{
						Resolver: &net.Resolver{
							PreferGo: false,
						},
					},
					cache: gache.New(),
				},
				checkFunc: func(gotConn net.Conn, gotErr error) error {
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dialer{
				der:   tt.field.der,
				cache: tt.field.cache,
			}

			conn, err := d.cachedDialer(tt.args.dctx, tt.args.network, tt.args.address)
			if err := tt.checkFunc(conn, err); err != nil {
				t.Error(err)
			}
		})
	}
}
