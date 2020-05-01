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
package search

import (
	"github.com/vdaas/vald/internal/client"
)

type SearchOption func(*search) error

var (
	defaultSearchOpts = []SearchOption{
		WithParallelDegree(100),
	}
)

func WithReader(r client.Reader) SearchOption {
	return func(s *search) error {
		s.r = r
		return nil
	}
}

func WithParallelDegree(p int) SearchOption {
	return func(s *search) error {
		s.p = p
		return nil
	}
}

func WithDataset(n string) SearchOption {
	return func(s *search) (err error) {
		s.n = n
		return nil
	}
}
