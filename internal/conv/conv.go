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
	"reflect"
	"unsafe"
)

// Btoa converts from byte slice to string.
func Btoa(b []byte) string {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: header.Data,
		Len:  header.Len,
	}))
}

// Atobs converts from string to byte slice.
func Atob(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: (*(*reflect.StringHeader)(unsafe.Pointer(&s))).Data,
		Len:  len(s),
		Cap:  len(s),
	}))
}

// F32stos converts from float32 slice to type string.
func F32stos(fs []float32) string {
	lf := 4 * len(fs)
	buf := (*(*[1]byte)(unsafe.Pointer(&(fs[0]))))[:]
	addr := unsafe.Pointer(&buf)
	(*(*int)(unsafe.Pointer(uintptr(addr) + uintptr(8)))) = lf
	(*(*int)(unsafe.Pointer(uintptr(addr) + uintptr(16)))) = lf
	return *(*string)(unsafe.Pointer(&buf))
}
