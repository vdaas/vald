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

package service

import (
	"context"
	"log"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/pkg/gateway/vald/config"
	"google.golang.org/grpc"
)

type Gateway interface {
	GetIPs() []string
	BroadCast(ctx context.Context, eg errgroup.Group, f func(client agent.AgentClient) error) error
}

type Agent interface {
	BroadCast(ctx context.Context, f func(client agent.AgentClient) error)
}

type vp struct {
}

func New(cfg *config.Data) (Gateway, error) {
	ctx := context.Background()
	opts := []grpc.DialOption{
		grpc.WithContextDialer(nil),
		grpc.WithBackoffConfig(grpc.BackoffConfig{}),
	}
	conn, err := grpc.DialContext(ctx, "host", opts...)
	if err != nil {
		log.Fatal("client connection error:", err)
	}
	defer conn.Close()
	client := agent.NewAgentClient(conn)
	_ = client
	return new(vp), nil
}

func (v *vp) BroadCast(ctx context.Context, eg errgroup.Group, f func(client agent.AgentClient) error) error {
	return nil
}

func (v *vp) GetIPs() []string {
	return nil
}
