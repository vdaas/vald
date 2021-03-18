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

// Package grpc provides grpc server logic
package grpc

import (
	"context"
	"fmt"

	"github.com/vdaas/vald/apis/grpc/v1/meta"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/db/nosql/cassandra"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/pkg/meta/cassandra/service"
)

type server struct {
	cassandra service.Cassandra
}

func New(opts ...Option) meta.MetaServer {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}

func isInternalError(err error) bool {
	return errors.As(err, &cassandra.ErrUnsupported) ||
		errors.As(err, &cassandra.ErrTooManyStmts) ||
		errors.As(err, &cassandra.ErrUseStmt) ||
		errors.As(err, &cassandra.ErrSessionClosed) ||
		errors.As(err, &cassandra.ErrNoKeyspace) ||
		errors.As(err, &cassandra.ErrNoMetadata) ||
		errors.As(err, &cassandra.ErrHostQueryFailed)
}

func isFailedPrecondition(err error) bool {
	return errors.As(err, &cassandra.ErrNoConnections) ||
		errors.As(err, &cassandra.ErrKeyspaceDoesNotExist) ||
		errors.As(err, &cassandra.ErrNoHosts) ||
		errors.As(err, &cassandra.ErrNoConnectionsStarted)
}

func getPreconditionFailureDetails(err error) (res []*errdetails.PreconditionFailureViolation) {
	res = make([]*errdetails.PreconditionFailureViolation, 0)

	if errors.As(err, &cassandra.ErrNoConnections) {
		res = append(res, &errdetails.PreconditionFailureViolation{
			Type:        "No connections",
			Subject:     "Cassandra",
			Description: cassandra.ErrNoConnections.Error(),
		})
	}

	if errors.As(err, &cassandra.ErrKeyspaceDoesNotExist) {
		res = append(res, &errdetails.PreconditionFailureViolation{
			Type:        "Keyspace does not exist",
			Subject:     "Cassandra",
			Description: cassandra.ErrKeyspaceDoesNotExist.Error(),
		})
	}

	if errors.As(err, &cassandra.ErrNoHosts) {
		res = append(res, &errdetails.PreconditionFailureViolation{
			Type:        "No hosts",
			Subject:     "Cassandra",
			Description: cassandra.ErrNoHosts.Error(),
		})
	}

	if errors.As(err, &cassandra.ErrNoConnectionsStarted) {
		res = append(res, &errdetails.PreconditionFailureViolation{
			Type:        "No connections started",
			Subject:     "Cassandra",
			Description: cassandra.ErrNoConnectionsStarted.Error(),
		})
	}

	return res
}

