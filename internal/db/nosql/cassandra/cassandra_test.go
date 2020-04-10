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

package cassandra

import (
	"context"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		want    Cassandra
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.opts...)
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

func Test_client_Open(t *testing.T) {
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
			if err := tt.c.Open(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("client.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_Close(t *testing.T) {
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
			if err := tt.c.Close(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("client.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_Query(t *testing.T) {
	type args struct {
		stmt  string
		names []string
	}
	tests := []struct {
		name string
		c    *client
		args args
		want *Queryx
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Query(tt.args.stmt, tt.args.names); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("client.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelect(t *testing.T) {
	type args struct {
		table   string
		columns []string
		cmps    []Cmp
	}
	tests := []struct {
		name      string
		args      args
		wantStmt  string
		wantNames []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStmt, gotNames := Select(tt.args.table, tt.args.columns, tt.args.cmps...)
			if gotStmt != tt.wantStmt {
				t.Errorf("Select() gotStmt = %v, want %v", gotStmt, tt.wantStmt)
			}
			if !reflect.DeepEqual(gotNames, tt.wantNames) {
				t.Errorf("Select() gotNames = %v, want %v", gotNames, tt.wantNames)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		table string
		cmps  []Cmp
	}
	tests := []struct {
		name string
		args args
		want *DeleteBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Delete(tt.args.table, tt.args.cmps...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	type args struct {
		table   string
		columns []string
	}
	tests := []struct {
		name string
		args args
		want *InsertBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Insert(tt.args.table, tt.args.columns...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		table string
	}
	tests := []struct {
		name string
		args args
		want *UpdateBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Update(tt.args.table); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBatch(t *testing.T) {
	tests := []struct {
		name string
		want *BatchBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Batch(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Batch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEq(t *testing.T) {
	type args struct {
		column string
	}
	tests := []struct {
		name string
		args args
		want Cmp
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Eq(tt.args.column); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Eq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIn(t *testing.T) {
	type args struct {
		column string
	}
	tests := []struct {
		name string
		args args
		want Cmp
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := In(tt.args.column); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("In() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		column string
	}
	tests := []struct {
		name string
		args args
		want Cmp
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.column); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrapErrorWithKeys(t *testing.T) {
	type args struct {
		err  error
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapErrorWithKeys(tt.args.err, tt.args.keys...); (err != nil) != tt.wantErr {
				t.Errorf("WrapErrorWithKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
