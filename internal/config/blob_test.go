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
	"go.uber.org/goleak"
)

func TestBlobStorageType_String(t *testing.T) {
	type want struct {
		want string
	}
	type test struct {
		name       string
		bst        BlobStorageType
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return s3 when the bst is S3",
			bst:  S3,
			want: want{
				want: "s3",
			},
		},
		{
			name: "return unknown when the bst is empty",
			want: want{
				want: "unknown",
			},
		},
		{
			name: "return unknown when the bst is invalid storage type",
			bst:  BlobStorageType(100),
			want: want{
				want: "unknown",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := test.bst.String()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestAtoBST(t *testing.T) {
	type args struct {
		bst string
	}
	type want struct {
		want BlobStorageType
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, BlobStorageType) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got BlobStorageType) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return S3 when the bst is s3",
			args: args{
				bst: "s3",
			},
			want: want{
				want: S3,
			},
		},
		{
			name: "return S3 when the bst is S3",
			args: args{
				bst: "S3",
			},
			want: want{
				want: S3,
			},
		},
		{
			name: "return 0 when the bst is empty",
			want: want{
				want: 0,
			},
		},
		{
			name: "return 0 when the bst is invalid storage type",
			args: args{
				bst: "storage",
			},
			want: want{
				want: 0,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := AtoBST(test.args.bst)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestBlob_Bind(t *testing.T) {
	type fields struct {
		StorageType string
		Bucket      string
		S3          *S3Config
	}
	type want struct {
		want *Blob
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Blob) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *Blob) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return Blob when the bind successes and the S3Config is nil",
				fields: fields{
					StorageType: "s3",
					Bucket:      "test.vald",
				},
				want: want{
					want: &Blob{
						StorageType: "s3",
						Bucket:      "test.vald",
						S3:          new(S3Config),
					},
				},
			}
		}(),
		func() test {
			s3 := new(S3Config)
			return test{
				name: "return Blob when the bind successes and the S3Config is not nil",
				fields: fields{
					StorageType: "s3",
					Bucket:      "test.vald",
					S3:          s3,
				},
				want: want{
					want: &Blob{
						StorageType: "s3",
						Bucket:      "test.vald",
						S3:          s3,
					},
				},
			}
		}(),
		func() test {
			m := map[string]string{
				"STORAGE_TYPE": "s3",
				"BUCKET":       "test.vald",
			}

			return test{
				name: "return Blob when the bind successes and the data is loaded from the environment variable",
				fields: fields{
					StorageType: "_STORAGE_TYPE_",
					Bucket:      "_BUCKET_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()

					for key, val := range m {
						if err := os.Setenv(key, val); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()

					for key := range m {
						if err := os.Unsetenv(key); err != nil {
							t.Fatal(err)
						}
					}
				},
				want: want{
					want: &Blob{
						StorageType: "s3",
						Bucket:      "test.vald",
						S3:          new(S3Config),
					},
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			b := &Blob{
				StorageType: test.fields.StorageType,
				Bucket:      test.fields.Bucket,
				S3:          test.fields.S3,
			}

			got := b.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestS3Config_Bind(t *testing.T) {
	type fields struct {
		Endpoint                   string
		Region                     string
		AccessKey                  string
		SecretAccessKey            string
		Token                      string
		MaxRetries                 int
		ForcePathStyle             bool
		UseAccelerate              bool
		UseARNRegion               bool
		UseDualStack               bool
		EnableSSL                  bool
		EnableParamValidation      bool
		Enable100Continue          bool
		EnableContentMD5Validation bool
		EnableEndpointDiscovery    bool
		EnableEndpointHostPrefix   bool
		MaxPartSize                string
		MaxChunkSize               string
	}
	type want struct {
		want *S3Config
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *S3Config) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *S3Config) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return S3Config when the bind successes",
				fields: fields{
					Endpoint:                   "https://test.us-west-2.amazonaws.com",
					Region:                     "us-west-2",
					AccessKey:                  "access_key",
					SecretAccessKey:            "secret_access_key",
					Token:                      "token",
					MaxRetries:                 0,
					ForcePathStyle:             false,
					UseAccelerate:              false,
					UseARNRegion:               false,
					UseDualStack:               false,
					EnableSSL:                  false,
					EnableParamValidation:      false,
					Enable100Continue:          false,
					EnableContentMD5Validation: false,
					EnableEndpointDiscovery:    false,
					EnableEndpointHostPrefix:   false,
					MaxPartSize:                "32mb",
					MaxChunkSize:               "42mb",
				},
				want: want{
					want: &S3Config{
						Endpoint:        "https://test.us-west-2.amazonaws.com",
						Region:          "us-west-2",
						AccessKey:       "access_key",
						SecretAccessKey: "secret_access_key",
						Token:           "token",
						MaxPartSize:     "32mb",
						MaxChunkSize:    "42mb",
					},
				},
			}
		}(),
		func() test {
			m := map[string]string{
				"ENDPOINT":          "https://test.us-west-2.amazonaws.com",
				"REGION":            "us-west-2",
				"ACCESS_KEY":        "access_key",
				"SECRET_ACCESS_KEY": "secret_access_key",
				"TOKEN":             "token",
				"MAX_PART_SIZE":     "32mb",
				"MAX_CHUNK_SIZE":    "42mb",
			}
			return test{
				name: "return S3Config when the bind successes and the data is loaded from the environment variable",
				fields: fields{
					Endpoint:                   "_ENDPOINT_",
					Region:                     "_REGION_",
					AccessKey:                  "_ACCESS_KEY_",
					SecretAccessKey:            "_SECRET_ACCESS_KEY_",
					Token:                      "_TOKEN_",
					MaxRetries:                 0,
					ForcePathStyle:             false,
					UseAccelerate:              false,
					UseARNRegion:               false,
					UseDualStack:               false,
					EnableSSL:                  false,
					EnableParamValidation:      false,
					Enable100Continue:          false,
					EnableContentMD5Validation: false,
					EnableEndpointDiscovery:    false,
					EnableEndpointHostPrefix:   false,
					MaxPartSize:                "_MAX_PART_SIZE_",
					MaxChunkSize:               "_MAX_CHUNK_SIZE_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()

					for key, val := range m {
						if err := os.Setenv(key, val); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()

					for key := range m {
						if err := os.Unsetenv(key); err != nil {
							t.Fatal(err)
						}
					}
				},
				want: want{
					want: &S3Config{
						Endpoint:        "https://test.us-west-2.amazonaws.com",
						Region:          "us-west-2",
						AccessKey:       "access_key",
						SecretAccessKey: "secret_access_key",
						Token:           "token",
						MaxPartSize:     "32mb",
						MaxChunkSize:    "42mb",
					},
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			s := &S3Config{
				Endpoint:                   test.fields.Endpoint,
				Region:                     test.fields.Region,
				AccessKey:                  test.fields.AccessKey,
				SecretAccessKey:            test.fields.SecretAccessKey,
				Token:                      test.fields.Token,
				MaxRetries:                 test.fields.MaxRetries,
				ForcePathStyle:             test.fields.ForcePathStyle,
				UseAccelerate:              test.fields.UseAccelerate,
				UseARNRegion:               test.fields.UseARNRegion,
				UseDualStack:               test.fields.UseDualStack,
				EnableSSL:                  test.fields.EnableSSL,
				EnableParamValidation:      test.fields.EnableParamValidation,
				Enable100Continue:          test.fields.Enable100Continue,
				EnableContentMD5Validation: test.fields.EnableContentMD5Validation,
				EnableEndpointDiscovery:    test.fields.EnableEndpointDiscovery,
				EnableEndpointHostPrefix:   test.fields.EnableEndpointHostPrefix,
				MaxPartSize:                test.fields.MaxPartSize,
				MaxChunkSize:               test.fields.MaxChunkSize,
			}

			got := s.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
