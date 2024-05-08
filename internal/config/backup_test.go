//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

func TestBackupManager_Bind(t *testing.T) {
	type fields struct {
		Client *GRPCClient
	}
	type want struct {
		want *BackupManager
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *BackupManager) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *BackupManager) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			addrs := []string{
				"10.40.3.342",
				"10.40.98.17",
				"10.40.84.215",
			}
			healthcheck := "30s"
			connectionPool := &ConnectionPool{
				ResolveDNS:           true,
				EnableRebalance:      true,
				RebalanceDuration:    "5m",
				Size:                 100,
				OldConnCloseDuration: "3m",
			}
			backoffOpts := &Backoff{
				InitialDuration:  "5m",
				BackoffTimeLimit: "10m",
				MaximumDuration:  "15m",
				JitterLimit:      "3m",
				BackoffFactor:    3,
				RetryCount:       100,
				EnableErrorLog:   true,
			}
			callOpts := &CallOption{
				WaitForReady:          true,
				MaxRetryRPCBufferSize: 100,
				MaxRecvMsgSize:        1000,
				MaxSendMsgSize:        1000,
			}
			dialOpts := &DialOption{
				WriteBufferSize:             10000,
				ReadBufferSize:              10000,
				InitialWindowSize:           100,
				InitialConnectionWindowSize: 100,
				MaxMsgSize:                  1000,
				BackoffMaxDelay:             "3m",
				BackoffBaseDelay:            "1m",
				BackoffJitter:               100,
				BackoffMultiplier:           10,
				MinimumConnectionTimeout:    "5m",
				EnableBackoff:               true,
				Insecure:                    true,
				Timeout:                     "5m",
				Net:                         &Net{},
				Keepalive: &GRPCClientKeepalive{
					Time:                "100s",
					Timeout:             "300s",
					PermitWithoutStream: true,
				},
			}
			tls := &TLS{
				Enabled: true,
				Cert:    "cert",
				Key:     "key",
				CA:      "ca",
			}
			return test{
				name: "return BackupManager when the b.Client is not nil",
				fields: fields{
					Client: &GRPCClient{
						Addrs:               addrs,
						HealthCheckDuration: healthcheck,
						ConnectionPool:      connectionPool,
						Backoff:             backoffOpts,
						CallOption:          callOpts,
						DialOption:          dialOpts,
						TLS:                 tls,
					},
				},
				want: want{
					want: &BackupManager{
						Client: &GRPCClient{
							Addrs:               addrs,
							HealthCheckDuration: healthcheck,
							ConnectionPool:      connectionPool,
							Backoff:             backoffOpts,
							CallOption:          callOpts,
							DialOption:          dialOpts,
							TLS:                 tls,
						},
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return BackupManager when the b.Client is nil",
				fields: fields{},
				want: want{
					want: &BackupManager{
						Client: &GRPCClient{
							DialOption: &DialOption{
								Insecure: true,
							},
						},
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
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			b := &BackupManager{
				Client: test.fields.Client,
			}

			got := b.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
