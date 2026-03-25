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
)

type rawKVClient struct {
	cli *rawkv.Client
}

func NewRawKVClient(ctx context.Context, addrs []string) (ManagedMetadataClient, error) {
	cli, err := rawkv.NewClient(ctx, addrs, tikvcfg.DefaultConfig().Security)
	if err != nil {
		return nil, err
	}
	return &rawKVClient{
		cli: cli,
	}, nil
}

func (c *rawKVClient) Close() error {
	return c.cli.Close()
}

func (c *rawKVClient) Get(ctx context.Context, key []byte) ([]byte, error) {
	return c.cli.Get(ctx, key)
}

func (c *rawKVClient) Put(ctx context.Context, key, val []byte) error {
	return c.cli.Put(ctx, key, val)
}

func (c *rawKVClient) Delete(ctx context.Context, key []byte) error {
	return c.cli.Delete(ctx, key)
}
