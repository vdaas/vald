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

package portforward

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/test/goleak"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	k8stesting "k8s.io/client-go/testing"
	"k8s.io/streaming/pkg/httpstream"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

// mockClientSet wires the client-go fake clientset into the kclient.Client
// interface used by the forwarder. The forwarder only ever calls
// GetClientSet and GetRESTConfig on it, so the embedded crclient.Client is
// left nil: it exists solely to satisfy client.Client's k8s.Client method set
// at compile time.
type mockClientSet struct {
	crclient.Client
	cs  kubernetes.Interface
	cfg *rest.Config
}

func (m *mockClientSet) GetClientSet() kubernetes.Interface {
	return m.cs
}

func (m *mockClientSet) GetRESTConfig() *rest.Config {
	if m.cfg != nil {
		return m.cfg
	}
	return &rest.Config{Host: "https://127.0.0.1:6443"}
}

func newTestForwarder(t *testing.T) Forwarder {
	t.Helper()
	pf, err := New(
		WithClient(&mockClientSet{cs: fake.NewClientset()}),
		WithNamespace("default"),
		WithServiceName("vald-lb-gateway"),
		WithAddress("127.0.0.1"),
		WithPorts(map[uint16]uint16{18081: 8081}),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return pf
}

// TestDialPodSPDY_DoesNotMutateSharedClient guards against the SPDY
// transport being assigned to the caller-provided (possibly process-shared)
// http.Client: the session must build its own client instead.
func TestDialPodSPDY_DoesNotMutateSharedClient(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.NotFoundHandler())
	t.Cleanup(srv.Close)

	cfg := &rest.Config{Host: srv.URL}
	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		t.Fatalf("kubernetes.NewForConfig() error = %v", err)
	}

	shared := &http.Client{Timeout: 3 * time.Second}
	pf := &portForward{
		client:     &mockClientSet{cs: cs, cfg: cfg},
		namespace:  "test-ns",
		httpClient: shared,
	}

	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	sc, err := pf.dialPodSPDY(ctx, "vald-agent-0")
	if sc != nil {
		if cerr := sc.Close(); cerr != nil {
			t.Errorf("sc.Close() error = %v", cerr)
		}
	}
	if err == nil {
		t.Fatal("dialPodSPDY() error = nil, want non-nil for 404 test server")
	}
	if shared.Transport != nil {
		t.Errorf("shared http.Client transport mutated to %T, want nil", shared.Transport)
	}
	if shared.Timeout != 3*time.Second {
		t.Errorf("shared http.Client timeout mutated to %v", shared.Timeout)
	}
}

const (
	podA = "pod-a"
	podB = "pod-b"
	podC = "pod-c"
)

// TestNormalize verifies the shared slice normalization helper used by
// updateTargets and New: the result must be sorted, fully
// deduplicated (including non-adjacent duplicates), and clipped to length.
func TestNormalize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "nil input stays nil",
			in:   nil,
			want: nil,
		},
		{
			name: "empty input stays empty",
			in:   []string{},
			want: []string{},
		},
		{
			name: "unsorted input is sorted",
			in:   []string{podC, podA, podB},
			want: []string{podA, podB, podC},
		},
		{
			name: "adjacent duplicates are removed",
			in:   []string{podA, podA, podB},
			want: []string{podA, podB},
		},
		{
			name: "non-adjacent duplicates are removed",
			in:   []string{podB, podA, podB, podA},
			want: []string{podA, podB},
		},
		{
			name: "excess capacity is clipped",
			in:   append(make([]string, 0, 16), podB, podA, podA),
			want: []string{podA, podB},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := normalize(test.in)
			if !slices.Equal(got, test.want) {
				t.Errorf("normalize() = %v, want %v", got, test.want)
			}
			if cap(got) != len(got) {
				t.Errorf("normalize() cap = %d, want %d (clipped to length)", cap(got), len(got))
			}
		})
	}
}

