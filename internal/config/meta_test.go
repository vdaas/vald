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
	"strconv"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestMeta_Bind(t *testing.T) {
	t.Parallel()
	type fields struct {
		Host                      string
		Port                      uint16
		Client                    *GRPCClient
		EnableCache               bool
		CacheExpiration           string
		ExpiredCacheCheckDuration string
	}
	type want struct {
		want *Meta
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Meta) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *Meta) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			host := "vald-meta.vald.svc.cluster.local"
			port := uint16(8081)
			enableCache := true
			cacheExpiration := "24h"
			expiredCacheCheckDuration := "1m"
			return test{
				name: "return Meta when all parameters are not nil or empty",
				fields: fields{
					Host: host,
					Port: port,
					Client: &GRPCClient{
						DialOption: &DialOption{
							Insecure: true,
						},
					},
					EnableCache:               enableCache,
					CacheExpiration:           cacheExpiration,
					ExpiredCacheCheckDuration: expiredCacheCheckDuration,
				},
				want: want{
					want: &Meta{
						Host: host,
						Port: port,
						Client: &GRPCClient{
							Addrs: []string{
								host + ":" + strconv.FormatUint(uint64(port), 10),
							},
							ConnectionPool: &ConnectionPool{},
							DialOption: &DialOption{
								Insecure: true,
							},
							TLS: &TLS{
								Enabled: false,
							},
						},
						EnableCache:               enableCache,
						CacheExpiration:           cacheExpiration,
						ExpiredCacheCheckDuration: expiredCacheCheckDuration,
					},
				},
			}
		}(),
		func() test {
			host := "vald-meta.vald.svc.cluster.local"
			port := uint16(8081)
			enableCache := true
			cacheExpiration := "24h"
			expiredCacheCheckDuration := "1m"
			return test{
				name: "return Meta when Client is nil and others are not empty",
				fields: fields{
					Host:                      host,
					Port:                      port,
					EnableCache:               enableCache,
					CacheExpiration:           cacheExpiration,
					ExpiredCacheCheckDuration: expiredCacheCheckDuration,
				},
				want: want{
					want: &Meta{
						Host: host,
						Port: port,
						Client: &GRPCClient{
							Addrs: []string{
								host + ":" + strconv.FormatUint(uint64(port), 10),
							},
							DialOption: &DialOption{
								Insecure: true,
							},
						},
						EnableCache:               enableCache,
						CacheExpiration:           cacheExpiration,
						ExpiredCacheCheckDuration: expiredCacheCheckDuration,
					},
				},
			}
		}(),
		func() test {
			p := map[string]string{
				"HOST":                         "vald-meta.vald.svc.cluster.local",
				"CACHE_EXPIRATION":             "24h",
				"EXPIRED_CACHE_CHECK_DURATION": "1m",
			}
			port := uint16(8081)
			enableCache := true
			return test{
				name: "return Meta when some parameters are set as environment value",
				fields: fields{
					Host:                      "_HOST_",
					Port:                      port,
					EnableCache:               enableCache,
					CacheExpiration:           "_CACHE_EXPIRATION_",
					ExpiredCacheCheckDuration: "_EXPIRED_CACHE_CHECK_DURATION_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range p {
						if err := os.Setenv(k, v); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					for k := range p {
						if err := os.Unsetenv(k); err != nil {
							t.Fatal(err)
						}
					}
				},
				want: want{
					want: &Meta{
						Host: "vald-meta.vald.svc.cluster.local",
						Port: port,
						Client: &GRPCClient{
							Addrs: []string{
								"vald-meta.vald.svc.cluster.local" + ":" + strconv.FormatUint(uint64(port), 10),
							},
							DialOption: &DialOption{
								Insecure: true,
							},
						},
						EnableCache:               enableCache,
						CacheExpiration:           "24h",
						ExpiredCacheCheckDuration: "1m",
					},
				},
			}
		}(),
		func() test {
			return test{
				name:   "return Meta when all parameters are nil or empty",
				fields: fields{},
				want: want{
					want: &Meta{
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
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			m := &Meta{
				Host:                      test.fields.Host,
				Port:                      test.fields.Port,
				Client:                    test.fields.Client,
				EnableCache:               test.fields.EnableCache,
				CacheExpiration:           test.fields.CacheExpiration,
				ExpiredCacheCheckDuration: test.fields.ExpiredCacheCheckDuration,
			}

			got := m.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
