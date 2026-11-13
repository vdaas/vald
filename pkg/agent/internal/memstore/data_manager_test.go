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
package memstore

import (
	"context"
	"slices"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/pkg/agent/internal/kvs"
	"github.com/vdaas/vald/pkg/agent/internal/vqueue"
)

// TestExistsGetObject_doNotWriteKVS pins that Exists and GetObject are pure reads: an insert that is
// superseded by a pending delete (dts == its) must not bump the kvs timestamp, and their
// observable results stay the same.
func TestExistsGetObject_doNotWriteKVS(t *testing.T) {
	t.Parallel()

	const (
		uuid = "uuid"
		oid  = uint32(7)
		kts  = int64(1)
		its  = int64(2)
	)
	kvVec := []float32{1, 2, 3}
	tests := []struct {
		call func(kv kvs.BidiMap, vq vqueue.Queue) error
		name string
	}{
		{
			name: "Exists",
			call: func(kv kvs.BidiMap, vq vqueue.Queue) error {
				got, ok := Exists(kv, vq, uuid)
				if !ok || got != oid {
					return errors.Errorf("Exists() = (%d, %v), want (%d, true)", got, ok, oid)
				}
				return nil
			},
		},
		{
			name: "GetObject",
			call: func(kv kvs.BidiMap, vq vqueue.Queue) error {
				vec, ts, err := GetObject(kv, vq, uuid, func(got uint32) ([]float32, error) {
					if got != oid {
						return nil, errors.Errorf("getVectorFn oid = %d, want %d", got, oid)
					}
					return kvVec, nil
				})
				if err != nil || ts != kts || !slices.Equal(vec, kvVec) {
					return errors.Errorf("GetObject() = (%v, %d, %v), want (%v, %d, nil)", vec, ts, err, kvVec, kts)
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			kv := kvs.New()
			kv.Set(uuid, oid, kts)
			vq, err := vqueue.New()
			if err != nil {
				t.Fatal(err)
			}
			// update pattern with equal timestamps: the insert is not visible (vqok=false, its != 0).
			if err := vq.PushInsert(uuid, []float32{4, 5, 6}, its); err != nil {
				t.Fatal(err)
			}
			if err := vq.PushDelete(uuid, its); err != nil {
				t.Fatal(err)
			}

			if err := tt.call(kv, vq); err != nil {
				t.Error(err)
			}
			gotOid, gotTs, ok := kv.Get(uuid)
			if !ok || gotOid != oid || gotTs != kts {
				t.Errorf("kvs after %s = (%d, %d, %v), want (%d, %d, true)", tt.name, gotOid, gotTs, ok, oid, kts)
			}
			if _, gotTs, ok := kv.GetInverse(oid); !ok || gotTs != kts {
				t.Errorf("inverse kvs after %s = (%d, %v), want (%d, true)", tt.name, gotTs, ok, kts)
			}
		})
	}
}

// afterGetKV runs afterGet once, right after the first Get, to interleave a concurrent
// kvs writer (CreateIndex or another UpdateTimestamp) between UpdateTimestamp's read and write.
type afterGetKV struct {
	kvs.BidiMap
	afterGet func(kv kvs.BidiMap)
}

func (k *afterGetKV) Get(key string) (uint32, int64, bool) {
	oid, ts, ok := k.BidiMap.Get(key)
	if f := k.afterGet; f != nil {
		k.afterGet = nil
		f(k.BidiMap)
	}
	return oid, ts, ok
}

// TestUpdateTimestamp_kvsWriteRacingCreateIndex pins that UpdateTimestamp never re-creates a
// kvs mapping that CreateIndex removed after UpdateTimestamp read it: kv.Set used to map the
// uuid back to an oid already removed from the index. A concurrent timestamp change alone
// must still succeed instead of reporting NotFound.
func TestUpdateTimestamp_kvsWriteRacingCreateIndex(t *testing.T) {
	t.Parallel()

	const (
		uuid = "uuid"
		oid  = uint32(7)
		kts  = int64(1)
		ts   = int64(10)
	)
	tests := []struct {
		queue    func(vq vqueue.Queue) error
		afterGet func(kv kvs.BidiMap)
		name     string
		wantErr  bool
		wantKV   bool
	}{
		{
			name: "kvs only, deleted by CreateIndex after the read",
			queue: func(vqueue.Queue) error {
				return nil
			},
			afterGet: func(kv kvs.BidiMap) { kv.Delete(uuid) },
			wantErr:  true,
		},
		{
			name: "kvs with a queued delete, applied by CreateIndex after the read",
			queue: func(vq vqueue.Queue) error {
				return vq.PushDelete(uuid, 2)
			},
			afterGet: func(kv kvs.BidiMap) { kv.Delete(uuid) },
			wantErr:  true,
		},
		{
			name: "kvs with a cancelled insert, deleted by CreateIndex after the read",
			queue: func(vq vqueue.Queue) error {
				if err := vq.PushInsert(uuid, []float32{1}, 2); err != nil {
					return err
				}
				return vq.PushDelete(uuid, 3)
			},
			afterGet: func(kv kvs.BidiMap) { kv.Delete(uuid) },
			wantErr:  true,
		},
		{
			name: "concurrent timestamp change still succeeds",
			queue: func(vqueue.Queue) error {
				return nil
			},
			afterGet: func(kv kvs.BidiMap) { kv.Set(uuid, oid, kts+1) },
			wantKV:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			base := kvs.New()
			t.Cleanup(func() {
				if err := base.Close(); err != nil {
					t.Error(err)
				}
			})
			base.Set(uuid, oid, kts)
			vq, err := vqueue.New()
			if err != nil {
				t.Fatal(err)
			}
			if err = tt.queue(vq); err != nil {
				t.Fatal(err)
			}

			err = UpdateTimestamp(&afterGetKV{BidiMap: base, afterGet: tt.afterGet}, vq, uuid, ts, false)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdateTimestamp() error = %v, wantErr %v", err, tt.wantErr)
			}
			gotOid, gotTs, ok := base.Get(uuid)
			if ok != tt.wantKV {
				t.Fatalf("kvs entry exists = %v (oid %d, ts %d), want %v", ok, gotOid, gotTs, tt.wantKV)
			}
			_, invTs, invOk := base.GetInverse(oid)
			if invOk != tt.wantKV {
				t.Fatalf("inverse kvs entry exists = %v, want %v", invOk, tt.wantKV)
			}
			if !tt.wantKV {
				if base.Len() != 0 {
					t.Errorf("kvs Len = %d, want 0", base.Len())
				}
				return
			}
			if gotOid != oid || gotTs != ts || invTs != ts {
				t.Errorf("kvs = (%d, %d), inverse ts %d, want (%d, %d)", gotOid, gotTs, invTs, oid, ts)
			}
		})
	}
}

