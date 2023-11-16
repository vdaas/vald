// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package service

// Option represents the functional option for index.
type Option func(_ *rotator) error

var defaultOpts = []Option{}

func WithReplicaId(id int) Option {
	return func(r *rotator) error {
		r.replicaid = id
		return nil
	}
}

func WithNamespace(ns string) Option {
	return func(r *rotator) error {
		r.namespace = ns
		return nil
	}
}

func WithDeploymentPrefix(dp string) Option {
	return func(r *rotator) error {
		r.deploymentPrefix = dp
		return nil
	}
}

func WithSnapshotPrefix(sp string) Option {
	return func(r *rotator) error {
		r.snapshotPrefix = sp
		return nil
	}
}

func WithPvcPrefix(pp string) Option {
	return func(r *rotator) error {
		r.pvcPrefix = pp
		return nil
	}
}
