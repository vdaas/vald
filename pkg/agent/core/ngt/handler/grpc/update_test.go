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
package grpc

import (
	"context"
	"math"
	"slices"
	"strconv"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/data/request"
	"github.com/vdaas/vald/internal/test/data/vector"
)

// tsNGT records which multi update/delete service method the handler called
// and with which timestamp; other methods panic through the nil embedded interface.
type tsNGT struct {
	stubNGT
	method string
	ts     int64
}

func (n *tsNGT) UpdateMultiple(map[string][]float32) error {
	n.method = "UpdateMultiple"
	return nil
}

func (n *tsNGT) UpdateMultipleWithTime(_ map[string][]float32, t int64) error {
	n.method, n.ts = "UpdateMultipleWithTime", t
	return nil
}

func (n *tsNGT) DeleteMultiple(...string) error {
	n.method = "DeleteMultiple"
	return nil
}

func (n *tsNGT) DeleteMultipleWithTime(_ []string, t int64) error {
	n.method, n.ts = "DeleteMultipleWithTime", t
	return nil
}

// Test_server_multi_configTimestamp pins that MultiUpdate and MultiRemove pass
// the first request's Config.Timestamp to the service like their single
// counterparts (Update, Remove) do, and fall back to server time when it is 0.
func Test_server_multi_configTimestamp(t *testing.T) {
	t.Parallel()
	const dim = 3
	vec := []float32{1, 2, 3}
	updateReqs := func(ts int64) *payload.Update_MultiRequest {
		return &payload.Update_MultiRequest{Requests: []*payload.Update_Request{
			{Vector: &payload.Object_Vector{Id: "a", Vector: vec}, Config: &payload.Update_Config{Timestamp: ts}},
			{Vector: &payload.Object_Vector{Id: "b", Vector: vec}, Config: &payload.Update_Config{Timestamp: ts + 1}},
		}}
	}
	removeReqs := func(ts int64) *payload.Remove_MultiRequest {
		return &payload.Remove_MultiRequest{Requests: []*payload.Remove_Request{
			{Id: &payload.Object_ID{Id: "a"}, Config: &payload.Remove_Config{Timestamp: ts}},
			{Id: &payload.Object_ID{Id: "b"}, Config: &payload.Remove_Config{Timestamp: ts + 1}},
		}}
	}
	tests := []struct {
		call       func(context.Context, *server) error
		name       string
		wantMethod string
		wantTS     int64
	}{
		{
			name: "MultiUpdate with Config.Timestamp",
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiUpdate(ctx, updateReqs(42))
				return err
			},
			wantMethod: "UpdateMultipleWithTime",
			wantTS:     42,
		},
		{
			name: "MultiUpdate without Config.Timestamp",
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiUpdate(ctx, updateReqs(0))
				return err
			},
			wantMethod: "UpdateMultiple",
		},
		{
			name: "MultiRemove with Config.Timestamp",
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiRemove(ctx, removeReqs(42))
				return err
			},
			wantMethod: "DeleteMultipleWithTime",
			wantTS:     42,
		},
		{
			name: "MultiRemove without Config.Timestamp",
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiRemove(ctx, removeReqs(0))
				return err
			},
			wantMethod: "DeleteMultiple",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			n := &tsNGT{dim: dim}
			s := &server{ngt: n}
			if err := test.call(t.Context(), s); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if n.method != test.wantMethod || n.ts != test.wantTS {
				t.Errorf("called %s(ts=%d), want %s(ts=%d)", n.method, n.ts, test.wantMethod, test.wantTS)
			}
		})
	}
}

