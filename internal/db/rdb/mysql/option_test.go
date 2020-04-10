//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

package mysql

import (
	"context"
	"crypto/tls"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/net"
)

func TestWithTimezone(t *testing.T) {
	type args struct {
		tz string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTimezone(tt.args.tz); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTimezone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCharset(t *testing.T) {
	type args struct {
		cs string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCharset(tt.args.cs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCharset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDB(t *testing.T) {
	type args struct {
		db string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDB(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithHost(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithHost(tt.args.host); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPort(t *testing.T) {
	type args struct {
		port int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithPort(tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithUser(t *testing.T) {
	type args struct {
		user string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithUser(tt.args.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPass(t *testing.T) {
	type args struct {
		pass string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithPass(tt.args.pass); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPass() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithInitialPingTimeLimit(t *testing.T) {
	type args struct {
		lim string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithInitialPingTimeLimit(tt.args.lim); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithInitialPingTimeLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithInitialPingDuration(t *testing.T) {
	type args struct {
		dur string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithInitialPingDuration(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithInitialPingDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithConnectionLifeTimeLimit(t *testing.T) {
	type args struct {
		dur string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithConnectionLifeTimeLimit(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithConnectionLifeTimeLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMaxIdleConns(t *testing.T) {
	type args struct {
		conns int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMaxIdleConns(tt.args.conns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMaxIdleConns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMaxOpenConns(t *testing.T) {
	type args struct {
		conns int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMaxOpenConns(tt.args.conns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMaxOpenConns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTLSConfig(t *testing.T) {
	type args struct {
		cfg *tls.Config
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTLSConfig(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTLSConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDialer(t *testing.T) {
	type args struct {
		der func(ctx context.Context, addr, port string) (net.Conn, error)
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDialer(tt.args.der); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDialer() = %v, want %v", got, tt.want)
			}
		})
	}
}
