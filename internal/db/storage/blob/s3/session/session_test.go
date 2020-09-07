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

package session

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want Session
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Session) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Session) error {
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
		           opts: nil,
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
		           opts: nil,
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

			got := New(test.args.opts...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_sess_Session(t *testing.T) {
	type fields struct {
		endpoint                   string
		region                     string
		accessKey                  string
		secretAccessKey            string
		token                      string
		maxRetries                 int
		forcePathStyle             bool
		useAccelerate              bool
		useARNRegion               bool
		useDualStack               bool
		enableSSL                  bool
		enableParamValidation      bool
		enable100Continue          bool
		enableContentMD5Validation bool
		enableEndpointDiscovery    bool
		enableEndpointHostPrefix   bool
	}
	type want struct {
		want *session.Session
		err  error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *session.Session, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *session.Session, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
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
		           endpoint: "",
		           region: "",
		           accessKey: "",
		           secretAccessKey: "",
		           token: "",
		           maxRetries: 0,
		           forcePathStyle: false,
		           useAccelerate: false,
		           useARNRegion: false,
		           useDualStack: false,
		           enableSSL: false,
		           enableParamValidation: false,
		           enable100Continue: false,
		           enableContentMD5Validation: false,
		           enableEndpointDiscovery: false,
		           enableEndpointHostPrefix: false,
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
		           endpoint: "",
		           region: "",
		           accessKey: "",
		           secretAccessKey: "",
		           token: "",
		           maxRetries: 0,
		           forcePathStyle: false,
		           useAccelerate: false,
		           useARNRegion: false,
		           useDualStack: false,
		           enableSSL: false,
		           enableParamValidation: false,
		           enable100Continue: false,
		           enableContentMD5Validation: false,
		           enableEndpointDiscovery: false,
		           enableEndpointHostPrefix: false,
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
			s := &sess{
				endpoint:                   test.fields.endpoint,
				region:                     test.fields.region,
				accessKey:                  test.fields.accessKey,
				secretAccessKey:            test.fields.secretAccessKey,
				token:                      test.fields.token,
				maxRetries:                 test.fields.maxRetries,
				forcePathStyle:             test.fields.forcePathStyle,
				useAccelerate:              test.fields.useAccelerate,
				useARNRegion:               test.fields.useARNRegion,
				useDualStack:               test.fields.useDualStack,
				enableSSL:                  test.fields.enableSSL,
				enableParamValidation:      test.fields.enableParamValidation,
				enable100Continue:          test.fields.enable100Continue,
				enableContentMD5Validation: test.fields.enableContentMD5Validation,
				enableEndpointDiscovery:    test.fields.enableEndpointDiscovery,
				enableEndpointHostPrefix:   test.fields.enableEndpointHostPrefix,
			}

			got, err := s.Session()
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
