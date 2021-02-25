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
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/pkg/meta/redis/service"
)

type server struct {
	redis service.Redis
}

func New(opts ...Option) meta.MetaServer {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *server) GetMeta(ctx context.Context, key *payload.Meta_Key) (*payload.Meta_Val, error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-redis.GetMeta")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	val, err := s.redis.Get(ctx, key.GetKey())
	if err != nil {
		if errors.IsErrRedisNotFound(err) {
			log.Warnf("[GetMeta]\tnot found\t%v\t%s", key.GetKey(), err.Error())
			err = status.WrapWithNotFound(fmt.Sprintf("GetMeta API: not found: key %s", key.GetKey()), err,
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMeta",
					ResourceName: key.GetKey(),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil, err
		}
		log.Errorf("[GetMeta]\tunknown error\t%+v", err)
		err = status.WrapWithUnknown(fmt.Sprintf("GetMeta API: unknown error occurred: key %s", key.GetKey()), err,
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMeta",
				ResourceName: key.GetKey(),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeUnknown(err.Error()))
		}
		return nil, err
	}
	return &payload.Meta_Val{
		Val: val,
	}, nil
}

func (s *server) GetMetas(ctx context.Context, keys *payload.Meta_Keys) (mv *payload.Meta_Vals, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-redis.GetMetas")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	mv = new(payload.Meta_Vals)
	mv.Vals, err = s.redis.GetMultiple(ctx, keys.GetKeys()...)
	if err != nil {
		if errors.IsErrRedisNotFound(err) {
			log.Warnf("[GetMetas]\tnot found\t%v\t%s", keys.GetKeys(), err.Error())
			err = status.WrapWithNotFound(fmt.Sprintf("GetMetas API: not found: keys %#v", keys.GetKeys()), err,
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetas",
					ResourceName: fmt.Sprintf("%#v", keys.GetKeys()),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return mv, err
		}
		log.Errorf("[GetMetas]\tunknown error\t%+v", err)
		err = status.WrapWithUnknown(fmt.Sprintf("GetMetas API: unknown error occurred: keys %#v", keys.GetKeys()), err,
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetas",
				ResourceName: fmt.Sprintf("%#v", keys.GetKeys()),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeUnknown(err.Error()))
		}
		return mv, err
	}
	return mv, nil
}

func (s *server) GetMetaInverse(ctx context.Context, val *payload.Meta_Val) (*payload.Meta_Key, error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-redis.GetMetaInverse")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	key, err := s.redis.GetInverse(ctx, val.GetVal())
	if err != nil {
		if errors.IsErrRedisNotFound(err) {
			log.Warnf("[GetMetaInverse]\tnot found\t%v\t%s", val.GetVal(), err.Error())
			err = status.WrapWithNotFound(fmt.Sprintf("GetMetaInverse API: not found: val %s", val.GetVal()), err,
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetaInverse",
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil, err
		}
		log.Errorf("[GetMetaInverse]\tunknown error\t%+v", err)
		err = status.WrapWithUnknown(fmt.Sprintf("GetMetaInverse API: unknown error occurred: val %s", val.GetVal()), err,
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetaInverse",
				ResourceName: val.GetVal(),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeUnknown(err.Error()))
		}
		return nil, err
	}
	return &payload.Meta_Key{
		Key: key,
	}, nil
}

func (s *server) GetMetasInverse(ctx context.Context, vals *payload.Meta_Vals) (mk *payload.Meta_Keys, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-redis.GetMetasInverse")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	mk = new(payload.Meta_Keys)
	mk.Keys, err = s.redis.GetInverseMultiple(ctx, vals.GetVals()...)
	if err != nil {
		if errors.IsErrRedisNotFound(err) {
			log.Warnf("[GetMetasInverse]\tnot found\t%v\t%s", vals.GetVals(), err.Error())
			err = status.WrapWithNotFound(fmt.Sprintf("GetMetasInverse API: not found: vals %#v", vals.GetVals()), err,
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetasInverse",
					ResourceName: fmt.Sprintf("%#v", vals.GetVals()),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return mk, err
		}
		log.Errorf("[GetMetasInverse]\tunknown error\t%+v", err)
		err = status.WrapWithUnknown(fmt.Sprintf("GetMetasInverse API: unknown error occurred: vals %#v", vals.GetVals()), err,
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.GetMetasInverse",
				ResourceName: fmt.Sprintf("%#v", vals.GetVals()),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeUnknown(err.Error()))
		}
		return mk, err
	}
	return mk, nil
}