// Test_server_MultiUpdate_perIDErrors pins that MultiUpdate reports the same
// per-ID failures as Update (NotFound for a missing ID, InvalidArgument for an
// empty ID, AlreadyExists for an identical vector) instead of returning OK,
// while the valid IDs of the batch are still updated.
func Test_server_MultiUpdate_perIDErrors(t *testing.T) {
	t.Parallel()
	const (
		dim       = 3
		insertNum = 10
		existID   = "exist"
		otherID   = "other"
	)
	existVec := []float32{1, 2, 3}
	otherVec := []float32{4, 5, 6}
	newVec := []float32{7, 8, 9}
	cfg := &config.NGT{
		Dimension:        dim,
		DistanceType:     ngt.L2.String(),
		ObjectType:       ngt.Float.String(),
		CreationEdgeSize: 60,
		SearchEdgeSize:   20,
		KVSDB:            &config.KVSDB{Concurrency: 10},
		VQueue: &config.VQueue{
			InsertBufferPoolSize: 1000,
			DeleteBufferPoolSize: 1000,
		},
	}
	req := func(id string, vec []float32) *payload.Update_Request {
		return &payload.Update_Request{Vector: &payload.Object_Vector{Id: id, Vector: vec}}
	}
	tests := []struct {
		name     string
		reqs     []*payload.Update_Request
		wantCode codes.Code
	}{
		{
			name:     "all IDs updated",
			reqs:     []*payload.Update_Request{req(existID, newVec)},
			wantCode: codes.OK,
		},
		{
			name:     "non-existent ID",
			reqs:     []*payload.Update_Request{req(existID, newVec), req("missing", newVec)},
			wantCode: codes.NotFound,
		},
		{
			name:     "empty ID",
			reqs:     []*payload.Update_Request{req(existID, newVec), req("", newVec)},
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "identical vector",
			reqs:     []*payload.Update_Request{req(existID, newVec), req(otherID, otherVec)},
			wantCode: codes.AlreadyExists,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			eg, ctx := errgroup.New(t.Context())
			n, err := newIndexedNGTService(ctx, eg, request.Float, vector.Gaussian, insertNum, nil, cfg, nil,
				[]string{existID, otherID}, [][]float32{existVec, otherVec})
			if err != nil {
				t.Fatal(err)
			}
			s, err := New(WithErrGroup(eg), WithNGT(n))
			if err != nil {
				t.Fatal(err)
			}
			res, err := s.MultiUpdate(ctx, &payload.Update_MultiRequest{Requests: test.reqs})
			if test.wantCode == codes.OK {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got := len(res.GetLocations()); got != len(test.reqs) {
					t.Errorf("got %d locations, want %d", got, len(test.reqs))
				}
			} else {
				st, ok := status.FromError(err)
				if err == nil || !ok || st.Code() != test.wantCode {
					t.Fatalf("got error %v, want code %v", err, test.wantCode)
				}
			}
			// the valid ID of the batch is updated in every case.
			got, _, err := n.GetObject(existID)
			if err != nil {
				t.Fatalf("GetObject(%s) error: %v", existID, err)
			}
			if !slices.Equal(got, newVec) {
				t.Errorf("GetObject(%s) = %v, want %v", existID, got, newVec)
			}
		})
	}
}

// countedError is a per-ID failure that counts how often its message is built.
type countedError struct {
	calls *atomic.Int64
	msg   string
}

func (e countedError) Error() string {
	e.calls.Add(1)
	return e.msg
}

// multiErrNGT fails every multi insert/update/delete with err; other methods
// panic through the nil embedded interface.
type multiErrNGT struct {
	stubNGT
}

func (n *multiErrNGT) InsertMultiple(map[string][]float32) error { return n.err }

func (n *multiErrNGT) UpdateMultiple(map[string][]float32) error { return n.err }

func (n *multiErrNGT) DeleteMultiple(...string) error { return n.err }

