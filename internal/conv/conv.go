// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package conv

import (
	"io"
	"reflect"
	"strings"
	"unsafe"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// Btoa converts from byte slice to string.
func Btoa(b []byte) (s string) {
	// skipcq: GSC-G103
	slh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	// skipcq: GSC-G103
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh.Data = slh.Data
	sh.Len = slh.Len
	return s
}

// Atobs converts from string to byte slice.
func Atob(s string) (b []byte) {
	// skipcq: GSC-G103
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	// skipcq: GSC-G103
	slh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	slh.Data = sh.Data
	slh.Len = sh.Len
	slh.Cap = sh.Len
	return b
}

// F32stos converts from float32 slice to type string.
func F32stos(fs []float32) (s string) {
	lf := 4 * len(fs)
	// skipcq: GSC-G103
	buf := (*(*[1]byte)(unsafe.Pointer(&(fs[0]))))[:]
	// skipcq: GSC-G103
	addr := unsafe.Pointer(&buf)
	// skipcq: GSC-G103
	(*(*int)(unsafe.Pointer(uintptr(addr) + uintptr(8)))) = lf
	// skipcq: GSC-G103
	(*(*int)(unsafe.Pointer(uintptr(addr) + uintptr(16)))) = lf
	return Btoa(buf)
}

// Utf8ToSjis converts a UTF8 string to sjis string.
func Utf8ToSjis(s string) (string, error) {
	return encode(strings.NewReader(s), japanese.ShiftJIS.NewEncoder())
}

// Utf8ToEucjp converts a UTF8 string to eucjp string.
func Utf8ToEucjp(s string) (string, error) {
	return encode(strings.NewReader(s), japanese.EUCJP.NewEncoder())
}

func encode(r io.Reader, t transform.Transformer) (string, error) {
	b, err := io.ReadAll(transform.NewReader(r, t))
	if err != nil {
		return "", err
	}
	return string(b), nil
}
