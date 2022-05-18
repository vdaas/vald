//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package net provides net functionality for vald's network connection
package net

import (
	"context"
	"math"
	"net"
	"strconv"
	"sync"
	"syscall"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"inet.af/netaddr"
)

type (
	// Addr is an alias of net.Addr.
	Addr = net.Addr

	// Conn is an alias of net.Conn.
	Conn = net.Conn

	// ListenConfig is an alias of net.ListenConfig.
	ListenConfig = net.ListenConfig

	// Listener is an alias of net.Listener.
	Listener = net.Listener

	// Resolver is an alias of net.Resolver.
	Resolver = net.Resolver

	// UDPConn is an alias of net.UDPConn.
	UDPConn = net.UDPConn

	// TCPListener is an alias of net.TCPListener.
	TCPListener = net.TCPListener

	// UnixListener is an alias of net.UnixListener.
	UnixListener = net.UnixListener
)

	// NetworkType represents a network type such as TCP, TCP6, etc.
	NetworkType uint
)

const (
	localIPv4   = "127.0.0.1"
	localIPv6   = "::1"
	localHost   = "localhost"
	defaultPort = 80

	Unknown NetworkType = iota
	UNIX
	UNIXGRAM
	UNIXPACKET
	ICMP
	ICMP6
	IGMP
	TCP
	TCP4
	TCP6
	UDP
	UDP4
	UDP6
)

var (
	// DefaultResolver is an alias of net.DefaultResolver.
	DefaultResolver = net.DefaultResolver

	// Listen is an alias of net.Listen.
	Listen = net.Listen

	// IPv4 is an alias of net.IPv4.
	IPv4 = net.IPv4
)

// NetworkTypeFromString returns the corresponding network type from string.
func NetworkTypeFromString(str string) NetworkType {
	switch strings.ToLower(str) {
	case UNIX.String():
		return UNIX
	case UNIXGRAM.String():
		return UNIXGRAM
	case UNIXPACKET.String():
		return UNIXPACKET
	case TCP.String():
		return TCP
	case TCP4.String():
		return TCP4
	case TCP6.String():
		return TCP6
	case UDP.String():
		return UDP
	case UDP4.String():
		return UDP4
	case UDP6.String():
		return UDP6
	case ICMP.String():
		return ICMP
	case ICMP6.String():
		return ICMP6
	case IGMP.String():
		return IGMP
	}
	return Unknown
}

// String returns the string of the network type.
func (n NetworkType) String() string {
	switch n {
	case UNIX:
		return "unix"
	case UNIXGRAM:
		return "unixgram"
	case UNIXPACKET:
		return "unixpacket"
	case TCP:
		return "tcp"
	case TCP4:
		return "tcp4"
	case TCP6:
		return "tcp6"
	case UDP:
		return "udp"
	case UDP4:
		return "udp4"
	case UDP6:
		return "udp6"
	case ICMP:
		return "icmp"
	case IGMP:
		return "igmp"
	case ICMP6:
		return "ipv6-icmp"
	}
	return "unknown"
}

// IsLocal returns if the host is the localhost address.
func IsLocal(host string) bool {
	return host == localHost ||
		host == localIPv4 ||
		host == localIPv6
}

// IsUDP returns if the network type is the udp or udp4 or udp6.
func IsUDP(network string) bool {
	rip := NetworkTypeFromString(network)
	return rip == UDP ||
		rip == UDP4 ||
		rip == UDP6
}

// IsTCP returns if the network type is the tcp or tcp4 or tcp6.
func IsTCP(network string) bool {
	rip := NetworkTypeFromString(network)
	return rip == TCP ||
		rip == TCP4 ||
		rip == TCP6
}

// Parse parses the hostname, IPv4 or IPv6 address and return the hostname/IP, port number,
// whether the address is local IP and IPv4 or IPv6, and any parsing error occurred.
// The address should contains the port number, otherwise an error will return.
func Parse(addr string) (host string, port uint16, isLocal, isIPv4, isIPv6 bool, err error) {
	host, port, err = SplitHostPort(addr)
	if err != nil {
		log.Warnf("failed to parse addr %s\terror: %v", addr, err)
		host = addr
	}

	ip, nerr := netaddr.ParseIP(host)
	if nerr != nil {
		log.Debugf("host: %s,\tport: %d,\tip: %#v,\terror: %v", host, port, ip, nerr)
	}

	// return host and port and flags
	return host, port,
		// check is local ip or not
		IsLocal(host) || ip.IsLoopback(),
		// check is IPv4 or not
		// ic < 2,
		nerr == nil && ip.Is4(),
		// check is IPv6 or not
		// ic >= 2,
		nerr == nil && (ip.Is6() || ip.Is4in6()),
		// Split error
		err
}

// DialContext is a wrapper function of the net.Dial function.
func DialContext(ctx context.Context, network, addr string) (conn Conn, err error) {
	if DefaultResolver.Dial == nil {
		return net.Dial(network, addr)
	}
	return DefaultResolver.Dial(ctx, network, addr)
}

// JoinHostPort joins the host/IP address and the port number,
func JoinHostPort(host string, port uint16) string {
	return net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10))
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
	p, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil || p > math.MaxUint16 {
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

	var der net.Dialer
	for i := start; i >= start && i <= end; i++ {
		port := i
		eg.Go(safety.RecoverFunc(func() error {
			select {
			case <-egctx.Done():
				return egctx.Err()
			default:
				addr := JoinHostPort(host, port)
				conn, err := der.DialContext(ctx, TCP.String(), addr)
				if err != nil {
					return nil
				}

				mu.Lock()
				ports = append(ports, port)
				mu.Unlock()

				if err = conn.Close(); err != nil {
					log.Warnf("failed to close scan connection for %s, error: %v", addr, err)
				}
				return nil
			}
		}))
	}

	if err = eg.Wait(); err != nil {
		return nil, err
	}

	if len(ports) == 0 {
		return nil, errors.ErrNoPortAvailable(host, start, end)
	}

	return ports, nil
}

// LoadLocalIP returns local ip address
func LoadLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Warn(err)
		return ""
	}
	for _, address := range addrs {
		if ipn, ok := address.(*net.IPNet); ok {
			if ip, ok := netaddr.FromStdIPNet(ipn); ok && ip.Valid() && ip.IP().IsLoopback() &&
				(ip.IP().Is4() || ip.IP().Is6() || ip.IP().Is4in6()) {
				return ip.IP().String()
			}
		}
	}
	return ""
}
