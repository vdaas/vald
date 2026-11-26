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

// package portforward provides a persistent port forwarding daemon for Kubernetes services.
package portforward

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	kclient "github.com/vdaas/vald/internal/k8s/client"
	k8s "github.com/vdaas/vald/internal/k8s/resource"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	watch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"k8s.io/streaming/pkg/httpstream"
)

// healthCheckInterval is the polling interval Start uses while waiting for
// the first port forward session to become healthy.
const healthCheckInterval = 100 * time.Millisecond

// echBufferSize is the buffer size of the daemon error channel: one slot for
// each of the two daemon goroutines (endpoints watcher and connection loop)
// so a single pending error from each never blocks.
const echBufferSize = 2

const (
	decimalBase        = 10
	failureModeAll     = "all"
	failureModeAny     = "any"
	failureModeCap     = 2
	localhostHost      = "localhost"
	localhostIPv4      = "127.0.0.1"
	localhostIPv6      = "::1"
	networkClosedError = "use of closed network connection"
)

// Forwarder defines the interface for a persistent port forwarding daemon.
type Forwarder interface {
	// Start launches the port forward daemon and returns an error channel (named "ech")
	// to report runtime issues.
	Start(ctx context.Context) (<-chan error, error)
	// Stop gracefully terminates the port forwarding daemon.
	Stop() error
}

// portForward is the concrete implementation of the Forwarder interface.
// It holds all configuration and state required to run the persistent port forward daemon.
type portForward struct {
	// Client provides access to the Kubernetes API.
	client kclient.Client

	// EndpointsClient is used to watch the Endpoints resource.
	endpointsClient k8s.EndpointClient

	// Backoff settings for the connection loop.
	backoff backoff.Backoff

	// errgroup is used to manage the lifecycle of the daemon goroutines.
	eg errgroup.Group

	// Ports maps local ports to target pod ports.
	ports map[uint16]uint16

	// HTTP client used for SPDY transport.
	httpClient *http.Client

	// cancel cancels the overall port forward daemon context.
	cancel context.CancelFunc

	// ech is the error channel used to report errors during runtime.
	ech chan error

	// targetsChanged notifies the connection loop when endpoints pods change.
	targetsChanged chan struct{}

	// dialSession opens a SPDY portforward connection to a pod (overrideable in tests).
	dialSession func(ctx context.Context, podName string) (httpstream.Connection, error)

	// Namespace where the service and pods reside.
	namespace string
	// ServiceName is the target service name used to fetch endpoints.
	serviceName string
	// Addresses are the local bind addresses for the port forward.
	addresses []string

	// targets holds the current list of available pod names (extracted from Endpoints).
	targets []string

	// boundAddrs holds the currently active listener addresses (host:port).
	boundAddrs []string

	// mu protects access to targets and boundAddrs.
	mu sync.RWMutex

	// stopOnce makes Stop idempotent: the error channel must be closed at
	// most once and Stop may be called before Start.
	stopOnce sync.Once

	spdyMu       sync.Mutex
	spdyRT       http.RoundTripper
	spdyUpgrader spdy.Upgrader

	healthy atomic.Bool

	// current is used for efficient round-robin selection.
	current atomic.Uint64 // using atomic operations for concurrent safety
}

const streamCopyBufSize = 32 * 1024

var copyBufPool = sync.Pool[*[]byte]{ //nolint:gochecknoglobals // stream copy buffers are reused across forwarded connections
	New: func() *[]byte {
		buf := make([]byte, streamCopyBufSize)
		return &buf
	},
}

type (
	listenAddress struct {
		address     string
		protocol    string
		failureMode string
	}

	portPair struct {
		local  uint16
		remote uint16
	}

	podSession struct {
		streamConn httpstream.Connection
		podName    string
		requestID  atomic.Uint64
	}
)

