//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package metrics

import "github.com/vdaas/vald/internal/errors"

// mergeableShard represents a shard that can be merged.
// T is the interface type that the shard implements (e.g., Histogram, TDigest).
type mergeableShard[T any] interface {
	Merge(T) error
}

// mergeShards merges the shards from source into target.
// S is the concrete shard type (e.g., *histogram).
// T is the interface type accepted by Merge (e.g., Histogram).
// S must implement mergeableShard[T].
// Note: S is expected to be assignable to T (or implement T), but we can't enforce this
// relationship strictly with Go 1.18+ generics constraints without using T as a constraint for S,
// which is not allowed. We use a runtime type assertion instead.
func mergeShards[T any, S mergeableShard[T]](target []S, source []S) error {
	if len(target) != len(source) {
		return errors.New("incompatible shards: count mismatch")
	}
	for i := range target {
		// Cast source[i] to T. This is safe as long as S implements T.
		src, ok := any(source[i]).(T)
		if !ok {
			return errors.New("shard does not implement interface T")
		}
		if err := target[i].Merge(src); err != nil {
			return err
		}
	}
	return nil
}
