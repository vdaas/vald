#
# Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
	proto/gen/code \
	proto/gen/api/docs

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
	$(GOBIN)/buf \
	$(GOBIN)/protoc-gen-doc

.PHONY: proto/clean/deps
## uninstall all protobuf dependencies
proto/clean/deps:
	rm -rf $(GOBIN)/buf
	rm -rf $(GOBIN)/protoc-gen-doc

$(ROOTDIR)/apis/proto/v1/rpc/errdetails/error_details.proto:
	curl -fsSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/rpc/error_details.proto -o $(ROOTDIR)/apis/proto/v1/rpc/errdetails/error_details.proto
	sed -i -e "s/package google.rpc/package rpc.v1/" $(ROOTDIR)/apis/proto/v1/rpc/errdetails/error_details.proto
	sed -i -e "s%google.golang.org/genproto/googleapis/rpc/errdetails;errdetails%$(GOPKG)/apis/grpc/v1/rpc/errdetails%" $(ROOTDIR)/apis/proto/v1/rpc/errdetails/error_details.proto
	sed -i -e "s/com.google.rpc/org.vdaas.vald.api.v1.rpc/" $(ROOTDIR)/apis/proto/v1/rpc/errdetails/error_details.proto

.PHONY: proto/gen/code
## generate proto code
proto/gen/code: \
	$(PROTOS) \
	proto/deps
	@$(call green, "generating pb.go and swagger.json files and documents for API v1...")
	buf format -w
	buf generate
	make proto/replace

.PHONY: proto/gen/api/docs
## generate proto api docs
proto/gen/api/docs: \
	proto/gen/api/docs/payload \
	$(PROTO_VALD_API_DOCS) \
	$(PROTO_MIRROR_API_DOCS)

.PHONY: proto/gen/api/docs/payload
## generate proto api payload docs
proto/gen/api/docs/payload: $(ROOTDIR)/apis/docs/v1/payload.md.tmpl

$(ROOTDIR)/apis/docs/v1/payload.md.tmpl: $(ROOTDIR)/apis/proto/v1/payload/payload.proto $(ROOTDIR)/apis/docs/v1/payload.tmpl
	@$(call green,"generating payload v1...")
	buf generate --template=apis/docs/buf.gen.payload.yaml

$(ROOTDIR)/apis/docs/v1/%.md: $(ROOTDIR)/apis/proto/v1/vald/%.proto $(ROOTDIR)/apis/docs/v1/payload.md.tmpl $(ROOTDIR)/apis/docs/v1/doc.tmpl
	@$(call green,"generating documents for API v1...")
	@$(call gen-api-document,$@,$(subst $(ROOTDIR)/,,$<))

$(ROOTDIR)/apis/docs/v1/mirror.md: $(ROOTDIR)/apis/proto/v1/mirror/mirror.proto $(ROOTDIR)/apis/docs/v1/payload.md.tmpl $(ROOTDIR)/apis/docs/v1/doc.tmpl
	@$(call green,"generating documents for API v1...")
	@$(call gen-api-document,$@,$(subst $(ROOTDIR)/,,$<))

.PHONY: proto/replace
## replace generated proto code
proto/replace: \
	files
	@cat $(ROOTDIR)/.gitfiles | grep -E '^(\./)?apis/grpc/.*\.go$$' | xargs -I {} -P$(CORES) bash -c ' \
	echo "Replacing gRPC Go {}" && \
	sed -i -E "s%google.golang.org/grpc/codes%$(GOPKG)/internal/net/grpc/codes%g" {} && \
	sed -i -E "s%google.golang.org/grpc/status%$(GOPKG)/internal/net/grpc/status%g" {} && \
	sed -i -E "s%\"io\"%\"$(GOPKG)/internal/io\"%g" {} && \
	sed -i -E "s%\"sync\"%\"$(GOPKG)/internal/sync\"%g" {} && \
	sed -i -E "s%interface\{\}%any%g" {} && \
	sed -i -E "s%For_%For%g" {}'
	@echo "Proto file Replace complete."
