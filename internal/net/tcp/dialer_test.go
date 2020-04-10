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
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/net"
)

func TestNewDialer(t *testing.T) {
	type args struct {
		opts []DialerOption
	}
	tests := []struct {
		name    string
		args    args
		wantDer Dialer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDer, err := NewDialer(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDialer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDer, tt.wantDer) {
				t.Errorf("NewDialer() = %v, want %v", gotDer, tt.wantDer)
			}
		})
	}
}

func Test_dialer_GetDialer(t *testing.T) {
	tests := []struct {
		name string
		d    *dialer
		want func(ctx context.Context, network, addr string) (net.Conn, error)
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetDialer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dialer.GetDialer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dialer_lookup(t *testing.T) {
	type args struct {
		ctx  context.Context
		addr string
	}
	tests := []struct {
		name    string
		d       *dialer
		args    args
		wantIps map[int]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIps, err := tt.d.lookup(tt.args.ctx, tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("dialer.lookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIps, tt.wantIps) {
				t.Errorf("dialer.lookup() = %v, want %v", gotIps, tt.wantIps)
			}
		})
	}
}

func Test_dialer_StartDialerCache(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		d    *dialer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.StartDialerCache(tt.args.ctx)
		})
	}
}

func Test_dialer_DialContext(t *testing.T) {
	type args struct {
		ctx     context.Context
		network string
		address string
	}
	tests := []struct {
		name    string
		d       *dialer
		args    args
		want    net.Conn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.DialContext(tt.args.ctx, tt.args.network, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("dialer.DialContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dialer.DialContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dialer_cachedDialer(t *testing.T) {
	type args struct {
		dctx    context.Context
		network string
		addr    string
	}
	tests := []struct {
		name     string
		d        *dialer
		args     args
		wantConn net.Conn
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConn, err := tt.d.cachedDialer(tt.args.dctx, tt.args.network, tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("dialer.cachedDialer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotConn, tt.wantConn) {
				t.Errorf("dialer.cachedDialer() = %v, want %v", gotConn, tt.wantConn)
			}
		})
	}
}
