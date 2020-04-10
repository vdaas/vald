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

// Package tcp provides tcp option
package tcp

import (
	"crypto/tls"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/cache"
)

func TestWithCache(t *testing.T) {
	type args struct {
		c cache.Cache
	}
	tests := []struct {
		name string
		args args
		want DialerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCache(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDNSRefreshDuration(t *testing.T) {
	type args struct {
		dur string
	}
	tests := []struct {
		name string
		args args
		want DialerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDNSRefreshDuration(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDNSRefreshDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDNSCacheExpiration(t *testing.T) {
	type args struct {
		dur string
	}
	tests := []struct {
		name string
		args args
		want DialerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDNSCacheExpiration(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDNSCacheExpiration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDialerTimeout(t *testing.T) {
	type args struct {
		dur string
	}
	tests := []struct {
		name string
		args args
		want DialerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDialerTimeout(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDialerTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDialerKeepAlive(t *testing.T) {
	type args struct {
		dur string
	}
	tests := []struct {
		name string
		args args
		want DialerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDialerKeepAlive(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDialerKeepAlive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTLS(t *testing.T) {
	type args struct {
		cfg *tls.Config
	}
	tests := []struct {
		name string
		args args
		want DialerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTLS(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTLS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithEnableDNSCache(t *testing.T) {
	tests := []struct {
		name string
		want DialerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithEnableDNSCache(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithEnableDNSCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisableDNSCache(t *testing.T) {
	tests := []struct {
		name string
		want DialerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDisableDNSCache(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisableDNSCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithEnableDialerDualStack(t *testing.T) {
	tests := []struct {
		name string
		want DialerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithEnableDialerDualStack(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithEnableDialerDualStack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisableDialerDualStack(t *testing.T) {
	tests := []struct {
		name string
		want DialerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDisableDialerDualStack(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisableDialerDualStack() = %v, want %v", got, tt.want)
			}
		})
	}
}
