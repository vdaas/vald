// +build darwin,!linux,!windows,!wasm,!js

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

// Package control provides network socket option
package control

import (
	"syscall"
)

var SetsockoptInt = syscall.SetsockoptInt

const (
	SOL_SOCKET  = syscall.SOL_SOCKET
	IPPROTO_TCP = syscall.IPPROTO_TCP
	SOL_IP      = syscall.SOL_IP
	SOL_IPV6    = syscall.SOL_IPV6

	SO_REUSEADDR = syscall.SO_REUSEADDR
	SO_REUSEPORT = syscall.SO_REUSEPORT
	SO_KEEPALIVE = syscall.SO_KEEPALIVE

	TCP_NODELAY      = syscall.TCP_NODELAY
	TCP_CORK         = syscall.TCP_CORK
	TCP_QUICKACK     = syscall.TCP_QUICKACK
	TCP_DEFER_ACCEPT = syscall.TCP_DEFER_ACCEPT
	TCP_KEEPINTVL    = syscall.TCP_KEEPINTVL
	TCP_KEEPIDLE     = syscall.TCP_KEEPIDLE

	IP_TRANSPARENT   = syscall.IP_TRANSPARENT
	IPV6_TRANSPARENT = syscall.IPV6_TRANSPARENT

	IP_RECVORIGDSTADDR   = syscall.IP_RECVORIGDSTADDR
	IPV6_RECVORIGDSTADDR = syscall.IPV6_RECVORIGDSTADDR
)
