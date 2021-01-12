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
)

func TestGateway_Bind(t *testing.T) {
	type fields struct {
		AgentPort      int
		AgentName      string
		AgentNamespace string
		AgentDNS       string
		NodeName       string
		IndexReplica   int
		Discoverer     *DiscovererClient
		Meta           *Meta
		BackupManager  *BackupManager
		EgressFilter   *EgressFilter
	}
	type want struct {
		want *Gateway
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Gateway) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *Gateway) error {
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
		           AgentPort: 0,
		           AgentName: "",
		           AgentNamespace: "",
		           AgentDNS: "",
		           NodeName: "",
		           IndexReplica: 0,
		           Discoverer: DiscovererClient{},
		           Meta: Meta{},
		           BackupManager: BackupManager{},
		           EgressFilter: EgressFilter{},
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
		           AgentPort: 0,
		           AgentName: "",
		           AgentNamespace: "",
		           AgentDNS: "",
		           NodeName: "",
		           IndexReplica: 0,
		           Discoverer: DiscovererClient{},
		           Meta: Meta{},
		           BackupManager: BackupManager{},
		           EgressFilter: EgressFilter{},
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
			g := &Gateway{
				AgentPort:      test.fields.AgentPort,
				AgentName:      test.fields.AgentName,
				AgentNamespace: test.fields.AgentNamespace,
				AgentDNS:       test.fields.AgentDNS,
				NodeName:       test.fields.NodeName,
				IndexReplica:   test.fields.IndexReplica,
				Discoverer:     test.fields.Discoverer,
				Meta:           test.fields.Meta,
				BackupManager:  test.fields.BackupManager,
				EgressFilter:   test.fields.EgressFilter,
			}

			got := g.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
