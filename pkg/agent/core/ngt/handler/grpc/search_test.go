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
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/data/request"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

const searchErrorsTestID = "uuid"

// stubNGT returns a fixed error from every search method; other methods panic through the nil embedded interface.
type stubNGT struct {
	service.NGT
	err error
	dim int
}

func (n *stubNGT) GetDimensionSize() int { return n.dim }

func (n *stubNGT) Search(
	context.Context, []float32, uint32, float32, float32, int32,
) (*payload.Search_Response, error) {
	return nil, n.err
}

func (n *stubNGT) SearchByID(
	context.Context, string, uint32, float32, float32, int32,
) ([]float32, *payload.Search_Response, error) {
	return nil, nil, n.err
}

func (n *stubNGT) LinearSearch(
	context.Context, []float32, uint32,
) (*payload.Search_Response, error) {
	return nil, n.err
}

func (n *stubNGT) LinearSearchByID(
	context.Context, string, uint32,
) ([]float32, *payload.Search_Response, error) {
	return nil, nil, n.err
}

// Test_server_search_concurrentCoreError calls each search handler concurrently with an NGT core error.
// The errors.As target used to be a package-level variable, which is a data race under -race.
func Test_server_search_concurrentCoreError(t *testing.T) {
	t.Parallel()

	const callers = 32
	vec := []float32{0, 1}
	tests := []struct {
		call    func(context.Context, *server) error
		name    string
		wantMsg string
	}{
		{
			name: "Search", //nolint:goconst // subtest name equals the handler name, as in the other cases
			call: func(ctx context.Context, s *server) error {
				_, err := s.Search(ctx, &payload.Search_Request{Vector: vec})
				return err
			},
			wantMsg: "Search API failed to process search request due to ngt core process returned error",
		},
		{
			name: "SearchByID",
			call: func(ctx context.Context, s *server) error {
				_, err := s.SearchByID(ctx, &payload.Search_IDRequest{Id: searchErrorsTestID})
				return err
			},
			wantMsg: "SearchByID API failed to process search request due to ngt core process returned error",
		},
		{
			name: "LinearSearch",
			call: func(ctx context.Context, s *server) error {
				_, err := s.LinearSearch(ctx, &payload.Search_Request{Vector: vec})
				return err
			},
			wantMsg: "LinearSearch API failed to process search request due to ngt core process returned error",
		},
		{
			name: "LinearSearchByID",
			call: func(ctx context.Context, s *server) error {
				_, err := s.LinearSearchByID(ctx, &payload.Search_IDRequest{Id: searchErrorsTestID})
				return err
			},
			wantMsg: "LinearSearchByID API failed to process search request due to ngt core process returned error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// the core returns NGTError by value (errors.NewNGTError), which the As target must match.
			s := &server{ngt: &stubNGT{err: errors.NewNGTError("boom"), dim: len(vec)}}

			start := make(chan struct{})
			errs := make(chan error, callers)
			for range callers {
				go func() {
					<-start
					errs <- test.call(context.Background(), s)
				}()
			}
			close(start)

			timeout := time.After(10 * time.Second)
			for range callers {
				select {
				case err := <-errs:
					st, ok := status.FromError(err)
					if !ok {
						t.Errorf("got non-status error: %v", err)
						continue
					}
					if st.Code() != codes.Internal {
						t.Errorf("code = %v, want %v", st.Code(), codes.Internal)
					}
					if st.Message() != test.wantMsg {
						t.Errorf("message = %q, want %q", st.Message(), test.wantMsg)
					}
				case <-timeout:
					t.Fatal("timed out waiting for search callers")
				}
			}
		})
	}
}

// fanOutProbe records how many calls to hold overlap. A call stays in flight
// until more than limit calls overlap (then every call is released at once, so
// an unbounded fan-out is seen without waiting) or a short grace period ends
// (so a bounded fan-out still makes progress).
type fanOutProbe struct {
	over   chan struct{}
	mu     sync.Mutex
	limit  int
	cur    int
	peak   int
	calls  int
	closed bool
}

