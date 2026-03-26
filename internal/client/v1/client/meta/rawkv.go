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

package meta

import (
	"context"

	tikvcfg "github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/rawkv"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/trace"
)

const apiName = "vald/internal/client/v1/client/meta"

type rawKVClient struct {
	cli *rawkv.Client
}

func NewRawKVClient(ctx context.Context, addrs []string) (ManagedMetadataClient, error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/NewRawKVClient")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	cli, err := rawkv.NewClient(ctx, addrs, tikvcfg.DefaultConfig().Security)
	if err != nil {
		err = errors.ErrNewTiKVRawClientFailed(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return &rawKVClient{
		cli: cli,
	}, nil
}

func (c *rawKVClient) Close() error {
	err := c.cli.Close()
	if err != nil {
		return errors.ErrTiKVRawClientCloseOperationFailed(err)
	}
	return nil
}

func (c *rawKVClient) Get(ctx context.Context, key []byte) ([]byte, error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Get")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	val, err := c.cli.Get(ctx, key)
	if err != nil {
		err = errors.ErrTiKVGetOperationFailed(key, err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return val, nil
}

func (c *rawKVClient) Put(ctx context.Context, key, val []byte) error {
	ctx, span := trace.StartSpan(ctx, apiName+"/Put")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	err := c.cli.Put(ctx, key, val)
	if err != nil {
		err = errors.ErrTiKVSetOperationFailed(key, val, err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
	}
	return err
}

func (c *rawKVClient) Delete(ctx context.Context, key []byte) error {
	ctx, span := trace.StartSpan(ctx, apiName+"/Delete")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	err := c.cli.Delete(ctx, key)
	if err != nil {
		err = errors.ErrTiKVDeleteOperationFailed(key, err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
	}
	return err
}
