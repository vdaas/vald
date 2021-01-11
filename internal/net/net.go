//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package net provides net functionality for grpc
package net

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

const (
	localIPv4   = "127.0.0.1"
	localIPv6   = "::1"
	localHost   = "localhost"
	defaultPort = 80
)

type (
	// Conn is an alias of net.Conn.
	Conn = net.Conn

	// Dialer is an alias of net.Dialer.
	Dialer = net.Dialer

	// ListenConfig is an alias of net.ListenConfig.
	ListenConfig = net.ListenConfig

	// Listener is an alias of net.Listener.
	Listener = net.Listener

	// Resolver is an alias of net.Resolver.
	Resolver = net.Resolver
)

// DefaultResolver is an alias of net.DefaultResolver.
var DefaultResolver = net.DefaultResolver

// Listen is a wrapper function of the net.Listen function.
func Listen(network, address string) (Listener, error) {
	return net.Listen(network, address)
}

// IsLocal returns if the host is the localhost address.
func IsLocal(host string) bool {
	return host == localHost ||
		host == localIPv4 ||
		host == localIPv6
}

// Dial is a wrapper function of the net.Dial function.
func Dial(network, addr string) (conn Conn, err error) {
	return net.Dial(network, addr)
}

// Parse parses the hostname, IPv4 or IPv6 address and return the hostname/IP, port number,
// whether the address is IP, and any parsing error occurred.
// The address should contains the port number, otherwise an error will return.
func Parse(addr string) (host string, port uint16, isIP bool, err error) {
	host, port, err = SplitHostPort(addr)
	isIP = IsIPv6(host) || IsIPv4(host)
	if err != nil {
		return host, defaultPort, isIP, err
	}
	return host, port, isIP, err
}

// IsIPv6 returns weather the address is IPv6 address.
func IsIPv6(addr string) bool {
	return net.ParseIP(addr) != nil && strings.Count(addr, ":") >= 2
}

// IsIPv4 returns weather the address is IPv4 address.
func IsIPv4(addr string) bool {
	return net.ParseIP(addr) != nil && strings.Count(addr, ":") < 2
}

// SplitHostPort splits the address, and return the host/IP address and the port number,
// and any error occurred.
// If it is the loopback address, it will return the loopback address and corresponding port number.
// IPv6 loopback address is not supported yet.
// For more information, please read https://github.com/vdaas/vald/projects/3#card-43504189
func SplitHostPort(hostport string) (host string, port uint16, err error) {
	if !strings.HasPrefix(hostport, "::") && strings.HasPrefix(hostport, ":") {
		hostport = localIPv4 + hostport
	}
	host, portStr, err := net.SplitHostPort(hostport)
	if err != nil {
		host = hostport
		port = defaultPort
	}
	p, err := strconv.Atoi(portStr)
	if err != nil {
		port = defaultPort
	} else {
		port = uint16(p)
	}
	return host, port, err
}

// ScanPorts scans the given range of port numbers from the host (inclusively),
// and return the list of ports that can be connected through TCP, or any error occurred.
func ScanPorts(ctx context.Context, start, end uint16, host string) (ports []uint16, err error) {
	if start > end {
		start, end = end, start
	}

	var rl syscall.Rlimit
	if err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rl); err != nil {
		return nil, err
	}
	eg, egctx := errgroup.New(ctx)
	eg.Limitation(int(rl.Max) / 2)

	var mu sync.Mutex

	for i := start; i >= start && i <= end; i++ {
		port := i
		eg.Go(func() error {
			select {
			case <-egctx.Done():
				return egctx.Err()
			default:
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
				if err != nil {
					log.Warn(err)
					return nil
				}

				mu.Lock()
				ports = append(ports, port)
				mu.Unlock()

				if err = conn.Close(); err != nil {
					log.Warn(err)
				}
				return nil
			}
		})
	}

	if err = eg.Wait(); err != nil {
		return nil, err
	}

	if len(ports) == 0 {
		return nil, errors.ErrNoPortAvailable(host, start, end)
	}

	return ports, nil
}

func LoadLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Warn(err)
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