// Test_server_multi_perIDErrorClassification pins the gRPC code and message
// MultiInsert, MultiUpdate and MultiRemove build from a joined per-ID error,
// and that classifying it builds each member's message O(1) times: checking
// every ID with errors.Is rebuilt the whole joined message per ID, O(N^2).
func Test_server_multi_perIDErrorClassification(t *testing.T) {
	t.Parallel()
	const (
		dim = 3
		n   = 200
	)
	vec := []float32{1, 2, 3}
	ids := make([]string, 0, n)
	for i := range n {
		ids = append(ids, "id-"+strconv.Itoa(i))
	}
	idList := "[" + strings.Join(ids, " ") + "]"
	// leaves builds one counted failure per ID; msg maps an ID to its message.
	leaves := func(calls *atomic.Int64, msg func(id string) string) []error {
		errs := make([]error, 0, n)
		for _, id := range ids {
			errs = append(errs, countedError{calls: calls, msg: msg(id)})
		}
		return errs
	}
	notFound := func(id string) string { return errors.ErrObjectIDNotFound(id).Error() }
	alreadyExists := func(id string) string { return errors.ErrUUIDAlreadyExists(id).Error() }
	update := func(ctx context.Context, s *server) error {
		reqs := &payload.Update_MultiRequest{Requests: make([]*payload.Update_Request, 0, n)}
		for _, id := range ids {
			reqs.Requests = append(reqs.Requests, &payload.Update_Request{Vector: &payload.Object_Vector{Id: id, Vector: vec}})
		}
		_, err := s.MultiUpdate(ctx, reqs)
		return err
	}
	insert := func(ctx context.Context, s *server) error {
		reqs := &payload.Insert_MultiRequest{Requests: make([]*payload.Insert_Request, 0, n)}
		for _, id := range ids {
			reqs.Requests = append(reqs.Requests, &payload.Insert_Request{Vector: &payload.Object_Vector{Id: id, Vector: vec}})
		}
		_, err := s.MultiInsert(ctx, reqs)
		return err
	}
	remove := func(ctx context.Context, s *server) error {
		reqs := &payload.Remove_MultiRequest{Requests: make([]*payload.Remove_Request, 0, n)}
		for _, id := range ids {
			reqs.Requests = append(reqs.Requests, &payload.Remove_Request{Id: &payload.Object_ID{Id: id}})
		}
		_, err := s.MultiRemove(ctx, reqs)
		return err
	}
	tests := []struct {
		call     func(context.Context, *server) error
		err      func(calls *atomic.Int64) error
		name     string
		wantMsg  string // empty: not checked (built from map iteration order)
		wantCode codes.Code
	}{
		{
			name:     "MultiUpdate not found",
			call:     update,
			err:      func(c *atomic.Int64) error { return errors.Join(leaves(c, notFound)...) },
			wantCode: codes.NotFound,
			wantMsg:  "MultiUpdate API uuids " + idList + " not found",
		},
		{
			name: "MultiUpdate not found matched case-insensitively like errors.Is",
			call: update,
			err: func(c *atomic.Int64) error {
				return errors.Join(leaves(c, func(id string) string { return strings.ToUpper(notFound(id)) })...)
			},
			wantCode: codes.NotFound,
			wantMsg:  "MultiUpdate API uuids " + idList + " not found",
		},
		{
			name: "MultiUpdate not found in a nested and wrapped join",
			call: update,
			err: func(c *atomic.Int64) error {
				l := leaves(c, notFound)
				return errors.Join(l[0], errors.Wrap(errors.Join(l[1:]...), "insert failed"))
			},
			wantCode: codes.NotFound,
			wantMsg:  "MultiUpdate API uuids " + idList + " not found",
		},
		{
			name: "MultiUpdate invalid dimension",
			call: update,
			err: func(c *atomic.Int64) error {
				return errors.Join(leaves(c, func(string) string {
					return errors.ErrInvalidDimensionSize(dim, dim).Error()
				})...)
			},
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "MultiUpdate already exists",
			call:     update,
			err:      func(c *atomic.Int64) error { return errors.Join(leaves(c, alreadyExists)...) },
			wantCode: codes.AlreadyExists,
			wantMsg:  "MultiUpdate API uuids " + idList + " already exists",
		},
		{
			name: "MultiUpdate unknown error",
			call: update,
			err: func(c *atomic.Int64) error {
				return errors.Join(leaves(c, func(id string) string { return "boom " + id })...)
			},
			wantCode: codes.Internal,
			wantMsg:  "MultiUpdate API failed",
		},
		{
			name:     "MultiInsert already exists",
			call:     insert,
			err:      func(c *atomic.Int64) error { return errors.Join(leaves(c, alreadyExists)...) },
			wantCode: codes.AlreadyExists,
			wantMsg:  "MultiInsert API uuids " + idList + " already exists",
		},
		{
			name:     "MultiRemove not found",
			call:     remove,
			err:      func(c *atomic.Int64) error { return errors.Join(leaves(c, notFound)...) },
			wantCode: codes.NotFound,
			wantMsg:  "MultiRemove API uuids " + idList + " not found",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var calls atomic.Int64
			s := &server{ngt: &multiErrNGT{err: test.err(&calls), dim: dim}}
			calls.Store(0)
			err := test.call(t.Context(), s)
			st, ok := status.FromError(err)
			if err == nil || !ok || st.Code() != test.wantCode {
				t.Fatalf("got error %v, want code %v", err, test.wantCode)
			}
			if test.wantMsg != "" && st.Message() != test.wantMsg {
				t.Errorf("message = %q, want %q", st.Message(), test.wantMsg)
			}
			// a fixed number of passes over the members, not one per ID.
			if got := calls.Load(); got > 20*n {
				t.Errorf("member messages built %d times for %d IDs, want <= %d", got, n, 20*n)
			}
		})
	}
}

