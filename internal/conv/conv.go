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
package conv

import (
	"io/ioutil"
	"reflect"
	"strings"
	"unsafe"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// Btoa converts from byte slice to string.
func Btoa(b []byte) (s string) {
	slh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh.Data = slh.Data
	sh.Len = slh.Len
	return s
}

// Atobs converts from string to byte slice.
func Atob(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	slh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	slh.Data = sh.Data
	slh.Len = sh.Len
	slh.Cap = sh.Len
	return b
}

// F32stos converts from float32 slice to type string.
func F32stos(fs []float32) (s string) {
	lf := 4 * len(fs)
	buf := (*(*[1]byte)(unsafe.Pointer(&(fs[0]))))[:]
	addr := unsafe.Pointer(&buf)
	(*(*int)(unsafe.Pointer(uintptr(addr) + uintptr(8)))) = lf
	(*(*int)(unsafe.Pointer(uintptr(addr) + uintptr(16)))) = lf
	return Btoa(buf)
}

func Utf8ToSjis(s string) string {
	b, _ := ioutil.ReadAll(transform.NewReader(strings.NewReader(s), japanese.ShiftJIS.NewEncoder()))
	return string(b)
}

func Utf8ToEucjp(s string) string {
	b, _ := ioutil.ReadAll(transform.NewReader(strings.NewReader(s), japanese.EUCJP.NewEncoder()))
	return string(b)
}
