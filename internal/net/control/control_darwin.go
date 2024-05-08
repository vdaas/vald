//go:build darwin && !linux && !windows && !wasm && !js

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
	"syscall"

	"golang.org/x/sys/unix"
)

var SetsockoptInt = unix.SetsockoptInt

const (
	SOL_SOCKET  = syscall.SOL_SOCKET
	IPPROTO_TCP = syscall.IPPROTO_TCP
	SOL_IP      = 0
	SOL_IPV6    = 0

	SO_REUSEADDR = unix.SO_REUSEADDR
	SO_REUSEPORT = unix.SO_REUSEPORT
	SO_KEEPALIVE = unix.SO_KEEPALIVE

	TCP_NODELAY          = unix.TCP_NODELAY
	TCP_CORK             = 0
	TCP_QUICKACK         = 0
	TCP_DEFER_ACCEPT     = 0
	TCP_KEEPINTVL        = unix.TCP_KEEPINTVL
	TCP_KEEPIDLE         = 0
	TCP_FASTOPEN         = unix.TCP_FASTOPEN
	TCP_FASTOPEN_CONNECT = 0

	IP_TRANSPARENT   = 0
	IPV6_TRANSPARENT = 0

	IP_RECVORIGDSTADDR   = 0
	IPV6_RECVORIGDSTADDR = 0
)
