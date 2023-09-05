//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package malloc provides Go wrapper of malloc_info
package malloc

/*
#include <errno.h>
#include <malloc.h>
#include <stdio.h>
#include <string.h>

char* strerror2() {
	return strerror(errno);
}

*/
import "C"

import (
	"encoding/xml"
	"unsafe"

	"github.com/vdaas/vald/internal/errors"
)

type Size struct {
	From  int `xml:"from,attr"`
	To    int `xml:"to,attr"`
	Total int `xml:"total,attr"`
	Count int `xml:"count,attr"`
}

type Sizes struct {
	Size     []*Size `xml:"size"`
	Unsorted Size    `xml:"unsorted"`
}

type Total struct {
	Type  string `xml:"type,attr"`
	Count int    `xml:"count,attr"`
	Size  int    `xml:"size,attr"`
}

type System struct {
	Type string `xml:"type,attr"`
	Size int    `xml:"size,attr"`
}

type Aspace struct {
	Type string `xml:"type,attr"`
	Size int    `xml:"size,attr"`
}

type Heap struct {
	Nr     string    `xml:"nr,attr"`
	Sizes  Sizes     `xml:"sizes"`
	Total  []*Total  `xml:"total"`
	System []*System `xml:"system"`
	Aspace []*Aspace `xml:"aspace"`
}

type MallocInfo struct {
	Version string    `xml:"version,attr"`
	Heap    []*Heap   `xml:"heap"`
	Total   []*Total  `xml:"total"`
	System  []*System `xml:"system"`
	Aspace  []*Aspace `xml:"aspace"`
}

func convert(body string) (*MallocInfo, error) {
	var m MallocInfo
	if err := xml.Unmarshal([]byte(body), &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func GetMallocInfo() (*MallocInfo, error) {
	var ptr *C.char
	var size C.size_t
	in := C.open_memstream(&ptr, &size)
	defer func() {
		C.fclose(in)
		C.free(unsafe.Pointer(ptr))
	}()

	ret := C.malloc_info(0, in)
	switch ret {
	case 0:
		C.fflush(in)
		return convert(C.GoStringN(ptr, C.int(size)))
	case -1:
		return nil, errors.New(C.GoString(C.strerror2()))
	default:
		return nil, errors.ErrUnexpectedReturnCode(int(ret))
	}
}