// New creates a Forwarder instance with default backoff settings.
func New(opts ...Option) (Forwarder, error) {
	pf := new(portForward)
	for _, opt := range append(defaultOptions, opts...) {
		opt(pf)
	}
	if pf.client == nil {
		return nil, errors.ErrKubernetesClientNotFound
	}
	if pf.namespace == "" {
		return nil, errors.ErrUndefinedNamespace
	}
	if pf.serviceName == "" {
		return nil, errors.ErrUndefinedService
	}
	if len(pf.addresses) == 0 {
		return nil, errors.ErrPortForwardAddressNotFound
	}
	// Normalize a private copy once: WithAddress keeps the caller's slice, and
	// normalizing it in place on every reconnect would reorder the caller's
	// data and, because slices.Compact zeroes the removed tail, leave "" in
	// pf.addresses for the next session when duplicates were given.
	pf.addresses = normalize(slices.Clone(pf.addresses))
	if len(pf.ports) == 0 {
		return nil, errors.ErrPortForwardPortPairNotFound
	}
	pf.endpointsClient = k8s.Endpoints(pf.client, pf.namespace)

	if pf.httpClient == nil {
		pf.httpClient = http.DefaultClient
	}
	pf.targetsChanged = make(chan struct{}, 1)
	return pf, nil
}

// normalize sorts ss in place, removes duplicates, and clips excess capacity,
// returning the canonical slice form shared by every consumer in this package
// (round-robin target selection and the port forward session arguments).
// Sorting before Compact guarantees non-adjacent duplicates are also removed.
func normalize(ss []string) []string {
	slices.Sort(ss)
	return slices.Clip(slices.Compact(ss))
}

// updateTargets safely replaces the current target pod list, resets the
// round-robin counter, and reports whether the sorted pod list changed.
func (pf *portForward) updateTargets(pods []string) bool {
	pods = normalize(pods)
	pf.mu.Lock()
	changed := !slices.Equal(pf.targets, pods)
	pf.targets = pods
	pf.mu.Unlock()
	if changed {
		pf.current.Store(0)
	}
	return changed
}

func (pf *portForward) snapshotTargets() []string {
	pf.mu.RLock()
	defer pf.mu.RUnlock()
	return slices.Clone(pf.targets)
}

func (pf *portForward) selectSession(sessions []*podSession) *podSession {
	idx := pf.current.Add(1)
	return sessions[(idx-1)%uint64(len(sessions))]
}

func (pf *portForward) notifyTargetsChanged() {
	if pf.targetsChanged == nil {
		return
	}
	select {
	case pf.targetsChanged <- struct{}{}:
	default:
	}
}

func (pf *portForward) setBoundAddrs(addrs []string) {
	pf.mu.Lock()
	pf.boundAddrs = addrs
	pf.mu.Unlock()
}

func (pf *portForward) getBoundAddrs() []string {
	pf.mu.RLock()
	defer pf.mu.RUnlock()
	return slices.Clone(pf.boundAddrs)
}