func newFanOutProbe(limit int) *fanOutProbe {
	return &fanOutProbe{over: make(chan struct{}), limit: limit}
}

func (p *fanOutProbe) hold() {
	p.mu.Lock()
	p.cur++
	p.calls++
	p.peak = max(p.peak, p.cur)
	if p.cur > p.limit && !p.closed {
		p.closed = true
		close(p.over)
	}
	p.mu.Unlock()
	timer := time.NewTimer(20 * time.Millisecond)
	select {
	case <-p.over:
	case <-timer.C:
	}
	timer.Stop()
	p.mu.Lock()
	p.cur--
	p.mu.Unlock()
}

func (p *fanOutProbe) stats() (peak, calls int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.peak, p.calls
}

// probeNGT answers every search through the probe; the vector whose first
// element is 0 and the id "id-0" panic instead when panics is set.
type probeNGT struct {
	p *fanOutProbe
	stubNGT
	panics bool
}

func (n *probeNGT) search(first float32) (*payload.Search_Response, error) {
	n.p.hold()
	if n.panics && first == 0 {
		panic("probe panic")
	}
	return &payload.Search_Response{}, nil
}

func (n *probeNGT) searchByID(uuid string) ([]float32, *payload.Search_Response, error) {
	n.p.hold()
	if n.panics && uuid == "id-0" {
		panic("probe panic")
	}
	return nil, &payload.Search_Response{}, nil
}

func (n *probeNGT) Search(
	_ context.Context, vec []float32, _ uint32, _, _ float32, _ int32,
) (*payload.Search_Response, error) {
	return n.search(vec[0])
}

func (n *probeNGT) SearchByID(
	_ context.Context, uuid string, _ uint32, _, _ float32, _ int32,
) ([]float32, *payload.Search_Response, error) {
	return n.searchByID(uuid)
}

func (n *probeNGT) LinearSearch(
	_ context.Context, vec []float32, _ uint32,
) (*payload.Search_Response, error) {
	return n.search(vec[0])
}

func (n *probeNGT) LinearSearchByID(
	_ context.Context, uuid string, _ uint32,
) ([]float32, *payload.Search_Response, error) {
	return n.searchByID(uuid)
}

