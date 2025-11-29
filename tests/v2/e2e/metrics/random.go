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

import (
	"crypto/rand"
	"encoding/binary"
	"math/big"
)

func secureUint64() uint64 {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return binary.LittleEndian.Uint64(b[:])
}

func secureUint64N(n uint64) uint64 {
	if n == 0 {
		return 0
	}
	max := new(big.Int).SetUint64(n)
	v, _ := rand.Int(rand.Reader, max)
	return v.Uint64()
}

func secureIntN(n int) int {
	if n <= 0 {
		return 0
	}
	max := big.NewInt(int64(n))
	v, _ := rand.Int(rand.Reader, max)
	return int(v.Int64())
}
