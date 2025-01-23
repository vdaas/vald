//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// Package health provides generic functionality for grpc health checks.
package health

import (
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Register register the generic gRPC health check server implementation to the srv.
func Register(srv *grpc.Server) {
	hsrv := health.NewServer()
	healthpb.RegisterHealthServer(srv, hsrv)
	for api := range srv.GetServiceInfo() {
		hsrv.SetServingStatus(api, healthpb.HealthCheckResponse_SERVING)
		log.Debug("gRPC health check server registered for service:\t" + api)
	}
	hsrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
}