// TestUpdateTimestamp_queuedInsertAndDelete pins UpdateTimestamp for an indexed object with
// both a queued insert and a queued delete that is not older than it. Raising the kvs
// timestamp to ts used to make CreateIndex skip the queued delete (newest wins), so a newer
// delete was lost and the superseded indexed vector was served again.
func TestUpdateTimestamp_queuedInsertAndDelete(t *testing.T) {
	t.Parallel()

	const (
		uuid = "uuid"
		oid  = uint32(7)
		kts  = int64(50)
	)
	qvec := []float32{1, 2}
	tests := []struct {
		// wantErr is the expected error; nil means success.
		wantErr error
		name    string
		its     int64
		dts     int64
		ts      int64
		// wantIts is the queued insert timestamp afterwards.
		wantIts int64
		force   bool
	}{
		{
			name:    "newer delete cancels the insert: pending deletion is not found",
			its:     101,
			dts:     200,
			ts:      300,
			wantErr: errors.ErrObjectNotFound(nil, uuid),
			wantIts: 101,
		},
		{
			name:    "force on a pending deletion is not found",
			its:     101,
			dts:     200,
			ts:      150,
			force:   true,
			wantErr: errors.ErrObjectNotFound(nil, uuid),
			wantIts: 101,
		},
		{
			name:    "delete tying the insert: the queued vector is re-queued at ts",
			its:     100,
			dts:     100,
			ts:      300,
			wantIts: 300,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			kv := kvs.New()
			t.Cleanup(func() {
				if err := kv.Close(); err != nil {
					t.Error(err)
				}
			})
			kv.Set(uuid, oid, kts)
			vq, err := vqueue.New()
			if err != nil {
				t.Fatal(err)
			}
			if err = vq.PushInsert(uuid, qvec, tt.its); err != nil {
				t.Fatal(err)
			}
			if err = vq.PushDelete(uuid, tt.dts); err != nil {
				t.Fatal(err)
			}

			err = UpdateTimestamp(kv, vq, uuid, tt.ts, tt.force)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("UpdateTimestamp() error = %v, want nil", err)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf("UpdateTimestamp() error = %v, want %v", err, tt.wantErr)
			}

			// The kvs timestamp must not pass the queued delete, or CreateIndex skips it.
			gotOid, gotKts, ok := kv.Get(uuid)
			if !ok || gotOid != oid || gotKts != kts {
				t.Errorf("kvs = (%d, %d, %v), want (%d, %d, true)", gotOid, gotKts, ok, oid, kts)
			}
			vec, its, dts, _ := vq.GetVectorWithTimestamp(uuid)
			if !slices.Equal(vec, qvec) || its != tt.wantIts {
				t.Errorf("queued insert = (%v, %d), want (%v, %d)", vec, its, qvec, tt.wantIts)
			}
			if dts != tt.dts {
				t.Errorf("queued delete ts = %d, want %d", dts, tt.dts)
			}
		})
	}
}

