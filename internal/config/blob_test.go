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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           bst: "",
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
		           args: args {
		           bst: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *Blob) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           StorageType: "",
		           Bucket: "",
		           S3: S3Config{},
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
		           StorageType: "",
		           Bucket: "",
		           S3: S3Config{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
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
	}
	type want struct {
		want *S3Config
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *S3Config) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *S3Config) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           Endpoint: "",
		           Region: "",
		           AccessKey: "",
		           SecretAccessKey: "",
		           Token: "",
		           MaxRetries: 0,
		           ForcePathStyle: false,
		           UseAccelerate: false,
		           UseARNRegion: false,
		           UseDualStack: false,
		           EnableSSL: false,
		           EnableParamValidation: false,
		           Enable100Continue: false,
		           EnableContentMD5Validation: false,
		           EnableEndpointDiscovery: false,
		           EnableEndpointHostPrefix: false,
		           MaxPartSize: "",
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
		           Endpoint: "",
		           Region: "",
		           AccessKey: "",
		           SecretAccessKey: "",
		           Token: "",
		           MaxRetries: 0,
		           ForcePathStyle: false,
		           UseAccelerate: false,
		           UseARNRegion: false,
		           UseDualStack: false,
		           EnableSSL: false,
		           EnableParamValidation: false,
		           Enable100Continue: false,
		           EnableContentMD5Validation: false,
		           EnableEndpointDiscovery: false,
		           EnableEndpointHostPrefix: false,
		           MaxPartSize: "",
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
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
			}

			got := s.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