func (s *server) GetMeta(ctx context.Context, key *payload.Meta_Key) (*payload.Meta_Val, error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-cassandra.GetMeta")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	val, err := s.cassandra.Get(key.GetKey())
	if err != nil {
		switch {
		case errors.IsErrCassandraNotFound(err):
			log.Warnf("[GetMeta]\tnot found\t%v\t%s", key.GetKey(), err.Error())
			err = status.WrapWithNotFound(fmt.Sprintf("GetMeta API: not found: key %s", key.GetKey()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(key),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMeta",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil, err
		case errors.IsErrCassandraUnavailable(err):
			log.Warnf("[GetMeta]\tunavailable\t%+v", err)
			err = status.WrapWithUnavailable(fmt.Sprintf("GetMeta API: Cassandra unavailable: key %s", key.GetKey()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(key),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMeta",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnavailable(err.Error()))
			}
			return nil, err
		case isInternalError(err):
			log.Warnf("[GetMeta]\tinternal error\t%+v", err)
			err = status.WrapWithInternal(fmt.Sprintf("GetMeta API: internal error occurred: key %s", key.GetKey()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(key),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMeta",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, err
		case isFailedPrecondition(err):
			log.Warnf("[GetMeta]\tfailed precondition\t%+v", err)
			err = status.WrapWithFailedPrecondition(fmt.Sprintf("GetMeta API: failed precondition: key %s", key.GetKey()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(key),
				},
				&errdetails.PreconditionFailure{
					Violations: getPreconditionFailureDetails(err),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMeta",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeFailedPrecondition(err.Error()))
			}
			return nil, err
		default:
			log.Errorf("[GetMeta]\tunknown error\t%+v", err)
			err = status.WrapWithUnknown(fmt.Sprintf("GetMeta API: unknown error occurred: key %s", key.GetKey()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(key),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMeta",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnknown(err.Error()))
			}
			return nil, err
		}
	}
	return &payload.Meta_Val{
		Val: val,
	}, nil
}

func (s *server) GetMetas(ctx context.Context, keys *payload.Meta_Keys) (mv *payload.Meta_Vals, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-cassandra.GetMetas")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	mv = new(payload.Meta_Vals)
	mv.Vals, err = s.cassandra.GetMultiple(keys.GetKeys()...)
	if err != nil {
		switch {
		case errors.IsErrCassandraNotFound(err):
			log.Warnf("[GetMetas]\tnot found\t%v\t%s", keys.GetKeys(), err.Error())
			err = status.WrapWithNotFound(fmt.Sprintf("GetMetas API: not found: keys %#v", keys.GetKeys()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(keys),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetas",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return mv, err
		case errors.IsErrCassandraUnavailable(err):
			log.Warnf("[GetMetas]\tunavailable\t%+v", err)
			err = status.WrapWithUnavailable(fmt.Sprintf("GetMetas API: Cassandra unavailable: keys %#v", keys.GetKeys()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(keys),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetas",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnavailable(err.Error()))
			}
			return mv, err
		case isInternalError(err):
			log.Warnf("[GetMetas]\tinternal error\t%+v", err)
			err = status.WrapWithInternal(fmt.Sprintf("GetMetas API: internal error occurred: keys %#v", keys.GetKeys()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(keys),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetas",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, err
		case isFailedPrecondition(err):
			log.Warnf("[GetMetas]\tfailed precondition\t%+v", err)
			err = status.WrapWithFailedPrecondition(fmt.Sprintf("GetMetas API: failed precondition: keys %#v", keys.GetKeys()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(keys),
				},
				&errdetails.PreconditionFailure{
					Violations: getPreconditionFailureDetails(err),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetas",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeFailedPrecondition(err.Error()))
			}
			return nil, err
		default:
			log.Errorf("[GetMetas]\tunknown error\t%+v", err)
			err = status.WrapWithUnknown(fmt.Sprintf("GetMetas API: unknown error occurred: keys %#v", keys.GetKeys()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(keys),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetas",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnknown(err.Error()))
			}
			return mv, err
		}
	}
	return mv, nil
}

func (s *server) GetMetaInverse(ctx context.Context, val *payload.Meta_Val) (*payload.Meta_Key, error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-cassandra.GetMetaInverse")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	key, err := s.cassandra.GetInverse(val.GetVal())
	if err != nil {
		switch {
		case errors.IsErrCassandraNotFound(err):
			log.Warnf("[GetMetaInverse]\tnot found\t%v\t%s", val.GetVal(), err.Error())
			err = status.WrapWithNotFound(fmt.Sprintf("GetMetaInverse API: not found: val %s", val.GetVal()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(val),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetaInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil, err
		case errors.IsErrCassandraUnavailable(err):
			log.Warnf("[GetMetaInverse]\tunavailable\t%+v", err)
			err = status.WrapWithUnavailable(fmt.Sprintf("GetMetaInverse API: Cassandra unavailable: val %s", val.GetVal()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(val),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetaInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnavailable(err.Error()))
			}
			return nil, err
		case isInternalError(err):
			log.Warnf("[GetMetaInverse]\tinternal error\t%+v", err)
			err = status.WrapWithInternal(fmt.Sprintf("GetMetaInverse API: internal error occurred: val %s", val.GetVal()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(val),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetaInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, err
		case isFailedPrecondition(err):
			log.Warnf("[GetMetaInverse]\tfailed precondition\t%+v", err)
			err = status.WrapWithFailedPrecondition(fmt.Sprintf("GetMetaInverse API: failed precondition: val %s", val.GetVal()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(val),
				},
				&errdetails.PreconditionFailure{
					Violations: getPreconditionFailureDetails(err),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetaInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeFailedPrecondition(err.Error()))
			}
			return nil, err
		default:
			log.Errorf("[GetMetaInverse]\tunknown error\t%+v", err)
			err = status.WrapWithUnknown(fmt.Sprintf("GetMetaInverse API: unknown error occurred: val %s", val.GetVal()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(val),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetaInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnknown(err.Error()))
			}
			return nil, err
		}
	}
	return &payload.Meta_Key{
		Key: key,
	}, nil
}

func (s *server) GetMetasInverse(ctx context.Context, vals *payload.Meta_Vals) (mk *payload.Meta_Keys, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-cassandra.GetMetasInverse")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	mk = new(payload.Meta_Keys)
	mk.Keys, err = s.cassandra.GetInverseMultiple(vals.GetVals()...)
	if err != nil {
		switch {
		case errors.IsErrCassandraNotFound(err):
			log.Warnf("[GetMetasInverse]\tnot found\t%v\t%s", vals.GetVals(), err.Error())
			err = status.WrapWithNotFound(fmt.Sprintf("GetMetasInverse API: not found: vals %#v", vals.GetVals()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(vals),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetasInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return mk, err
		case errors.IsErrCassandraUnavailable(err):
			log.Warnf("[GetMetasInverse]\tunavailable\t%+v", err)
			err = status.WrapWithUnavailable(fmt.Sprintf("GetMetasInverse API: Cassandra unavailable: vals %#v", vals.GetVals()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(vals),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetasInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnavailable(err.Error()))
			}
			return mk, err
		case isInternalError(err):
			log.Warnf("[GetMetasInverse]\tinternal error\t%+v", err)
			err = status.WrapWithInternal(fmt.Sprintf("GetMetasInverse API: internal error occurred: vals %#v", vals.GetVals()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(vals),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetasInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, err
		case isFailedPrecondition(err):
			log.Warnf("[GetMetasInverse]\tfailed precondition\t%+v", err)
			err = status.WrapWithFailedPrecondition(fmt.Sprintf("GetMetasInverse API: failed precondition: vals %#v", vals.GetVals()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(vals),
				},
				&errdetails.PreconditionFailure{
					Violations: getPreconditionFailureDetails(err),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetasInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeFailedPrecondition(err.Error()))
			}
			return nil, err
		default:
			log.Errorf("[GetMetasInverse]\tunknown error\t%+v", err)
			err = status.WrapWithUnknown(fmt.Sprintf("GetMetasInverse API: unknown error occurred: vals %#v", vals.GetVals()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(vals),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetasInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnknown(err.Error()))
			}
			return mk, err
		}
	}
	return mk, nil
}

func (s *server) SetMeta(ctx context.Context, kv *payload.Meta_KeyVal) (_ *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-cassandra.SetMeta")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = s.cassandra.Set(kv.GetKey(), kv.GetVal())
	if err != nil {
		switch {
		case errors.IsErrCassandraUnavailable(err):
			log.Warnf("[SetMeta]\tunavailable\t%+v", err)
			err = status.WrapWithUnavailable(fmt.Sprintf("SetMeta API: Cassandra unavailable: key %s, val %s", kv.GetKey(), kv.GetVal()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(kv),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.SetMeta",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnavailable(err.Error()))
			}
			return nil, err
		case isInternalError(err):
			log.Warnf("[SetMeta]\tinternal error\t%+v", err)
			err = status.WrapWithInternal(fmt.Sprintf("SetMeta API: internal error occurred: key %s val %s", kv.GetKey(), kv.GetVal()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(kv),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.SetMeta",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, err
		case isFailedPrecondition(err):
			log.Warnf("[SetMeta]\tfailed precondition\t%+v", err)
			err = status.WrapWithFailedPrecondition(fmt.Sprintf("SetMeta API: failed precondition: key %s val %s", kv.GetKey(), kv.GetVal()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(kv),
				},
				&errdetails.PreconditionFailure{
					Violations: getPreconditionFailureDetails(err),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.SetMeta",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeFailedPrecondition(err.Error()))
			}
			return nil, err
		default:
			log.Errorf("[SetMeta]\tunknown error\t%+v", err)
			err = status.WrapWithInternal(fmt.Sprintf("SetMeta API: failed to store: key %s val %s", kv.GetKey(), kv.GetVal()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(kv),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.SetMeta",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, err
		}
	}
	return new(payload.Empty), nil
}

func (s *server) SetMetas(ctx context.Context, kvs *payload.Meta_KeyVals) (_ *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-cassandra.SetMetas")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	query := make(map[string]string, len(kvs.GetKvs())/2)
	for _, kv := range kvs.GetKvs() {
		query[kv.GetKey()] = kv.GetVal()
	}
	err = s.cassandra.SetMultiple(query)
	if err != nil {
		switch {
		case errors.IsErrCassandraUnavailable(err):
			log.Warnf("[SetMetas]\tunavailable\t%+v", err)
			err = status.WrapWithUnavailable(fmt.Sprintf("SetMetas API: Cassandra unavailable: kvs %#v", query), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(kvs),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.SetMetas",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnavailable(err.Error()))
			}
			return nil, err
		case isInternalError(err):
			log.Warnf("[SetMetas]\tinternal error\t%+v", err)
			err = status.WrapWithInternal(fmt.Sprintf("SetMetas API: internal error occurred: kvs %#v", query), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(kvs),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.SetMetas",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, err
		case isFailedPrecondition(err):
			log.Warnf("[SetMetas]\tfailed precondition\t%+v", err)
			err = status.WrapWithFailedPrecondition(fmt.Sprintf("SetMetas API: failed precondition: kvs %#v", query), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(kvs),
				},
				&errdetails.PreconditionFailure{
					Violations: getPreconditionFailureDetails(err),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.SetMetas",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeFailedPrecondition(err.Error()))
			}
			return nil, err
		default:
			log.Errorf("[SetMetas]\tunknown error\t%+v", err)
			err = status.WrapWithInternal(fmt.Sprintf("SetMetas API: failed to store: kvs %#v", query), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(kvs),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.SetMetas",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, err
		}
	}
	return new(payload.Empty), nil
}

func (s *server) DeleteMeta(ctx context.Context, key *payload.Meta_Key) (*payload.Meta_Val, error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-cassandra.DeleteMeta")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	val, err := s.cassandra.Delete(key.GetKey())
	if err != nil {
		switch {
		case errors.IsErrCassandraNotFound(err):
			log.Warnf("[DeleteMeta]\tnot found\t%v\t%s", key.GetKey(), err.Error())
			err = status.WrapWithNotFound(fmt.Sprintf("DeleteMeta API: not found: key %s", key.GetKey()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(key),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMeta",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil, err
		case errors.IsErrCassandraUnavailable(err):
			log.Warnf("[DeleteMeta]\tunavailable\t%+v", err)
			err = status.WrapWithUnavailable(fmt.Sprintf("DeleteMeta API: Cassandra unavailable: key %s", key.GetKey()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(key),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMeta",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnavailable(err.Error()))
			}
			return nil, err
		case isInternalError(err):
			log.Warnf("[DeleteMeta]\tinternal error\t%+v", err)
			err = status.WrapWithInternal(fmt.Sprintf("DeleteMeta API: internal error occurred: key %s", key.GetKey()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(key),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMeta",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, err
		case isFailedPrecondition(err):
			log.Warnf("[DeleteMeta]\tfailed precondition\t%+v", err)
			err = status.WrapWithFailedPrecondition(fmt.Sprintf("DeleteMeta API: failed precondition: key %s", key.GetKey()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(key),
				},
				&errdetails.PreconditionFailure{
					Violations: getPreconditionFailureDetails(err),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMeta",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeFailedPrecondition(err.Error()))
			}
			return nil, err
		default:
			log.Errorf("[DeleteMeta]\tunknown error\t%+v", err)
			err = status.WrapWithUnknown(fmt.Sprintf("DeleteMeta API: unknown error occurred: key %s", key.GetKey()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(key),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMeta",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnknown(err.Error()))
			}
			return nil, err
		}
	}
	return &payload.Meta_Val{
		Val: val,
	}, nil
}

func (s *server) DeleteMetas(ctx context.Context, keys *payload.Meta_Keys) (mv *payload.Meta_Vals, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-cassandra.DeleteMetas")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	mv = new(payload.Meta_Vals)
	mv.Vals, err = s.cassandra.DeleteMultiple(keys.GetKeys()...)
	if err != nil {
		switch {
		case errors.IsErrCassandraNotFound(err):
			log.Warnf("[DeleteMetas]\tnot found\t%v\t%s", keys.GetKeys(), err.Error())
			err = status.WrapWithNotFound(fmt.Sprintf("DeleteMetas API: not found: keys %#v", keys.GetKeys()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(keys),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetas",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return mv, err
		case errors.IsErrCassandraUnavailable(err):
			log.Warnf("[DeleteMetas]\tunavailable\t%+v", err)
			err = status.WrapWithUnavailable(fmt.Sprintf("DeleteMetas API: Cassandra unavailable: keys %#v", keys.GetKeys()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(keys),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetas",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnavailable(err.Error()))
			}
			return nil, err
		case isInternalError(err):
			log.Warnf("[DeleteMetas]\tinternal error\t%+v", err)
			err = status.WrapWithInternal(fmt.Sprintf("DeleteMetas API: internal error occurred: keys %#v", keys.GetKeys()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(keys),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetas",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, err
		case isFailedPrecondition(err):
			log.Warnf("[DeleteMetas]\tfailed precondition\t%+v", err)
			err = status.WrapWithFailedPrecondition(fmt.Sprintf("DeleteMetas API: failed precondition: keys %#v", keys.GetKeys()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(keys),
				},
				&errdetails.PreconditionFailure{
					Violations: getPreconditionFailureDetails(err),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetas",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeFailedPrecondition(err.Error()))
			}
			return nil, err
		default:
			log.Errorf("[DeleteMetas]\tunknown error\t%+v", err)
			err = status.WrapWithUnknown(fmt.Sprintf("DeleteMetas API: unknown error occurred: keys %#v", keys.GetKeys()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(keys),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetas",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnknown(err.Error()))
			}
			return mv, err
		}
	}
	return mv, nil
}

func (s *server) DeleteMetaInverse(ctx context.Context, val *payload.Meta_Val) (*payload.Meta_Key, error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-cassandra.DeleteMetaInverse")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	key, err := s.cassandra.DeleteInverse(val.GetVal())
	if err != nil {
		switch {
		case errors.IsErrCassandraNotFound(err):
			log.Warnf("[DeleteMetaInverse]\tnot found\t%v\t%s", val.GetVal(), err.Error())
			err = status.WrapWithNotFound(fmt.Sprintf("DeleteMetaInverse API: not found: val %s", val.GetVal()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(val),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetaInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil, err
		case errors.IsErrCassandraUnavailable(err):
			log.Warnf("[DeleteMetaInverse]\tunavailable\t%+v", err)
			err = status.WrapWithUnavailable(fmt.Sprintf("DeleteMetaInverse API: Cassandra unavailable: val %s", val.GetVal()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(val),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetaInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnavailable(err.Error()))
			}
			return nil, err
		case isInternalError(err):
			log.Warnf("[DeleteMetaInverse]\tinternal error\t%+v", err)
			err = status.WrapWithInternal(fmt.Sprintf("DeleteMetaInverse API: internal error occurred: val %s", val.GetVal()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(val),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetaInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, err
		case isFailedPrecondition(err):
			log.Warnf("[DeleteMetaInverse]\tfailed precondition\t%+v", err)
			err = status.WrapWithFailedPrecondition(fmt.Sprintf("DeleteMetaInverse API: failed precondition: val %s", val.GetVal()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(val),
				},
				&errdetails.PreconditionFailure{
					Violations: getPreconditionFailureDetails(err),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetaInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeFailedPrecondition(err.Error()))
			}
			return nil, err
		default:
			log.Errorf("[DeleteMetaInverse]\tunknown error\t%+v", err)
			err = status.WrapWithUnknown(fmt.Sprintf("DeleteMetaInverse API: unknown error occurred: val %s", val.GetVal()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(val),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetaInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnknown(err.Error()))
			}
			return nil, err
		}
	}
	return &payload.Meta_Key{
		Key: key,
	}, nil
}

func (s *server) DeleteMetasInverse(ctx context.Context, vals *payload.Meta_Vals) (mk *payload.Meta_Keys, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-cassandra.DeleteMetasInverse")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	mk = new(payload.Meta_Keys)
	mk.Keys, err = s.cassandra.DeleteInverseMultiple(vals.GetVals()...)
	if err != nil {
		switch {
		case errors.IsErrCassandraNotFound(err):
			log.Warnf("[DeleteMetasInverse]\tnot found\t%v\t%s", vals.GetVals(), err.Error())
			err = status.WrapWithNotFound(fmt.Sprintf("DeleteMetasInverse API: not found: vals %#v", vals.GetVals()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(vals),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetasInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return mk, err
		case errors.IsErrCassandraUnavailable(err):
			log.Warnf("[DeleteMetasInverse]\tunavailable\t%+v", err)
			err = status.WrapWithUnavailable(fmt.Sprintf("DeleteMetasInverse API: Cassandra unavailable: vals %#v", vals.GetVals()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(vals),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetasInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnavailable(err.Error()))
			}
			return nil, err
		case isInternalError(err):
			log.Warnf("[DeleteMetasInverse]\tinternal error\t%+v", err)
			err = status.WrapWithInternal(fmt.Sprintf("DeleteMetasInverse API: internal error occurred: vals %#v", vals.GetVals()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(vals),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetasInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, err
		case isFailedPrecondition(err):
			log.Warnf("[DeleteMetasInverse]\tfailed precondition\t%+v", err)
			err = status.WrapWithFailedPrecondition(fmt.Sprintf("DeleteMetasInverse API: failed precondition: vals %#v", vals.GetVals()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(vals),
				},
				&errdetails.PreconditionFailure{
					Violations: getPreconditionFailureDetails(err),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetasInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeFailedPrecondition(err.Error()))
			}
			return nil, err
		default:
			log.Errorf("[DeleteMetasInverse]\tunknown error\t%+v", err)
			err = status.WrapWithUnknown(fmt.Sprintf("DeleteMetasInverse API: unknown error occurred: vals %#v", vals.GetVals()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(vals),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetasInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeUnknown(err.Error()))
			}
			return mk, err
		}
	}
	return mk, nil
}
