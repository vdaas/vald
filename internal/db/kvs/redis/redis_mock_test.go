package redis

import redis "github.com/go-redis/redis/v7"

type MockRedis struct {
	TxPipelineFunc func() redis.Pipeliner
	PingFunc       func() *StatusCmd
	CloseFunc      func() error
	GetFunc        func(string) *redis.StringCmd
	MGetFunc       func(...string) *redis.SliceCmd
	DelFunc        func(keys ...string) *redis.IntCmd
}

var _ = (*MockRedis)(nil)

func (m *MockRedis) TxPipeline() redis.Pipeliner {
	return m.TxPipelineFunc()
}

func (m *MockRedis) Ping() *StatusCmd {
	return m.PingFunc()
}

func (m *MockRedis) Close() error {
	return m.CloseFunc()
}

func (m *MockRedis) Get(key string) *redis.StringCmd {
	return m.GetFunc(key)
}

func (m *MockRedis) MGet(keys ...string) *redis.SliceCmd {
	return m.MGetFunc(keys...)
}

func (m *MockRedis) Del(keys ...string) *redis.IntCmd {
	return m.DelFunc(keys...)
}
