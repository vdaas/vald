//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package strings

import (
	"bytes"
	"strings"
	"sync"
)

type (
	Builder = strings.Builder
)

var (
	Contains       = strings.Contains
	Count          = strings.Count
	EqualFold      = strings.EqualFold
	HasPrefix      = strings.HasPrefix
	HasSuffix      = strings.HasSuffix
	Index          = strings.Index
	IndexAny       = strings.IndexAny
	NewReader      = strings.NewReader
	NewReplacer    = strings.NewReplacer
	Replace        = strings.Replace
	ReplaceAll     = strings.ReplaceAll
	Split          = strings.Split
	SplitAfter     = strings.SplitAfter
	SplitAfterN    = strings.SplitAfterN
	SplitN         = strings.SplitN
	ToLower        = strings.ToLower
	ToLowerSpecial = strings.ToLowerSpecial
	ToUpper        = strings.ToUpper
	ToUpperSpecial = strings.ToUpperSpecial
	Trim           = strings.Trim
	TrimLeft       = strings.TrimLeft
	TrimPrefix     = strings.TrimPrefix
	TrimRight      = strings.TrimRight
	TrimSpace      = strings.TrimSpace
	TrimSuffix     = strings.TrimSuffix

	bufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 1024))
		},
	}
)

func Join(elems []string, sep string) (str string) {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return elems[0]
	}
	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	b := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(b)
	defer b.Reset()
	b.WriteString(elems[0])
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(s)
	}
	return b.String()
}
