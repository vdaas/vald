package errors

import (
	"math"
	"testing"
	"time"
)

func TestErrFailedInitDialer(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrFailedInitDialer error",
			want: want{
				want: New("failed to init dialer"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrFailedInitDialer
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrInvalidDNSConfig(t *testing.T) {
	type args struct {
		dnsRefreshDur time.Duration
		dnsCacheExp   time.Duration
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrInvalidDNSConfig when dnsRefreshDur is 5 minute, dnsCacheExp is 4 minute",
			args: args{
				dnsRefreshDur: 5 * time.Minute,
				dnsCacheExp:   4 * time.Minute,
			},
			want: want{
				want: New("dnsRefreshDuration  > dnsCacheExp, 5m0s, 4m0s"),
			},
		},
		{
			name: "return an ErrInvalidDNSConfig when dnsRefreshDur is 0, dnsCacheExp is 0",
			args: args{
				dnsRefreshDur: 0,
				dnsCacheExp:   0,
			},
			want: want{
				want: New("dnsRefreshDuration  > dnsCacheExp, 0s, 0s"),
			},
		},
		{
			name: "return an ErrInvalidDNSConfig when dnsRefreshDur and dnsCacheExp are the minimum number of int64",
			args: args{
				dnsRefreshDur: time.Duration(math.MinInt64),
				dnsCacheExp:   time.Duration(math.MinInt64),
			},
			want: want{
				want: Errorf("dnsRefreshDuration  > dnsCacheExp, %s, %s", time.Duration(math.MinInt64), time.Duration(math.MinInt64)),
			},
		},
		{
			name: "return an ErrInvalidDNSConfig when dnsRefreshDur and dnsCacheExp are the maximum number of int64",
			args: args{
				dnsRefreshDur: time.Duration(math.MaxInt64),
				dnsCacheExp:   time.Duration(math.MaxInt64),
			},
			want: want{
				want: Errorf("dnsRefreshDuration  > dnsCacheExp, %s, %s", time.Duration(math.MaxInt64), time.Duration(math.MaxInt64)),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrInvalidDNSConfig(test.args.dnsRefreshDur, test.args.dnsCacheExp)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrNoPortAvailable(t *testing.T) {
	type args struct {
		host  string
		start uint16
		end   uint16
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return an ErrNoPortAvailable when host is localhost and start is 65534 and end is 65535",
			args: args{
				host:  "localhost",
				start: 65534,
				end:   65535,
			},
			want: want{
				want: New("no port available for Host: localhost\tbetween 65534 ~ 65535"),
			},
		},
		{
			name: "return an ErrNoPortAvailable when host is empty and start is 0 and end is 0",
			args: args{
				host:  "",
				start: 0,
				end:   0,
			},
			want: want{
				want: New("no port available for Host: \tbetween 0 ~ 0"),
			},
		},
		{
			name: "return an ErrNoPortAvailable when host is empty and start is 65534 and end is 65535",
			args: args{
				host:  "",
				start: 65534,
				end:   65535,
			},
			want: want{
				want: New("no port available for Host: \tbetween 65534 ~ 65535"),
			},
		},
		{
			name: "return an ErrNoPortAvailable when host is localhost and start is 0 and end is 65535",
			args: args{
				host:  "localhost",
				start: 0,
				end:   65535,
			},
			want: want{
				want: New("no port available for Host: localhost\tbetween 0 ~ 65535"),
			},
		},
		{
			name: "return an ErrNoPortAvailable when host is localhost and start is 65534 and end is 0",
			args: args{
				host:  "localhost",
				start: 65534,
				end:   0,
			},
			want: want{
				want: New("no port available for Host: localhost\tbetween 65534 ~ 0"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ErrNoPortAvailable(test.args.host, test.args.start, test.args.end)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
