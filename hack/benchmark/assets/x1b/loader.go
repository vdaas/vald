// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package x1b

import (
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"github.com/vdaas/vald/internal/errors"
)

const (
	headerSize = 4
)

var (
	ErrOutOfBounds         = errors.New("out of bounds")
	ErrUnsupportedFileType = errors.New("unsupported file type")
)

type BillionScaleVectors interface {
	Load(i int) (interface{}, error)
	Dimension() int
	Size() int
	Close() error
}

type Uint8Vectors interface {
	BillionScaleVectors
	LoadUint8(i int) ([]uint8, error)
}

type FloatVectors interface {
	BillionScaleVectors
	LoadFloat32(i int) ([]float32, error)
}

type Int32Vectors interface {
	BillionScaleVectors
	LoadInt32(i int) ([]int32, error)
}

type file struct {
	mem   []byte
	dim   int
	size  int
	block int
}

type bvecs struct {
	*file
}

type fvecs struct {
	*file
}

type ivecs struct {
	*file
}

func openFile(fname string, elementSize int) (f *file, err error) {
	fp, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := fp.Close(); e != nil {
			err = errors.Wrap(err, e.Error())
		}
	}()

	fi, err := fp.Stat()
	if err != nil {
		return nil, err
	}

	mem, err := syscall.Mmap(int(fp.Fd()), 0, int(fi.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	// skipcq: GSC-G103
	dim := int(*(*int32)(unsafe.Pointer(&mem[0])))
	block := headerSize + dim*elementSize
	return &file{
		mem:   mem,
		dim:   dim,
		size:  len(mem) / block,
		block: block,
	}, nil
}

func (f *file) Close() error {
	return syscall.Munmap(f.mem)
}

func (f *file) load(i int) ([]byte, error) {
	if i >= f.size {
		return nil, ErrOutOfBounds
	}

	return f.mem[i*f.block+headerSize : (i+1)*f.block], nil
}

func (f *file) Dimension() int {
	return f.dim
}

func (f *file) Size() int {
	return f.size
}

func (bv *bvecs) LoadUint8(i int) ([]uint8, error) {
	buf, err := bv.load(i)
	if err != nil {
		return nil, err
	}
	// skipcq: GSC-G103
	return ((*[1 << 26]uint8)(unsafe.Pointer(&buf[0])))[:bv.dim:bv.dim], nil
}

func (bv *bvecs) Load(i int) (interface{}, error) {
	return bv.LoadUint8(i)
}

func (fv *fvecs) LoadFloat32(i int) ([]float32, error) {
	buf, err := fv.load(i)
	if err != nil {
		return nil, err
	}
	// skipcq: GSC-G103
	return ((*[1 << 26]float32)(unsafe.Pointer(&buf[0])))[:fv.dim:fv.dim], nil
}

func (fv *fvecs) Load(i int) (interface{}, error) {
	return fv.LoadFloat32(i)
}

func (iv *ivecs) LoadInt32(i int) ([]int32, error) {
	buf, err := iv.load(i)
	if err != nil {
		return nil, err
	}
	// skipcq: GSC-G103
	return ((*[1 << 26]int32)(unsafe.Pointer(&buf[0])))[:iv.dim:iv.dim], nil
}

func (iv *ivecs) Load(i int) (interface{}, error) {
	return iv.LoadInt32(i)
}

func NewUint8Vectors(fname string) (Uint8Vectors, error) {
	f, err := openFile(fname, 1)
	if err != nil {
		return nil, err
	}
	return &bvecs{f}, nil
}

func NewFloatVectors(fname string) (FloatVectors, error) {
	f, err := openFile(fname, 4)
	if err != nil {
		return nil, err
	}
	return &fvecs{f}, nil
}

func NewInt32Vectors(fname string) (Int32Vectors, error) {
	f, err := openFile(fname, 4)
	if err != nil {
		return nil, err
	}
	return &ivecs{f}, nil
}

func Open(fname string) (BillionScaleVectors, error) {
	switch filepath.Ext(fname) {
	case ".bvecs":
		return NewUint8Vectors(fname)
	case ".fvecs":
		return NewFloatVectors(fname)
	case ".ivecs":
		return NewInt32Vectors(fname)
	default:
		return nil, ErrUnsupportedFileType
	}
}
