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

// Package grpc provides grpc server logic
package grpc

import (
	"context"
	"fmt"

	"github.com/vdaas/vald/apis/grpc/meta"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/pkg/meta/redis/service"
)

type server struct {
	redis service.Redis
}

func New(opts ...Option) meta.MetaServer {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}
	return s
}

func (s *server) GetMeta(ctx context.Context, key *payload.Meta_Key) (*payload.Meta_Val, error) {
	val, err := s.redis.Get(key.GetKey())
	if err != nil {
		if errors.IsErrRedisNotFound(errors.UnWrapAll(err)) {
			log.Warnf("[GetMeta]\tnot found\t%v\t%+v", key.GetKey(), err)
			return nil, status.WrapWithNotFound(fmt.Sprintf("GetMeta API key %s not found", key.GetKey()), err, info.Get())
		}
		log.Errorf("[GetMeta]\tunknown error\t%+v", err)
		return nil, status.WrapWithUnknown(fmt.Sprintf("GetMeta API unknown error occurred key %s", key.GetKey()), err, info.Get())
	}
	return &payload.Meta_Val{
		Val: val,
	}, nil
}

func (s *server) GetMetas(ctx context.Context, keys *payload.Meta_Keys) (mv *payload.Meta_Vals, err error) {
	mv = new(payload.Meta_Vals)
	mv.Vals, err = s.redis.GetMultiple(keys.GetKeys()...)
	if err != nil {
		if errors.IsErrRedisNotFound(errors.UnWrapAll(err)) {
			log.Warnf("[GetMetas]\tnot found\t%v\t%+v", keys.GetKeys(), err)
			return mv, status.WrapWithNotFound(fmt.Sprintf("GetMetas API Redis entry keys %#v not found", keys.GetKeys()), err, info.Get())
		}
		log.Errorf("[GetMetas]\tunknown error\t%+v", err)
		return mv, status.WrapWithUnknown(fmt.Sprintf("GetMetas API Redis entry keys %#v unknown error occurred", keys.GetKeys()), err, info.Get())
	}
	return mv, nil
}

func (s *server) GetMetaInverse(ctx context.Context, val *payload.Meta_Val) (*payload.Meta_Key, error) {
	key, err := s.redis.GetInverse(val.GetVal())
	if err != nil {
		if errors.IsErrRedisNotFound(errors.UnWrapAll(err)) {
			log.Warnf("[GetMetaInverse]\tnot found\t%v\t%+v", val.GetVal(), err)
			return nil, status.WrapWithNotFound(fmt.Sprintf("GetMetaInverse API val %s not found", val.GetVal()), err, info.Get())
		}
		log.Errorf("[GetMetaInverse]\tunknown error\t%+v", err)
		return nil, status.WrapWithUnknown(fmt.Sprintf("GetMetaInverse API val %s unknown error occurred", val.GetVal()), err, info.Get())
	}
	return &payload.Meta_Key{
		Key: key,
	}, nil
}

func (s *server) GetMetasInverse(ctx context.Context, vals *payload.Meta_Vals) (mk *payload.Meta_Keys, err error) {
	mk = new(payload.Meta_Keys)
	mk.Keys, err = s.redis.GetInverseMultiple(vals.GetVals()...)
	if err != nil {
		if errors.IsErrRedisNotFound(errors.UnWrapAll(err)) {
			log.Warnf("[GetMetasInverse]\tnot found\t%v\t%+v", vals.GetVals(), err)
			return mk, status.WrapWithNotFound(fmt.Sprintf("GetMetasInverse API vals %#v not found", vals.GetVals()), err, info.Get())
		}
		log.Errorf("[GetMetasInverse]\tunknown error\t%+v", err)
		return mk, status.WrapWithUnknown(fmt.Sprintf("GetMetasInverse API vals %#v unknown error occurred", vals.GetVals()), err, info.Get())
	}
	return mk, nil
}

func (s *server) SetMeta(ctx context.Context, kv *payload.Meta_KeyVal) (_ *payload.Empty, err error) {
	err = s.redis.Set(kv.GetKey(), kv.GetVal())
	if err != nil {
		log.Errorf("[SetMeta]\tunknown error\t%+v", err)
		return nil, status.WrapWithInternal(fmt.Sprintf("SetMeta API key %s val %s failed to store", kv.GetKey(), kv.GetVal()), err, info.Get())
	}
	return new(payload.Empty), nil
}