// TestPortForward_StopBeforeStart guards Stop against nil panics when Start
// was never called.
func TestPortForward_StopBeforeStart(t *testing.T) {
	t.Parallel()

	pf := newTestForwarder(t)
	if err := pf.Stop(); err != nil {
		t.Errorf("Stop() before Start error = %v, want nil", err)
	}
	if err := pf.Stop(); err != nil {
		t.Errorf("second Stop() error = %v, want nil", err)
	}
}

// TestPortForward_StartStopIdempotent verifies that Start returns once the
// context expires (no busy-wait hang) and that Stop is idempotent afterwards
// (no double close of the error channel).
func TestPortForward_StartStopIdempotent(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	pf := newTestForwarder(t)
	ech, err := pf.Start(ctx)
	if ech == nil {
		t.Fatal("Start() ech = nil, want non-nil channel")
	}
	if err == nil {
		t.Error("Start() error = nil, want context error because no pod ever becomes healthy")
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		// First Stop drains the goroutines and may surface the context
		// error from the daemon goroutines; the second must be a no-op.
		if serr := pf.Stop(); serr != nil {
			t.Logf("first Stop() error = %v", serr)
		}
		if serr := pf.Stop(); serr != nil {
			t.Errorf("second Stop() error = %v, want nil", serr)
		}
	}()
	select {
	case <-done:
	case <-time.After(30 * time.Second):
		t.Fatal("Stop() did not return: goroutines leaked or deadlocked")
	}
}