// Start launches the port forwarding daemon.
// It starts two goroutines:
//  1. runEndpointWatcher: continuously monitors Endpoints and updates the target pod list.
//  2. runConnectionLoop: repeatedly establishes port forwarding using round-robin selection with exponential backoff.
//
// skipcq: GO-R1005
func (pf *portForward) Start(ctx context.Context) (<-chan error, error) {
	// Create a cancelable context for the entire daemon.
	ctx, pf.cancel = context.WithCancel(ctx)

	// Initialize the error channel (named "ech").
	pf.ech = make(chan error, echBufferSize)
	if pf.targetsChanged == nil {
		pf.targetsChanged = make(chan struct{}, 1)
	}

	if pf.eg == nil {
		pf.eg, ctx = errgroup.New(ctx)
	}

	// Perform an initial update of the target pod list.
	pf.loadTargets(ctx)

	// Start the endpoints watcher goroutine.
	pf.eg.Go(safety.RecoverFunc(func() (err error) {
		// Create a watch on the Endpoints resource using the short syntax.
		watcher, err := pf.endpointsWatcher(ctx)
		if err != nil {
			return err
		}
		// The closure reads the watcher variable at defer-execution time so
		// the latest (possibly restarted) watcher is stopped, not the first
		// one. watcher is only reassigned by this goroutine, so no race.
		defer func() {
			if watcher != nil {
				watcher.Stop()
			}
		}()

		// Process events from the watcher channel.
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case ev, ok := <-watcher.ResultChan():
				if !ok {
					watcher.Stop()
					log.Error("endpoints watcher channel closed, restarting watcher")
					watcher, err = pf.endpointsWatcher(ctx)
					if err != nil {
						pf.notifyError(ctx, err)
						return err
					}
				} else {
					var changed bool
					if ep, isEP := ev.Object.(*v1.Endpoints); isEP && ep != nil { //nolint:staticcheck // corev1.Endpoints matches EndpointsClient watch type
						if pods := extractPodsFromEndpoints(ep); len(pods) > 0 {
							changed = pf.updateTargets(pods)
						} else if ev.Type == watch.Deleted {
							changed = pf.loadTargets(ctx)
						}
					} else {
						changed = pf.loadTargets(ctx)
					}
					if changed {
						// On target list changes, notify the connection loop to refresh SPDY sessions.
						pf.notifyTargetsChanged()
					}
				}
			}
		}
	}))

	pf.eg.Go(safety.RecoverFunc(func() (err error) {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if pf.backoff != nil {
					_, err = pf.backoff.Do(ctx, func(ctx context.Context) (any, bool, error) {
						return nil, true, pf.portForwardToService(ctx)
					})
				} else {
					err = pf.portForwardToService(ctx)
				}
				if err != nil {
					if errors.IsNot(err, context.Canceled, context.DeadlineExceeded) {
						log.Errorf("port forward connection loop ended with error: %v", err)
					}
					pf.notifyError(ctx, err)
				}
			}
		}
	}))

	// Poll the healthy flag on a ticker instead of a busy-wait so waiting for
	// the first established session does not burn a CPU core.
	ticker := time.NewTicker(healthCheckInterval)
	defer ticker.Stop()
	for {
		if pf.healthy.Load() {
			return pf.ech, nil
		}
		select {
		case <-ctx.Done():
			// report why the daemon context ended (e.g. a failed endpoints
			// watch cancels the errgroup with its error), not just Canceled.
			return pf.ech, context.Cause(ctx)
		case <-ticker.C:
		}
	}
}

// notifyError reports err on the error channel without blocking: when nobody
// drains ech and the buffer is full the error is only logged, so the
// reconnection loop never stalls on a full channel.
func (pf *portForward) notifyError(ctx context.Context, err error) {
	select {
	case <-ctx.Done():
	case pf.ech <- err:
	default:
		log.Warnf("port forward error channel is full, dropping error: %v", err)
	}
}

// Stop gracefully terminates the port forwarding daemon by canceling the context
// and waiting for all goroutines to finish. Stop is idempotent and safe to
// call before Start.
func (pf *portForward) Stop() (err error) {
	pf.stopOnce.Do(func() {
		if pf.cancel != nil {
			pf.cancel()
		}
		if pf.eg != nil {
			err = pf.eg.Wait()
		}
		pf.spdyMu.Lock()
		if pf.spdyRT != nil {
			closeIdleTransport(pf.spdyRT)
			pf.spdyRT = nil
			pf.spdyUpgrader = nil
		}
		pf.spdyMu.Unlock()
		if pf.ech != nil {
			close(pf.ech)
		}
	})
	return err
}

func (pf *portForward) endpointsWatcher(ctx context.Context) (w watch.Interface, err error) {
	w, err = pf.endpointsClient.Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", pf.serviceName),
	})
	if err != nil {
		log.Errorf("failed to watch endpoints for service %s: %v", pf.serviceName, err)
		return nil, err
	}
	return w, err
}

func extractPodsFromEndpoints(endpoints *v1.Endpoints) []string { //nolint:staticcheck // corev1.Endpoints matches EndpointsClient.Get return type
	if endpoints == nil {
		return nil
	}
	pods := make([]string, 0, len(endpoints.Subsets))
	for _, subset := range endpoints.Subsets {
		for _, addr := range subset.Addresses {
			if addr.TargetRef != nil && addr.TargetRef.Kind == "Pod" {
				pods = append(pods, addr.TargetRef.Name)
			}
		}
	}
	return pods
}

