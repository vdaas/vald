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

// Package config providers configuration type and load configuration logic
package config

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/net/tcp"
)

func TestDNS_Bind(t *testing.T) {
	tests := []struct {
		name string
		d    *DNS
		want *DNS
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Bind(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DNS.Bind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDialer_Bind(t *testing.T) {
	tests := []struct {
		name string
		d    *Dialer
		want *Dialer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Bind(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dialer.Bind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCP_Bind(t *testing.T) {
	tests := []struct {
		name string
		t    *TCP
		want *TCP
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Bind(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCP.Bind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCP_Opts(t *testing.T) {
	tests := []struct {
		name string
		t    *TCP
		want []tcp.DialerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Opts(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCP.Opts() = %v, want %v", got, tt.want)
			}
		})
	}
}
