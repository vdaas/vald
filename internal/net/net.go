//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package net

import (
	"context"
	"math"
	"net"
	"net/netip"
	"strconv"
	"syscall"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"golang.org/x/net/idna"
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
	if str == "" {
		return Unknown
	}
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

// Parse takes a network address string in various formats, such as hostname,
// IP address, or A record, splits it into host and port, and analyzes its properties.
//
// The input addr can be in the following formats:
// - "host:port" (e.g., "example.com:80", "192.0.2.1:443", "[2001:db8::1]:8080")
// - "host" (e.g., "example.com", "192.0.2.1", "[2001:db8::1]")
//
// If the host part is not an IP address literal, DNS resolution is attempted.
// During this process, the hostname is converted to Punycode to mitigate IDN homograph attacks.
//
// Returns:
//   - host: The IP address literal or hostname.
//   - port: The port number. Returns 0 if no port is specified.
//   - isLocal: True if the address is local (loopback, private, or link-local).
//   - isIPv4: True if the address is an IPv4 address.
//   - isIPv6: True if the address is an IPv6 address.
//   - err: Non-nil if an error occurs during parsing or resolution. The error is structured
//     and can be inspected using errors.As.
func Parse(
	ctx context.Context, addr string,
) (host string, port uint16, isLocal, isIPv4, isIPv6 bool, err error) {
	// 1. Split host and port.
	host, port, err = SplitHostPort(addr)
	if err != nil && !errors.Is(err, errors.Errorf("address %s: missing port in address", addr)) {
		log.Warnf("failed to parse addr %s\terror: %v", addr, err)
	}
	if host == "" {
		host = addr
	}
	if port == 0 {
		port = defaultPort
	}

	// 2. Resolve and analyze the host part.
	// Try to parse it as an IP address literal.
	ipAddr, err := netip.ParseAddr(host)
	if err != nil {
		// DNS resolution.
		ips, err := lookupNetIP(ctx, DefaultResolver, host)
		if err != nil || len(ips) == 0 {
			return "", 0, false, false, false, errors.Wrapf(err, "failed to lookup ip for %s", host)
		}
		// If multiple IPs are returned, analyze the first one.
		ipAddr = ips[0]
	}

	// 3. Classify the IP address.
	// Normalize IPv4-mapped IPv6 addresses.
	if ipAddr.Is4In6() {
		ipAddr = ipAddr.Unmap()
	}

	isIPv4 = ipAddr.Is4()
	isIPv6 = ipAddr.Is6()
	isLocal = ipAddr.IsLoopback() || ipAddr.IsPrivate() || ipAddr.IsLinkLocalUnicast() || IsLocal(host)

	return host, port, isLocal, isIPv4, isIPv6, nil
}

func lookupNetIP(ctx context.Context, rsv *net.Resolver, host string) ([]netip.Addr, error) {
	// Mitigate IDN homograph attacks.
	asciiHost, err := idna.Lookup.ToASCII(host)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid idna request addr %s", host)
	}
	if rsv == nil {
		rsv = DefaultResolver
	}
	return rsv.LookupNetIP(ctx, "ip", asciiHost)
}

// DialContext is a wrapper function of the net.Dial function.
func DialContext(ctx context.Context, network, addr string) (conn Conn, err error) {
	if DefaultResolver.Dial == nil {
		return net.Dial(network, addr)
	}
	return DefaultResolver.Dial(ctx, network, addr)
}

// JoinHostPort joins the host/IP address and the port number,.
func JoinHostPort(host string, port uint16) string {
	return net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10))
}

// SplitHostPort splits the address, and return the host/IP address and the port number,
// and any error occurred.
// If it is the loopback address, it will return the loopback address and corresponding port number.
// Falls back to default port and local IP if not provided.
// Example: ":8080" → "127.0.0.1", 8080.
func SplitHostPort(hostport string) (host string, port uint16, err error) {
	if !strings.HasPrefix(hostport, "::") && strings.HasPrefix(hostport, ":") {
		hostport = localIPv4 + hostport
	}
	var portStr string
	host, portStr, err = net.SplitHostPort(hostport)
	if err != nil {
		if host == "" {
			host = hostport
		}
		port = defaultPort
	}
	if len(portStr) > 0 {
		var p uint64
		p, err = strconv.ParseUint(portStr, 10, 16)
		if err != nil || p > math.MaxUint16 {
			port = defaultPort
		} else {
			port = uint16(p)
		}
	}
	if len(host) == 0 {
		host = "localhost"
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

	concurrency := int(rl.Max) / 2

	log.Debugf("starting to scan available ports from %d to %d, concurrency %d", start, end, concurrency)

	eg, egctx := errgroup.New(ctx)
	eg.SetLimit(concurrency)

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

	log.Debugf("finished to scan available ports %v", ports)

	return ports, nil
}

// LoadLocalIP returns local ip address.
func LoadLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Warn(err)
		return ""
	}
	for _, address := range addrs {
		if ipn, ok := address.(*net.IPNet); ok {
			if ip, err := netip.ParsePrefix(ipn.String()); err == nil && ip.IsValid() && ip.Addr().IsLoopback() &&
				(ip.Addr().Is4() || ip.Addr().Is6() || ip.Addr().Is4In6()) {
				return ip.Addr().String()
			}
		}
	}
	return ""
}
