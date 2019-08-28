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

// Package model defines object structure
package model

import (
	"github.com/vdaas/vald/apis/grpc/agent"
	"google.golang.org/grpc"
)

type Agent struct {
	Client   agent.AgentClient
	Conn     *grpc.ClientConn
	IP       string
	Name     string
	CPU      float64
	Mem      float64
	HostIP   string
	HostName string
	HostCPU  float64
	HostMem  float64
}

type Agents []Agent

func (a Agents) Len() int {
	return len(a)
}

func (a Agents) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Agents) Less(i, j int) bool {
	return a[i].HostCPU < a[j].HostCPU &&
		a[i].HostMem < a[j].HostMem &&
		a[i].CPU < a[j].CPU &&
		a[i].Mem < a[j].Mem
}
