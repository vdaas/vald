// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package ngt provides implementation of Go API for https://github.com/yahoojapan/NGT
package ngt

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"reflect"
	"testing"
)

const (
	index    = "./assets/index"
	poolSize = 2
)

func TestCreate(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "tmpdir")
	if err != nil {
		t.Errorf("Unexpected error: TestCreate(%v)", err)
	}
	defer os.RemoveAll(tmpdir)

	ngt, err := New(
		WithIndexPath(tmpdir),
		WithObjectType(Uint8),
		WithDimension(6),
	)
	defer ngt.Close()
	if err != nil {
		t.Errorf("Unexpected error: TestCreate(%v)", err)
	}
}

func TestInsert(t *testing.T) {
	tests := []struct {
		vector []float64
		want   uint
	}{
		{[]float64{1, 0, 0, 0, 0, 0}, 1},
		{[]float64{0, 1, 0, 0, 0, 0}, 2},
		{[]float64{0, 0, 1, 0, 0, 0}, 3},
		{[]float64{0, 0, 0, 1, 0, 0}, 4},
		{[]float64{0, 0, 0, 0, 1, 0}, 5},
		{[]float64{0, 0, 0, 0, 0, 1}, 6},
		{[]float64{1, 1, 0, 0, 0, 0}, 7},
	}

	tmpdir, err := ioutil.TempDir("", "tmpdir")
	if err != nil {
		t.Errorf("Unexpected error: TestInsert(%v)", err)
	}
	defer os.RemoveAll(tmpdir)

	ngt, err := New(
		WithIndexPath(tmpdir),
		WithObjectType(Uint8),
		WithDimension(6),
	)
	defer ngt.Close()
	if err != nil {
		t.Errorf("Unexpected error: TestInsert(%v)", err)
	}

	for _, tt := range tests {
		id, err := ngt.Insert(tt.vector)
		if err != nil {
			t.Fatal(err)
		}
		if id != tt.want {
			t.Errorf("TestInsert(%v): %v, wanted: %v", tt.vector, id, tt.want)
		}
	}
}

func TestInsertCommit(t *testing.T) {
	tests := []struct {
		vector []float64
		want   uint
	}{
		{[]float64{1, 0, 0, 0, 0, 0}, 1},
		{[]float64{0, 1, 0, 0, 0, 0}, 2},
		{[]float64{0, 0, 1, 0, 0, 0}, 3},
		{[]float64{0, 0, 0, 1, 0, 0}, 4},
		{[]float64{0, 0, 0, 0, 1, 0}, 5},
		{[]float64{0, 0, 0, 0, 0, 1}, 6},
		{[]float64{1, 1, 0, 0, 0, 0}, 7},
	}

	tmpdir, err := ioutil.TempDir("", "tmpdir")
	if err != nil {
		t.Errorf("Unexpected error: TestInsertCommit(%v)", err)
	}
	defer os.RemoveAll(tmpdir)

	ngt, err := New(
		WithIndexPath(tmpdir),
		WithObjectType(Uint8),
		WithDimension(6),
	)
	defer ngt.Close()
	if err != nil {
		t.Errorf("Unexpected error: TestInsertCommit(%v)", err)
	}

	for _, tt := range tests {
		id, err := ngt.InsertCommit(tt.vector, 2)
		if err != nil {
			t.Errorf("Unexpected error: TestInsertCommit(%v)", err)
		}
		if id != tt.want {
			t.Errorf("TestInsertCommit(%v): %v, wanted: %v", tt.vector, id, tt.want)
		}
	}
}

func TestBulkInsert(t *testing.T) {
	tests := []struct {
		vectors [][]float64
		wants   []uint
	}{
		{
			[][]float64{
				{1, 0, 0, 0, 0, 0},
				{0, 1, 0, 0, 0, 0},
				{0, 0, 1, 0, 0, 0},
				{0, 0, 0, 1, 0, 0},
				{0, 0, 0, 0, 1, 0},
				{0, 0, 0, 0, 0, 1},
				{1, 1, 0, 0, 0, 0},
			},
			[]uint{1, 2, 3, 4, 5, 6, 7},
		},
	}

	tmpdir, err := ioutil.TempDir("", "tmpdir")
	if err != nil {
		t.Errorf("Unexpected error: TestBulkInsert(%v)", err)
	}
	defer os.RemoveAll(tmpdir)

	ngt, err := New(
		WithIndexPath(tmpdir),
		WithObjectType(Uint8),
		WithDimension(6),
	)
	defer ngt.Close()
	if err != nil {
		t.Errorf("Unexpected error: TestBulkInsert(%v)", err)
	}

	for _, tt := range tests {
		ids, errs := ngt.BulkInsert(tt.vectors)
		if len(errs) > 0 {
			t.Errorf("Unexpected error: TestBulkInsert(%v)", errs)
		}
		if !reflect.DeepEqual(ids, tt.wants) {
			t.Errorf("TestBulkInsert(%v): %v, wanted: %v", tt.vectors, ids, tt.wants)
		}
	}
}