// When the indexed object is newer than both queued entries, CreateIndex skips
// them (newest wins), so UpdateTimestamp must update the indexed object and drop
// the stale insert; re-queuing that insert at ts would let it replace the newer
// indexed object at the next CreateIndex.
func TestUpdateTimestamp_indexedNewerThanQueued(t *testing.T) {
	t.Parallel()

	const (
		uuid = "uuid"
		oid  = uint32(7)
		kts  = int64(200)
		ts   = int64(300)
	)
	tests := []struct {
		name string
		its  int64
		dts  int64
	}{
		{name: "delete tying the stale insert", its: 100, dts: 100},
		{name: "delete newer than the stale insert", its: 100, dts: 150},
		{name: "stale insert without a delete", its: 100},
		{name: "stale insert tying the indexed object", its: kts},
		{name: "stale insert newer than its delete", its: 150, dts: 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			kv := kvs.New()
			t.Cleanup(func() {
				if err := kv.Close(); err != nil {
					t.Error(err)
				}
			})
			kv.Set(uuid, oid, kts)
			vq, err := vqueue.New()
			if err != nil {
				t.Fatal(err)
			}
			if err = vq.PushInsert(uuid, []float32{1, 2}, tt.its); err != nil {
				t.Fatal(err)
			}
			if tt.dts != 0 {
				if err = vq.PushDelete(uuid, tt.dts); err != nil {
					t.Fatal(err)
				}
			}

			if err = UpdateTimestamp(kv, vq, uuid, ts, false); err != nil {
				t.Fatalf("UpdateTimestamp() error = %v, want nil", err)
			}
			gotOid, gotKts, ok := kv.Get(uuid)
			if !ok || gotOid != oid || gotKts != ts {
				t.Errorf("kvs = (%d, %d, %v), want (%d, %d, true)", gotOid, gotKts, ok, oid, ts)
			}
			if vec, its, dts, _ := vq.GetVectorWithTimestamp(uuid); vec != nil || its != 0 || dts != 0 {
				t.Errorf("stale entries still queued: insert (%v, %d), delete %d", vec, its, dts)
			}
		})
	}
}

type lenKV struct {
	kvs.BidiMap
	uuids []string
}

func (k lenKV) Len() uint64 { return uint64(len(k.uuids)) }

func (k lenKV) Range(_ context.Context, f func(string, uint32, int64) bool) {
	for i, uuid := range k.uuids {
		if !f(uuid, uint32(i+1), 1) {
			return
		}
	}
}

type lenVQ struct {
	vqueue.Queue
	ivq, dvq int
}

func (q lenVQ) IVQLen() int { return q.ivq }

func (q lenVQ) DVQLen() int { return q.dvq }

func (lenVQ) Range(context.Context, func(string, []float32, int64) bool) {}

func (lenVQ) DVExists(string) (int64, bool) { return 0, false }