// loadTargets retrieves the current Endpoints for the service,
// extracts the associated pod names, and updates the internal targets.
func (pf *portForward) loadTargets(ctx context.Context) bool {
	endpoints, err := pf.endpointsClient.Get(ctx, pf.serviceName, metav1.GetOptions{})
	if err != nil {
		log.Errorf("failed to get endpoints for service %s: %v", pf.serviceName, err)
		return false
	}
	pods := extractPodsFromEndpoints(endpoints)
	if len(pods) == 0 {
		log.Errorf("no pods found in endpoints for service %s", pf.serviceName)
		return false
	}
	return pf.updateTargets(pods)
}

func parseListenAddresses(addrs []string) ([]listenAddress, error) {
	if len(addrs) == 0 {
		return nil, errors.ErrPortForwardAddressNotFound
	}
	parsed := make(map[string]listenAddress, len(addrs))
	for _, addr := range addrs {
		if addr == localhostHost {
			if _, ok := parsed[localhostIPv4]; !ok {
				parsed[localhostIPv4] = listenAddress{address: localhostIPv4, protocol: "tcp4", failureMode: failureModeAll}
			}
			if _, ok := parsed[localhostIPv6]; !ok {
				parsed[localhostIPv6] = listenAddress{address: localhostIPv6, protocol: "tcp6", failureMode: failureModeAll}
			}
			continue
		}
		ip := net.ParseIP(addr)
		if ip == nil {
			return nil, errors.Errorf("%s is not a valid IP", addr)
		}
		if ip.To4() != nil {
			parsed[addr] = listenAddress{address: addr, protocol: "tcp4", failureMode: failureModeAny}
		} else {
			parsed[addr] = listenAddress{address: addr, protocol: "tcp6", failureMode: failureModeAny}
		}
	}
	res := make([]listenAddress, 0, len(parsed))
	for _, v := range parsed {
		res = append(res, v)
	}
	slices.SortFunc(res, func(a, b listenAddress) int {
		return strings.Compare(a.address, b.address)
	})
	return res, nil
}

func parsePortPairs(ports map[uint16]uint16) ([]portPair, error) {
	if len(ports) == 0 {
		return nil, errors.ErrPortForwardPortPairNotFound
	}
	pairs := make([]portPair, 0, len(ports))
	for local, remote := range ports {
		if remote == 0 {
			return nil, errors.New("remote port must be > 0")
		}
		pairs = append(pairs, portPair{local: local, remote: remote})
	}
	slices.SortFunc(pairs, func(a, b portPair) int {
		return cmp.Or(cmp.Compare(a.local, b.local), cmp.Compare(a.remote, b.remote))
	})
	return pairs, nil
}

type ctxRoundTripper struct {
	ctx context.Context //nolint:containedctx // scopes SPDY dial requests to caller context
	rt  http.RoundTripper
}

func (c ctxRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return c.rt.RoundTrip(req.WithContext(c.ctx))
}

func isClosedNetErr(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), networkClosedError)
}

func closeIdleTransport(rt http.RoundTripper) {
	if ic, ok := rt.(interface{ CloseIdleConnections() }); ok {
		ic.CloseIdleConnections()
	}
}

func closeListeners(bound []boundListener) {
	for _, b := range bound {
		if cerr := b.listener.Close(); cerr != nil && !isClosedNetErr(cerr) {
			log.Debugf("failed to close portforward listener %s: %v", b.listener.Addr(), cerr)
		}
	}
}

func closeSessions(sessions []*podSession) {
	for _, s := range sessions {
		if cerr := s.streamConn.Close(); cerr != nil && !isClosedNetErr(cerr) {
			log.Debugf("failed to close SPDY session for pod %s: %v", s.podName, cerr)
		}
	}
}