func TestBulkInsertCommit(t *testing.T) {
	tests := []struct {
		vectors [][]float64
		wants   []uint
	}{
		{
			[][]float64{
				{1, 0, 0, 0, 0, 0},
				{0, 1, 0, 0, 0, 0},
				{0, 0, 1, 0, 0, 0},
				{0, 0, 0, 1, 0, 0},
				{0, 0, 0, 0, 1, 0},
				{0, 0, 0, 0, 0, 1},
				{1, 1, 0, 0, 0, 0},
			},
			[]uint{1, 2, 3, 4, 5, 6, 7},
		},
	}

	tmpdir, err := ioutil.TempDir("", "tmpdir")
	if err != nil {
		t.Errorf("Unexpected error: TestBulkInsert(%v)", err)
	}
	defer os.RemoveAll(tmpdir)

	ngt, err := New(
		WithIndexPath(tmpdir),
		WithObjectType(Uint8),
		WithDimension(6),
	)
	defer ngt.Close()
	if err != nil {
		t.Errorf("Unexpected error: TestBulkInsertCommit(%v)", err)
	}

	for _, tt := range tests {
		ids, errs := ngt.BulkInsertCommit(tt.vectors, 2)
		if len(errs) > 0 {
			t.Errorf("Unexpected error: TestBulkInsertCommit(%v)", errs)
		}
		if !reflect.DeepEqual(ids, tt.wants) {
			t.Errorf("TestBulkInsertCommit(%v): %v, wanted: %v", tt.vectors, ids, tt.wants)
		}
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		vector []float64
		want   SearchResult
	}{
		{[]float64{1, 0, 0, 0, 0, 0}, SearchResult{1, 0, nil}},
		{[]float64{0, 1, 0, 0, 0, 0}, SearchResult{2, 0, nil}},
		{[]float64{0, 0, 1, 0, 0, 0}, SearchResult{3, 0, nil}},
		{[]float64{0, 0, 0, 1, 0, 0}, SearchResult{4, 0, nil}},
		{[]float64{0, 0, 0, 0, 1, 0}, SearchResult{5, 0, nil}},
		{[]float64{1, 1, 0, 0, 0, 0}, SearchResult{6, 0, nil}},
	}

	ngt, err := Load(
		WithIndexPath(index),
		WithObjectType(Uint8),
		WithDimension(6),
	)
	defer ngt.Close()
	if err != nil {
		t.Errorf("Unexpected error: TestSearch(%v)", err)
	}

	for _, tt := range tests {
		result, err := ngt.Search(tt.vector, 1, 0.1, -1.0)
		if err != nil {
			t.Errorf("Unexpected error: TestSearch(%v)", err)
		}
		if result[0].ID != tt.want.ID || result[0].Distance != tt.want.Distance {
			t.Errorf("TestSearch(%v): %v, wanted: %v", tt.vector, result, tt.want)
		}
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		id   uint
		want error
	}{
		{1, nil},
		{2, nil},
		{3, nil},
		{4, nil},
		{5, nil},
		{6, nil},
	}
	tmpdir, err := ioutil.TempDir("", "tmpdir")
	if err != nil {
		t.Errorf("Unexpected error: TestRemove(%v)", err)
	}
	defer os.RemoveAll(tmpdir)

	if err := exec.Command("cp", "-r", index, tmpdir).Run(); err != nil {
		t.Errorf("Unexpected error: TestRemove(%v)", err)
	}

	ngt, err := Load(
		WithIndexPath(path.Join(tmpdir, "index")),
	)
	defer ngt.Close()
	if err != nil {
		t.Errorf("Unexpected error: TestRemove(%v)", err)
	}
	for _, tt := range tests {
		if err := ngt.Remove(tt.id); err != tt.want {
			t.Errorf("TestRemove(%v): %v, wanted: %v", tt.id, err, tt.want)
		}
	}
}

func TestGetVector(t *testing.T) {
	tests := []struct {
		id   uint
		want []float32
	}{
		{1, []float32{1, 0, 0, 0, 0, 0}},
		{2, []float32{0, 1, 0, 0, 0, 0}},
		{3, []float32{0, 0, 1, 0, 0, 0}},
		{4, []float32{0, 0, 0, 1, 0, 0}},
		{5, []float32{0, 0, 0, 0, 1, 0}},
		{6, []float32{1, 1, 0, 0, 0, 0}},
	}
	ngt, err := Load(
		WithIndexPath(index),
		WithObjectType(Uint8),
		WithDimension(6),
	)
	defer ngt.Close()
	if err != nil {
		t.Errorf("Unexpected error: TestGetVector(%v)", err)
	}
	for _, tt := range tests {
		vec, err := ngt.GetVector(tt.id)
		if err != nil {
			t.Errorf("Unexpected error: TestGetVector(%v)", err)
		}
		if !reflect.DeepEqual(vec, tt.want) {
			t.Errorf("TestGetVector(%v): %v, wanted: %v", tt.id, vec, tt.want)
		}
	}
}