func (s *server) SetMetas(ctx context.Context, kvs *payload.Meta_KeyVals) (_ *payload.Empty, err error) {
	query := make(map[string]string, len(kvs.GetKvs())/2)
	for _, kv := range kvs.GetKvs() {
		query[kv.GetKey()] = kv.GetVal()
	}
	err = s.redis.SetMultiple(query)
	if err != nil {
		log.Errorf("[SetMetas]\tunknown error\t%+v", err)
		return nil, status.WrapWithInternal(fmt.Sprintf("SetMetas API failed to store %#v", query), err, info.Get())
	}
	return new(payload.Empty), nil
}

func (s *server) DeleteMeta(ctx context.Context, key *payload.Meta_Key) (*payload.Meta_Val, error) {
	val, err := s.redis.Delete(key.GetKey())
	if err != nil {
		if errors.IsErrRedisNotFound(errors.UnWrapAll(err)) {
			log.Warnf("[DeleteMeta]\tnot found\t%v\t%+v", key.GetKey(), err)
			return nil, status.WrapWithNotFound(fmt.Sprintf("DeleteMeta API key %s not found", key.GetKey()), err, info.Get())
		}
		log.Errorf("[DeleteMeta]\tunknown error\t%+v", err)
		return nil, status.WrapWithUnknown(fmt.Sprintf("DeleteMeta API unknown error occurred key %s", key.GetKey()), err, info.Get())
	}
	return &payload.Meta_Val{
		Val: val,
	}, nil
}

func (s *server) DeleteMetas(ctx context.Context, keys *payload.Meta_Keys) (mv *payload.Meta_Vals, err error) {
	mv = new(payload.Meta_Vals)
	mv.Vals, err = s.redis.DeleteMultiple(keys.GetKeys()...)
	if err != nil {
		if errors.IsErrRedisNotFound(errors.UnWrapAll(err)) {
			log.Warnf("[DeleteMetas]\tnot found\t%v\t%+v", keys.GetKeys(), err)
			return mv, status.WrapWithNotFound(fmt.Sprintf("DeleteMetas API Redis entry keys %#v not found", keys.GetKeys()), err, info.Get())
		}
		log.Errorf("[DeleteMetas]\tunknown error\t%+v", err)
		return mv, status.WrapWithUnknown(fmt.Sprintf("DeleteMetas API Redis entry keys %#v unknown error occurred", keys.GetKeys()), err, info.Get())
	}
	return mv, nil
}

func (s *server) DeleteMetaInverse(ctx context.Context, val *payload.Meta_Val) (*payload.Meta_Key, error) {
	key, err := s.redis.DeleteInverse(val.GetVal())
	if err != nil {
		if errors.IsErrRedisNotFound(errors.UnWrapAll(err)) {
			log.Warnf("[DeleteMetaInverse]\tnot found\t%v\t%+v", val.GetVal(), err)
			return nil, status.WrapWithNotFound(fmt.Sprintf("DeleteMetaInverse API val %s not found", val.GetVal()), err, info.Get())
		}
		log.Errorf("[DeleteMetaInverse]\tunknown error\t%+v", err)
		return nil, status.WrapWithUnknown(fmt.Sprintf("DeleteMetaInverse API val %s unknown error occurred", val.GetVal()), err, info.Get())
	}
	return &payload.Meta_Key{
		Key: key,
	}, nil
}

func (s *server) DeleteMetasInverse(ctx context.Context, vals *payload.Meta_Vals) (mk *payload.Meta_Keys, err error) {
	mk = new(payload.Meta_Keys)
	mk.Keys, err = s.redis.DeleteInverseMultiple(vals.GetVals()...)
	if err != nil {
		if errors.IsErrRedisNotFound(errors.UnWrapAll(err)) {
			log.Warnf("[DeleteMetasInverse]\tnot found\t%v\t%+v", vals.GetVals(), err)
			return mk, status.WrapWithNotFound(fmt.Sprintf("DeleteMetasInverse API vals %#v not found", vals.GetVals()), err, info.Get())
		}
		log.Errorf("[DeleteMetasInverse]\tunknown error\t%+v", err)
		return mk, status.WrapWithUnknown(fmt.Sprintf("DeleteMetasInverse API vals %#v unknown error occurred", vals.GetVals()), err, info.Get())
	}
	return mk, nil
}
