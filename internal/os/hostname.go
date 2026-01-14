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

package os

import (
	"net"
	"os"

	"github.com/vdaas/vald/internal/strings"
)

type (
	DirEntry     = os.DirEntry
	FileInfo     = os.FileInfo
	FileMode     = os.FileMode
	File         = os.File
	LinkError    = os.LinkError
	PathError    = os.PathError
	ProcAttr     = os.ProcAttr
	Process      = os.Process
	ProcessState = os.ProcessState
	Root         = os.Root
	Signal       = os.Signal
	SyscallError = os.SyscallError
)

const (
	unknownHost = "unknown-host"

	O_RDONLY = os.O_RDONLY
	O_WRONLY = os.O_WRONLY
	O_RDWR   = os.O_RDWR
	O_APPEND = os.O_APPEND
	O_CREATE = os.O_CREATE
	O_EXCL   = os.O_EXCL
	O_SYNC   = os.O_SYNC
	O_TRUNC  = os.O_TRUNC

	ModeDir        = os.ModeDir
	ModeAppend     = os.ModeAppend
	ModeExclusive  = os.ModeExclusive
	ModeTemporary  = os.ModeTemporary
	ModeSymlink    = os.ModeSymlink
	ModeDevice     = os.ModeDevice
	ModeNamedPipe  = os.ModeNamedPipe
	ModeSocket     = os.ModeSocket
	ModeSetuid     = os.ModeSetuid
	ModeSetgid     = os.ModeSetgid
	ModeCharDevice = os.ModeCharDevice
	ModeSticky     = os.ModeSticky
	ModeIrregular  = os.ModeIrregular
	ModeType       = os.ModeType
	ModePerm       = os.ModePerm

	PathSeparator     = os.PathSeparator
	PathListSeparator = os.PathListSeparator
	DevNull           = os.DevNull
)

var (
	Chdir           = os.Chdir
	Chmod           = os.Chmod
	Chown           = os.Chown
	Chtimes         = os.Chtimes
	Clearenv        = os.Clearenv
	CopyFS          = os.CopyFS
	Create          = os.Create
	CreateTemp      = os.CreateTemp
	DirFS           = os.DirFS
	Environ         = os.Environ
	Executable      = os.Executable
	Exit            = os.Exit
	ExpandEnv       = os.ExpandEnv
	Expand          = os.Expand
	FindProcess     = os.FindProcess
	Getegid         = os.Getegid
	Getenv          = os.Getenv
	Geteuid         = os.Geteuid
	Getgid          = os.Getgid
	Getgroups       = os.Getgroups
	Getpagesize     = os.Getpagesize
	Getpid          = os.Getpid
	Getppid         = os.Getppid
	Getuid          = os.Getuid
	Getwd           = os.Getwd
	IsExist         = os.IsExist
	IsNotExist      = os.IsNotExist
	IsPathSeparator = os.IsPathSeparator
	IsPermission    = os.IsPermission
	IsTimeout       = os.IsTimeout
	Lchown          = os.Lchown
	Link            = os.Link
	LookupEnv       = os.LookupEnv
	Lstat           = os.Lstat
	MkdirAll        = os.MkdirAll
	Mkdir           = os.Mkdir
	MkdirTemp       = os.MkdirTemp
	NewFile         = os.NewFile
	NewSyscallError = os.NewSyscallError
	OpenFile        = os.OpenFile
	OpenInRoot      = os.OpenInRoot
	Open            = os.Open
	OpenRoot        = os.OpenRoot
	Pipe            = os.Pipe
	ReadDir         = os.ReadDir
	ReadFile        = os.ReadFile
	Readlink        = os.Readlink
	RemoveAll       = os.RemoveAll
	Remove          = os.Remove
	Rename          = os.Rename
	SameFile        = os.SameFile
	Setenv          = os.Setenv
	StartProcess    = os.StartProcess
	Stat            = os.Stat
	Symlink         = os.Symlink
	TempDir         = os.TempDir
	Truncate        = os.Truncate
	Unsetenv        = os.Unsetenv
	UserCacheDir    = os.UserCacheDir
	UserConfigDir   = os.UserConfigDir
	UserHomeDir     = os.UserHomeDir
	WriteFile       = os.WriteFile

	ErrInvalid          = os.ErrInvalid
	ErrPermission       = os.ErrPermission
	ErrExist            = os.ErrExist
	ErrNotExist         = os.ErrNotExist
	ErrClosed           = os.ErrClosed
	ErrNoDeadline       = os.ErrNoDeadline
	ErrDeadlineExceeded = os.ErrDeadlineExceeded
	ErrProcessDone      = os.ErrProcessDone

	Stdin  = os.Stdin
	Stdout = os.Stdout
	Stderr = os.Stderr

	Args = os.Args

	hostname = func() string {
		h, err := os.Hostname()
		if err != nil {
			addrs, err := net.InterfaceAddrs()
			if err != nil {
				return unknownHost
			}
			ips := make([]string, 0, len(addrs))
			for _, addr := range addrs {
				if ipn, ok := addr.(*net.IPNet); ok && !ipn.IP.IsLoopback() {
					ips = append(ips, ipn.IP.String())
				}
			}
			if len(ips) == 0 {
				return unknownHost
			}
			return strings.Join(ips, ",\t")
		}
		return h
	}()
)

func Hostname() (hn string, err error) {
	if hostname != "" {
		return hostname, nil
	}
	return os.Hostname()
}
