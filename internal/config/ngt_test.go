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

	"github.com/vdaas/vald/internal/errors"
)

func TestNGT_Bind(t *testing.T) {
	type fields struct {
		IndexPath              string
		Dimension              int
		BulkInsertChunkSize    int
		DistanceType           string
		ObjectType             string
		CreationEdgeSize       int
		SearchEdgeSize         int
		AutoIndexDurationLimit string
		AutoIndexCheckDuration string
		AutoIndexLength        int
		EnableInMemoryMode     bool
	}
	type want struct {
		want *NGT
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *NGT) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *NGT) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           IndexPath: "",
		           Dimension: 0,
		           BulkInsertChunkSize: 0,
		           DistanceType: "",
		           ObjectType: "",
		           CreationEdgeSize: 0,
		           SearchEdgeSize: 0,
		           AutoIndexDurationLimit: "",
		           AutoIndexCheckDuration: "",
		           AutoIndexLength: 0,
		           EnableInMemoryMode: false,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           IndexPath: "",
		           Dimension: 0,
		           BulkInsertChunkSize: 0,
		           DistanceType: "",
		           ObjectType: "",
		           CreationEdgeSize: 0,
		           SearchEdgeSize: 0,
		           AutoIndexDurationLimit: "",
		           AutoIndexCheckDuration: "",
		           AutoIndexLength: 0,
		           EnableInMemoryMode: false,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			n := &NGT{
				IndexPath:              test.fields.IndexPath,
				Dimension:              test.fields.Dimension,
				BulkInsertChunkSize:    test.fields.BulkInsertChunkSize,
				DistanceType:           test.fields.DistanceType,
				ObjectType:             test.fields.ObjectType,
				CreationEdgeSize:       test.fields.CreationEdgeSize,
				SearchEdgeSize:         test.fields.SearchEdgeSize,
				AutoIndexDurationLimit: test.fields.AutoIndexDurationLimit,
				AutoIndexCheckDuration: test.fields.AutoIndexCheckDuration,
				AutoIndexLength:        test.fields.AutoIndexLength,
				EnableInMemoryMode:     test.fields.EnableInMemoryMode,
			}

			got := n.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
