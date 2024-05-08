// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
package strings

import (
	"bytes"
	"strings"
	"syscall"

	"github.com/vdaas/vald/internal/sync"
)

type (
	Builder  = strings.Builder
	Reader   = strings.Reader
	Replacer = strings.Replacer
)

var (
	Clone          = strings.Clone
	Compare        = strings.Compare
	Contains       = strings.Contains
	ContainsAny    = strings.ContainsAny
	ContainsFunc   = strings.ContainsFunc
	ContainsRune   = strings.ContainsRune
	Count          = strings.Count
	Cut            = strings.Cut
	CutPrefix      = strings.CutPrefix
	CutSuffix      = strings.CutSuffix
	EqualFold      = strings.EqualFold
	Fields         = strings.Fields
	FieldsFunc     = strings.FieldsFunc
	HasPrefix      = strings.HasPrefix
	HasSuffix      = strings.HasSuffix
	Index          = strings.Index
	IndexAny       = strings.IndexAny
	IndexByte      = strings.IndexByte
	IndexFunc      = strings.IndexFunc
	IndexRune      = strings.IndexRune
	LastIndex      = strings.LastIndex
	LastIndexAny   = strings.LastIndexAny
	LastIndexByte  = strings.LastIndexByte
	LastIndexFunc  = strings.LastIndexFunc
	Map            = strings.Map
	Repeat         = strings.Repeat
	Replace        = strings.Replace
	ReplaceAll     = strings.ReplaceAll
	Split          = strings.Split
	SplitAfter     = strings.SplitAfter
	SplitAfterN    = strings.SplitAfterN
	SplitN         = strings.SplitN
	Title          = strings.Title
	ToLower        = strings.ToLower
	ToLowerSpecial = strings.ToLowerSpecial
	ToTitle        = strings.ToTitle
	ToTitleSpecial = strings.ToTitleSpecial
	ToUpper        = strings.ToUpper
	ToUpperSpecial = strings.ToUpperSpecial
	ToValidUTF8    = strings.ToValidUTF8
	Trim           = strings.Trim
	TrimFunc       = strings.TrimFunc
	TrimLeft       = strings.TrimLeft
	TrimLeftFunc   = strings.TrimLeftFunc
	TrimPrefix     = strings.TrimPrefix
	TrimRight      = strings.TrimRight
	TrimRightFunc  = strings.TrimRightFunc
	TrimSpace      = strings.TrimSpace
	TrimSuffix     = strings.TrimSuffix
	NewReader      = strings.NewReader
	NewReplacer    = strings.NewReplacer

	bufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, syscall.Getpagesize()))
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
	b.Grow(n)
	b.WriteString(elems[0])
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(s)
	}
	return b.String()
}
