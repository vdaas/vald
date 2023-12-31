// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package redis

import (
	"context"

	redis "github.com/go-redis/redis/v8"
)

type MockRedis struct {
	TxPipelineFunc func() redis.Pipeliner
	PingFunc       func() *StatusCmd
	CloseFunc      func() error
	GetFunc        func(string) *redis.StringCmd
	MGetFunc       func(...string) *redis.SliceCmd
	DelFunc        func(keys ...string) *redis.IntCmd
}

var _ Redis = (*MockRedis)(nil)

func (m *MockRedis) TxPipeline() redis.Pipeliner {
	return m.TxPipelineFunc()
}

func (m *MockRedis) Ping(context.Context) *StatusCmd {
	return m.PingFunc()
}

func (m *MockRedis) Close() error {
	return m.CloseFunc()
}

func (m *MockRedis) Get(_ context.Context, key string) *redis.StringCmd {
	return m.GetFunc(key)
}

func (m *MockRedis) MGet(_ context.Context, keys ...string) *redis.SliceCmd {
	return m.MGetFunc(keys...)
}

func (m *MockRedis) Del(_ context.Context, keys ...string) *redis.IntCmd {
	return m.DelFunc(keys...)
}

type dummyHook struct {
	name string
}

func (*dummyHook) BeforeProcess(ctx context.Context, _ Cmder) (context.Context, error) {
	return ctx, nil
}

func (*dummyHook) AfterProcess(context.Context, Cmder) error {
	return nil
}

func (*dummyHook) BeforeProcessPipeline(ctx context.Context, _ []Cmder) (context.Context, error) {
	return ctx, nil
}

func (*dummyHook) AfterProcessPipeline(context.Context, []Cmder) error {
	return nil
}

func dummyWithFunc(err error) Option {
	return func(*redisClient) error {
		return err
	}
}
