// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package urlopener

import (
	"net/http"
	"testing"

	"cloud.google.com/go/storage"
)

const (
	testServiceAccountJSON = `{"type":"service_account","client_email":"sa@example.iam.gserviceaccount.com","private_key":"test-private-key","private_key_id":"kid","token_uri":"https://oauth2.googleapis.com/token"}`
	testAuthorizedUserJSON = `{"type":"authorized_user","client_id":"cid","client_secret":"secret","refresh_token":"refresh"}`
)

func Test_signingOptions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantID  string
		wantKey string
		json    []byte
		wantErr bool
	}{
		{name: "nil JSON (metadata server / Workload Identity)", json: nil},
		{name: "authorized_user has no signing key", json: []byte(testAuthorizedUserJSON)},
		{
			name:    "service_account fills signing fields",
			json:    []byte(testServiceAccountJSON),
			wantID:  "sa@example.iam.gserviceaccount.com",
			wantKey: "test-private-key",
		},
		{name: "malformed JSON returns error", json: []byte(`{`), wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := signingOptions(tc.json, storage.ScopeReadWrite)
			if (err != nil) != tc.wantErr {
				t.Fatalf("signingOptions() error = %v, wantErr %v", err, tc.wantErr)
			}
			if got.GoogleAccessID != tc.wantID {
				t.Errorf("GoogleAccessID = %q, want %q", got.GoogleAccessID, tc.wantID)
			}
			if string(got.PrivateKey) != tc.wantKey {
				t.Errorf("PrivateKey = %q, want %q", got.PrivateKey, tc.wantKey)
			}
		})
	}
}

func Test_urlOpener_URLOpener_credentialTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		json   string
		wantID string
	}{
		{name: "authorized_user credentials are accepted", json: testAuthorizedUserJSON},
		{
			name:   "service_account credentials keep signing fields",
			json:   testServiceAccountJSON,
			wantID: "sa@example.iam.gserviceaccount.com",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			uo, err := New(WithHTTPClient(&http.Client{}), WithCredentialsJSON(tc.json))
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			got, err := uo.URLOpener(t.Context())
			if err != nil {
				t.Fatalf("URLOpener() error = %v", err)
			}
			if got.Client == nil {
				t.Fatal("URLOpener().Client is nil")
			}
			if got.Options.GoogleAccessID != tc.wantID {
				t.Errorf("GoogleAccessID = %q, want %q", got.Options.GoogleAccessID, tc.wantID)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want URLOpener
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, URLOpener, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got URLOpener, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           opts:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           opts:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_urlOpener_URLOpener(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client              *http.Client
// 		credentialsFilePath string
// 		credentialsJSON     string
// 	}
// 	type want struct {
// 		wantGuo *gcsblob.URLOpener
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *gcsblob.URLOpener, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotGuo *gcsblob.URLOpener, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotGuo, w.wantGuo) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotGuo, w.wantGuo)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           credentialsFilePath:"",
// 		           credentialsJSON:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           credentialsFilePath:"",
// 		           credentialsJSON:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			uo := &urlOpener{
// 				client:              test.fields.client,
// 				credentialsFilePath: test.fields.credentialsFilePath,
// 				credentialsJSON:     test.fields.credentialsJSON,
// 			}
//
// 			gotGuo, err := uo.URLOpener(test.args.ctx)
// 			if err := checkFunc(test.want, gotGuo, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
