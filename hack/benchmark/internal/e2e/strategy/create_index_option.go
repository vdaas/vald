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

// Package strategy provides strategy for e2e testing functions
package strategy

import "github.com/vdaas/vald/internal/client/v1/client"

type CreateIndexOption func(*createIndex)

var defaultCreateIndexOptions = []CreateIndexOption{
	WithCreateIndexPoolSize(10000),
}

func WithCreateIndexPoolSize(size int) CreateIndexOption {
	return func(ci *createIndex) {
		if size > 0 {
			ci.poolSize = uint32(size)
		}
	}
}

func WithCreateIndexClient(c client.Indexer) CreateIndexOption {
	return func(ci *createIndex) {
		if c != nil {
			ci.Indexer = c
		}
	}
}
