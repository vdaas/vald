//go:build e2e

//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// package portforward provides a persistent port forwarding daemon for Kubernetes services.
package portforward

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
	k8s "github.com/vdaas/vald/tests/v2/e2e/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	watch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
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
	client k8s.Client

	// EndpointsClient is used to watch the Endpoints resource.
	eclient k8s.EndpointClient

	// Backoff settings for the connection loop.
	backoff backoff.Backoff

	// errgroup is used to manage the lifecycle of the daemon goroutines.
	eg errgroup.Group

	// Namespace where the service and pods reside.
	namespace string
	// ServiceName is the target service name used to fetch endpoints.
	serviceName string
	// Addresses are the local bind addresses for the port forward.
	addresses []string
	// Ports maps local ports to target pod ports.
	ports map[uint16]uint16

	// HTTP client used for SPDY transport.
	httpClient *http.Client

	// targets holds the current list of available pod names (extracted from Endpoints).
	targets []string
	// current is used for efficient round-robin selection.
	current atomic.Uint64 // using atomic operations for concurrent safety

	// cancel cancels the overall port forward daemon context.
	cancel context.CancelFunc

	// ech is the error channel used to report errors during runtime.
	ech chan error

	// mu protects access to the targets slice.
	mu sync.RWMutex

	healthy atomic.Bool
}

// NewForwarder creates a new instance of a Forwarder with default backoff settings.
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
	if len(pf.ports) == 0 {
		return nil, errors.ErrPortForwardPortPairNotFound
	}
	pf.eclient = k8s.Endpoints(pf.client, pf.namespace)

	if pf.httpClient == nil {
		pf.httpClient = http.DefaultClient
	}
	return pf, nil
}

// updateTargets safely replaces the current target pod list and resets the round-robin counter.
func (pf *portForward) updateTargets(pods []string) {
	pods = slices.Clip(pods)
	slices.Sort(pods)
	pods = slices.Compact(pods)
	pf.mu.Lock()
	pf.targets = pods
	pf.mu.Unlock()
	pf.current.Store(0)
}

// getNextPod returns a pod name using a counter-based round-robin strategy.
// It uses an atomic counter with a modulus operation to avoid slice rotation.
func (pf *portForward) getNextPod() (string, error) {
	idx := pf.current.Add(1)
	pf.mu.RLock()
	defer pf.mu.RUnlock()
	if len(pf.targets) == 0 {
		return "", fmt.Errorf("no available pods")
	}
	pod := pf.targets[int(idx-1)%len(pf.targets)]
	return pod, nil
}

// Start launches the port forwarding daemon.
// It starts two goroutines:
//  1. runEndpointWatcher: continuously monitors Endpoints and updates the target pod list.
//  2. runConnectionLoop: repeatedly establishes port forwarding using round-robin selection with exponential backoff.
func (pf *portForward) Start(ctx context.Context) (<-chan error, error) {
	// Create a cancelable context for the entire daemon.
	ctx, pf.cancel = context.WithCancel(ctx)

	// Initialize the error channel (named "ech").
	pf.ech = make(chan error, 2)

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
		defer watcher.Stop()

		// Process events from the watcher channel.
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case _, ok := <-watcher.ResultChan():
				if !ok {
					watcher.Stop()
					log.Error("endpoints watcher channel closed, restarting watcher")
					watcher, err = pf.endpointsWatcher(ctx)
					if err != nil {
						select {
						case <-ctx.Done():
							return ctx.Err()
						case pf.ech <- err:
						}
						return err
					}
				} else {
					// On any event, update the target pod list.
					pf.loadTargets(ctx)
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
					select {
					case <-ctx.Done():
						return ctx.Err()
					case pf.ech <- err:
					}
				}
			}
		}
	}))

	for {
		select {
		case <-ctx.Done():
			return pf.ech, ctx.Err()
		default:
			if pf.healthy.Load() {
				return pf.ech, nil
			}
		}
	}
}

// Stop gracefully terminates the port forwarding daemon by canceling the context
// and waiting for all goroutines to finish.
func (pf *portForward) Stop() (err error) {
	if pf.cancel != nil {
		pf.cancel()
	}
	err = pf.eg.Wait()
	close(pf.ech)
	return err
}

func (pf *portForward) endpointsWatcher(ctx context.Context) (w watch.Interface, err error) {
	w, err = pf.eclient.Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", pf.serviceName),
	})
	if err != nil {
		log.Errorf("failed to watch endpoints for service %s: %v", pf.serviceName, err)
		return nil, err
	}
	return w, err
}

// loadTargets retrieves the current Endpoints for the service,
// extracts the associated pod names, and updates the internal targets.
func (pf *portForward) loadTargets(ctx context.Context) {
	endpoints, err := pf.eclient.Get(ctx, pf.serviceName, metav1.GetOptions{})
	if err != nil {
		log.Errorf("failed to get endpoints for service %s: %v", pf.serviceName, err)
		return
	}
	pods := make([]string, 0, len(endpoints.Subsets))
	for _, subset := range endpoints.Subsets {
		for _, addr := range subset.Addresses {
			if addr.TargetRef != nil && addr.TargetRef.Kind == "Pod" {
				pods = append(pods, addr.TargetRef.Name)
			}
		}
	}
	if len(pods) == 0 {
		log.Errorf("no pods found in endpoints for service %s", pf.serviceName)
		return
	}
	pf.updateTargets(pods)
}

