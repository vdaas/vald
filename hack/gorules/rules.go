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
package gorules

import (
	"github.com/quasilyte/go-ruleguard/dsl"
)

func CheckPayloadObjectAccess(m dsl.Matcher) {
	m.Import("github.com/vdaas/vald/apis/grpc/v1/payload")

	m.Match(
		`$*_ = $*_, $x.$y, $*_`,
		`$*_, $x.$y, $*_`,
		`$x.$y != $_`,
		`$x.$y == $_`,
		`$x.$y.$_`,
		`$x.$y[$_].$_`,
	).Where(!m["y"].Text.Matches(`Get.+`) &&
		m["x"].Type.Implements(`payload.Payload`)).
		At(m["y"]).
		Report("Avoid to access struct fields directly").
		Suggest(`Get$y()`)
}

// from Ruleguard by example
// https://go-ruleguard.github.io/by-example/matching-comments
func PrintFmt(m dsl.Matcher) {
	// Here we scan a $s string for %s, %d and %v.
	m.Match(`fmt.Println($s, $*_)`,
		`fmt.Print($s, $*_)`).
		Where(m["s"].Text.Matches(`%[sdv]`)).
		Report("found formatting directive in non-formatting call")
}
