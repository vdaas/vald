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

// Package discoverer
package discoverer

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/net/grpc"
)

func TestWithOnDiscoverFunc(t *testing.T) {
	type args struct {
		f func(ctx context.Context, c Client, addrs []string) error
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
			if got := WithOnDiscoverFunc(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithOnDiscoverFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithOnConnectFunc(t *testing.T) {
	type args struct {
		f func(ctx context.Context, c Client, addr string) error
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
			if got := WithOnConnectFunc(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithOnConnectFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithOnDisconnectFunc(t *testing.T) {
	type args struct {
		f func(ctx context.Context, c Client, addr string) error
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
			if got := WithOnDisconnectFunc(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithOnDisconnectFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscovererClient(t *testing.T) {
	type args struct {
		gc grpc.Client
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
			if got := WithDiscovererClient(tt.args.gc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDiscovererClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscovererAddr(t *testing.T) {
	type args struct {
		addr string
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
			if got := WithDiscovererAddr(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDiscovererAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscovererHostPort(t *testing.T) {
	type args struct {
		host string
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
			if got := WithDiscovererHostPort(tt.args.host, tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDiscovererHostPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscoverDuration(t *testing.T) {
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
			if got := WithDiscoverDuration(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDiscoverDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithOptions(t *testing.T) {
	type args struct {
		opts []grpc.Option
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
			if got := WithOptions(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithAutoConnect(t *testing.T) {
	type args struct {
		flg bool
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
			if got := WithAutoConnect(tt.args.flg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithAutoConnect() = %v, want %v", got, tt.want)
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

func TestWithNamespace(t *testing.T) {
	type args struct {
		ns string
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
			if got := WithNamespace(tt.args.ns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithNamespace() = %v, want %v", got, tt.want)
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

func TestWithServiceDNSARecord(t *testing.T) {
	type args struct {
		a string
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
			if got := WithServiceDNSARecord(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithServiceDNSARecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithNodeName(t *testing.T) {
	type args struct {
		nn string
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
			if got := WithNodeName(tt.args.nn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithNodeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithErrGroup(t *testing.T) {
	type args struct {
		eg errgroup.Group
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
			if got := WithErrGroup(tt.args.eg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithErrGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
