// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

import "strings"

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

	oldnew = []string{
		" ", "",
		",", "",
		"-", "",
		".", "",
		"/", "",
		":", "",
		";", "",
		"=", "",
		"\\", "",
		"_", "",
		"|", "",
		"~", "",
		"\n", "",
		"\r", "",
		"\t", "",
	}
	reps = strings.NewReplacer(oldnew...)
	_    = reps.Replace(" ")
)

// Join is strings.Join: it sizes one Builder exactly, so a pooled buffer only
// added a Get/Put and a copy in String().
func Join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

func TrimForCompare[S ~string](str S) S {
	return S(reps.Replace(ToLower(string(str))))
}
