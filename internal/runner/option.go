//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package runner

import "github.com/vdaas/vald/internal/config"

type Option[T any] func(*runner[T])

func WithName[T any](name string) Option[T] {
	return func(r *runner[T]) {
		if name != "" {
			r.name = name
		}
	}
}

func WithVersion[T any](ver, maxVer, minVer string) Option[T] {
	return func(r *runner[T]) {
		if ver != "" {
			r.version = ver
		}
		if maxVer != "" {
			r.maxVersion = maxVer
		}
		if minVer != "" {
			r.minVersion = minVer
		}
	}
}

func WithConfigLoader[T any](f func(string) (T, *config.GlobalConfig, error)) Option[T] {
	return func(r *runner[T]) {
		if f != nil {
			r.loadConfig = f
		}
	}
}

func WithDaemonInitializer[T any](f func(T) (Runner, error)) Option[T] {
	return func(r *runner[T]) {
		if f != nil {
			r.initializeDaemon = f
		}
	}
}