// Test_server_Multi_boundedFanOut pins that the Multi* search handlers run at
// most streamConcurrency items at once, still answer every item in request
// order, and report a panicking item as an error instead of dropping it.
func Test_server_Multi_boundedFanOut(t *testing.T) {
	t.Parallel()
	const (
		n     = 6
		limit = 2
	)
	searchReqs := func() *payload.Search_MultiRequest {
		reqs := &payload.Search_MultiRequest{Requests: make([]*payload.Search_Request, 0, n)}
		for i := range n {
			reqs.Requests = append(reqs.Requests, &payload.Search_Request{
				Vector: []float32{float32(i), 1},
				Config: &payload.Search_Config{RequestId: "req-" + strconv.Itoa(i)},
			})
		}
		return reqs
	}
	idReqs := func() *payload.Search_MultiIDRequest {
		reqs := &payload.Search_MultiIDRequest{Requests: make([]*payload.Search_IDRequest, 0, n)}
		for i := range n {
			reqs.Requests = append(reqs.Requests, &payload.Search_IDRequest{
				Id:     "id-" + strconv.Itoa(i),
				Config: &payload.Search_Config{RequestId: "req-" + strconv.Itoa(i)},
			})
		}
		return reqs
	}
	calls := []struct {
		call func(context.Context, *server) (*payload.Search_Responses, error)
		name string
	}{
		{
			name: "MultiSearch",
			call: func(ctx context.Context, s *server) (*payload.Search_Responses, error) {
				return s.MultiSearch(ctx, searchReqs())
			},
		},
		{
			name: "MultiSearchByID",
			call: func(ctx context.Context, s *server) (*payload.Search_Responses, error) {
				return s.MultiSearchByID(ctx, idReqs())
			},
		},
		{
			name: "MultiLinearSearch",
			call: func(ctx context.Context, s *server) (*payload.Search_Responses, error) {
				return s.MultiLinearSearch(ctx, searchReqs())
			},
		},
		{
			name: "MultiLinearSearchByID",
			call: func(ctx context.Context, s *server) (*payload.Search_Responses, error) {
				return s.MultiLinearSearchByID(ctx, idReqs())
			},
		},
	}
	for _, c := range calls {
		for _, panics := range []bool{false, true} {
			t.Run(c.name+"/panics="+strconv.FormatBool(panics), func(t *testing.T) {
				t.Parallel()
				p := newFanOutProbe(limit)
				// eg is the shared group the handlers used to fan out on.
				eg, _ := errgroup.New(t.Context())
				s := &server{
					ngt:               &probeNGT{dim: 2, p: p, panics: panics},
					eg:                eg,
					streamConcurrency: limit,
				}
				type result struct {
					res *payload.Search_Responses
					err error
				}
				done := make(chan result, 1)
				go func() {
					res, err := c.call(t.Context(), s)
					done <- result{res: res, err: err}
				}()
				var got result
				select {
				case got = <-done:
				case <-time.After(10 * time.Second):
					t.Fatal("timed out waiting for the handler")
				}
				if peak, _ := p.stats(); peak > limit {
					t.Errorf("peak concurrency = %d, want <= %d", peak, limit)
				}
				if panics {
					if got.err == nil || !strings.Contains(got.err.Error(), "probe panic") {
						t.Errorf("err = %v, want the recovered panic", got.err)
					}
					// The panicking item must not make the group skip the queued ones.
					if _, calls := p.stats(); calls != n {
						t.Errorf("searched %d items, want all %d", calls, n)
					}
					return
				}
				if got.err != nil {
					t.Fatalf("unexpected error: %v", got.err)
				}
				if _, calls := p.stats(); calls != n {
					t.Errorf("searched %d items, want %d", calls, n)
				}
				if len(got.res.GetResponses()) != n {
					t.Fatalf("got %d responses, want %d", len(got.res.GetResponses()), n)
				}
				for i, r := range got.res.GetResponses() {
					if want := "req-" + strconv.Itoa(i); r.GetRequestId() != want {
						t.Errorf("responses[%d].RequestId = %q, want %q", i, r.GetRequestId(), want)
					}
				}
			})
		}
	}
}

// listNGT lists ids and fails DeleteWithTime for each of them with a per-id
// not-found error; other methods panic through the nil embedded interface.
type listNGT struct {
	stubNGT
	ids []string
}

func (n *listNGT) ListObjectFunc(
	_ context.Context, f func(uuid string, oid uint32, timestamp int64) bool,
) {
	for i, id := range n.ids {
		if !f(id, uint32(i), int64(i)) {
			return
		}
	}
}

func (*listNGT) DeleteWithTime(uuid string, _ int64) error {
	return errors.ErrObjectIDNotFound(uuid)
}

