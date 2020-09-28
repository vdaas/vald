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
	"net/http"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/comparator"
	"go.uber.org/goleak"
)

var (
	atop = func(s string) *string {
		return &s
	}
	btop = func(b bool) *bool {
		return &b
	}
	itop = func(i int) *int {
		return &i
	}

	handlerListComparator = func(x, y request.HandlerList) bool {
		return reflect.DeepEqual(x, y)
	}
	stringPointerComparator = func(x, y *string) bool {
		return reflect.DeepEqual(x, y)
	}

	awsConfigComparatorOpts = []comparator.Option{
		comparator.AllowUnexported(aws.Config{}),
	}
	awsConfigComparator = func(x, y aws.Config) bool {
		return comparator.Diff(x, y, awsConfigComparatorOpts...) == ""
	}

	sessionComparatorOpts = []comparator.Option{
		comparator.IgnoreUnexported(session.Session{}),
		comparator.Comparer(handlerListComparator),
		comparator.Comparer(stringPointerComparator),

		comparator.IgnoreFields(aws.Config{}, "EndpointResolver", "Credentials", "Logger"),
		comparator.IgnoreTypes(request.Handlers{}),
	}
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

		if diff := comparator.Diff(w.want, got, sessionComparatorOpts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "set endpoint success",
			fields: fields{
				endpoint: "127.0.0.1",
				region:   "jp",
			},
			want: want{
				want: &session.Session{
					Config: &aws.Config{
						Endpoint:                      atop("127.0.0.1"),
						Region:                        atop("jp"),
						DisableSSL:                    btop(true),
						HTTPClient:                    &http.Client{},
						LogLevel:                      aws.LogLevel(aws.LogLevelType(uint(0))),
						MaxRetries:                    itop(0),
						DisableParamValidation:        btop(true),
						S3ForcePathStyle:              btop(false),
						S3Disable100Continue:          btop(true),
						S3UseAccelerate:               btop(false),
						S3DisableContentMD5Validation: btop(true),
						S3UseARNRegion:                btop(false),
						UseDualStack:                  btop(false),
						EnableEndpointDiscovery:       btop(false),
						DisableEndpointHostPrefix:     btop(true),
						STSRegionalEndpoint:           endpoints.LegacySTSEndpoint,
						S3UsEast1RegionalEndpoint:     endpoints.LegacyS3UsEast1Endpoint,
					},
				},
			},
		},
		{
			name: "set cred success",
			fields: fields{
				accessKey:       "abc",
				secretAccessKey: "def",
			},
			want: want{
				want: &session.Session{
					Config: &aws.Config{
						Region: atop(""),
						Credentials: func() *credentials.Credentials {
							creds := credentials.NewStaticCredentials(
								"abc",
								"def",
								"",
							)
							_, _ = creds.Get()
							return creds
						}(),
						DisableSSL:                    btop(true),
						HTTPClient:                    &http.Client{},
						LogLevel:                      aws.LogLevel(aws.LogLevelType(uint(0))),
						MaxRetries:                    itop(0),
						DisableParamValidation:        btop(true),
						S3ForcePathStyle:              btop(false),
						S3Disable100Continue:          btop(true),
						S3UseAccelerate:               btop(false),
						S3DisableContentMD5Validation: btop(true),
						S3UseARNRegion:                btop(false),
						UseDualStack:                  btop(false),
						EnableEndpointDiscovery:       btop(false),
						DisableEndpointHostPrefix:     btop(true),
						STSRegionalEndpoint:           endpoints.LegacySTSEndpoint,
						S3UsEast1RegionalEndpoint:     endpoints.LegacyS3UsEast1Endpoint,
					},
				},
			},
		},
		{
			name: "set maxretries success",
			fields: fields{
				maxRetries: 999,
			},
			want: want{
				want: &session.Session{
					Config: &aws.Config{
						Region:                        atop(""),
						Credentials:                   nil,
						DisableSSL:                    btop(true),
						HTTPClient:                    &http.Client{},
						LogLevel:                      aws.LogLevel(aws.LogLevelType(uint(0))),
						MaxRetries:                    itop(999),
						DisableParamValidation:        btop(true),
						S3ForcePathStyle:              btop(false),
						S3Disable100Continue:          btop(true),
						S3UseAccelerate:               btop(false),
						S3DisableContentMD5Validation: btop(true),
						S3UseARNRegion:                btop(false),
						UseDualStack:                  btop(false),
						EnableEndpointDiscovery:       btop(false),
						DisableEndpointHostPrefix:     btop(true),
						STSRegionalEndpoint:           endpoints.LegacySTSEndpoint,
						S3UsEast1RegionalEndpoint:     endpoints.LegacyS3UsEast1Endpoint,
					},
				},
			},
		},
		{
			name: "set forcepathstyle success",
			fields: fields{
				forcePathStyle: true,
			},
			want: want{
				want: &session.Session{
					Config: &aws.Config{
						Region:                        atop(""),
						Credentials:                   nil,
						DisableSSL:                    btop(true),
						HTTPClient:                    &http.Client{},
						LogLevel:                      aws.LogLevel(aws.LogLevelType(uint(0))),
						MaxRetries:                    itop(0),
						DisableParamValidation:        btop(true),
						S3ForcePathStyle:              btop(true),
						S3Disable100Continue:          btop(true),
						S3UseAccelerate:               btop(false),
						S3DisableContentMD5Validation: btop(true),
						S3UseARNRegion:                btop(false),
						UseDualStack:                  btop(false),
						EnableEndpointDiscovery:       btop(false),
						DisableEndpointHostPrefix:     btop(true),
						STSRegionalEndpoint:           endpoints.LegacySTSEndpoint,
						S3UsEast1RegionalEndpoint:     endpoints.LegacyS3UsEast1Endpoint,
					},
				},
			},
		},
		{
			name: "set useAccelerate success",
			fields: fields{
				useAccelerate: true,
			},
			want: want{
				want: &session.Session{
					Config: &aws.Config{
						Region:                        atop(""),
						Credentials:                   nil,
						DisableSSL:                    btop(true),
						HTTPClient:                    &http.Client{},
						LogLevel:                      aws.LogLevel(aws.LogLevelType(uint(0))),
						MaxRetries:                    itop(0),
						DisableParamValidation:        btop(true),
						S3ForcePathStyle:              btop(false),
						S3Disable100Continue:          btop(true),
						S3UseAccelerate:               btop(true),
						S3DisableContentMD5Validation: btop(true),
						S3UseARNRegion:                btop(false),
						UseDualStack:                  btop(false),
						EnableEndpointDiscovery:       btop(false),
						DisableEndpointHostPrefix:     btop(true),
						STSRegionalEndpoint:           endpoints.LegacySTSEndpoint,
						S3UsEast1RegionalEndpoint:     endpoints.LegacyS3UsEast1Endpoint,
					},
				},
			},
		},
		{
			name: "set useARNRegion success",
			fields: fields{
				useARNRegion: true,
			},
			want: want{
				want: &session.Session{
					Config: &aws.Config{
						Region:                        atop(""),
						Credentials:                   nil,
						DisableSSL:                    btop(true),
						HTTPClient:                    &http.Client{},
						LogLevel:                      aws.LogLevel(aws.LogLevelType(uint(0))),
						MaxRetries:                    itop(0),
						DisableParamValidation:        btop(true),
						S3ForcePathStyle:              btop(false),
						S3Disable100Continue:          btop(true),
						S3UseAccelerate:               btop(false),
						S3DisableContentMD5Validation: btop(true),
						S3UseARNRegion:                btop(true),
						UseDualStack:                  btop(false),
						EnableEndpointDiscovery:       btop(false),
						DisableEndpointHostPrefix:     btop(true),
						STSRegionalEndpoint:           endpoints.LegacySTSEndpoint,
						S3UsEast1RegionalEndpoint:     endpoints.LegacyS3UsEast1Endpoint,
					},
				},
			},
		},
		{
			name: "set useDualStack success",
			fields: fields{
				useDualStack: true,
			},
			want: want{
				want: &session.Session{
					Config: &aws.Config{
						Region:                        atop(""),
						Credentials:                   nil,
						DisableSSL:                    btop(true),
						HTTPClient:                    &http.Client{},
						LogLevel:                      aws.LogLevel(aws.LogLevelType(uint(0))),
						MaxRetries:                    itop(0),
						DisableParamValidation:        btop(true),
						S3ForcePathStyle:              btop(false),
						S3Disable100Continue:          btop(true),
						S3UseAccelerate:               btop(false),
						S3DisableContentMD5Validation: btop(true),
						S3UseARNRegion:                btop(false),
						UseDualStack:                  btop(true),
						EnableEndpointDiscovery:       btop(false),
						DisableEndpointHostPrefix:     btop(true),
						STSRegionalEndpoint:           endpoints.LegacySTSEndpoint,
						S3UsEast1RegionalEndpoint:     endpoints.LegacyS3UsEast1Endpoint,
					},
				},
			},
		},
		{
			name: "set enableEndpointDiscovery success",
			fields: fields{
				enableEndpointDiscovery: true,
			},
			want: want{
				want: &session.Session{
					Config: &aws.Config{
						Region:                        atop(""),
						Credentials:                   nil,
						DisableSSL:                    btop(true),
						HTTPClient:                    &http.Client{},
						LogLevel:                      aws.LogLevel(aws.LogLevelType(uint(0))),
						MaxRetries:                    itop(0),
						DisableParamValidation:        btop(true),
						S3ForcePathStyle:              btop(false),
						S3Disable100Continue:          btop(true),
						S3UseAccelerate:               btop(false),
						S3DisableContentMD5Validation: btop(true),
						S3UseARNRegion:                btop(false),
						UseDualStack:                  btop(false),
						EnableEndpointDiscovery:       btop(true),
						DisableEndpointHostPrefix:     btop(true),
						STSRegionalEndpoint:           endpoints.LegacySTSEndpoint,
						S3UsEast1RegionalEndpoint:     endpoints.LegacyS3UsEast1Endpoint,
					},
				},
			},
		},
		{
			name: "set enableSSL success",
			fields: fields{
				enableSSL: true,
			},
			want: want{
				want: &session.Session{
					Config: &aws.Config{
						Region:      atop(""),
						Credentials: nil,
						//DisableSSL:                    btop(false),
						HTTPClient:                    &http.Client{},
						LogLevel:                      aws.LogLevel(aws.LogLevelType(uint(0))),
						MaxRetries:                    itop(0),
						DisableParamValidation:        btop(true),
						S3ForcePathStyle:              btop(false),
						S3Disable100Continue:          btop(true),
						S3UseAccelerate:               btop(false),
						S3DisableContentMD5Validation: btop(true),
						S3UseARNRegion:                btop(false),
						UseDualStack:                  btop(false),
						EnableEndpointDiscovery:       btop(false),
						DisableEndpointHostPrefix:     btop(true),
						STSRegionalEndpoint:           endpoints.LegacySTSEndpoint,
						S3UsEast1RegionalEndpoint:     endpoints.LegacyS3UsEast1Endpoint,
					},
				},
			},
		},
		{
			name: "set EnableParamValdiation success",
			fields: fields{
				enableParamValidation: true,
			},
			want: want{
				want: &session.Session{
					Config: &aws.Config{
						Region:                        atop(""),
						Credentials:                   nil,
						DisableSSL:                    btop(true),
						HTTPClient:                    &http.Client{},
						LogLevel:                      aws.LogLevel(aws.LogLevelType(uint(0))),
						MaxRetries:                    itop(0),
						S3ForcePathStyle:              btop(false),
						S3Disable100Continue:          btop(true),
						S3UseAccelerate:               btop(false),
						S3DisableContentMD5Validation: btop(true),
						S3UseARNRegion:                btop(false),
						UseDualStack:                  btop(false),
						EnableEndpointDiscovery:       btop(false),
						DisableEndpointHostPrefix:     btop(true),
						STSRegionalEndpoint:           endpoints.LegacySTSEndpoint,
						S3UsEast1RegionalEndpoint:     endpoints.LegacyS3UsEast1Endpoint,
						//DisableParamValidation:        btop(true),
					},
				},
			},
		},
		{
			name: "set Enable100Conitnue success",
			fields: fields{
				enable100Continue: true,
			},
			want: want{
				want: &session.Session{
					Config: &aws.Config{
						Region:           atop(""),
						Credentials:      nil,
						DisableSSL:       btop(true),
						HTTPClient:       &http.Client{},
						LogLevel:         aws.LogLevel(aws.LogLevelType(uint(0))),
						MaxRetries:       itop(0),
						S3ForcePathStyle: btop(false),
						//S3Disable100Continue:          btop(true),
						S3UseAccelerate:               btop(false),
						S3DisableContentMD5Validation: btop(true),
						S3UseARNRegion:                btop(false),
						UseDualStack:                  btop(false),
						EnableEndpointDiscovery:       btop(false),
						DisableEndpointHostPrefix:     btop(true),
						STSRegionalEndpoint:           endpoints.LegacySTSEndpoint,
						S3UsEast1RegionalEndpoint:     endpoints.LegacyS3UsEast1Endpoint,
						DisableParamValidation:        btop(true),
					},
				},
			},
		},
		{
			name: "set ContentMD5Validation success",
			fields: fields{
				enableContentMD5Validation: true,
			},
			want: want{
				want: &session.Session{
					Config: &aws.Config{
						Region:               atop(""),
						Credentials:          nil,
						DisableSSL:           btop(true),
						HTTPClient:           &http.Client{},
						LogLevel:             aws.LogLevel(aws.LogLevelType(uint(0))),
						MaxRetries:           itop(0),
						S3ForcePathStyle:     btop(false),
						S3Disable100Continue: btop(true),
						S3UseAccelerate:      btop(false),
						// S3DisableContentMD5Validation: btop(true),
						S3UseARNRegion:            btop(false),
						UseDualStack:              btop(false),
						EnableEndpointDiscovery:   btop(false),
						DisableEndpointHostPrefix: btop(true),
						STSRegionalEndpoint:       endpoints.LegacySTSEndpoint,
						S3UsEast1RegionalEndpoint: endpoints.LegacyS3UsEast1Endpoint,
						DisableParamValidation:    btop(true),
					},
				},
			},
		},
		{
			name: "set endpointHostPrefix success",
			fields: fields{
				enableEndpointHostPrefix: true,
			},
			want: want{
				want: &session.Session{
					Config: &aws.Config{
						Region:                        atop(""),
						Credentials:                   nil,
						DisableSSL:                    btop(true),
						HTTPClient:                    &http.Client{},
						LogLevel:                      aws.LogLevel(aws.LogLevelType(uint(0))),
						MaxRetries:                    itop(0),
						S3ForcePathStyle:              btop(false),
						S3Disable100Continue:          btop(true),
						S3UseAccelerate:               btop(false),
						S3DisableContentMD5Validation: btop(true),
						S3UseARNRegion:                btop(false),
						UseDualStack:                  btop(false),
						EnableEndpointDiscovery:       btop(false),
						//DisableEndpointHostPrefix:     btop(true),
						STSRegionalEndpoint:       endpoints.LegacySTSEndpoint,
						S3UsEast1RegionalEndpoint: endpoints.LegacyS3UsEast1Endpoint,
						DisableParamValidation:    btop(true),
					},
				},
			},
		},
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
