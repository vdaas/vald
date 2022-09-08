//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/test/goleak"
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
			name: "return cloud_storage when the bst is CloudStorage",
			bst:  CloudStorage,
			want: want{
				want: "cloud_storage",
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

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := test.bst.String()
			if err := checkFunc(test.want, got); err != nil {
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
			name: "return CloudStorage when the bst is cloud_storage",
			args: args{
				bst: "cloud_storage",
			},
			want: want{
				want: CloudStorage,
			},
		},
		{
			name: "return CloudStorage when the bst is CLOUD_storage",
			args: args{
				bst: "CLOUD_storage",
			},
			want: want{
				want: CloudStorage,
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

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := AtoBST(test.args.bst)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestBlob_Bind(t *testing.T) {
	type fields struct {
		StorageType  string
		Bucket       string
		S3           *S3Config
		CloudStorage *CloudStorageConfig
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
				name: "return Blob when the bind successes and the S3Config and CloudStorage is nil",
				fields: fields{
					StorageType: "s3",
					Bucket:      "test.vald",
				},
				want: want{
					want: &Blob{
						StorageType:  "s3",
						Bucket:       "test.vald",
						S3:           new(S3Config),
						CloudStorage: new(CloudStorageConfig),
					},
				},
			}
		}(),
		func() test {
			s3 := &S3Config{
				Endpoint: "https://test.vald",
			}
			cloudStorage := &CloudStorageConfig{
				URL:    "gs://test.vald",
				Client: new(CloudStorageClient),
			}
			return test{
				name: "return Blob when the bind successes and the S3Config CloudStorageConfig is not nil",
				fields: fields{
					StorageType:  "s3",
					Bucket:       "test.vald",
					S3:           s3,
					CloudStorage: cloudStorage,
				},
				want: want{
					want: &Blob{
						StorageType:  "s3",
						Bucket:       "test.vald",
						S3:           s3,
						CloudStorage: cloudStorage,
					},
				},
			}
		}(),
		func() test {
			envPrefix := "BLOB_BIND_"
			m := map[string]string{
				envPrefix + "STORAGE_TYPE": "s3",
				envPrefix + "BUCKET":       "test.vald",
			}

			return test{
				name: "return Blob when the bind successes and the data is loaded from the environment variable",
				fields: fields{
					StorageType: "_" + envPrefix + "STORAGE_TYPE_",
					Bucket:      "_" + envPrefix + "BUCKET_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()

					for key, val := range m {
						t.Setenv(key, val)
					}
				},
				want: want{
					want: &Blob{
						StorageType:  "s3",
						Bucket:       "test.vald",
						S3:           new(S3Config),
						CloudStorage: new(CloudStorageConfig),
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			b := &Blob{
				StorageType:  test.fields.StorageType,
				Bucket:       test.fields.Bucket,
				S3:           test.fields.S3,
				CloudStorage: test.fields.CloudStorage,
			}

			got := b.Bind()
			if err := checkFunc(test.want, got); err != nil {
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
			envPrefix := "S3CONFIG_BIND_"
			m := map[string]string{
				envPrefix + "ENDPOINT":          "https://test.us-west-2.amazonaws.com",
				envPrefix + "REGION":            "us-west-2",
				envPrefix + "ACCESS_KEY":        "access_key",
				envPrefix + "SECRET_ACCESS_KEY": "secret_access_key",
				envPrefix + "TOKEN":             "token",
				envPrefix + "MAX_PART_SIZE":     "32mb",
				envPrefix + "MAX_CHUNK_SIZE":    "42mb",
			}
			return test{
				name: "return S3Config when the bind successes and the data is loaded from the environment variable",
				fields: fields{
					Endpoint:                   "_" + envPrefix + "ENDPOINT_",
					Region:                     "_" + envPrefix + "REGION_",
					AccessKey:                  "_" + envPrefix + "ACCESS_KEY_",
					SecretAccessKey:            "_" + envPrefix + "SECRET_ACCESS_KEY_",
					Token:                      "_" + envPrefix + "TOKEN_",
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
					MaxPartSize:                "_" + envPrefix + "MAX_PART_SIZE_",
					MaxChunkSize:               "_" + envPrefix + "MAX_CHUNK_SIZE_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()

					for key, val := range m {
						t.Setenv(key, val)
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

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
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
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestCloudStorageConfig_Bind(t *testing.T) {
	type fields struct {
		URL                     string
		Client                  *CloudStorageClient
		WriteBufferSize         int
		WriteCacheControl       string
		WriteContentDisposition string
		WriteContentEncoding    string
		WriteContentLanguage    string
		WriteContentType        string
	}
	type want struct {
		want *CloudStorageConfig
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *CloudStorageConfig) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *CloudStorageConfig) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			fields := fields{
				URL: "gs://test.vald",
				Client: &CloudStorageClient{
					CredentialsFilePath: "/var/cred",
					CredentialsJSON:     "{\"type\": \"json\"}",
				},
				WriteBufferSize:         256,
				WriteCacheControl:       "no-cache",
				WriteContentDisposition: "attachment",
				WriteContentEncoding:    "uint8",
				WriteContentLanguage:    "en-US",
				WriteContentType:        "text/plain",
			}
			return test{
				name:   "return CloudStorageConfig when the CloudStorageClient is not nil",
				fields: fields,
				want: want{
					want: &CloudStorageConfig{
						URL:                     fields.URL,
						Client:                  fields.Client,
						WriteBufferSize:         fields.WriteBufferSize,
						WriteCacheControl:       fields.WriteCacheControl,
						WriteContentDisposition: fields.WriteContentDisposition,
						WriteContentEncoding:    fields.WriteContentEncoding,
						WriteContentLanguage:    fields.WriteContentLanguage,
						WriteContentType:        fields.WriteContentType,
					},
				},
			}
		}(),
		func() test {
			fields := fields{
				URL:                     "gs://test.vald",
				WriteBufferSize:         256,
				WriteCacheControl:       "no-cache",
				WriteContentDisposition: "attachment",
				WriteContentEncoding:    "uint8",
				WriteContentLanguage:    "en-US",
				WriteContentType:        "text/plain",
			}
			return test{
				name:   "return CloudStorageConfig when the CloudStorageClient is nil",
				fields: fields,
				want: want{
					want: &CloudStorageConfig{
						URL:                     fields.URL,
						Client:                  new(CloudStorageClient),
						WriteBufferSize:         fields.WriteBufferSize,
						WriteCacheControl:       fields.WriteCacheControl,
						WriteContentDisposition: fields.WriteContentDisposition,
						WriteContentEncoding:    fields.WriteContentEncoding,
						WriteContentLanguage:    fields.WriteContentLanguage,
						WriteContentType:        fields.WriteContentType,
					},
				},
			}
		}(),
		func() test {
			envPrefix := "CLOUDSTORAGECONFIG_BIND_"
			m := map[string]string{
				envPrefix + "URL":                          "gs://test.vald",
				envPrefix + "CLIENT_CREDENTIALS_FILE_PATH": "/var/cred",
				envPrefix + "CLIENT_CREDENTIALS_JSON":      "{\"type\": \"json\"}",
				envPrefix + "WRITE_CACHE_CONTROL":          "no-cache",
				envPrefix + "WRITE_CONTENT_DISPOSITION":    "attachment",
				envPrefix + "WRITE_CONTENT_ENCODING":       "uint8",
				envPrefix + "WRITE_CONTENT_LANGUAGE":       "en-US",
				envPrefix + "WRITE_CONTENT_TYPE":           "text/plain",
			}
			return test{
				name: "return CloudStorageConfig when the data is loaded from the environment variable",
				fields: fields{
					URL: "_" + envPrefix + "URL_",
					Client: &CloudStorageClient{
						CredentialsFilePath: "_" + envPrefix + "CLIENT_CREDENTIALS_FILE_PATH_",
						CredentialsJSON:     "_" + envPrefix + "CLIENT_CREDENTIALS_JSON_",
					},
					WriteBufferSize:         256,
					WriteCacheControl:       "_" + envPrefix + "WRITE_CACHE_CONTROL_",
					WriteContentDisposition: "_" + envPrefix + "WRITE_CONTENT_DISPOSITION_",
					WriteContentEncoding:    "_" + envPrefix + "WRITE_CONTENT_ENCODING_",
					WriteContentLanguage:    "_" + envPrefix + "WRITE_CONTENT_LANGUAGE_",
					WriteContentType:        "_" + envPrefix + "WRITE_CONTENT_TYPE_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range m {
						t.Setenv(k, v)
					}
				},
				want: want{
					want: &CloudStorageConfig{
						URL: "gs://test.vald",
						Client: &CloudStorageClient{
							CredentialsFilePath: "/var/cred",
							CredentialsJSON:     "{\"type\": \"json\"}",
						},
						WriteBufferSize:         256,
						WriteCacheControl:       "no-cache",
						WriteContentDisposition: "attachment",
						WriteContentEncoding:    "uint8",
						WriteContentLanguage:    "en-US",
						WriteContentType:        "text/plain",
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			c := &CloudStorageConfig{
				URL:                     test.fields.URL,
				Client:                  test.fields.Client,
				WriteBufferSize:         test.fields.WriteBufferSize,
				WriteCacheControl:       test.fields.WriteCacheControl,
				WriteContentDisposition: test.fields.WriteContentDisposition,
				WriteContentEncoding:    test.fields.WriteContentEncoding,
				WriteContentLanguage:    test.fields.WriteContentLanguage,
				WriteContentType:        test.fields.WriteContentType,
			}

			got := c.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