func (s *server) SetMeta(ctx context.Context, kv *payload.Meta_KeyVal) (_ *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-redis.SetMeta")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = s.redis.Set(ctx, kv.GetKey(), kv.GetVal())
	if err != nil {
		log.Errorf("[SetMeta]\tunknown error\t%+v", err)
		err = status.WrapWithInternal(fmt.Sprintf("SetMeta API: failed to store: key %s val %s", kv.GetKey(), kv.GetVal()), err,
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.SetMeta",
				ResourceName: kv.GetKey(),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return new(payload.Empty), nil
}

func (s *server) SetMetas(ctx context.Context, kvs *payload.Meta_KeyVals) (_ *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-redis.SetMetas")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	query := make(map[string]string, len(kvs.GetKvs())/2)
	for _, kv := range kvs.GetKvs() {
		query[kv.GetKey()] = kv.GetVal()
	}
	err = s.redis.SetMultiple(ctx, query)
	if err != nil {
		log.Errorf("[SetMetas]\tunknown error\t%+v", err)
		err = status.WrapWithInternal(fmt.Sprintf("SetMetas API: failed to store: %#v", query), err,
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.SetMetas",
				ResourceName: fmt.Sprintf("%#v", query),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return new(payload.Empty), nil
}

func (s *server) DeleteMeta(ctx context.Context, key *payload.Meta_Key) (*payload.Meta_Val, error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-redis.DeleteMeta")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	val, err := s.redis.Delete(ctx, key.GetKey())
	if err != nil {
		if errors.IsErrRedisNotFound(err) {
			log.Warnf("[DeleteMeta]\tnot found\t%v\t%s", key.GetKey(), err.Error())
			err = status.WrapWithNotFound(fmt.Sprintf("DeleteMeta API: not found: key %s", key.GetKey()), err,
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMeta",
					ResourceName: key.GetKey(),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil, err
		}
		log.Errorf("[DeleteMeta]\tunknown error\t%+v", err)
		err = status.WrapWithUnknown(fmt.Sprintf("DeleteMeta API: unknown error occurred: key %s", key.GetKey()), err,
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMeta",
				ResourceName: key.GetKey(),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeUnknown(err.Error()))
		}
		return nil, err
	}
	return &payload.Meta_Val{
		Val: val,
	}, nil
}

func (s *server) DeleteMetas(ctx context.Context, keys *payload.Meta_Keys) (mv *payload.Meta_Vals, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-redis.DeleteMetas")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	mv = new(payload.Meta_Vals)
	mv.Vals, err = s.redis.DeleteMultiple(ctx, keys.GetKeys()...)
	if err != nil {
		if errors.IsErrRedisNotFound(err) {
			log.Warnf("[DeleteMetas]\tnot found\t%v\t%s", keys.GetKeys(), err.Error())
			err = status.WrapWithNotFound(fmt.Sprintf("DeleteMetas API: not found: keys %#v", keys.GetKeys()), err,
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetas",
					ResourceName: fmt.Sprintf("%#v", keys.GetKeys()),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return mv, err
		}
		log.Errorf("[DeleteMetas]\tunknown error\t%+v", err)
		err = status.WrapWithUnknown(fmt.Sprintf("DeleteMetas API: unknown error occurred: keys %#v", keys.GetKeys()), err,
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetas",
				ResourceName: fmt.Sprintf("%#v", keys.GetKeys()),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeUnknown(err.Error()))
		}
		return mv, err
	}
	return mv, nil
}

func (s *server) DeleteMetaInverse(ctx context.Context, val *payload.Meta_Val) (*payload.Meta_Key, error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-redis.DeleteMetaInverse")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	key, err := s.redis.DeleteInverse(ctx, val.GetVal())
	if err != nil {
		if errors.IsErrRedisNotFound(err) {
			log.Warnf("[DeleteMetaInverse]\tnot found\t%v\t%s", val.GetVal(), err.Error())
			err = status.WrapWithNotFound(fmt.Sprintf("DeleteMetaInverse API: not found: val %s", val.GetVal()), err,
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetaInverse",
					ResourceName: val.GetVal(),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil, err
		}
		log.Errorf("[DeleteMetaInverse]\tunknown error\t%+v", err)
		err = status.WrapWithUnknown(fmt.Sprintf("DeleteMetaInverse API: unknown error occurred: val %s", val.GetVal()), err,
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetaInverse",
				ResourceName: val.GetVal(),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeUnknown(err.Error()))
		}
		return nil, err
	}
	return &payload.Meta_Key{
		Key: key,
	}, nil
}

func (s *server) DeleteMetasInverse(ctx context.Context, vals *payload.Meta_Vals) (mk *payload.Meta_Keys, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/meta-redis.DeleteMetasInverse")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	mk = new(payload.Meta_Keys)
	mk.Keys, err = s.redis.DeleteInverseMultiple(ctx, vals.GetVals()...)
	if err != nil {
		if errors.IsErrRedisNotFound(err) {
			log.Warnf("[DeleteMetasInverse]\tnot found\t%v\t%s", vals.GetVals(), err.Error())
			err = status.WrapWithNotFound(fmt.Sprintf("DeleteMetasInverse API: not found: vals %#v", vals.GetVals()), err,
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetasInverse",
					ResourceName: fmt.Sprintf("%#v", vals.GetVals()),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return mk, err
		}
		log.Errorf("[DeleteMetasInverse]\tunknown error\t%+v", err)
		err = status.WrapWithUnknown(fmt.Sprintf("DeleteMetasInverse API: unknown error occurred: vals %#v", vals.GetVals()), err,
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta/Meta.DeleteMetasInverse",
				ResourceName: fmt.Sprintf("%#v", vals.GetVals()),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeUnknown(err.Error()))
		}
		return mk, err
	}
	return mk, nil
}
