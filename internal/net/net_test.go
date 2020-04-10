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

// Package net provides net functionality for grpc
package net

import (
	"context"
	"reflect"
	"testing"
)

func TestListen(t *testing.T) {
	type args struct {
		network string
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    Listener
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Listen(tt.args.network, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("Listen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Listen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsLocal(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLocal(tt.args.host); got != tt.want {
				t.Errorf("IsLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name     string
		args     args
		wantHost string
		wantPort string
		wantIsIP bool
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotPort, gotIsIP, err := Parse(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHost != tt.wantHost {
				t.Errorf("Parse() gotHost = %v, want %v", gotHost, tt.wantHost)
			}
			if gotPort != tt.wantPort {
				t.Errorf("Parse() gotPort = %v, want %v", gotPort, tt.wantPort)
			}
			if gotIsIP != tt.wantIsIP {
				t.Errorf("Parse() gotIsIP = %v, want %v", gotIsIP, tt.wantIsIP)
			}
		})
	}
}

func TestIsIPv6(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIPv6(tt.args.addr); got != tt.want {
				t.Errorf("IsIPv6() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIPv4(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIPv4(tt.args.addr); got != tt.want {
				t.Errorf("IsIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitHostPort(t *testing.T) {
	type args struct {
		hostport string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := SplitHostPort(tt.args.hostport)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitHostPort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SplitHostPort() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SplitHostPort() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestScanPorts(t *testing.T) {
	type args struct {
		ctx   context.Context
		start uint16
		end   uint16
		host  string
	}
	tests := []struct {
		name      string
		args      args
		wantPorts []uint16
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPorts, err := ScanPorts(tt.args.ctx, tt.args.start, tt.args.end, tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanPorts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPorts, tt.wantPorts) {
				t.Errorf("ScanPorts() = %v, want %v", gotPorts, tt.wantPorts)
			}
		})
	}
}
