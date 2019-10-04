//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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

	ErrInvalidRequest = New("invalid request")

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

	// gRPC

	ErrAgentClientNotConnected = New("agent client not connected")

	//DB

	ErrAddrsNotFound = New("addrs not found")

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

	ErrUncommittedIndexExists = func(num uint64) error {
		return Errorf("%d indexes are not commited", num)
	}

	ErrUncommittedIndexNotFound = New("uncommitted indexes are not found")

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
