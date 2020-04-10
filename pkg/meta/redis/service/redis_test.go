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

// Package service manages the main logic of server.
package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/config"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg *config.Redis
	}
	tests := []struct {
		name    string
		args    args
		want    Redis
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_Disconnect(t *testing.T) {
	tests := []struct {
		name    string
		c       *client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Disconnect(); (err != nil) != tt.wantErr {
				t.Errorf("client.Disconnect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_Connect(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Connect(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("client.Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("client.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_GetMultiple(t *testing.T) {
	type args struct {
		keys []string
	}
	tests := []struct {
		name     string
		c        *client
		args     args
		wantVals []string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVals, err := tt.c.GetMultiple(tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.GetMultiple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVals, tt.wantVals) {
				t.Errorf("client.GetMultiple() = %v, want %v", gotVals, tt.wantVals)
			}
		})
	}
}

func Test_client_GetInverse(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.GetInverse(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.GetInverse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("client.GetInverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_GetInverseMultiple(t *testing.T) {
	type args struct {
		vals []string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.GetInverseMultiple(tt.args.vals...)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.GetInverseMultiple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("client.GetInverseMultiple() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_appendPrefix(t *testing.T) {
	type args struct {
		prefix string
		key    string
	}
	tests := []struct {
		name string
		c    *client
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.appendPrefix(tt.args.prefix, tt.args.key); got != tt.want {
				t.Errorf("client.appendPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_get(t *testing.T) {
	type args struct {
		prefix string
		key    string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.c.get(tt.args.prefix, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("client.get() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func Test_client_getMulti(t *testing.T) {
	type args struct {
		prefix string
		keys   []string
	}
	tests := []struct {
		name     string
		c        *client
		args     args
		wantVals []string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVals, err := tt.c.getMulti(tt.args.prefix, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.getMulti() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVals, tt.wantVals) {
				t.Errorf("client.getMulti() = %v, want %v", gotVals, tt.wantVals)
			}
		})
	}
}

func Test_client_Set(t *testing.T) {
	type args struct {
		key string
		val string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Set(tt.args.key, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("client.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_SetMultiple(t *testing.T) {
	type args struct {
		kvs map[string]string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.SetMultiple(tt.args.kvs); (err != nil) != tt.wantErr {
				t.Errorf("client.SetMultiple() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Delete(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("client.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_DeleteMultiple(t *testing.T) {
	type args struct {
		keys []string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.DeleteMultiple(tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.DeleteMultiple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("client.DeleteMultiple() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_DeleteInverse(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.DeleteInverse(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.DeleteInverse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("client.DeleteInverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_DeleteInverseMultiple(t *testing.T) {
	type args struct {
		vals []string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.DeleteInverseMultiple(tt.args.vals...)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.DeleteInverseMultiple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("client.DeleteInverseMultiple() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_delete(t *testing.T) {
	type args struct {
		pfx    string
		pfxInv string
		key    string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		wantVal string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := tt.c.delete(tt.args.pfx, tt.args.pfxInv, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("client.delete() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func Test_client_deleteMulti(t *testing.T) {
	type args struct {
		pfx    string
		pfxInv string
		keys   []string
	}
	tests := []struct {
		name     string
		c        *client
		args     args
		wantVals []string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVals, err := tt.c.deleteMulti(tt.args.pfx, tt.args.pfxInv, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.deleteMulti() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVals, tt.wantVals) {
				t.Errorf("client.deleteMulti() = %v, want %v", gotVals, tt.wantVals)
			}
		})
	}
}