// Test_server_multi_joinedErrors pins the error the agent's Multi* search
// handlers and RemoveByTimestamp return now that per-item failures are
// collected in a slice and joined once after all items are processed: N
// distinct failures give one flat joined error with N members, in a status
// carrying the members' code.
func Test_server_multi_joinedErrors(t *testing.T) {
	t.Parallel()
	const n = 8
	vec := []float32{1, 2, 3}
	ids := make([]string, 0, n)
	for i := range n {
		ids = append(ids, "id-"+strconv.Itoa(i))
	}
	searchReqs := func() *payload.Search_MultiRequest {
		reqs := &payload.Search_MultiRequest{Requests: make([]*payload.Search_Request, 0, n)}
		for _, id := range ids {
			reqs.Requests = append(reqs.Requests, &payload.Search_Request{Vector: vec, Config: &payload.Search_Config{RequestId: id}})
		}
		return reqs
	}
	idReqs := func() *payload.Search_MultiIDRequest {
		reqs := &payload.Search_MultiIDRequest{Requests: make([]*payload.Search_IDRequest, 0, n)}
		for _, id := range ids {
			reqs.Requests = append(reqs.Requests, &payload.Search_IDRequest{Id: id})
		}
		return reqs
	}
	tests := []struct {
		call     func(context.Context, *server) error
		name     string
		wantCode codes.Code
	}{
		{
			name: "MultiSearch",
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiSearch(ctx, searchReqs())
				return err
			},
			wantCode: codes.NotFound,
		},
		{
			name: "MultiSearchByID",
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiSearchByID(ctx, idReqs())
				return err
			},
			wantCode: codes.NotFound,
		},
		{
			name: "MultiLinearSearch",
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiLinearSearch(ctx, searchReqs())
				return err
			},
			wantCode: codes.NotFound,
		},
		{
			name: "MultiLinearSearchByID",
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiLinearSearchByID(ctx, idReqs())
				return err
			},
			wantCode: codes.NotFound,
		},
		{
			name: "RemoveByTimestamp",
			call: func(ctx context.Context, s *server) error {
				_, err := s.RemoveByTimestamp(ctx, &payload.Remove_TimestampRequest{})
				return err
			},
			wantCode: codes.NotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			eg, ectx := errgroup.New(ctx)
			s := &server{
				ngt: &listNGT{err: errors.ErrEmptySearchResult, dim: len(vec), ids: ids},
				eg:  eg,
			}
			done := make(chan error, 1)
			go func() { done <- test.call(ectx, s) }()
			var err error
			select {
			case err = <-done:
			case <-ctx.Done():
				t.Fatal("timed out waiting for the handler")
			}
			if err == nil {
				t.Fatal("returned nil, want the per-item failures")
			}
			if st, ok := status.FromError(err); !ok || st == nil || st.Code() != test.wantCode {
				t.Errorf("status code = %v (ok=%t), want %v: %v", st.Code(), ok, test.wantCode, err)
			}
			var u interface{ Unwrap() []error }
			if !errors.As(err, &u) {
				t.Fatalf("got %T, want a joined error: %v", err, err)
			}
			if got := len(u.Unwrap()); got != n {
				t.Errorf("members = %d, want %d (flat): %v", got, n, err)
			}
			for _, id := range ids {
				if !strings.Contains(err.Error(), " "+id+"'s ") && !strings.Contains(err.Error(), " "+id+" ") {
					t.Errorf("error %q does not mention %q", err.Error(), id)
				}
			}
		})
	}
}

