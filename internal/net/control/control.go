//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

// SocketController represents the socket listener controller.
type SocketController interface {
	GetControl() func(network, addr string, c syscall.RawConn) (err error)
}

type control struct {
	reusePort                bool
	reuseAddr                bool
	tcpFastOpen              bool
	tcpNoDelay               bool
	tcpCork                  bool
	tcpQuickAck              bool
	tcpDeferAccept           bool
	ipTransparent            bool
	ipRecoverDestinationAddr bool
	keepAlive                int
}

// SocketFlag represents the flag to enable specific feature for the socket listener.
type SocketFlag uint

const (
	ReusePort SocketFlag = 1 << iota
	ReuseAddr
	TCPFastOpen
	TCPNoDelay
	TCPCork
	TCPQuickAck
	TCPDeferAccept
	IPTransparent
	IPRecoverDestinationAddr
)

// New returns the socket controller.
func New(flag SocketFlag, keepAlive int) SocketController {
	return &control{
		reusePort:                flag&ReusePort == ReusePort,
		reuseAddr:                flag&ReuseAddr == ReuseAddr,
		tcpFastOpen:              flag&TCPFastOpen == TCPFastOpen,
		tcpNoDelay:               flag&TCPNoDelay == TCPNoDelay,
		tcpCork:                  flag&TCPCork == TCPCork,
		tcpQuickAck:              flag&TCPQuickAck == TCPQuickAck,
		tcpDeferAccept:           flag&TCPDeferAccept == TCPDeferAccept,
		ipTransparent:            flag&IPTransparent == IPTransparent,
		ipRecoverDestinationAddr: flag&IPRecoverDestinationAddr == IPRecoverDestinationAddr,
		keepAlive:                keepAlive,
	}
}

func boolint(b bool) int {
	if b {
		return 1
	}
	return 0
}

func isTCP(network string) bool {
	switch network {
	case "tcp", "tcp4", "tcp6":
		return true
	default:
		return false
	}
}

// GetControl returns the controller function for the socket listener.
func (ctrl *control) GetControl() func(network, addr string, c syscall.RawConn) (err error) {
	if ctrl == nil {
		return func(network, address string, c syscall.RawConn) (err error) { return nil }
	}
	return ctrl.controlFunc
}

func (ctrl *control) controlFunc(network, address string, c syscall.RawConn) (err error) {
	return c.Control(func(fd uintptr) {
		log.Debugf("controlling socket for %s://%s, config %#v", network, address, ctrl)
		f := int(fd)
		var ierr error
		if SO_REUSEPORT != 0 {
			ierr = SetsockoptInt(f, SOL_SOCKET, SO_REUSEPORT, boolint(ctrl.reusePort))
			if ierr != nil {
				err = errors.Wrap(err, ierr.Error())
			}
		}
		if SO_REUSEADDR != 0 {
			ierr = SetsockoptInt(f, SOL_SOCKET, SO_REUSEADDR, boolint(ctrl.reuseAddr))
			if ierr != nil {
				err = errors.Wrap(err, ierr.Error())
			}
		}
		if isTCP(network) {
			if TCP_FASTOPEN != 0 {
				ierr = SetsockoptInt(f, IPPROTO_TCP, TCP_FASTOPEN, boolint(ctrl.tcpFastOpen))
				if ierr != nil {
					err = errors.Wrap(err, ierr.Error())
				}
			}
			if TCP_FASTOPEN_CONNECT != 0 {
				ierr = SetsockoptInt(f, IPPROTO_TCP, TCP_FASTOPEN_CONNECT, boolint(ctrl.tcpFastOpen))
				if ierr != nil {
					err = errors.Wrap(err, ierr.Error())
				}
			}
			if TCP_NODELAY != 0 {
				ierr = SetsockoptInt(f, IPPROTO_TCP, TCP_NODELAY, boolint(ctrl.tcpNoDelay))
				if ierr != nil {
					err = errors.Wrap(err, ierr.Error())
				}
			}
			if TCP_CORK != 0 {
				ierr = SetsockoptInt(f, IPPROTO_TCP, TCP_CORK, boolint(ctrl.tcpCork))
				if ierr != nil {
					err = errors.Wrap(err, ierr.Error())
				}
			}
			if TCP_QUICKACK != 0 {
				ierr = SetsockoptInt(f, IPPROTO_TCP, TCP_QUICKACK, boolint(ctrl.tcpQuickAck))
				if ierr != nil {
					err = errors.Wrap(err, ierr.Error())
				}
			}
			if TCP_DEFER_ACCEPT != 0 {
				ierr = SetsockoptInt(f, IPPROTO_TCP, TCP_DEFER_ACCEPT, boolint(ctrl.tcpDeferAccept))
				if ierr != nil {
					err = errors.Wrap(err, ierr.Error())
				}
			}
		}
		var sol, trans, rda int
		switch network {
		case "tcp", "tcp4", "udp", "udp4":
			sol, trans, rda = SOL_IP, IP_TRANSPARENT, IP_RECVORIGDSTADDR
		case "tcp6", "udp6":
			sol, trans, rda = SOL_IPV6, IPV6_TRANSPARENT, IPV6_RECVORIGDSTADDR
		}
		if sol != 0 && trans != 0 {
			ierr = SetsockoptInt(f, sol, trans, boolint(ctrl.ipTransparent))
			if ierr != nil {
				err = errors.Wrap(err, ierr.Error())
			}
		}
		if sol != 0 && rda != 0 {
			ierr = SetsockoptInt(f, sol, rda, boolint(ctrl.ipRecoverDestinationAddr))
			if ierr != nil {
				err = errors.Wrap(err, ierr.Error())
			}
		}

		if SO_KEEPALIVE != 0 {
			ierr = SetsockoptInt(f, SOL_SOCKET, SO_KEEPALIVE, boolint(ctrl.keepAlive > 0))
			if ierr != nil {
				err = errors.Wrap(err, ierr.Error())
			}
		}
		if ctrl.keepAlive > 0 && isTCP(network) {
			if TCP_KEEPINTVL != 0 {
				ierr = SetsockoptInt(f, IPPROTO_TCP, TCP_KEEPINTVL, ctrl.keepAlive)
				if ierr != nil {
					err = errors.Wrap(err, ierr.Error())
				}
			}
			if TCP_KEEPIDLE != 0 {
				ierr = SetsockoptInt(f, IPPROTO_TCP, TCP_KEEPIDLE, ctrl.keepAlive)
				if ierr != nil {
					err = errors.Wrap(err, ierr.Error())
				}
			}
		}
	})
}
