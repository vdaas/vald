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

package metadata

import (
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/os"
)

func TestLoad_nilDecode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "returns ErrMetadataFileEmpty when the file contains only null",
			content: "null",
		},
		{
			name:    "returns ErrMetadataFileEmpty when the file contains only whitespace",
			content: " \n\t",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(t.TempDir(), AgentMetadataFileName)
			if err := os.WriteFile(path, []byte(tc.content), fs.ModePerm); err != nil {
				t.Fatalf("failed to write metadata: %v", err)
			}
			got, err := Load(path)
			if !errors.Is(err, errors.ErrMetadataFileEmpty) {
				t.Errorf("got error %v, want %v", err, errors.ErrMetadataFileEmpty)
			}
			if got != nil {
				t.Errorf("got metadata %#v, want nil", got)
			}
		})
	}
}

func TestStore_atomic(t *testing.T) {
	t.Parallel()
	tests := []struct {
		check func(t *testing.T, path string)
		name  string
	}{
		{
			name: "replaces an existing metadata file without rewriting it in place",
			check: func(t *testing.T, path string) {
				t.Helper()
				old := &Metadata{NGT: &NGT{IndexCount: 1}}
				if err := Store(path, old); err != nil {
					t.Fatalf("failed to store initial metadata: %v", err)
				}
				// A hard link keeps pointing at the old inode. An in-place O_TRUNC rewrite
				// changes that inode (so a crash mid-write leaves it empty); a rename does not.
				link := filepath.Join(filepath.Dir(path), "old-metadata-link")
				if err := os.Link(path, link); err != nil {
					t.Fatalf("failed to link metadata: %v", err)
				}
				if err := Store(path, &Metadata{NGT: &NGT{IndexCount: 2}}); err != nil {
					t.Fatalf("failed to store metadata: %v", err)
				}
				got, err := Load(path)
				if err != nil {
					t.Fatalf("failed to load metadata: %v", err)
				}
				if got == nil || got.NGT == nil || got.NGT.IndexCount != 2 {
					t.Fatalf("got metadata %#v, want index count 2", got)
				}
				prev, err := Load(link)
				if err != nil {
					t.Fatalf("failed to load previous metadata: %v", err)
				}
				if prev == nil || prev.NGT == nil || prev.NGT.IndexCount != 1 {
					t.Fatalf("previous metadata was modified in place: %#v", prev)
				}
				assertNoTmpFile(t, path)
			},
		},
		{
			name: "creates a new metadata file",
			check: func(t *testing.T, path string) {
				t.Helper()
				if err := Store(path, &Metadata{IsInvalid: true}); err != nil {
					t.Fatalf("failed to store metadata: %v", err)
				}
				got, err := Load(path)
				if err != nil {
					t.Fatalf("failed to load metadata: %v", err)
				}
				if got == nil || !got.IsInvalid {
					t.Fatalf("got metadata %#v, want IsInvalid", got)
				}
				assertNoTmpFile(t, path)
			},
		},
		{
			name: "removes the temporary file when the rename fails",
			check: func(t *testing.T, path string) {
				t.Helper()
				// A non-empty directory at path makes the final rename fail.
				if err := os.MkdirAll(filepath.Join(path, "child"), fs.ModePerm); err != nil {
					t.Fatalf("failed to create directory: %v", err)
				}
				if err := Store(path, &Metadata{}); err == nil {
					t.Fatal("Store succeeded, want an error")
				}
				assertNoTmpFile(t, path)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.check(t, filepath.Join(t.TempDir(), AgentMetadataFileName))
		})
	}
}

func assertNoTmpFile(t *testing.T, path string) {
	t.Helper()
	_, err := os.Stat(path + tmpFileSuffix)
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("temporary metadata file is left behind or cannot be checked: %v", err)
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestLoad(t *testing.T) {
// 	type args struct {
// 		path string
// 	}
// 	type want struct {
// 		wantMeta *Metadata
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *Metadata, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotMeta *Metadata, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotMeta, w.wantMeta) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotMeta, w.wantMeta)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           path:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           path:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotMeta, err := Load(test.args.path)
// 			if err := checkFunc(test.want, gotMeta, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestStore(t *testing.T) {
// 	type args struct {
// 		path string
// 		meta *Metadata
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           path:"",
// 		           meta:Metadata{},
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           path:"",
// 		           meta:Metadata{},
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			err := Store(test.args.path, test.args.meta)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
