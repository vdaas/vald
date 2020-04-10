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
package mock

import "testing"

func TestLogger_Debug(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	tests := []struct {
		name string
		l    *Logger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Debug(tt.args.vals...)
		})
	}
}

func TestLogger_Debugf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	tests := []struct {
		name string
		l    *Logger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Debugf(tt.args.format, tt.args.vals...)
		})
	}
}

func TestLogger_Info(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	tests := []struct {
		name string
		l    *Logger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Info(tt.args.vals...)
		})
	}
}

func TestLogger_Infof(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	tests := []struct {
		name string
		l    *Logger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Infof(tt.args.format, tt.args.vals...)
		})
	}
}

func TestLogger_Warn(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	tests := []struct {
		name string
		l    *Logger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Warn(tt.args.vals...)
		})
	}
}

func TestLogger_Warnf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	tests := []struct {
		name string
		l    *Logger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Warnf(tt.args.format, tt.args.vals...)
		})
	}
}

func TestLogger_Error(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	tests := []struct {
		name string
		l    *Logger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Error(tt.args.vals...)
		})
	}
}

func TestLogger_Errorf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	tests := []struct {
		name string
		l    *Logger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Errorf(tt.args.format, tt.args.vals...)
		})
	}
}

func TestLogger_Fatal(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	tests := []struct {
		name string
		l    *Logger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Fatal(tt.args.vals...)
		})
	}
}

func TestLogger_Fatalf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	tests := []struct {
		name string
		l    *Logger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Fatalf(tt.args.format, tt.args.vals...)
		})
	}
}
