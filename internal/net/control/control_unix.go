//go:build linux && !windows && !wasm && !js && !darwin

//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package control provides network socket option
package control

import (
	"golang.org/x/sys/unix"
)

var SetsockoptInt = unix.SetsockoptInt

const (
	SOL_SOCKET  = unix.SOL_SOCKET
	IPPROTO_TCP = unix.IPPROTO_TCP
	SOL_IP      = unix.SOL_IP
	SOL_IPV6    = unix.SOL_IPV6

	SO_REUSEADDR = unix.SO_REUSEADDR
	SO_REUSEPORT = unix.SO_REUSEPORT
	SO_KEEPALIVE = unix.SO_KEEPALIVE

	TCP_NODELAY      = unix.TCP_NODELAY
	TCP_CORK         = unix.TCP_CORK
	TCP_QUICKACK     = unix.TCP_QUICKACK
	TCP_DEFER_ACCEPT = unix.TCP_DEFER_ACCEPT
	TCP_KEEPINTVL    = unix.TCP_KEEPINTVL
	TCP_KEEPIDLE     = unix.TCP_KEEPIDLE
	// from linux/include/uapi/linux/tcp.h
	TCP_FASTOPEN         = unix.TCP_FASTOPEN
	TCP_FASTOPEN_CONNECT = unix.TCP_FASTOPEN_CONNECT

	IP_TRANSPARENT   = unix.IP_TRANSPARENT
	IPV6_TRANSPARENT = unix.IPV6_TRANSPARENT

	IP_RECVORIGDSTADDR   = unix.IP_RECVORIGDSTADDR
	IPV6_RECVORIGDSTADDR = unix.IPV6_RECVORIGDSTADDR
)
