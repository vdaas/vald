// +build windows

//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package tcp provides tcp option
package tcp

import (
	"syscall"

	"github.com/vdaas/vald/internal/errors"
	"golang.org/x/sys/windows"
)

// Control controls raw network connection.
// And Control calls c.Control invokes f on the underlying connection's file
// descriptor or handle.
func Control(network, address string, c syscall.RawConn) (err error) {
	cerr := c.Control(func(fd uintptr) {
		err = windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_REUSEADDR, 1)
	})
	if cerr != nil {
		return errors.Wrap(err, cerr.Error())
	}
	return
}