// TestNew_NormalizesPrivateAddressCopy guards against the session code
// normalizing the caller's address slice in place on every reconnect:
// slices.Compact zeroes the removed tail, so duplicate addresses left an empty
// address in pf.addresses that broke every later session. New must normalize
// a private copy once and leave the caller's slice untouched.
func TestNew_NormalizesPrivateAddressCopy(t *testing.T) {
	t.Parallel()

	const v4, v6 = "127.0.0.1", "::1"
	addrs := []string{v4, v6, v4}
	orig := slices.Clone(addrs)
	f, err := New(
		WithClient(&mockClientSet{cs: fake.NewClientset()}),
		WithNamespace("default"),
		WithServiceName("vald-lb-gateway"),
		WithAddress(addrs...),
		WithPorts(map[uint16]uint16{18081: 8081}),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	pf, ok := f.(*portForward)
	if !ok {
		t.Fatalf("New() returned %T, want *portForward", f)
	}
	if want := []string{v4, v6}; !slices.Equal(pf.addresses, want) {
		t.Errorf("pf.addresses = %q, want %q", pf.addresses, want)
	}
	if !slices.Equal(addrs, orig) {
		t.Errorf("caller's address slice mutated to %q, want %q", addrs, orig)
	}
}

// TestPortForward_StartReturnsCause guards against Start reporting a bare
// context.Canceled when the daemon errgroup is cancelled by a goroutine
// failure: the endpoints watch error must be returned as the cause.
func TestPortForward_StartReturnsCause(t *testing.T) {
	t.Parallel()

	watchErr := errors.New("endpoints watch forbidden")
	cs := fake.NewClientset()
	cs.PrependWatchReactor("endpoints", func(k8stesting.Action) (bool, watch.Interface, error) {
		return true, nil, watchErr
	})
	f, err := New(
		WithClient(&mockClientSet{cs: cs}),
		WithNamespace("default"),
		WithServiceName("vald-lb-gateway"),
		WithAddress("127.0.0.1"),
		WithPorts(map[uint16]uint16{18081: 8081}),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	type result struct {
		err error
	}
	res := make(chan result, 1)
	go func() {
		_, serr := f.Start(ctx)
		res <- result{err: serr}
	}()
	select {
	case r := <-res:
		if !errors.Is(r.err, watchErr) {
			t.Errorf("Start() error = %v, want %v", r.err, watchErr)
		}
	case <-time.After(20 * time.Second):
		cancel()
		t.Error("Start() did not return after the endpoints watch failed")
		select {
		case <-res:
		case <-time.After(10 * time.Second):
			t.Fatal("Start() did not return after cancel")
		}
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		if serr := f.Stop(); serr != nil {
			t.Logf("Stop() error = %v", serr)
		}
	}()
	select {
	case <-done:
	case <-time.After(30 * time.Second):
		t.Fatal("Stop() did not return: goroutines leaked or deadlocked")
	}
}

type fakeStream struct {
	conn    net.Conn
	headers http.Header
	onDone  func()
	once    sync.Once
	id      uint32
}

func (s *fakeStream) Read(p []byte) (int, error) {
	return s.conn.Read(p)
}

func (s *fakeStream) Write(p []byte) (int, error) {
	return s.conn.Write(p)
}

func (s *fakeStream) Close() error {
	// For error streams, Close() closes only the client write direction; for data streams, closing triggers EOF.
	if s.headers.Get(v1.StreamType) == v1.StreamTypeError {
		return nil
	}
	return s.Reset()
}

func (s *fakeStream) Reset() error {
	var err error
	s.once.Do(func() {
		err = s.conn.Close()
		if s.onDone != nil {
			s.onDone()
		}
	})
	return err
}

func (s *fakeStream) Headers() http.Header {
	return s.headers
}

func (s *fakeStream) Identifier() uint32 {
	return s.id
}

type fakeStreamConn struct {
	tb        testing.TB
	closeCh   chan bool
	errReads  map[string]net.Conn
	podName   string
	wg        sync.WaitGroup
	closeOnce sync.Once
	mu        sync.Mutex
	nextID    atomic.Uint32
}

func newFakeStreamConn(tb testing.TB, podName string) *fakeStreamConn {
	tb.Helper()
	return &fakeStreamConn{
		tb:       tb,
		podName:  podName,
		closeCh:  make(chan bool),
		errReads: make(map[string]net.Conn),
	}
}

func (c *fakeStreamConn) closeConn(conn net.Conn) {
	c.tb.Helper()
	if err := conn.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
		c.tb.Errorf("Close() error = %v", err)
	}
}

func (c *fakeStreamConn) CreateStream(headers http.Header) (httpstream.Stream, error) {
	select {
	case <-c.closeCh:
		return nil, errors.New("stream connection closed")
	default:
	}
	id := c.nextID.Add(1)
	reqID := headers.Get(v1.PortForwardRequestIDHeader)
	streamType := headers.Get(v1.StreamType)
	clientConn, serverConn := net.Pipe()
	if streamType == v1.StreamTypeError {
		c.mu.Lock()
		c.errReads[reqID] = serverConn
		c.mu.Unlock()
		return &fakeStream{
			headers: headers.Clone(),
			id:      id,
			conn:    clientConn,
			onDone: func() {
				c.closeConn(serverConn)
			},
		}, nil
	}
	c.wg.Go(func() {
		defer func() {
			c.closeConn(serverConn)
			c.mu.Lock()
			if errServer, ok := c.errReads[reqID]; ok {
				delete(c.errReads, reqID)
				c.mu.Unlock()
				c.closeConn(errServer)
			} else {
				c.mu.Unlock()
			}
		}()
		buf := make([]byte, 64)
		if _, err := serverConn.Read(buf); err != nil {
			return
		}
		if _, err := serverConn.Write([]byte("pong:" + c.podName)); err != nil {
			return
		}
	})
	return &fakeStream{
		headers: headers.Clone(),
		id:      id,
		conn:    clientConn,
	}, nil
}

func (c *fakeStreamConn) Close() error {
	c.closeOnce.Do(func() {
		close(c.closeCh)
		c.mu.Lock()
		for k, s := range c.errReads {
			c.closeConn(s)
			delete(c.errReads, k)
		}
		c.mu.Unlock()
	})
	c.wg.Wait()
	return nil
}

func (c *fakeStreamConn) CloseChan() <-chan bool {
	return c.closeCh
}

func (*fakeStreamConn) SetIdleTimeout(time.Duration) {}

func (*fakeStreamConn) RemoveStreams(...httpstream.Stream) {}

// TestPortForward_MultiPodRoundRobin verifies that a single local listener
// multiplexes accepted TCP connections round-robin across all pod SPDY sessions
// and shuts down cleanly without leaking goroutines.
func TestPortForward_MultiPodRoundRobin(t *testing.T) {
	t.Parallel()

	const podKind = "Pod"
	cs := fake.NewClientset(&v1.Endpoints{ //nolint:staticcheck // Endpoints matches portForward's current endpoint client.
		Name:      "vald-lb-gateway",
		Namespace: "default",
		Subsets: []v1.EndpointSubset{ //nolint:staticcheck // EndpointSubset matches portForward's current endpoint client.
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: "10.0.0.1",
						TargetRef: &v1.ObjectReference{
							Kind: podKind,
							Name: "vald-lb-gateway-0",
						},
					},
					{
						IP: "10.0.0.2",
						TargetRef: &v1.ObjectReference{
							Kind: podKind,
							Name: "vald-lb-gateway-1",
						},
					},
					{
						IP: "10.0.0.3",
						TargetRef: &v1.ObjectReference{
							Kind: podKind,
							Name: "vald-lb-gateway-2",
						},
					},
				},
			},
		},
	})

	f, err := New(
		WithClient(&mockClientSet{cs: cs}),
		WithNamespace("default"),
		WithServiceName("vald-lb-gateway"),
		WithAddress("127.0.0.1"),
		WithPorts(map[uint16]uint16{0: 8081}),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	pf := f.(*portForward)
	pf.dialSession = func(_ context.Context, podName string) (httpstream.Connection, error) {
		return newFakeStreamConn(t, podName), nil
	}

	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()

	if _, err := pf.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(func() {
		if serr := pf.Stop(); serr != nil && !errors.Is(serr, context.Canceled) {
			t.Logf("Stop() returned: %v", serr)
		}
	})

	addrs := pf.getBoundAddrs()
	if len(addrs) != 1 {
		t.Fatalf("expected 1 bound addr, got %v", addrs)
	}

	counts := make(map[string]int, 3)
	dialer := &net.Dialer{Timeout: 3 * time.Second}
	for range 9 {
		conn, derr := dialer.DialContext(ctx, "tcp", addrs[0])
		if derr != nil {
			t.Fatalf("DialContext() error = %v", derr)
		}
		if _, werr := conn.Write([]byte("ping")); werr != nil {
			if cerr := conn.Close(); cerr != nil {
				t.Errorf("conn.Close() error = %v", cerr)
			}
			t.Fatalf("Write() error = %v", werr)
		}
		resp, rerr := io.ReadAll(conn)
		if cerr := conn.Close(); cerr != nil {
			t.Errorf("conn.Close() error = %v", cerr)
		}
		if rerr != nil {
			t.Fatalf("ReadAll() error = %v", rerr)
		}
		counts[string(resp)]++
	}

	for _, wantPod := range []string{"pong:vald-lb-gateway-0", "pong:vald-lb-gateway-1", "pong:vald-lb-gateway-2"} {
		if counts[wantPod] != 3 {
			t.Errorf("pod %s hit count = %d, want 3 (counts=%v)", wantPod, counts[wantPod], counts)
		}
	}
}

