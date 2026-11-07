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
	"slices"
	"strconv"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/config"
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

// Test_errMatcher pins that errMatcher(err)(target) agrees with
// errors.Is(err, target) for message-matched targets, including the Unicode
// case folding of strings.EqualFold.
func Test_errMatcher(t *testing.T) {
	t.Parallel()
	join := errors.Join(errors.New("a"), errors.Wrap(errors.Join(errors.New("b"), errors.New("c")), "w"))
	tests := []struct {
		err    error
		target error
		name   string
	}{
		{name: "leaf of flat join", err: errors.Join(errors.New("a"), errors.New("b")), target: errors.New("b")},
		{name: "missing", err: errors.Join(errors.New("a"), errors.New("b")), target: errors.New("c")},
		{name: "nested leaf", err: join, target: errors.New("c")},
		{name: "wrapper layer", err: join, target: errors.New("w: b\nc")},
		{name: "joined message itself", err: errors.Join(errors.New("a"), errors.New("b")), target: errors.New("a\nb")},
		{name: "ASCII case", err: errors.New("uuid ABC not found"), target: errors.New("UUID abc NOT found")},
		{name: "long s folds to s", err: errors.New("\u017f"), target: errors.New("S")},
		{name: "Kelvin folds to k", err: errors.New("\u212a"), target: errors.New("k")},
		{name: "different length", err: errors.New("ab"), target: errors.New("a")},
		{name: "invalid UTF-8 is the replacement rune", err: errors.New("\xff"), target: errors.New("\ufffd")},
		{name: "non-letter", err: errors.New("[1]"), target: errors.New("{1}")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			want := errors.Is(test.err, test.target)
			if got := errMatcher(test.err)(test.target); got != want {
				t.Errorf("errMatcher(%q)(%q) = %t, errors.Is = %t", test.err, test.target, got, want)
			}
		})
	}
}
