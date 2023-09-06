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
package operation

import "github.com/vdaas/vald/internal/client/v1/client"

type Option func(*operation)

func WithClient(c client.Client) Option {
	return func(o *operation) {
		o.client = c
	}
}

func WithIndexer(c client.Indexer) Option {
	return func(o *operation) {
		o.indexerC = c
	}
}
