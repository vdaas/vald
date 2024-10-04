#
# Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# You may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
.PHONY: proto/all
## build protobufs
proto/all: \
	proto/deps \
	proto/gen

.PHONY: proto/clean
## clean proto artifacts
proto/clean:
	find $(ROOTDIR)/apis/grpc -name "*.pb.go" | xargs -P$(CORES) rm -f
	find $(ROOTDIR)/apis/grpc -name "*.pb.json.go" | xargs -P$(CORES) rm -f
	rm -rf $(ROOTDIR)/apis/swagger $(ROOTDIR)/apis/docs

.PHONY: proto/paths/print
## print proto paths
proto/paths/print:
	@echo $(PROTO_PATHS)

.PHONY: proto/deps
## install protobuf dependencies
proto/deps: \
	$(GOBIN)/buf

.PHONY: proto/clean/deps
## uninstall all protobuf dependencies
proto/clean/deps:
	rm -rf $(GOBIN)/buf

$(GOBIN)/buf:
	$(call go-install, github.com/bufbuild/buf/cmd/buf)

$(ROOTDIR)/apis/proto/v1/rpc/errdetails/error_details.proto:
	curl -fsSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/rpc/error_details.proto -o $(ROOTDIR)/apis/proto/v1/rpc/errdetails/error_details.proto
	sed  -i -e "s/package google.rpc/package rpc.v1/" $(ROOTDIR)/apis/proto/v1/rpc/errdetails/error_details.proto
	sed  -i -e "s%google.golang.org/genproto/googleapis/rpc/errdetails;errdetails%$(GOPKG)/apis/grpc/v1/rpc/errdetails%" $(ROOTDIR)/apis/proto/v1/rpc/errdetails/error_details.proto
	sed  -i -e "s/com.google.rpc/org.vdaas.vald.api.v1.rpc/" $(ROOTDIR)/apis/proto/v1/rpc/errdetails/error_details.proto

proto/gen: \
	$(PROTOS) \
	proto/deps
	@$(call green, "generating pb.go and swagger.json files and documents for API v1...")
	buf format -w
	buf generate
	make proto/replace

proto/replace:
	find $(ROOTDIR)/apis/grpc/* -name '*.go' | xargs -P$(CORES) sed -i -E "s%google.golang.org/grpc/codes%$(GOPKG)/internal/net/grpc/codes%g"
	find $(ROOTDIR)/apis/grpc/* -name '*.go' | xargs -P$(CORES) sed -i -E "s%google.golang.org/grpc/status%$(GOPKG)/internal/net/grpc/status%g"
	find $(ROOTDIR)/apis/grpc/* -name '*.go' | xargs -P$(CORES) sed -i -E "s%\"io\"%\"$(GOPKG)/internal/io\"%g"
	find $(ROOTDIR)/apis/grpc/* -name '*.go' | xargs -P$(CORES) sed -i -E "s%\"sync\"%\"$(GOPKG)/internal/sync\"%g"
	find $(ROOTDIR)/apis/grpc/* -name '*.go' | xargs -P$(CORES) sed -i -E "s%interface\{\}%any%g"
	find $(ROOTDIR)/apis/grpc/* -name '*.go' | xargs -P$(CORES) sed -i -E "s%For_%For%g"