// portForwardToServicePod
func (pf *portForward) portForwardToService(ctx context.Context) (err error) {
	// Retrieve the next available pod.
	podName, err := pf.getNextPod()
	if err != nil || podName == "" {
		log.Errorf("Port forward connection failed: %v", err)
		return errors.ErrNoAvailablePods
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	log.Infof("Attempting port forward to pod: %s on %v:%v", podName, pf.addresses, pf.ports)
	// Create an inner context for this port forward session.
	stop, ech, err := PortforwardExtended(ctx, pf.client, pf.namespace, podName, pf.addresses, pf.ports, pf.httpClient)
	if err != nil {
		log.Errorf("Failed to establish port forward to pod %s: %v", podName, err)
		if stop != nil {
			stop()
		}
		return err
	}
	defer stop()

	pf.healthy.Store(true)
	defer pf.healthy.Store(false)

	log.Infof("successfully established port forward to pod %s", podName)

	// Wait for the port forward session to end or the context to be cancelled.
	select {
	case err = <-ech:
		if err != nil {
			log.Errorf("Port forward session ended with error on pod %s: %v", podName, err)
			return err
		}
		log.Infof("Port forward session ended normally on pod %s", podName)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// PortforwardExtended establishes port forwarding for a specific pod.
// It is used internally by the connection loop.
func PortforwardExtended(
	ctx context.Context,
	c k8s.Client,
	namespace, podName string,
	addresses []string,
	ports map[uint16]uint16,
	hc *http.Client,
) (cancel context.CancelFunc, errorChan <-chan error, err error) {
	if c == nil {
		return cancel, nil, errors.ErrKubernetesClientNotFound
	}
	// Create a cancelable context for the port forward session.
	ctx, cancel = context.WithCancel(ctx)

	//tctx, tcancel := context.WithTimeout(ctx, time.Second*30)
	//_, ok, err := k8s.WaitForStatus(tctx, k8s.Pod(c, namespace), podName, k8s.StatusAvailable)
	//tcancel()
	//if !ok || err != nil {
	//	return cancel, nil, errors.Join(err, errors.ErrPodIsNotRunning(namespace, podName))
	//}
	//log.Debugf("pod %s is running", podName)

	if hc == nil {
		hc = http.DefaultClient
	}

	// Set up the SPDY round tripper required for port forwarding.
	transport, upgrader, err := spdy.RoundTripperFor(c.GetRESRConfig())
	if err != nil {
		return cancel, nil, err
	}
	hc.Transport = transport

	if addresses == nil {
		return cancel, nil, errors.ErrPortForwardAddressNotFound
	}

	// Build port pairs in the format "local:target".
	portPairs := make([]string, 0, len(ports))
	for local, target := range ports {
		portPairs = append(portPairs, fmt.Sprintf("%d:%d", local, target))
	}
	if len(portPairs) == 0 {
		return cancel, nil, errors.ErrPortForwardPortPairNotFound
	}
	slices.Sort(portPairs)
	portPairs = slices.Clip(slices.Compact(portPairs))
	slices.Sort(addresses)
	addresses = slices.Clip(slices.Compact(addresses))

	// Create a channel to signal when the port forwarder is ready.
	readyChan := make(chan struct{})

	// Construct the URL for the pod's portforward subresource.
	// Create a new port forwarder instance.
	pf, err := portforward.NewOnAddresses(
		spdy.NewDialer(upgrader, hc, http.MethodPost, c.GetClientSet().CoreV1().RESTClient().Post().
			Resource("pods").
			Namespace(namespace).
			Name(podName).
			SubResource("portforward").URL()),
		addresses, portPairs, ctx.Done(), readyChan, os.Stdout, os.Stderr,
	)
	if err != nil {
		log.Errorf("failed to create port forwarder, addresses: %v, portPairs: %v, error: %v", addresses, portPairs, err)
		return cancel, nil, err
	}

	// Prepare the error channel (named "ech") to report errors.
	ech := make(chan error, 1)
	errgroup.Go(safety.RecoverFunc(func() (err error) {
		defer cancel()
		defer close(ech)
		log.Debugf("port forwarder starting on %v:%v", addresses, portPairs)
		// ForwardPorts blocks until the session ends.
		if err = pf.ForwardPorts(); err != nil {
			select {
			case <-ctx.Done():
			case ech <- err:
			}
		}
		return nil
	}))

	// Wait until the port forwarder signals readiness or context cancellation.
	select {
	case <-ctx.Done():
		return cancel, ech, ctx.Err()
	case <-readyChan:
		log.Debugf("port forwarder ready for pod %s on %v", podName, portPairs)
		return cancel, ech, nil
	}
}
