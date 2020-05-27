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
package usecase

import "strings"

type Operation uint8

const (
	Unknown Operation = iota
	Insert
	Search
)

func Atoo(s string) Operation {
	switch strings.ToLower(s) {
	case "insert":
		return Insert
	case "search":
		return Search
	default:
		return Unknown
	}
}

func (o Operation) String() string {
	switch o {
	case Insert:
		return "insert"
	case Search:
		return "search"
	default:
		return "unknown operation"
	}
}