// TestUUIDs_capacityHint pins that a delete queue longer than kvs+insert queue
// (kvs is decremented before the delete entry is popped) does not wrap the
// uint64 capacity hint and panic in make.
func TestUUIDs_capacityHint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		kv   []string
		ivq  int
		dvq  int
	}{
		{name: "more pending deletes than objects", kv: []string{"a"}, ivq: 0, dvq: 2},
		{name: "pending deletes and empty kvs", kv: nil, ivq: 0, dvq: 1},
		{name: "fewer pending deletes than objects", kv: []string{"a", "b"}, ivq: 1, dvq: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got []string
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("UUIDs panicked: %v", r)
					}
				}()
				got = UUIDs(t.Context(), lenKV{uuids: tt.kv}, lenVQ{ivq: tt.ivq, dvq: tt.dvq})
			}()
			if !slices.Equal(got, tt.kv) {
				t.Errorf("UUIDs() = %v, want %v", got, tt.kv)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestExists(t *testing.T) {
// 	type args struct {
// 		kv   kvs.BidiMap
// 		vq   vqueue.Queue
// 		uuid string
// 	}
// 	type want struct {
// 		wantOid uint32
// 		wantOk  bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, uint32, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotOid uint32, gotOk bool) error {
// 		if !reflect.DeepEqual(gotOid, w.wantOid) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOid, w.wantOid)
// 		}
// 		if !reflect.DeepEqual(gotOk, w.wantOk) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           kv:nil,
// 		           vq:nil,
// 		           uuid:"",
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
// 		           kv:nil,
// 		           vq:nil,
// 		           uuid:"",
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
// 			gotOid, gotOk := Exists(test.args.kv, test.args.vq, test.args.uuid)
// 			if err := checkFunc(test.want, gotOid, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestGetObject(t *testing.T) {
// 	type args struct {
// 		kv          kvs.BidiMap
// 		vq          vqueue.Queue
// 		uuid        string
// 		getVectorFn func(oid uint32) ([]float32, error)
// 	}
// 	type want struct {
// 		wantVec       []float32
// 		wantTimestamp int64
// 		err           error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []float32, int64, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec []float32, gotTimestamp int64, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
// 		}
// 		if !reflect.DeepEqual(gotTimestamp, w.wantTimestamp) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTimestamp, w.wantTimestamp)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           kv:nil,
// 		           vq:nil,
// 		           uuid:"",
// 		           getVectorFn:nil,
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
// 		           kv:nil,
// 		           vq:nil,
// 		           uuid:"",
// 		           getVectorFn:nil,
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
// 			gotVec, gotTimestamp, err := GetObject(test.args.kv, test.args.vq, test.args.uuid, test.args.getVectorFn)
// 			if err := checkFunc(test.want, gotVec, gotTimestamp, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestListObjectFunc(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		kv  kvs.BidiMap
// 		vq  vqueue.Queue
// 		f   func(uuid string, oid uint32, ts int64) bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           kv:nil,
// 		           vq:nil,
// 		           f:nil,
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
// 		           ctx:nil,
// 		           kv:nil,
// 		           vq:nil,
// 		           f:nil,
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
// 			ListObjectFunc(test.args.ctx, test.args.kv, test.args.vq, test.args.f)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestUUIDs(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		kv  kvs.BidiMap
// 		vq  vqueue.Queue
// 	}
// 	type want struct {
// 		wantUuids []string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotUuids []string) error {
// 		if !reflect.DeepEqual(gotUuids, w.wantUuids) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotUuids, w.wantUuids)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           kv:nil,
// 		           vq:nil,
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
// 		           ctx:nil,
// 		           kv:nil,
// 		           vq:nil,
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
// 			gotUuids := UUIDs(test.args.ctx, test.args.kv, test.args.vq)
// 			if err := checkFunc(test.want, gotUuids); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestUpdateTimestamp(t *testing.T) {
// 	type args struct {
// 		kv          kvs.BidiMap
// 		vq          vqueue.Queue
// 		uuid        string
// 		ts          int64
// 		force bool
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
// 		           kv:nil,
// 		           vq:nil,
// 		           uuid:"",
// 		           ts:0,
// 		           force:false,
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
// 		           kv:nil,
// 		           vq:nil,
// 		           uuid:"",
// 		           ts:0,
// 		           force:false,
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
// 			err := UpdateTimestamp(test.args.kv, test.args.vq, test.args.uuid, test.args.ts, test.args.force)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