func Test_server_Search(t *testing.T) {
	t.Parallel()

	type args struct {
		req       *payload.Search_Request
		insertNum int
	}
	type fields struct {
		ngtCfg       *config.NGT
		overwriteVec [][]float32
		ngtOpts      []service.Option
		objectType   request.ObjectType
		distribution vector.Distribution
	}
	type want struct {
		resultSize int
		code       codes.Code
	}
	type test struct {
		args       args
		checkFunc  func(want, *payload.Search_Response, error) error
		beforeFunc func(*testing.T, context.Context, fields, args) (Server, error)
		afterFunc  func(args)
		name       string
		fields     fields
		want       want
	}

	const (
		defaultDimensionSize = 32
	)

	defaultInsertConfig := &payload.Insert_Config{
		SkipStrictExistCheck: true,
	}
	defaultBeforeFunc := func(t *testing.T, ctx context.Context, f fields, a args) (Server, error) {
		t.Helper()
		eg, ctx := errgroup.New(ctx)
		ngt, err := newIndexedNGTService(ctx, eg, f.objectType, f.distribution, a.insertNum, defaultInsertConfig, f.ngtCfg, f.ngtOpts, nil, f.overwriteVec)
		if err != nil {
			return nil, err
		}
		s, err := New(WithErrGroup(eg), WithNGT(ngt))
		if err != nil {
			return nil, err
		}
		return s, nil
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				errors.Errorf("got error cannot convert to Status: \"%#v\"", err)
			}
			if st.Code() != w.code {
				return errors.Errorf("got_code: \"%#v\",\n\t\t\t\twant: \"%#v\"", st.Code(), w.code)
			}
		}
		if gotSize := len(gotRes.GetResults()); gotSize != w.resultSize {
			res, _ := gotRes.MarshalJSON()
			return errors.Errorf("got size: \"%#v\",\n\t\t\t\twant size: \"%#v\"\ngotResults: \"%s\"", gotSize, w.resultSize, string(res))
		}
		return nil
	}

	ngtConfig := func(dim int, objectType string) *config.NGT {
		return &config.NGT{
			Dimension:        dim,
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
	}
	defaultSearch_Config := &payload.Search_Config{
		Num:     10,
		Radius:  -1,
		Epsilon: 0.1,
		Timeout: int64(time.Second),
	}
	genSameVecs := func(ot request.ObjectType, n int, dim int) [][]float32 {
		var vecs [][]float32
		var err error
		switch ot {
		case request.Float:
			vecs, err = vector.GenF32Vec(vector.Gaussian, 1, dim)
		case request.Uint8:
			vecs, err = vector.GenUint8Vec(vector.Gaussian, 1, dim)
		}
		if err != nil {
			t.Error(err)
		}

		res := make([][]float32, n)
		for i := range n {
			res[i] = vecs[0]
		}

		return res
	}

	/*
		Search test cases:
		- Equivalence Class Testing
			- case 1.1: success search vector from 1000 vectors (type: uint8)
			- case 1.2: success search vector from 1000 vectors (type: float32)
			- case 2.1: fail search with different dimension vector from 1000 vectors (type: uint8)
			- case 2.2: fail search with different dimension vector from 1000 vectors (type: float32)
		- Boundary Value Testing
			- case 1.1: success search with 0 value (min value) vector from 1000 vectors (type: uint8)
			- case 1.2: success search with +0 value vector from 1000 vectors (type: float32)
			- case 1.3: success search with -0 value vector from 1000 vectors (type: float32)
			- case 2.1: success search with max value vector from 1000 vectors (type: uint8)
			- case 2.2: success search with max value vector from 1000 vectors (type: float32)
			- case 3.1: success search with min value vector from 1000 vectors (type: float32)
			- case 4.1: fail search with NaN value vector from 1000 vectors (type: float32)
			- case 5.1: fail search with Inf value vector from 1000 vectors (type: float32)
			- case 6.1: fail search with -Inf value vector from 1000 vectors (type: float32)
			- case 7.1: fail search with 0 length vector from 1000 vectors (type: uint8)
			- case 7.2: fail search with 0 length vector from 1000 vectors (type: float32)
			- case 8.1: fail search with max dimension vector from 1000 vectors (type: uint8)
			- case 8.2: fail search with max dimension vector from 1000 vectors (type: float32)
			- case 9.1: fail search with nil vector from 1000 vectors (type: uint8)
			- case 9.2: fail search with nil vector from 1000 vectors (type: float32)
		- Decision Table Testing
			- case 1.1: success search with Search_Config.Num=10 from 5 different vectors (type: uint8)
			- case 1.2: success search with Search_Config.Num=10 from 5 different vectors (type: float32)
			- case 2.1: success search with Search_Config.Num=10 from 10 different vectors (type: uint8)
			- case 2.2: success search with Search_Config.Num=10 from 10 different vectors (type: float32)
			- case 3.1: success search with Search_Config.Num=10 from 20 different vectors (type: uint8)
			- case 3.2: success search with Search_Config.Num=10 from 20 different vectors (type: float32)
			- case 4.1: success search with Search_Config.Num=10 from 5 same vectors (type: uint8)
			- case 4.2: success search with Search_Config.Num=10 from 5 same vectors (type: float32)
			- case 5.1: success search with Search_Config.Num=10 from 10 same vectors (type: uint8)
			- case 5.2: success search with Search_Config.Num=10 from 10 same vectors (type: float32)
			- case 6.1: success search with Search_Config.Num=10 from 20 same vectors (type: uint8)
			- case 6.2: success search with Search_Config.Num=10 from 20 same vectors (type: float32)
	*/
	tests := []test{
		// Equivalence Class Testing
		{
			name: "Equivalence Class Testing case 1.1: success search vector (type: uint8)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Equivalence Class Testing case 1.2: success search vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Equivalence Class Testing case 2.1: fail search vector with different dimension (type: uint8)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize+1))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Equivalence Class Testing case 2.2: fail search vector with different dimension (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize+1)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},

		// Boundary Value Testing
		{
			name: "Boundary Value Testing case 1.1: success search with 0 value (min value) vector (type: uint8)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, float32(uint8(0))),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 1.2: success search with +0 value vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, +0.0),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 1.3: success search with -0 value vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, float32(math.Copysign(0, -1.0))),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.1: success search with max value vector (type: uint8)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, float32(math.MaxUint8)),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.2: success search with max value vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, math.MaxFloat32),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.1: success search with min value vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, -math.MaxFloat32),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 4.1: fail search with NaN value vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, float32(math.NaN())),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 5.1: fail search with Inf value vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, float32(math.Inf(+1.0))),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 6.1: fail search with -Inf value vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GenSameValueVec(defaultDimensionSize, float32(math.Inf(-1.0))),
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 7.1: fail search with 0 length vector (type: uint8)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: []float32{},
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, "uint8"),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 7.2: fail search with 0 length vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: []float32{},
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 8.1: fail search with max dimension vector (type: uint8)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, math.MaxInt32>>7))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 8.2: fail search with max dimension vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, math.MaxInt32>>7)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 9.1: fail search with nil vector (type: uint8)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: nil,
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 9.2: fail search with nil vector (type: float32)",
			args: args{
				insertNum: 1000,
				req: &payload.Search_Request{
					Vector: nil,
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},

		// Decision Table Testing
		{
			name: "Decision Table Testing case 1.1: success search with Search_Config.Num=10 from 5 different vectors (type: uint8)",
			args: args{
				insertNum: 5,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 5,
			},
		},
		{
			name: "Decision Table Testing case 1.2: success search with Search_Config.Num=10 from 5 different vectors (type: float32)",
			args: args{
				insertNum: 5,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 5,
			},
		},
		{
			name: "Decision Table Testing case 2.1: success search with Search_Config.Num=10 from 10 different vectors (type: uint8)",
			args: args{
				insertNum: 10,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 2.2: success search with Search_Config.Num=10 from 10 different vectors (type: float32)",
			args: args{
				insertNum: 10,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 3.1: success search with Search_Config.Num=10 from 20 different vectors (type: uint8)",
			args: args{
				insertNum: 20,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 3.2: success search with Search_Config.Num=10 from 20 different vectors (type: float32)",
			args: args{
				insertNum: 20,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Float,
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 4.1: success search with Search_Config.Num=10 from 5 same vectors (type: uint8)",
			args: args{
				insertNum: 5,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				overwriteVec: genSameVecs(request.Uint8, 5, defaultDimensionSize),
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 5,
			},
		},
		{
			name: "Decision Table Testing case 4.2: success search with Search_Config.Num=10 from 5 same vectors (type: float32)",
			args: args{
				insertNum: 5,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				overwriteVec: genSameVecs(request.Float, 5, defaultDimensionSize),
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 5,
			},
		},
		{
			name: "Decision Table Testing case 5.1: success search with Search_Config.Num=10 from 10 same vectors (type: uint8)",
			args: args{
				insertNum: 10,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				overwriteVec: genSameVecs(request.Uint8, 10, defaultDimensionSize),
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 5.2: success search with Search_Config.Num=10 from 10 same vectors (type: float32)",
			args: args{
				insertNum: 10,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				overwriteVec: genSameVecs(request.Float, 10, defaultDimensionSize),
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 6.1: success search with Search_Config.Num=10 from 20 same vectors (type: uint8)",
			args: args{
				insertNum: 20,
				req: &payload.Search_Request{
					Vector: vector.ConvertVectorsUint8ToFloat32(vector.GaussianDistributedUint8VectorGenerator(1, defaultDimensionSize))[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				overwriteVec: genSameVecs(request.Uint8, 20, defaultDimensionSize),
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Uint8.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
		{
			name: "Decision Table Testing case 6.2: success search with Search_Config.Num=10 from 20 same vectors (type: float32)",
			args: args{
				insertNum: 20,
				req: &payload.Search_Request{
					Vector: vector.GaussianDistributedFloat32VectorGenerator(1, defaultDimensionSize)[0],
					Config: defaultSearch_Config,
				},
			},
			fields: fields{
				distribution: vector.Gaussian,
				objectType:   request.Uint8,
				overwriteVec: genSameVecs(request.Float, 20, defaultDimensionSize),
				ngtCfg:       ngtConfig(defaultDimensionSize, ngt.Float.String()),
			},
			want: want{
				resultSize: 10,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			ctx := tt.Context()

			if test.beforeFunc == nil {
				test.beforeFunc = defaultBeforeFunc
			}
			s, err := test.beforeFunc(tt, ctx, test.fields, test.args)
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

			gotRes, err := s.Search(ctx, test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_SearchByID(t *testing.T) {
	t.Parallel()

	type args struct {
		indexID  string
		searchID string
	}
	type want struct {
		resultSize int
		code       codes.Code
	}
	type test struct {
		checkFunc  func(want, *payload.Search_Response, error) error
		beforeFunc func(*testing.T, context.Context, args) (Server, error)
		afterFunc  func(args)
		args       args
		name       string
		want       want
	}
	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				errors.Errorf("got error cannot convert to Status: \"%#v\"", err)
			}
			if st.Code() != w.code {
				return errors.Errorf("got_code: \"%#v\",\n\t\t\t\twant: \"%#v\"", st.Code(), w.code)
			}
		}
		if gotSize := len(gotRes.GetResults()); gotSize != w.resultSize {
			return errors.Errorf("got size: \"%#v\",\n\t\t\t\twant size: \"%#v\"", gotSize, w.resultSize)
		}
		return nil
	}

	const (
		insertNum = 1000
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
	defaultNgtConfig := &config.NGT{
		Dimension:        128,
		DistanceType:     ngt.L2.String(),
		ObjectType:       ngt.Float.String(),
		CreationEdgeSize: 60,
		SearchEdgeSize:   20,
		KVSDB: &config.KVSDB{
			Concurrency: 10,
		},
		VQueue: &config.VQueue{
			InsertBufferPoolSize: 1000,
			DeleteBufferPoolSize: 1000,
		},
		EnableInMemoryMode: true,
	}
	defaultInsertConfig := &payload.Insert_Config{
		SkipStrictExistCheck: true,
	}
	defaultBeforeFunc := func(t *testing.T, ctx context.Context, a args) (Server, error) {
		t.Helper()
		eg, ctx := errgroup.New(ctx)
		ngt, err := newIndexedNGTService(ctx, eg, request.Float, vector.Gaussian, insertNum, defaultInsertConfig, defaultNgtConfig, nil, []string{a.indexID}, nil)
		if err != nil {
			return nil, err
		}
		s, err := New(WithErrGroup(eg), WithNGT(ngt))
		if err != nil {
			return nil, err
		}
		return s, nil
	}
	defaultSearch_Config := &payload.Search_Config{
		Num:     10,
		Radius:  -1,
		Epsilon: 0.1,
		Timeout: 1000000000,
	}

	/*
		SearchByID test cases ( focus on ID(string), only test float32 ):
		- Equivalence Class Testing ( 1000 vectors inserted before a search )
			- case 1.1: success search vector
			- case 2.1: fail search with non-existent ID
		- Boundary Value Testing ( 1000 vectors inserted before a search )
			- case 1.1: fail search with ""
			- case 2.1: success search with ^@
			- case 2.2: success search with ^I
			- case 2.3: success search with ^J
			- case 2.4: success search with ^M
			- case 2.5: success search with ^[
			- case 2.6: success search with ^?
			- case 3.1: success search with utf-8 ID from utf-8 index
			- case 3.2: fail search with utf-8 ID from s-jis index
			- case 3.3: fail search with utf-8 ID from euc-jp index
			- case 3.4: fail search with s-jis ID from utf-8 index
			- case 3.5: success search with s-jis ID from s-jis index
			- case 3.6: fail search with s-jis ID from euc-jp index
			- case 3.4: fail search with euc-jp ID from utf-8 index
			- case 3.5: fail search with euc-jp ID from s-jis index
			- case 3.6: success search with euc-jp ID from euc-jp index
			- case 4.1: success search with 😀
		- Decision Table Testing
		    - NONE
	*/
	tests := []test{
		{
			name: "Equivalence Class Testing case 1.1: success search vector",
			args: args{
				indexID:  "test",
				searchID: "test",
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Equivalence Class Testing case 2.1: fail search with non-existent ID",
			args: args{
				indexID:  "test",
				searchID: "non-existent",
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 1.1: fail search with \"\"",
			args: args{
				indexID:  "test",
				searchID: "",
			},
			want: want{
				resultSize: 0,
				code:       codes.InvalidArgument,
			},
		},
		{
			name: "Boundary Value Testing case 2.1: success search with ^@",
			args: args{
				indexID:  string([]byte{0}),
				searchID: string([]byte{0}),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.2: success search with ^I",
			args: args{
				indexID:  "\t",
				searchID: "\t",
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.3: success search with ^J",
			args: args{
				indexID:  "\n",
				searchID: "\n",
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.4: success search with ^M",
			args: args{
				indexID:  "\r",
				searchID: "\r",
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.5: success search with ^[",
			args: args{
				indexID:  string([]byte{27}),
				searchID: string([]byte{27}),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 2.6: success search with ^?",
			args: args{
				indexID:  string([]byte{127}),
				searchID: string([]byte{127}),
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 3.1: success search with utf-8 ID from utf-8 index",
			args: args{
				indexID:  utf8Str,
				searchID: utf8Str,
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 3.2: fail search with utf-8 ID from s-jis index",
			args: args{
				indexID:  sjisStr,
				searchID: utf8Str,
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.3: fail search with utf-8 ID from euc-jp index",
			args: args{
				indexID:  eucjpStr,
				searchID: utf8Str,
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.4: fail search with s-jis ID from utf-8 index",
			args: args{
				indexID:  utf8Str,
				searchID: sjisStr,
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.5: success search with s-jis ID from s-jis index",
			args: args{
				indexID:  sjisStr,
				searchID: sjisStr,
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 3.6: fail search with s-jis ID from euc-jp index",
			args: args{
				indexID:  eucjpStr,
				searchID: sjisStr,
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.7: fail search with euc-jp ID from utf-8 index",
			args: args{
				indexID:  utf8Str,
				searchID: eucjpStr,
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.8: fail search with euc-jp ID from s-jis index",
			args: args{
				indexID:  sjisStr,
				searchID: eucjpStr,
			},
			want: want{
				resultSize: 0,
				code:       codes.NotFound,
			},
		},
		{
			name: "Boundary Value Testing case 3.9: success search with euc-jp ID from euc-jp index",
			args: args{
				indexID:  eucjpStr,
				searchID: eucjpStr,
			},
			want: want{
				resultSize: int(defaultSearch_Config.GetNum()),
			},
		},
		{
			name: "Boundary Value Testing case 4.1: success search with 😀",
			args: args{
				indexID:  "😀",
				searchID: "😀",
			},
			want: want{
				resultSize: 10,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			ctx := tt.Context()

			if test.beforeFunc == nil {
				test.beforeFunc = defaultBeforeFunc
			}
			s, err := test.beforeFunc(tt, ctx, test.args)
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

			req := &payload.Search_IDRequest{
				Id:     test.args.searchID,
				Config: defaultSearch_Config,
			}
			gotRes, err := s.SearchByID(ctx, req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_server_StreamSearch(t *testing.T) {
// 	type args struct {
// 		stream vald.Search_StreamSearchServer
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
// 			err := s.StreamSearch(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamSearchByID(t *testing.T) {
// 	type args struct {
// 		stream vald.Search_StreamSearchByIDServer
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
// 			err := s.StreamSearchByID(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiSearch(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Search_MultiRequest
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
// 		wantRes *payload.Search_Responses
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Responses, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
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
// 			gotRes, err := s.MultiSearch(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiSearchByID(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Search_MultiIDRequest
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
// 		wantRes *payload.Search_Responses
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Responses, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
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
// 			gotRes, err := s.MultiSearchByID(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
