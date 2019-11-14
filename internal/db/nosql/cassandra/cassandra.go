//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cassandra

import (
	"context"
)

type Cassandra interface {
	// TODO Define interface
	Close() error
}

type client struct {
}

func New(ctx context.Context, opts ...Option) (Cassandra, error) {
	c := new(client)

	for _, opt := range append(defaultOpts, opts...) {
		opt(c)
	}

	return nil, nil
}