func Test_server_Update(t *testing.T) {
	t.Parallel()

	type args struct {
		req         *payload.Update_Request
		indexID     string
		indexVector []float32
	}
	type fields struct {
		objectType string
	}
	type want struct {
		wantUUID string
		code     codes.Code
	}
	type test struct {
		checkFunc  func(want, *payload.Object_Location, error) error
		beforeFunc func(*testing.T, args) (Server, error)
		afterFunc  func(args)
		name       string
		fields     fields
		args       args
		want       want
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_Location, err error) error {
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				errors.Errorf("got error cannot convert to Status: \"%#v\"", err)
			}
			if st.Code() != w.code {
				return errors.Errorf("got code: \"%#v\",\n\t\t\t\twant code: \"%#v\"", st.Code(), w.code)
			}
		} else {
			if uuid := gotRes.GetUuid(); w.wantUUID != uuid {
				return errors.Errorf("got uuid: \"%#v\",\n\t\t\t\twant uuid: \"%#v\"", uuid, w.wantUUID)
			}
		}
		return nil
	}

	const (
		insertNum = 1000
		dimension = 128
	)

	utf8Str := "こんにちは"
	eucjpStr, err := conv.Utf8ToEucjp(utf8Str)
	if err != nil {
		t.Error(err)
	}
	sjisStr, err := conv.Utf8ToSjis(utf8Str)
	if err != nil {
		t.Error(err)
	}

	defaultUpdateConfig := &payload.Update_Config{
		SkipStrictExistCheck: true,
	}
	defaultInsertConfig := &payload.Insert_Config{
		SkipStrictExistCheck: true,
	}
	beforeFunc := func(t *testing.T, ctx context.Context, objectType string) func(*testing.T, args) (Server, error) {
		t.Helper()
		if objectType == "" {
			objectType = ngt.Float.String()
		}

		cfg := &config.NGT{
			Dimension:        dimension,
			DistanceType:     ngt.L2.String(),
			ObjectType:       objectType,
			CreationEdgeSize: 60,
			SearchEdgeSize:   20,
			KVSDB: &config.KVSDB{
				Concurrency: 10,
			},
			VQueue: &config.VQueue{
				InsertBufferPoolSize: 1000,
				DeleteBufferPoolSize: 1000,
			},
		}

		return func(t *testing.T, a args) (Server, error) {
			t.Helper()
			var overwriteVec [][]float32
			if a.indexVector != nil {
				overwriteVec = [][]float32{
					a.indexVector,
				}
			}

			eg, ctx := errgroup.New(ctx)
			ngt, err := newIndexedNGTService(ctx, eg, request.Float, vector.Gaussian, insertNum, defaultInsertConfig, cfg, nil, []string{a.indexID}, overwriteVec)
			if err != nil {
				return nil, err
			}
			s, err := New(WithErrGroup(eg), WithNGT(ngt))
			if err != nil {
				return nil, err
			}
			return s, nil
		}
	}

	/*
		Update test cases (only test float32 unless otherwise specified):
		- Equivalence Class Testing ( 1000 vectors inserted before a update )
			- case 1.1: success update one vector
			- case 2.1: fail update with non-existent ID
			- case 3.1: fail update with one different dimension vector (type: uint8)
			- case 3.2: fail update with one different dimension vector (type: float32)
		- Boundary Value Testing ( 1000 vectors inserted before a update )
			- case 1.1: fail update with "" as ID
			- case 2.1: success update with ^@ as ID
			- case 2.2: success update with ^I as ID
			- case 2.3: success update with ^J as ID
			- case 2.4: success update with ^M as ID
			- case 2.5: success update with ^[ as ID
			- case 2.6: success update with ^? as ID
			- case 3.1: success update with utf-8 ID from utf-8 index
			- case 3.2: fail update with utf-8 ID from s-jis index
			- case 3.3: fail update with utf-8 ID from euc-jp index
			- case 3.4: fail update with s-jis ID from utf-8 index
			- case 3.5: success update with s-jis ID from s-jis index
			- case 3.6: fail update with s-jis ID from euc-jp index
			- case 3.4: fail update with euc-jp ID from utf-8 index
			- case 3.5: fail update with euc-jp ID from s-jis index
			- case 3.6: success update with euc-jp ID from euc-jp index
			- case 4.1: success update with 😀 as ID
			- case 5.1: success update with one 0 value vector (type: uint8)
			- case 5.2: success update with one +0 value vector (type: float32)
			- case 5.3: success update with one -0 value vector (type: float32)
			- case 6.1: success update with one min value vector (type: uint8)
			- case 6.2: success update with one min value vector (type: float32)
			- case 7.1: success update with one max value vector (type: uint8)
			- case 7.2: success update with one max value vector (type: float32)
			- case 8.1: fail update with one NaN value vector (type: float32)
			- case 9.1: fail update with one +inf value vector (type: float32)
			- case 9.2: fail update with one -inf value vector (type: float32)
			- case 10.1: fail update with one nil vector
			- case 11.1: fail update with one empty vector
		- Decision Table Testing
			- case 1.1: fail update with one duplicated vector, duplicated ID and SkipStrictExistCheck is true
			- case 1.2: success update with one different vector, duplicated ID and SkipStrictExistsCheck is true
			- case 1.3: success update with one duplicated vector, different ID and SkipStrictExistCheck is true
			- case 2.1: fail update with one duplicated vector, duplicated ID and SkipStrictExistCheck is false
			- case 2.2: success update with one different vector, duplicated ID and SkipStrictExistsCheck is false
			- case 2.3: success update with one duplicated vector, different ID and SkipStrictExistCheck is false
	*/
	tests := []test{
		{
			name: "Equivalent Class Testing case 1.1: success update one vector",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Equivalent Class Testing case 2.1: fail update with non-existent ID",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "non-existent",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Equivalent Class Testing case 3.1: fail update with one different dimension vector (type: uint8)",
			fields: fields{
				objectType: ngt.Uint8.String(),
			},
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, dimension+1))[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Equivalent Class Testint case 3.2: fail update with one different dimension vector (type: float32)",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension+1)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				code: codes.InvalidArgument,
			},
		},

		{
			name: "Boundary Value Testing case 1.1: fail update with \"\" as ID",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 2.1: success update with ^@ as ID",
			args: args{
				indexID: string([]byte{0}),
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     string([]byte{0}),
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: string([]byte{0}),
			},
		},
		{
			name: "Boundary Value Testing case 2.2: success update with ^I as ID",
			args: args{
				indexID: "\t",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "\t",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: "\t",
			},
		},
		{
			name: "Boundary Value Testing case 2.3: success update with ^J as ID",
			args: args{
				indexID: "\n",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "\n",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: "\n",
			},
		},
		{
			name: "Boundary Value Testing case 2.4: success update with ^M as ID",
			args: args{
				indexID: "\r",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "\r",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: "\r",
			},
		},
		{
			name: "Boundary Value Testing case 2.5: success update with ^[ as ID",
			args: args{
				indexID: string([]byte{27}),
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     string([]byte{27}),
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: string([]byte{27}),
			},
		},
		{
			name: "Boundary Value Testing case 2.6: success update with ^? as ID",
			args: args{
				indexID: string([]byte{127}),
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     string([]byte{127}),
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: string([]byte{127}),
			},
		},
		{
			name: "Boundary Value Testing case 3.1: success update with utf-8 ID from utf-8 index",
			args: args{
				indexID: utf8Str,
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     utf8Str,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: utf8Str,
			},
		},
		{
			name: "Boundary Value Testing case 3.2: success update with utf-8 ID from s-jis index",
			args: args{
				indexID: sjisStr,
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     utf8Str,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.3: success update with utf-8 ID from euc-jp index",
			args: args{
				indexID: eucjpStr,
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     utf8Str,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.4: fail update with s-jis ID from utf-8 index",
			args: args{
				indexID: utf8Str,
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     sjisStr,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.5: success update with s-jis ID from s-jis index",
			args: args{
				indexID: sjisStr,
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     sjisStr,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: sjisStr,
			},
		},
		{
			name: "Boundary Value Testing case 3.6: fail update with s-jis ID from euc-jp index",
			args: args{
				indexID: eucjpStr,
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     sjisStr,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.7: fail update with euc-jp ID from utf-8 index",
			args: args{
				indexID: utf8Str,
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     eucjpStr,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.8: fail update with euc-jp ID from s-jis index",
			args: args{
				indexID: sjisStr,
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     eucjpStr,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				code: codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.9: success update with euc-jp ID from euc-jp index",
			args: args{
				indexID: eucjpStr,
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     eucjpStr,
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: eucjpStr,
			},
		},
		{
			name: "Boundary Value Testing case 4.1: success update with 😀 as ID",
			args: args{
				indexID: "😀",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "😀",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: "😀",
			},
		},
		{
			name: "Boundary Value Testing case 5.1: success update with one 0 value vector (type: uint8)",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(uint8(0))),
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 5.2: success update with one +0 value vector (type: float32)",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, 0),
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 5.3: success update with one -0 value vector (type: float32)",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.Copysign(0, -1.0))),
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 6.1: success update with one min value vector (type: uint8)",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(uint8(0))),
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 6.2: success update with one min value vector (type: float32)",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, -math.MaxFloat32),
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 7.1: success update with one max value vector (type: uint8)",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.MaxUint8)),
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 7.2: success update with one max value vector (type: float32)",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, math.MaxFloat32),
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Boundary Value Testing case 8.1: fail update with one NaN value vector (type: float32)",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.NaN())),
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 9.1: fail update with one +inf value vector (type: float32)",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.Inf(1.0))),
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 9.2: fail update with one -inf value vector (type: float32)",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GenSameValueVec(dimension, float32(math.Inf(-1.0))),
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 10.1: fail update with one nil vector",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: nil,
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				code: codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 11.1: fail update with one empty vector",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: []float32{},
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				code: codes.InvalidArgument,
			},
		},

		{
			name: "Decision Table Testing case 1.1: fail update with one duplicated vector, duplicated ID and SkipStrictExistCheck is true",
			args: func() args {
				vector := vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0]
				return args{
					indexID:     "test",
					indexVector: vector,
					req: &payload.Update_Request{
						Vector: &payload.Object_Vector{
							Id:     "test",
							Vector: vector,
						},
						Config: defaultUpdateConfig,
					},
				}
			}(),
			want: want{
				code: codes.AlreadyExists,
			},
		},
		{
			name: "Decision Table Testing case 1.2: success update with one different vector, duplicated ID and SkipStrictExistCheck is true",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: defaultUpdateConfig,
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Decision Table Testing case 1.3: success update with one duplicated vector, different ID and SkipStrictExistCheck is true",
			args: func() args {
				vector := vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0]
				return args{
					indexID:     "test",
					indexVector: vector,
					req: &payload.Update_Request{
						Vector: &payload.Object_Vector{
							Id:     "uuid-2", // the first uuid is overwritten, so use the second one
							Vector: vector,
						},
						Config: defaultUpdateConfig,
					},
				}
			}(),
			want: want{
				wantUUID: "uuid-2",
			},
		},
		{
			name: "Decision Table Testing case 2.1: fail update with one duplicated vector, duplicated ID and SkipStrictExistCheck is false",
			args: func() args {
				vector := vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0]
				return args{
					indexID:     "test",
					indexVector: vector,
					req: &payload.Update_Request{
						Vector: &payload.Object_Vector{
							Id:     "test",
							Vector: vector,
						},
						Config: &payload.Update_Config{
							SkipStrictExistCheck: false,
						},
					},
				}
			}(),
			want: want{
				code: codes.AlreadyExists,
			},
		},
		{
			name: "Decision Table Testing case 2.2: success update with one duplicated vector, duplicated ID and SkipStrictExistCheck is false",
			args: args{
				indexID: "test",
				req: &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     "test",
						Vector: vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0],
					},
					Config: &payload.Update_Config{
						SkipStrictExistCheck: false,
					},
				},
			},
			want: want{
				wantUUID: "test",
			},
		},
		{
			name: "Decision Table Testing case 2.3: success update with one duplicated vector, different ID and SkipStrictExistCheck is false",
			args: func() args {
				vector := vector.GaussianDistributedFloat32VectorGenerator(1, dimension)[0]
				return args{
					indexID:     "test",
					indexVector: vector,
					req: &payload.Update_Request{
						Vector: &payload.Object_Vector{
							Id:     "uuid-2", // the first uuid is overwritten, so use the second one
							Vector: vector,
						},
						Config: &payload.Update_Config{
							SkipStrictExistCheck: false,
						},
					},
				}
			}(),
			want: want{
				wantUUID: "uuid-2",
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			ctx := tt.Context()

			if test.beforeFunc == nil {
				test.beforeFunc = beforeFunc(tt, ctx, tc.fields.objectType)
			}
			s, err := test.beforeFunc(tt, test.args)
			if err != nil {
				tt.Errorf("error = %v", err)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			gotRes, err := s.Update(ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_server_StreamUpdate(t *testing.T) {
// 	type args struct {
// 		stream vald.Update_StreamUpdateServer
// 	}
// 	type fields struct {
// 		UnimplementedAgentServer agent.UnimplementedAgentServer
// 		UnimplementedValdServer  vald.UnimplementedValdServer
// 		ngt                      service.NGT
// 		eg                       errgroup.Group
// 		name                     string
// 		ip                       string
// 		streamConcurrency        int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
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
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           ngt:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           ngt:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 			s := &server{
// 				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
// 				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
// 				ngt:                      test.fields.ngt,
// 				eg:                       test.fields.eg,
// 				name:                     test.fields.name,
// 				ip:                       test.fields.ip,
// 				streamConcurrency:        test.fields.streamConcurrency,
// 			}
//
// 			err := s.StreamUpdate(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiUpdate(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Update_MultiRequest
// 	}
// 	type fields struct {
// 		UnimplementedAgentServer agent.UnimplementedAgentServer
// 		UnimplementedValdServer  vald.UnimplementedValdServer
// 		ngt                      service.NGT
// 		eg                       errgroup.Group
// 		name                     string
// 		ip                       string
// 		streamConcurrency        int
// 	}
// 	type want struct {
// 		wantRes *payload.Object_Locations
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           ngt:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           ngt:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 			s := &server{
// 				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
// 				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
// 				ngt:                      test.fields.ngt,
// 				eg:                       test.fields.eg,
// 				name:                     test.fields.name,
// 				ip:                       test.fields.ip,
// 				streamConcurrency:        test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.MultiUpdate(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_UpdateTimestamp(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Update_TimestampRequest
// 	}
// 	type fields struct {
// 		UnimplementedAgentServer agent.UnimplementedAgentServer
// 		UnimplementedValdServer  vald.UnimplementedValdServer
// 		ngt                      service.NGT
// 		eg                       errgroup.Group
// 		name                     string
// 		ip                       string
// 		streamConcurrency        int
// 	}
// 	type want struct {
// 		wantRes *payload.Object_Location
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           ngt:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedAgentServer:nil,
// 		           UnimplementedValdServer:nil,
// 		           ngt:nil,
// 		           eg:nil,
// 		           name:"",
// 		           ip:"",
// 		           streamConcurrency:0,
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
// 			s := &server{
// 				UnimplementedAgentServer: test.fields.UnimplementedAgentServer,
// 				UnimplementedValdServer:  test.fields.UnimplementedValdServer,
// 				ngt:                      test.fields.ngt,
// 				eg:                       test.fields.eg,
// 				name:                     test.fields.name,
// 				ip:                       test.fields.ip,
// 				streamConcurrency:        test.fields.streamConcurrency,
// 			}
//
// 			gotRes, err := s.UpdateTimestamp(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
