#
# Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	pbgo \
	pbdoc \
	swagger

.PHONY: pbgo
pbgo: $(PBGOS)

.PHONY: swagger
swagger: $(SWAGGERS)

.PHONY: pbdoc
pbdoc: $(PBDOCS)

.PHONY: proto/clean
## clean proto artifacts
proto/clean:
	rm -rf apis/grpc apis/swagger apis/docs

.PHONY: proto/paths/print
## print proto paths
proto/paths/print:
	@echo $(PROTO_PATHS)

.PHONY: proto/deps
## install protobuf dependencies
proto/deps: \
	$(GOBIN)/protoc-gen-doc \
	$(GOBIN)/protoc-gen-go \
	$(GOBIN)/protoc-gen-go-grpc \
	$(GOBIN)/protoc-gen-go-vtproto \
	$(GOBIN)/protoc-gen-grpc-gateway \
	$(GOBIN)/protoc-gen-swagger \
	$(GOBIN)/protoc-gen-validate \
	$(GOBIN)/prototool \
	$(GOBIN)/swagger \
	$(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate \
	$(GOPATH)/src/github.com/golang/protobuf \
	$(GOPATH)/src/github.com/googleapis/googleapis \
	$(GOPATH)/src/github.com/planetscale/vtprotobuf \
	$(GOPATH)/src/github.com/protocolbuffers/protobuf \
	$(GOPATH)/src/google.golang.org/genproto \
	$(GOPATH)/src/google.golang.org/protobuf

.PHONY: proto/clean/deps
## uninstall all protobuf dependencies
proto/clean/deps:
	rm -rf $(GOBIN)/protoc-gen-doc \
	$(GOBIN)/protoc-gen-go \
	$(GOBIN)/protoc-gen-go-grpc \
	$(GOBIN)/protoc-gen-go-vtproto \
	$(GOBIN)/protoc-gen-grpc-gateway \
	$(GOBIN)/protoc-gen-swagger \
	$(GOBIN)/protoc-gen-validate \
	$(GOBIN)/prototool \
	$(GOBIN)/swagger \
	$(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate \
	$(GOPATH)/src/github.com/golang/protobuf \
	$(GOPATH)/src/github.com/googleapis/googleapis \
	$(GOPATH)/src/github.com/planetscale/vtprotobuf \
	$(GOPATH)/src/github.com/protocolbuffers/protobuf \
	$(GOPATH)/src/google.golang.org/genproto \
	$(GOPATH)/src/google.golang.org/protobuf


$(GOPATH)/src/github.com/protocolbuffers/protobuf:
	git clone \
		--depth 1 \
		https://github.com/protocolbuffers/protobuf \
		$(GOPATH)/src/github.com/protocolbuffers/protobuf

$(GOPATH)/src/github.com/googleapis/googleapis:
	git clone \
		--depth 1 \
		https://github.com/googleapis/googleapis \
		$(GOPATH)/src/github.com/googleapis/googleapis

$(GOPATH)/src/github.com/golang/protobuf:
	git clone \
		--depth 1 \
		https://github.com/golang/protobuf \
		$(GOPATH)/src/github.com/golang/protobuf

$(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate:
	git clone \
		--depth 1 \
		https://github.com/envoyproxy/protoc-gen-validate \
		$(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate

$(GOPATH)/src/google.golang.org/protobuf:
	git clone \
		--depth 1 \
		https://go.googlesource.com/protobuf \
		$(GOPATH)/src/google.golang.org/protobuf

$(GOPATH)/src/github.com/planetscale/vtprotobuf:
	git clone \
		--depth 1 \
		https://github.com/planetscale/vtprotobuf \
		$(GOPATH)/src/github.com/planetscale/vtprotobuf

$(GOPATH)/src/google.golang.org/genproto:
	git clone \
		--depth 1 \
		https://github.com/googleapis/go-genproto \
		$(GOPATH)/src/google.golang.org/genproto

$(GOBIN)/protoc-gen-go:
	$(call go-install, google.golang.org/protobuf/cmd/protoc-gen-go)

$(GOBIN)/protoc-gen-go-grpc:
	$(call go-install, google.golang.org/grpc/cmd/protoc-gen-go-grpc)

$(GOBIN)/protoc-gen-grpc-gateway:
	$(call go-install, github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway)

$(GOBIN)/protoc-gen-swagger:
	$(call go-install, github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger)

$(GOBIN)/protoc-gen-validate:
	$(call go-install, github.com/envoyproxy/protoc-gen-validate)

$(GOBIN)/prototool:
	$(call go-install, github.com/uber/prototool/cmd/prototool)

$(GOBIN)/protoc-gen-doc:
	$(call go-install, github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc)

$(GOBIN)/protoc-gen-go-vtproto:
	$(call go-install, github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto)

$(GOBIN)/swagger:
	$(call go-install, github.com/go-swagger/go-swagger/cmd/swagger)

$(ROOTDIR)/apis/proto/v1/rpc/error_details.proto:
	curl -fsSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/rpc/error_details.proto -o $(ROOTDIR)/apis/proto/v1/rpc/error_details.proto
	sed  -i -e "s/package google.rpc/package rpc.v1/" $(ROOTDIR)/apis/proto/v1/rpc/error_details.proto
	sed  -i -e "s%google.golang.org/genproto/googleapis/rpc/errdetails;errdetails%github.com/vdaas/vald/apis/grpc/v1/rpc/errdetails%" $(ROOTDIR)/apis/proto/v1/rpc/error_details.proto
	sed  -i -e "s/com.google.rpc/org.vdaas.vald.api.v1.rpc/" $(ROOTDIR)/apis/proto/v1/rpc/error_details.proto

$(PBGOS): \
	$(PROTOS) \
	proto/deps
	@$(call green, "generating pb.go files...")
	$(call mkdir, $(dir $@))
	$(call proto-code-gen, $(patsubst apis/grpc/%.pb.go,apis/proto/%.proto,$@))
	find $(ROOTDIR)/apis/grpc/* -name '*.go' | xargs sed -i -E "s%google.golang.org/grpc/codes%github.com/vdaas/vald/internal/net/grpc/codes%g"
	find $(ROOTDIR)/apis/grpc/* -name '*.go' | xargs sed -i -E "s%google.golang.org/grpc/status%github.com/vdaas/vald/internal/net/grpc/status%g"
	find $(ROOTDIR)/apis/grpc/* -name '*.go' | xargs sed -i -E "s%\"io\"%\"github.com/vdaas/vald/internal/io\"%g"
	find $(ROOTDIR)/apis/grpc/* -name '*.go' | xargs sed -i -E "s%\"sync\"%\"github.com/vdaas/vald/internal/sync\"%g"
	find $(ROOTDIR)/apis/grpc/* -name '*.go' | xargs sed -i -E "s/Vector = &Object_Vector\{\}/Vector = Object_VectorFromVTPool\(\)/g"
	find $(ROOTDIR)/apis/grpc/* -name '*.go' | xargs sed -i -E "s/v := &Object_Vector\{\}/v := Object_VectorFromVTPool\(\)/g"

$(SWAGGERS): \
	$(PROTOS) \
	proto/deps
	@$(call green, "generating swagger.json files...")
	$(call mkdir, $(dir $@))
	$(call protoc-gen, $(patsubst apis/swagger/%.swagger.json,apis/proto/%.proto,$@), --swagger_out=json_names_for_fields=true:$(dir $@))

$(PBDOCS): \
	$(PROTOS) \
	proto/deps

apis/docs/v1/docs.md: $(PROTOS_V1)
	@$(call green, "generating documents for API v1...")
	$(call mkdir, $(dir $@))
	$(call protoc-gen, $(PROTOS_V1), --plugin=protoc-gen-doc=$(GOBIN)/protoc-gen-doc --doc_opt=markdown$(COMMA)docs.md --doc_out=$(dir $@))