// TestPortForward_MultiPodPartialDialFailure verifies that if one pod fails to
// establish a SPDY session, the forwarder still becomes healthy using remaining pods.
func TestPortForward_MultiPodPartialDialFailure(t *testing.T) {
	t.Parallel()

	const podKind = "Pod"
	cs := fake.NewClientset(&v1.Endpoints{ //nolint:staticcheck // Endpoints matches portForward's current endpoint client.
		Name:      "vald-lb-gateway",
		Namespace: "default",
		Subsets: []v1.EndpointSubset{ //nolint:staticcheck // EndpointSubset matches portForward's current endpoint client.
			{
				Addresses: []v1.EndpointAddress{
					{TargetRef: &v1.ObjectReference{Kind: podKind, Name: podA}},
					{TargetRef: &v1.ObjectReference{Kind: podKind, Name: podB}},
					{TargetRef: &v1.ObjectReference{Kind: podKind, Name: podC}},
				},
			},
		},
	})

	f, err := New(
		WithClient(&mockClientSet{cs: cs}),
		WithNamespace("default"),
		WithServiceName("vald-lb-gateway"),
		WithAddress("127.0.0.1"),
		WithPorts(map[uint16]uint16{0: 8081}),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	pf := f.(*portForward)
	pf.dialSession = func(_ context.Context, podName string) (httpstream.Connection, error) {
		if podName == podB {
			return nil, errors.New("simulated dial failure")
		}
		return newFakeStreamConn(t, podName), nil
	}

	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()

	if _, err := pf.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(func() {
		if err := pf.Stop(); err != nil {
			t.Errorf("pf.Stop() error = %v", err)
		}
	})

	addrs := pf.getBoundAddrs()
	if len(addrs) != 1 {
		t.Fatalf("expected 1 bound addr, got %v", addrs)
	}

	counts := make(map[string]int, 2)
	dialer := &net.Dialer{Timeout: 3 * time.Second}
	for range 4 {
		conn, derr := dialer.DialContext(ctx, "tcp", addrs[0])
		if derr != nil {
			t.Fatalf("DialContext() error = %v", derr)
		}
		if _, werr := conn.Write([]byte("ping")); werr != nil {
			if cerr := conn.Close(); cerr != nil {
				t.Errorf("conn.Close() error = %v", cerr)
			}
			t.Fatalf("Write() error = %v", werr)
		}
		resp, rerr := io.ReadAll(conn)
		if cerr := conn.Close(); cerr != nil {
			t.Errorf("conn.Close() error = %v", cerr)
		}
		if rerr != nil {
			t.Fatalf("ReadAll() error = %v", rerr)
		}
		counts[string(resp)]++
	}
	if counts["pong:pod-a"] != 2 || counts["pong:pod-c"] != 2 {
		t.Errorf("unexpected counts across surviving pods: %v", counts)
	}
}

func TestPortForward_DialPodSessions_CanceledCtx(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	var dialed []*fakeStreamConn
	pf := &portForward{
		dialSession: func(_ context.Context, podName string) (httpstream.Connection, error) {
			sc := newFakeStreamConn(t, podName)
			dialed = append(dialed, sc)
			cancel()
			return sc, nil
		},
	}

	sessions, err := pf.dialPodSessions(ctx, []string{podA, podB, podC})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("dialPodSessions() err = %v, want context.Canceled", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("dialPodSessions() sessions = %d, want 0", len(sessions))
	}
	if len(dialed) != 1 {
		t.Fatalf("dialSession called %d times, want 1 before ctx cancel", len(dialed))
	}
	select {
	case <-dialed[0].CloseChan():
	default:
		t.Fatal("already-dialed session was not closed when dialPodSessions aborted on canceled ctx")
	}
}

type nonClosingErrorStream struct {
	*fakeStream
}

func (*nonClosingErrorStream) Close() error { return nil }

type blockingErrStreamConn struct {
	tb      testing.TB
	closeCh chan bool
}

func newBlockingErrStreamConn(tb testing.TB) *blockingErrStreamConn {
	tb.Helper()
	return &blockingErrStreamConn{tb: tb, closeCh: make(chan bool)}
}

func (c *blockingErrStreamConn) CreateStream(headers http.Header) (httpstream.Stream, error) {
	clientConn, serverConn := net.Pipe()
	return &nonClosingErrorStream{
		fakeStream: &fakeStream{
			headers: headers.Clone(),
			conn:    clientConn,
			onDone: func() {
				if err := serverConn.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
					c.tb.Errorf("Close() error = %v", err)
				}
			},
		},
	}, nil
}