func (pf *portForward) getSPDY() (http.RoundTripper, spdy.Upgrader, error) {
	pf.spdyMu.Lock()
	defer pf.spdyMu.Unlock()
	if pf.spdyRT != nil && pf.spdyUpgrader != nil {
		return pf.spdyRT, pf.spdyUpgrader, nil
	}
	transport, upgrader, err := spdy.RoundTripperFor(pf.client.GetRESTConfig())
	if err != nil {
		return nil, nil, err
	}
	pf.spdyRT = transport
	pf.spdyUpgrader = upgrader
	return transport, upgrader, nil
}

func (pf *portForward) dialPodSPDY(
	ctx context.Context, podName string,
) (httpstream.Connection, error) {
	if pf.client == nil {
		return nil, errors.ErrKubernetesClientNotFound
	}
	transport, upgrader, err := pf.getSPDY()
	if err != nil {
		return nil, err
	}
	sc := &http.Client{Transport: ctxRoundTripper{ctx: ctx, rt: transport}}
	if pf.httpClient != nil {
		sc.CheckRedirect = pf.httpClient.CheckRedirect
		sc.Jar = pf.httpClient.Jar
		sc.Timeout = pf.httpClient.Timeout
	}
	url := pf.client.GetClientSet().CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(pf.namespace).
		Name(podName).
		SubResource("portforward").URL()
	streamConn, protocol, err := spdy.NewDialerForStreaming(upgrader, sc, http.MethodPost, url).Dial(portforward.PortForwardProtocolV1Name)
	if err != nil {
		return nil, err
	}
	if protocol != portforward.PortForwardProtocolV1Name {
		if cerr := streamConn.Close(); cerr != nil && !isClosedNetErr(cerr) {
			log.Debugf("failed to close SPDY stream connection for pod %s: %v", podName, cerr)
		}
		return nil, errors.Errorf("unable to negotiate protocol: client supports %q, server returned %q", portforward.PortForwardProtocolV1Name, protocol)
	}
	if err := ctx.Err(); err != nil {
		if cerr := streamConn.Close(); cerr != nil && !isClosedNetErr(cerr) {
			log.Debugf("failed to close SPDY stream connection for pod %s: %v", podName, cerr)
		}
		return nil, err
	}
	return streamConn, nil
}

func (pf *portForward) dialPodSessions(ctx context.Context, pods []string) ([]*podSession, error) {
	dialFn := pf.dialSession
	if dialFn == nil {
		dialFn = pf.dialPodSPDY
	}
	sessions := make([]*podSession, 0, len(pods))
	var errs []error
	for _, pod := range pods {
		if cerr := ctx.Err(); cerr != nil {
			closeSessions(sessions)
			return nil, cerr
		}
		sc, err := dialFn(ctx, pod)
		if err != nil {
			log.Warnf("failed to dial SPDY portforward session to pod %s: %v", pod, err)
			errs = append(errs, err)
			continue
		}
		sessions = append(sessions, &podSession{
			podName:    pod,
			streamConn: sc,
		})
	}
	if len(sessions) == 0 {
		if len(errs) > 0 {
			return nil, errors.Join(errs...)
		}
		return nil, errors.ErrNoAvailablePods
	}
	return sessions, nil
}

type boundListener struct {
	listener   net.Listener
	remotePort uint16
}

func bindListeners(
	ctx context.Context, addrs []listenAddress, pairs []portPair,
) ([]boundListener, []string, error) {
	var lc net.ListenConfig
	bound := make([]boundListener, 0, len(addrs)*len(pairs))
	boundAddrs := make([]string, 0, len(addrs)*len(pairs))
	for _, pp := range pairs {
		failCounters := make(map[string]int, failureModeCap)
		successCounters := make(map[string]int, failureModeCap)
		var listenErrs []error
		for _, addr := range addrs {
			l, err := lc.Listen(ctx, addr.protocol, net.JoinHostPort(addr.address, strconv.FormatUint(uint64(pp.local), decimalBase)))
			if err != nil {
				failCounters[addr.failureMode]++
				listenErrs = append(listenErrs, err)
				continue
			}
			successCounters[addr.failureMode]++
			bound = append(bound, boundListener{
				listener:   l,
				remotePort: pp.remote,
			})
			boundAddrs = append(boundAddrs, l.Addr().String())
		}
		if (successCounters[failureModeAll] == 0 && failCounters[failureModeAll] > 0) || failCounters[failureModeAny] > 0 {
			closeListeners(bound)
			return nil, nil, errors.Join(listenErrs...)
		}
	}
	if len(bound) == 0 {
		return nil, nil, errors.New("unable to listen on any of the requested ports")
	}
	return bound, boundAddrs, nil
}

