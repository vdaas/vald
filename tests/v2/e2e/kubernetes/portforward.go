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

// package kubernetes provides kubernetes client functions
package kubernetes

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

// Portforward establishes a port forward from a local port to a target port on a specified pod.
// It is a simple wrapper around PortforwardExtended with default settings.
func Portforward(
	ctx context.Context, c Client, namespace, podName string, localPort, targetPort uint16,
) (cancel context.CancelFunc, ech <-chan error, err error) {
	// Use "127.0.0.1" as the default bind address and http.DefaultClient for HTTP communication.
	return PortforwardExtended(ctx, c, namespace, podName, []string{"127.0.0.1"},
		map[uint16]uint16{localPort: targetPort}, http.DefaultClient)
}

// PortforwardExtended sets up port forwarding for a specific pod with extended configuration.
// It accepts a list of addresses to bind to and a mapping of local ports to target ports.
func PortforwardExtended(
	ctx context.Context,
	c Client,
	namespace, podName string,
	addresses []string,
	ports map[uint16]uint16,
	hc *http.Client,
) (context.CancelFunc, <-chan error, error) {
	// Create a cancelable context to manage the lifetime of the port forward.
	ctx, cancel := context.WithCancel(ctx)

	// Retrieve the pod object using the provided client.
	pclient := Pod(c, namespace)
	if pclient == nil {
		return cancel, nil, errors.ErrKubernetesClientNotFound
	}
	pod, err := pclient.Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return cancel, nil, err
	}

	// Check the pod's current state (e.g., running, pending) to ensure it is available.
	status, info, err := CheckResourceState(pod)
	if err != nil {
		return cancel, nil, err
	}
	log.Debug(info)
	if status != StatusAvailable {
		return cancel, nil, errors.ErrPodIsNotRunning(namespace, podName)
	}

	// Set up the SPDY round tripper, which is required for port forwarding communication.
	transport, upgrader, err := spdy.RoundTripperFor(c.GetRESRConfig())
	if err != nil {
		return cancel, nil, err
	}
	hc.Transport = transport

	// Ensure that bind addresses are provided.
	if addresses == nil {
		return cancel, nil, errors.ErrPortForwardAddressNotFound
	}

	// Build the list of port pairs in the format "localPort:targetPort".
	portPairs := make([]string, 0, len(ports))
	for local, target := range ports {
		portPairs = append(portPairs, fmt.Sprintf("%d:%d", local, target))
	}
	if len(portPairs) == 0 {
		return cancel, nil, errors.ErrPortForwardPortPairNotFound
	}

	// Create a channel that signals when the port forward is ready.
	readyChan := make(chan struct{})

	// Create a new port forwarder that will listen on the specified addresses and port pairs.
	pf, err := portforward.NewOnAddresses(
		spdy.NewDialer(upgrader, hc, http.MethodPost,
			// Construct the URL for the pod's "portforward" subresource.
			// This URL is used by the SPDY dialer to connect to the pod.
			c.GetClientSet().CoreV1().RESTClient().Post().
				Resource("pods").
				Namespace(namespace).
				Name(podName).
				SubResource("portforward").URL()),
		addresses, portPairs, ctx.Done(), readyChan, os.Stdout, os.Stderr,
	)
	if err != nil {
		return cancel, nil, err
	}

	// Prepare an error channel to report any errors during the port forwarding session.
	ech := make(chan error, 1)
	errgroup.Go(safety.RecoverFunc(func() (err error) {
		// Ensure the context is canceled and the error channel is closed when finished.
		defer cancel()
		defer close(ech)
		// Start forwarding ports. This call blocks until the port forward stops.
		if err = pf.ForwardPorts(); err != nil {
			ech <- err
		}
		return nil
	}))

	// Wait until the port forwarding setup is ready or the context is canceled.
	select {
	case <-ctx.Done():
		return cancel, ech, ctx.Err()
	case <-readyChan:
		return cancel, ech, nil
	}
}

// getPodNameFromService retrieves a Pod name from the endpoints of a service.
// It uses the shorthand Endpoints(c, namespace) to get the endpoints.
func getPodNameFromService(ctx context.Context, c Client, namespace, serviceName string) (string, error) {
	// Retrieve the Endpoints for the given service using the alternative syntax.
	endpoints, err := Endpoints(c, namespace).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	// Loop through the subsets and addresses to find a Pod reference.
	for _, subset := range endpoints.Subsets {
		for _, addr := range subset.Addresses {
			if addr.TargetRef != nil && addr.TargetRef.Kind == PodObjectKind {
				return addr.TargetRef.Name, nil
			}
		}
	}
	return "", fmt.Errorf("no pod found associated with service %s", serviceName)
}

// PortforwardWithService continuously monitors a service and attempts to maintain a port forwarding
// connection by re-establishing it to a new Pod if the current connection is lost.
// Note: Since TCP sessions cannot survive a pod replacement, the client must handle reconnection at the application level.
func PortforwardWithService(
	ctx context.Context,
	c Client,
	namespace, serviceName string,
	addresses []string,
	ports map[uint16]uint16,
	hc *http.Client,
) (cancel context.CancelFunc, errCh <-chan error, err error) {
	// Create a cancelable context for the re-connection loop.
	ctx, cancel = context.WithCancel(ctx)
	outerErrCh := make(chan error, 1)

	// Launch a goroutine that maintains the port forwarding connection.
	go func() {
		defer close(outerErrCh)
		for {
			// Attempt to retrieve a valid Pod name from the service endpoints.
			podName, err := getPodNameFromService(ctx, c, namespace, serviceName)
			if err != nil {
				log.Error(err)
				// If retrieval fails, wait 5 seconds before retrying.
				select {
				case <-time.After(5 * time.Second):
					continue
				case <-ctx.Done():
					return
				}
			}
			log.Infof("Selected pod %s for port forwarding", podName)

			stop, ech, err := PortforwardExtended(ctx, c, namespace, podName, addresses, ports, hc)
			if err != nil {
				log.Error(err)
				if stop != nil {
					stop()
				}
				// Wait 5 seconds before retrying if port forwarding setup fails.
				select {
				case <-time.After(5 * time.Second):
					continue
				case <-ctx.Done():
					return
				}
			}

			// Block until the current port forwarding session ends or the context is canceled.
			select {
			case err = <-ech:
				if err != nil {
					log.Errorf("Error during port forwarding to pod %s: %v", podName, err)
				} else {
					log.Infof("Port forwarding to pod %s ended normally", podName)
				}
				stop()
			case <-ctx.Done():
				stop()
				return
			}

			// Wait for 3 seconds before attempting to re-establish the connection.
			select {
			case <-time.After(3 * time.Second):
			case <-ctx.Done():
				return
			}
			log.Info("Attempting to re-establish port forwarding...")
		}
	}()

	return cancel, outerErrCh, nil
}