func (c *blockingErrStreamConn) Close() error {
	select {
	case <-c.closeCh:
	default:
		close(c.closeCh)
	}
	return nil
}

func (c *blockingErrStreamConn) CloseChan() <-chan bool { return c.closeCh }

func (*blockingErrStreamConn) SetIdleTimeout(time.Duration) {}

func (*blockingErrStreamConn) RemoveStreams(...httpstream.Stream) {}

func TestHandleForwardedConn_CancelResetsBlockingErrorStream(t *testing.T) {
	t.Parallel()
	localClient, localServer := net.Pipe()
	t.Cleanup(func() {
		if err := localClient.Close(); err != nil {
			t.Errorf("localClient.Close() error = %v", err)
		}
		if err := localServer.Close(); err != nil {
			t.Errorf("localServer.Close() error = %v", err)
		}
	})

	sess := &podSession{
		podName:    "pod-blocking",
		streamConn: newBlockingErrStreamConn(t),
	}
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan struct{})
	go func() {
		defer close(done)
		handleForwardedConn(ctx, localServer, sess, 8081)
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("handleForwardedConn hung waiting on unclosed errorStream after ctx cancel")
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want Forwarder
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args
// 		want
// 		checkFunc  func(want, Forwarder, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Forwarder, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           opts:nil,
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
// 		           opts:nil,
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
// 			got, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_normalize(t *testing.T) {
// 	type args struct {
// 		ss []string
// 	}
// 	type want struct {
// 		want []string
// 	}
// 	type test struct {
// 		name       string
// 		args
// 		want
// 		checkFunc  func(want, []string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []string) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ss:nil,
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
// 		           ss:nil,
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
// 			got := normalize(test.args.ss)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_portForward_updateTargets(t *testing.T) {
// 	type args struct {
// 		pods []string
// 	}
// 	type fields struct {
// 		client      kclient.Client
// 		endpointsClient     k8s.EndpointClient
// 		backoff     backoff.Backoff
// 		eg          errgroup.Group
// 		ports       map[uint16]uint16
// 		httpClient  *http.Client
// 		cancel      context.CancelFunc
// 		ech         chan error
// 		namespace   string
// 		serviceName string
// 		addresses   []string
// 		targets     []string
// 		healthy     atomic.Bool
// 		current     atomic.Uint64
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args
// 		fields
// 		want
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
// 		           pods:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           endpointsClient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 		           pods:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           endpointsClient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 			pf := &portForward{
// 				client:      test.fields.client,
// 				endpointsClient:     test.fields.endpointsClient,
// 				backoff:     test.fields.backoff,
// 				eg:          test.fields.eg,
// 				ports:       test.fields.ports,
// 				httpClient:  test.fields.httpClient,
// 				cancel:      test.fields.cancel,
// 				ech:         test.fields.ech,
// 				namespace:   test.fields.namespace,
// 				serviceName: test.fields.serviceName,
// 				addresses:   test.fields.addresses,
// 				targets:     test.fields.targets,
// 				healthy:     test.fields.healthy,
// 				current:     test.fields.current,
// 			}
//
// 			pf.updateTargets(test.args.pods)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_portForward_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client      kclient.Client
// 		endpointsClient     k8s.EndpointClient
// 		backoff     backoff.Backoff
// 		eg          errgroup.Group
// 		ports       map[uint16]uint16
// 		httpClient  *http.Client
// 		cancel      context.CancelFunc
// 		ech         chan error
// 		namespace   string
// 		serviceName string
// 		addresses   []string
// 		targets     []string
// 		healthy     atomic.Bool
// 		current     atomic.Uint64
// 	}
// 	type want struct {
// 		want <-chan error
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args
// 		fields
// 		want
// 		checkFunc  func(want, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
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
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           endpointsClient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           endpointsClient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 			pf := &portForward{
// 				client:      test.fields.client,
// 				endpointsClient:     test.fields.endpointsClient,
// 				backoff:     test.fields.backoff,
// 				eg:          test.fields.eg,
// 				ports:       test.fields.ports,
// 				httpClient:  test.fields.httpClient,
// 				cancel:      test.fields.cancel,
// 				ech:         test.fields.ech,
// 				namespace:   test.fields.namespace,
// 				serviceName: test.fields.serviceName,
// 				addresses:   test.fields.addresses,
// 				targets:     test.fields.targets,
// 				healthy:     test.fields.healthy,
// 				current:     test.fields.current,
// 			}
//
// 			got, err := pf.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_portForward_notifyError(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		err error
// 	}
// 	type fields struct {
// 		client      kclient.Client
// 		endpointsClient     k8s.EndpointClient
// 		backoff     backoff.Backoff
// 		eg          errgroup.Group
// 		ports       map[uint16]uint16
// 		httpClient  *http.Client
// 		cancel      context.CancelFunc
// 		ech         chan error
// 		namespace   string
// 		serviceName string
// 		addresses   []string
// 		targets     []string
// 		healthy     atomic.Bool
// 		current     atomic.Uint64
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args
// 		fields
// 		want
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
// 		           err:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           endpointsClient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 		           err:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           endpointsClient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 			pf := &portForward{
// 				client:      test.fields.client,
// 				endpointsClient:     test.fields.endpointsClient,
// 				backoff:     test.fields.backoff,
// 				eg:          test.fields.eg,
// 				ports:       test.fields.ports,
// 				httpClient:  test.fields.httpClient,
// 				cancel:      test.fields.cancel,
// 				ech:         test.fields.ech,
// 				namespace:   test.fields.namespace,
// 				serviceName: test.fields.serviceName,
// 				addresses:   test.fields.addresses,
// 				targets:     test.fields.targets,
// 				healthy:     test.fields.healthy,
// 				current:     test.fields.current,
// 			}
//
// 			pf.notifyError(test.args.ctx, test.args.err)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_portForward_Stop(t *testing.T) {
// 	type fields struct {
// 		client      kclient.Client
// 		endpointsClient     k8s.EndpointClient
// 		backoff     backoff.Backoff
// 		eg          errgroup.Group
// 		ports       map[uint16]uint16
// 		httpClient  *http.Client
// 		cancel      context.CancelFunc
// 		ech         chan error
// 		namespace   string
// 		serviceName string
// 		addresses   []string
// 		targets     []string
// 		healthy     atomic.Bool
// 		current     atomic.Uint64
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		fields
// 		want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
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
// 		       fields: fields {
// 		           client:nil,
// 		           endpointsClient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           client:nil,
// 		           endpointsClient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			pf := &portForward{
// 				client:      test.fields.client,
// 				endpointsClient:     test.fields.endpointsClient,
// 				backoff:     test.fields.backoff,
// 				eg:          test.fields.eg,
// 				ports:       test.fields.ports,
// 				httpClient:  test.fields.httpClient,
// 				cancel:      test.fields.cancel,
// 				ech:         test.fields.ech,
// 				namespace:   test.fields.namespace,
// 				serviceName: test.fields.serviceName,
// 				addresses:   test.fields.addresses,
// 				targets:     test.fields.targets,
// 				healthy:     test.fields.healthy,
// 				current:     test.fields.current,
// 			}
//
// 			err := pf.Stop()
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_portForward_endpointsWatcher(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client      kclient.Client
// 		endpointsClient     k8s.EndpointClient
// 		backoff     backoff.Backoff
// 		eg          errgroup.Group
// 		ports       map[uint16]uint16
// 		httpClient  *http.Client
// 		cancel      context.CancelFunc
// 		ech         chan error
// 		namespace   string
// 		serviceName string
// 		addresses   []string
// 		targets     []string
// 		healthy     atomic.Bool
// 		current     atomic.Uint64
// 	}
// 	type want struct {
// 		wantW watch.Interface
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args
// 		fields
// 		want
// 		checkFunc  func(want, watch.Interface, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotW watch.Interface, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotW, w.wantW) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotW, w.wantW)
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
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           endpointsClient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           endpointsClient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 			pf := &portForward{
// 				client:      test.fields.client,
// 				endpointsClient:     test.fields.endpointsClient,
// 				backoff:     test.fields.backoff,
// 				eg:          test.fields.eg,
// 				ports:       test.fields.ports,
// 				httpClient:  test.fields.httpClient,
// 				cancel:      test.fields.cancel,
// 				ech:         test.fields.ech,
// 				namespace:   test.fields.namespace,
// 				serviceName: test.fields.serviceName,
// 				addresses:   test.fields.addresses,
// 				targets:     test.fields.targets,
// 				healthy:     test.fields.healthy,
// 				current:     test.fields.current,
// 			}
//
// 			gotW, err := pf.endpointsWatcher(test.args.ctx)
// 			if err := checkFunc(test.want, gotW, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_portForward_loadTargets(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client      kclient.Client
// 		endpointsClient     k8s.EndpointClient
// 		backoff     backoff.Backoff
// 		eg          errgroup.Group
// 		ports       map[uint16]uint16
// 		httpClient  *http.Client
// 		cancel      context.CancelFunc
// 		ech         chan error
// 		namespace   string
// 		serviceName string
// 		addresses   []string
// 		targets     []string
// 		healthy     atomic.Bool
// 		current     atomic.Uint64
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args
// 		fields
// 		want
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
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           endpointsClient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           endpointsClient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 			pf := &portForward{
// 				client:      test.fields.client,
// 				endpointsClient:     test.fields.endpointsClient,
// 				backoff:     test.fields.backoff,
// 				eg:          test.fields.eg,
// 				ports:       test.fields.ports,
// 				httpClient:  test.fields.httpClient,
// 				cancel:      test.fields.cancel,
// 				ech:         test.fields.ech,
// 				namespace:   test.fields.namespace,
// 				serviceName: test.fields.serviceName,
// 				addresses:   test.fields.addresses,
// 				targets:     test.fields.targets,
// 				healthy:     test.fields.healthy,
// 				current:     test.fields.current,
// 			}
//
// 			pf.loadTargets(test.args.ctx)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_portForward_portForwardToService(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client      kclient.Client
// 		endpointsClient     k8s.EndpointClient
// 		backoff     backoff.Backoff
// 		eg          errgroup.Group
// 		ports       map[uint16]uint16
// 		httpClient  *http.Client
// 		cancel      context.CancelFunc
// 		ech         chan error
// 		namespace   string
// 		serviceName string
// 		addresses   []string
// 		targets     []string
// 		healthy     atomic.Bool
// 		current     atomic.Uint64
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args
// 		fields
// 		want
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
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           endpointsClient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           endpointsClient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 			pf := &portForward{
// 				client:      test.fields.client,
// 				endpointsClient:     test.fields.endpointsClient,
// 				backoff:     test.fields.backoff,
// 				eg:          test.fields.eg,
// 				ports:       test.fields.ports,
// 				httpClient:  test.fields.httpClient,
// 				cancel:      test.fields.cancel,
// 				ech:         test.fields.ech,
// 				namespace:   test.fields.namespace,
// 				serviceName: test.fields.serviceName,
// 				addresses:   test.fields.addresses,
// 				targets:     test.fields.targets,
// 				healthy:     test.fields.healthy,
// 				current:     test.fields.current,
// 			}
//
// 			err := pf.portForwardToService(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