func handleForwardedConn(ctx context.Context, conn net.Conn, sess *podSession, remotePort uint16) {
	var closeOnce sync.Once
	closeConn := func() {
		closeOnce.Do(func() {
			if cerr := conn.Close(); cerr != nil && !isClosedNetErr(cerr) {
				log.Debugf("failed to close forwarded local connection for pod %s: %v", sess.podName, cerr)
			}
		})
	}
	defer closeConn()
	stopAfter := context.AfterFunc(ctx, closeConn)
	defer stopAfter()

	reqID := sess.requestID.Add(1) - 1
	headers := http.Header{}
	headers.Set(v1.StreamType, v1.StreamTypeError)
	headers.Set(v1.PortHeader, strconv.FormatUint(uint64(remotePort), decimalBase))
	headers.Set(v1.PortForwardRequestIDHeader, strconv.FormatUint(reqID, decimalBase))
	errorStream, err := sess.streamConn.CreateStream(headers)
	if err != nil {
		log.Errorf("error creating error stream for pod %s: %v", sess.podName, err)
		return
	}
	if cerr := errorStream.Close(); cerr != nil && !isClosedNetErr(cerr) {
		log.Debugf("failed to close error stream write side for pod %s: %v", sess.podName, cerr)
	}
	defer sess.streamConn.RemoveStreams(errorStream)

	type readAllResult struct {
		err     error
		message []byte
	}
	errorChan := make(chan readAllResult, 1)
	var copyWG sync.WaitGroup
	copyWG.Go(func() {
		message, rerr := io.ReadAll(errorStream)
		errorChan <- readAllResult{
			message: message,
			err:     rerr,
		}
	})

	headers.Set(v1.StreamType, v1.StreamTypeData)
	dataStream, err := sess.streamConn.CreateStream(headers)
	if err != nil {
		log.Errorf("error creating data stream for pod %s: %v", sess.podName, err)
		if rerr := errorStream.Reset(); rerr != nil && !isClosedNetErr(rerr) {
			log.Debugf("failed to reset error stream for pod %s: %v", sess.podName, rerr)
		}
		copyWG.Wait()
		return
	}
	defer sess.streamConn.RemoveStreams(dataStream)

	localError := make(chan struct{})
	remoteDone := make(chan struct{})

	copyWG.Go(func() {
		defer close(remoteDone)
		buf := copyBufPool.Get()
		defer copyBufPool.Put(buf)
		if _, cerr := io.CopyBuffer(conn, dataStream, *buf); cerr != nil && ctx.Err() == nil && !isClosedNetErr(cerr) {
			log.Errorf("error copying from remote stream to local connection for pod %s: %v", sess.podName, cerr)
		}
	})

	copyWG.Go(func() {
		defer func() {
			if cerr := dataStream.Close(); cerr != nil && !isClosedNetErr(cerr) {
				log.Debugf("failed to close data stream for pod %s: %v", sess.podName, cerr)
			}
		}()
		buf := copyBufPool.Get()
		defer copyBufPool.Put(buf)
		if _, cerr := io.CopyBuffer(dataStream, conn, *buf); cerr != nil && ctx.Err() == nil && !isClosedNetErr(cerr) {
			log.Errorf("error copying from local connection to remote stream for pod %s: %v", sess.podName, cerr)
			close(localError)
		}
	})

	resetErrorStream := false
	select {
	case <-remoteDone:
	case <-localError:
		resetErrorStream = true
	case <-ctx.Done():
		resetErrorStream = true
	}

	if rerr := dataStream.Reset(); rerr != nil && !isClosedNetErr(rerr) {
		log.Debugf("failed to reset data stream for pod %s: %v", sess.podName, rerr)
	}
	if resetErrorStream {
		if rerr := errorStream.Reset(); rerr != nil && !isClosedNetErr(rerr) {
			log.Debugf("failed to reset error stream for pod %s: %v", sess.podName, rerr)
		}
	}
	closeConn()
	copyWG.Wait()

	errResult := <-errorChan
	switch {
	case errResult.err != nil:
		if !resetErrorStream && ctx.Err() == nil && !isClosedNetErr(errResult.err) {
			log.Errorf("error reading from error stream for pod %s: %v", sess.podName, errResult.err)
		}
		if cerr := sess.streamConn.Close(); cerr != nil && !isClosedNetErr(cerr) {
			log.Debugf("failed to close SPDY stream connection for pod %s: %v", sess.podName, cerr)
		}
	case len(errResult.message) > 0:
		log.Errorf("an error occurred forwarding to pod %s: %s", sess.podName, string(errResult.message))
		if cerr := sess.streamConn.Close(); cerr != nil && !isClosedNetErr(cerr) {
			log.Debugf("failed to close SPDY stream connection for pod %s: %v", sess.podName, cerr)
		}
	}
}

