//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	"os"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestNGT_Bind(t *testing.T) {
	type fields struct {
		IndexPath               string
		Dimension               int
		BulkInsertChunkSize     int
		DistanceType            string
		ObjectType              string
		CreationEdgeSize        int
		SearchEdgeSize          int
		AutoIndexDurationLimit  string
		AutoIndexCheckDuration  string
		AutoSaveIndexDuration   string
		AutoIndexLength         int
		InitialDelayMaxDuration string
		EnableInMemoryMode      bool
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return NGT when all fields contain no prefix/suffix symbol",
			fields: fields{
				IndexPath:               "config/ngt",
				Dimension:               1000,
				BulkInsertChunkSize:     100,
				DistanceType:            "l2",
				ObjectType:              "float",
				CreationEdgeSize:        3,
				SearchEdgeSize:          5,
				AutoIndexDurationLimit:  "1h",
				AutoIndexCheckDuration:  "30m",
				AutoSaveIndexDuration:   "30m",
				AutoIndexLength:         100,
				InitialDelayMaxDuration: "1h",
				EnableInMemoryMode:      false,
			},
			want: want{
				want: &NGT{
					IndexPath:               "config/ngt",
					Dimension:               1000,
					BulkInsertChunkSize:     100,
					DistanceType:            "l2",
					ObjectType:              "float",
					CreationEdgeSize:        3,
					SearchEdgeSize:          5,
					AutoIndexDurationLimit:  "1h",
					AutoIndexCheckDuration:  "30m",
					AutoSaveIndexDuration:   "30m",
					AutoIndexLength:         100,
					InitialDelayMaxDuration: "1h",
					EnableInMemoryMode:      false,
				},
			},
		},
		{
			name: "return NGT with environment variable when it contains `_` as prefix and suffix",
			fields: fields{
				IndexPath:               "_indexPath_",
				Dimension:               1000,
				BulkInsertChunkSize:     100,
				DistanceType:            "_distanceType_",
				ObjectType:              "_objectType_",
				CreationEdgeSize:        3,
				SearchEdgeSize:          5,
				AutoIndexDurationLimit:  "_autoIndexDurationLimit_",
				AutoIndexCheckDuration:  "_autoIndexCheckDuration_",
				AutoSaveIndexDuration:   "_autoSaveIndexDuration_",
				AutoIndexLength:         100,
				InitialDelayMaxDuration: "_initialDelayMaxDuration_",
				EnableInMemoryMode:      false,
			},
			beforeFunc: func() {
				_ = os.Setenv("indexPath", "config/ngt")
				_ = os.Setenv("distanceType", "l2")
				_ = os.Setenv("objectType", "float")
				_ = os.Setenv("autoIndexDurationLimit", "1h")
				_ = os.Setenv("autoIndexCheckDuration", "30m")
				_ = os.Setenv("autoSaveIndexDuration", "30m")
				_ = os.Setenv("initialDelayMaxDuration", "1h")
			},
			afterFunc: func() {
				_ = os.Unsetenv("indexPath")
				_ = os.Unsetenv("distanceType")
				_ = os.Unsetenv("objectType")
				_ = os.Unsetenv("autoIndexDurationLimit")
				_ = os.Unsetenv("autoIndexCheckDuration")
				_ = os.Unsetenv("autoSaveIndexDuration")
				_ = os.Unsetenv("initialDelayMaxDuration")
			},
			want: want{
				want: &NGT{
					IndexPath:               "config/ngt",
					Dimension:               1000,
					BulkInsertChunkSize:     100,
					DistanceType:            "l2",
					ObjectType:              "float",
					CreationEdgeSize:        3,
					SearchEdgeSize:          5,
					AutoIndexDurationLimit:  "1h",
					AutoIndexCheckDuration:  "30m",
					AutoSaveIndexDuration:   "30m",
					AutoIndexLength:         100,
					InitialDelayMaxDuration: "1h",
					EnableInMemoryMode:      false,
				},
			},
		},
		{
			name: "returns NGT when all fields are empty",
			want: want{
				want: new(NGT),
			},
		},
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
				IndexPath:               test.fields.IndexPath,
				Dimension:               test.fields.Dimension,
				BulkInsertChunkSize:     test.fields.BulkInsertChunkSize,
				DistanceType:            test.fields.DistanceType,
				ObjectType:              test.fields.ObjectType,
				CreationEdgeSize:        test.fields.CreationEdgeSize,
				SearchEdgeSize:          test.fields.SearchEdgeSize,
				AutoIndexDurationLimit:  test.fields.AutoIndexDurationLimit,
				AutoIndexCheckDuration:  test.fields.AutoIndexCheckDuration,
				AutoSaveIndexDuration:   test.fields.AutoSaveIndexDuration,
				AutoIndexLength:         test.fields.AutoIndexLength,
				InitialDelayMaxDuration: test.fields.InitialDelayMaxDuration,
				EnableInMemoryMode:      test.fields.EnableInMemoryMode,
			}

			got := n.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
