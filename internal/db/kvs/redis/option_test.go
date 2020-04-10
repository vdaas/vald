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

// Package redis provides implementation of Go API for redis interface
package redis

import (
	"context"
	"crypto/tls"
	"reflect"
	"testing"

	redis "github.com/go-redis/redis/v7"
	"github.com/vdaas/vald/internal/net"
)

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

func TestWithAddrs(t *testing.T) {
	type args struct {
		addrs []string
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
			if got := WithAddrs(tt.args.addrs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithAddrs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDB(t *testing.T) {
	type args struct {
		db int
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

func TestWithClusterSlots(t *testing.T) {
	type args struct {
		f func() ([]redis.ClusterSlot, error)
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
			if got := WithClusterSlots(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithClusterSlots() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDialTimeout(t *testing.T) {
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
			if got := WithDialTimeout(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDialTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithIdleCheckFrequency(t *testing.T) {
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
			if got := WithIdleCheckFrequency(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithIdleCheckFrequency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithIdleTimeout(t *testing.T) {
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
			if got := WithIdleTimeout(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithIdleTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithKeyPrefix(t *testing.T) {
	type args struct {
		prefix string
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
			if got := WithKeyPrefix(tt.args.prefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithKeyPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMaximumConnectionAge(t *testing.T) {
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
			if got := WithMaximumConnectionAge(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMaximumConnectionAge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithRedirectLimit(t *testing.T) {
	type args struct {
		maxRedirects int
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
			if got := WithRedirectLimit(tt.args.maxRedirects); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithRedirectLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithRetryLimit(t *testing.T) {
	type args struct {
		maxRetries int
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
			if got := WithRetryLimit(tt.args.maxRetries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithRetryLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMaximumRetryBackoff(t *testing.T) {
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
			if got := WithMaximumRetryBackoff(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMaximumRetryBackoff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMinimumIdleConnection(t *testing.T) {
	type args struct {
		minIdleConns int
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
			if got := WithMinimumIdleConnection(tt.args.minIdleConns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMinimumIdleConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMinimumRetryBackoff(t *testing.T) {
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
			if got := WithMinimumRetryBackoff(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMinimumRetryBackoff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithOnConnectFunction(t *testing.T) {
	type args struct {
		f func(*redis.Conn) error
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
			if got := WithOnConnectFunction(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithOnConnectFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithOnNewNodeFunction(t *testing.T) {
	type args struct {
		f func(*redis.Client)
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
			if got := WithOnNewNodeFunction(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithOnNewNodeFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPassword(t *testing.T) {
	type args struct {
		password string
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
			if got := WithPassword(tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPoolSize(t *testing.T) {
	type args struct {
		poolSize int
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
			if got := WithPoolSize(tt.args.poolSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPoolSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPoolTimeout(t *testing.T) {
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
			if got := WithPoolTimeout(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPoolTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithReadOnlyFlag(t *testing.T) {
	type args struct {
		readOnly bool
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
			if got := WithReadOnlyFlag(tt.args.readOnly); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithReadOnlyFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithReadTimeout(t *testing.T) {
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
			if got := WithReadTimeout(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithReadTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithRouteByLatencyFlag(t *testing.T) {
	type args struct {
		routeByLatency bool
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
			if got := WithRouteByLatencyFlag(tt.args.routeByLatency); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithRouteByLatencyFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithRouteRandomlyFlag(t *testing.T) {
	type args struct {
		routeRandomly bool
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
			if got := WithRouteRandomlyFlag(tt.args.routeRandomly); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithRouteRandomlyFlag() = %v, want %v", got, tt.want)
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

func TestWithWriteTimeout(t *testing.T) {
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
			if got := WithWriteTimeout(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithWriteTimeout() = %v, want %v", got, tt.want)
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
