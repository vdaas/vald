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

package defaults

import (
	"slices"

	"github.com/vdaas/vald/internal/k8s/resource"
)

func (in *Defaults) DeepCopyInto(out *Defaults) {
	*out = *in
	out.Logging = resource.CopyPtr(in.Logging)
	out.ServerConfig = resource.CopyPtrInto(in.ServerConfig)
	out.Observability = resource.CopyPtrInto(in.Observability)
}

func (in *GRPC) DeepCopyInto(out *GRPC) {
	*out = *in
	out.Server = resource.CopyPtrInto(in.Server)
}

func (in *GRPCServer) DeepCopyInto(out *GRPCServer) {
	*out = *in
	out.GRPC = resource.CopyPtrInto(in.GRPC)
}

func (in *InterceptorConfig) DeepCopyInto(out *InterceptorConfig) {
	*out = *in
	out.Interceptors = slices.Clone(in.Interceptors)
}

func (in *Observability) DeepCopyInto(out *Observability) {
	*out = *in
	out.Trace = resource.CopyPtr(in.Trace)
	out.OTLP = resource.CopyPtr(in.OTLP)
}

func (in *ServerConfig) DeepCopyInto(out *ServerConfig) {
	*out = *in
	out.Servers = resource.CopyPtrInto(in.Servers)
}

func (in *Servers) DeepCopyInto(out *Servers) {
	*out = *in
	out.GRPC = resource.CopyPtrInto(in.GRPC)
}
