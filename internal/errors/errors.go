// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package errors provides error types and function
package errors

import (
	"reflect"
	"runtime"
	"time"

	// "github.com/pkg/errors"
	"github.com/cockroachdb/errors"
)

var (
	ErrInvalidConfigVersion = func(cur, con string) error {
		return Errorf("invalid config version %s not satisfies version constraints %s", cur, con)
	}

	// ErrTLSDisabled is error variable, it's replesents config error that tls is disabled by config
	ErrTLSDisabled = New("tls feature is disabled")

	// ErrTLSCertOrKeyNotFound is error variable, it's replesents tls cert or key not found error
	ErrTLSCertOrKeyNotFound = New("cert or key file path not found")

	ErrCertificationFailed = New("certification failed")

	ErrTimeoutParseFailed = func(timeout string) error {
		return Errorf("invalid timeout value: %s\t:timeout parse error out put failed", timeout)
	}

	ErrServerNotFound = func(name string) error {
		return Errorf("server %s not found", name)
	}

	ErrOptionFailed = func(err error, ref reflect.Value) error {
		return Wrapf(err, "failed to setup option :\t%s",
			runtime.FuncForPC(ref.Pointer()).Name())
	}

	ErrArgumentParseFailed = func(err error) error {
		return Wrap(err, "argument parse failed")
	}

	ErrDaemonStopFailed = func(err error) error {
		return Wrap(err, "failed to stop daemon")
	}

	ErrBackoffTimeout = func(err error) error {
		return Wrap(err, "backoff timeout by limitation")
	}

	ErrInvalidTypeConversion = func(i interface{}, tgt interface{}) error {
		return Errorf("invalid type conversion %v to %v", reflect.TypeOf(i), reflect.TypeOf(tgt))
	}

	// HTTP

	ErrInvalidAPIConfig = New("invalid api config")

	ErrHandler = func(err error) error {
		return Wrap(err, "handler returned error")
	}

	ErrHandlerTimeout = func(err error, t time.Time) error {
		return Wrapf(err, "handler timeout %v", t)
	}

	ErrRequestBodyCloseAndFlush = func(err error) error {
		return Wrap(err, "request body flush & close failed")
	}

	ErrRequestBodyClose = func(err error) error {
		return Wrap(err, "request body close failed")
	}

	ErrRequestBodyFlush = func(err error) error {
		return Wrap(err, "request body flush failed")
	}

	ErrTransportRetryable = New("transport is retryable")

	//NGT

	ErrCreateProperty = func(err error) error {
		return Wrap(err, "failed to create property")
	}

	ErrIndexNotFound = New("index file not found")

	ErrUnsupportedObjectType = New("unsupported ObjectType")

	ErrUnsupportedDistanceType = New("unsupported DistanceType")

	ErrFailedToSetDistanceType = func(err error, distance string) error {
		return Wrap(err, "failed to set distance type "+distance)
	}

	ErrFailedToSetObjectType = func(err error, t string) error {
		return Wrap(err, "failed to set object type "+t)
	}

	ErrFailedToSetDimension = func(err error) error {
		return Wrap(err, "failed to set dimension")
	}

	ErrFailedToSetCreationEdgeSize = func(err error) error {
		return Wrap(err, "failed to set creation edge size")
	}

	ErrFailedToSetSearchEdgeSize = func(err error) error {
		return Wrap(err, "failed to set search edge size")
	}

	// ErrCAPINotImplemented raises using not implemented function in C API
	ErrCAPINotImplemented = New("not implemented in C API")

	ErrUUIDAlreadyExists = func(uuid string, oid uint32) error {
		return Errorf("ngt uuid %s object id %d already exists ", uuid, oid)
	}

	ErrUUIDNotFound = func(id uint32) error {
		return Errorf("ngt object id %d's metadata not found", id)
	}

	ErrObjectIDNotFound = func(uuid string) error {
		return Errorf("ngt uuid %s's object id not found", uuid)
	}

	ErrObjectNotFound = func(err error, uuid string) error {
		return Wrapf(err, "ngt uuid %s's object not found", uuid)
	}

	ErrRemoveRequestedBeforeIndexing = func(oid uint) error {
		return Errorf("object id %d is not indexed we cannnot remove it", oid)
	}

	// Runtime

	ErrPanicRecovered = func(err error, rec interface{}) error {
		return Wrap(err, Errorf("panic recovered: %v", rec).Error())
	}

	ErrPanicString = func(err error, msg string) error {
		return Wrap(err, Errorf("panic recovered: %v", msg).Error())
	}

	ErrRuntimeError = func(err error, r runtime.Error) error {
		return Wrap(err, Errorf("system paniced caused by runtime error: %v", r).Error())
	}

	New = func(msg string) error {
		if msg == "" {
			return nil
		}
		return errors.New(msg)
	}

	Wrap = func(err error, msg string) error {
		if err != nil {
			if msg != "" {
				return errors.Wrap(err, msg)
			}
			return err
		}
		return New(msg)
	}

	Wrapf = func(err error, format string, args ...interface{}) error {
		if err != nil {
			if format != "" && len(args) > 0 {
				return errors.Wrapf(err, format, args...)
			}
			return err
		}
		return Errorf(format, args...)
	}

	Cause = func(err error) error {
		if err != nil {
			return errors.Cause(err)
		}
		return nil
	}

	Errorf = func(format string, args ...interface{}) error {
		if format != "" && len(args) > 0 {
			return errors.Errorf(format, args...)
		}
		return nil
	}
)