func (pf *portForward) portForwardToService(ctx context.Context) (err error) {
	pods := pf.snapshotTargets()
	if len(pods) == 0 {
		log.Errorf("Port forward connection failed: %v", errors.ErrNoAvailablePods)
		return errors.ErrNoAvailablePods
	}

	addrs, err := parseListenAddresses(pf.addresses)
	if err != nil {
		return err
	}
	pairs, err := parsePortPairs(pf.ports)
	if err != nil {
		return err
	}

	sessionCtx, sessionCancel := context.WithCancel(ctx)
	defer sessionCancel()

	log.Infof("Attempting multi-pod port forward to pods %v on %v:%v", pods, pf.addresses, pf.ports)
	sessions, err := pf.dialPodSessions(sessionCtx, pods)
	if err != nil {
		log.Errorf("Failed to establish port forward to pods %v: %v", pods, err)
		return err
	}

	bound, boundAddrs, err := bindListeners(sessionCtx, addrs, pairs)
	if err != nil {
		closeSessions(sessions)
		return err
	}
	pf.setBoundAddrs(boundAddrs)

	var wg sync.WaitGroup
	defer func() {
		pf.healthy.Store(false)
		pf.setBoundAddrs(nil)
		sessionCancel()
		closeListeners(bound)
		closeSessions(sessions)
		wg.Wait()
	}()

	for _, b := range bound {
		bl := b
		wg.Go(func() {
			for {
				conn, aerr := bl.listener.Accept()
				if aerr != nil {
					if sessionCtx.Err() == nil && !isClosedNetErr(aerr) {
						log.Errorf("error accepting connection on %s: %v", bl.listener.Addr(), aerr)
					}
					return
				}
				sess := pf.selectSession(sessions)
				wg.Go(func() {
					handleForwardedConn(sessionCtx, conn, sess, bl.remotePort)
				})
			}
		})
	}

	sessionClosedCh := make(chan string, len(sessions))
	for _, s := range sessions {
		sess := s
		wg.Go(func() {
			select {
			case <-sessionCtx.Done():
			case <-sess.streamConn.CloseChan():
				select {
				case sessionClosedCh <- sess.podName:
				default:
				}
			}
		})
	}

	pf.healthy.Store(true)
	connectedPods := make([]string, 0, len(sessions))
	for _, s := range sessions {
		connectedPods = append(connectedPods, s.podName)
	}
	log.Infof("successfully established port forward to pods %v on %v", connectedPods, boundAddrs)

	select {
	case closedPod := <-sessionClosedCh:
		log.Errorf("Port forward session ended on pod %s", closedPod)
		return errors.Errorf("lost connection to pod %s", closedPod)
	case <-pf.targetsChanged:
		log.Infof("Endpoints targets changed for service %s, refreshing port forward sessions", pf.serviceName)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
