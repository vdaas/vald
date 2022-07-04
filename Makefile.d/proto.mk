#
# Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
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
	openapi \
	jsonschema

.PHONY: pbgo
pbgo: $(PBGOS)

.PHONY: openapi
openapi: $(OPENAPISPECS)

.PHONY: jsonschema
jsonschema: $(OPENAPIJSONSCHEMAS)

.PHONY: pbdoc
pbdoc: $(PBDOCS)

.PHONY: proto/clean
## clean proto artifacts
proto/clean:
	rm -rf apis/grpc apis/openapi apis/docs

.PHONY: proto/paths/print
## print proto paths
proto/paths/print:
	@echo $(PROTO_PATHS)

.PHONY: proto/deps
## install protobuf dependencies
proto/deps: \
	$(GOPATH)/bin/protoc-gen-doc \
	$(GOPATH)/bin/protoc-gen-go \
	$(GOPATH)/bin/protoc-gen-go-grpc \
	$(GOPATH)/bin/protoc-gen-go-vtproto \
	$(GOPATH)/bin/protoc-gen-grpc-gateway \
	$(GOPATH)/bin/protoc-gen-jsonschema \
	$(GOPATH)/bin/protoc-gen-openapi \
	$(GOPATH)/bin/protoc-gen-validate \
	$(GOPATH)/bin/prototool \
	$(GOPATH)/bin/swagger \
	$(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate \
	$(GOPATH)/src/github.com/golang/protobuf \
	$(GOPATH)/src/github.com/googleapis/googleapis \
	$(GOPATH)/src/github.com/planetscale/vtprotobuf \
	$(GOPATH)/src/github.com/protocolbuffers/protobuf \
	$(GOPATH)/src/github.com/vdaas/vald/apis/proto/v1 \
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

$(GOPATH)/bin/protoc-gen-go:
	$(call go-install, google.golang.org/protobuf/cmd/protoc-gen-go)

$(GOPATH)/bin/protoc-gen-go-grpc:
	$(call go-install, google.golang.org/grpc/cmd/protoc-gen-go-grpc)

$(GOPATH)/bin/protoc-gen-grpc-gateway:
	$(call go-install, github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway)

$(GOPATH)/bin/protoc-gen-openapi:
	$(call go-install, github.com/google/gnostic/cmd/protoc-gen-openapi)

$(GOPATH)/bin/protoc-gen-jsonschema:
	$(call go-install, github.com/google/gnostic/cmd/protoc-gen-jsonschema)

$(GOPATH)/bin/protoc-gen-validate:
	$(call go-install, github.com/envoyproxy/protoc-gen-validate)

$(GOPATH)/bin/prototool:
	$(call go-install, github.com/uber/prototool/cmd/prototool)

$(GOPATH)/bin/protoc-gen-doc:
	$(call go-install, github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc)

$(GOPATH)/bin/protoc-gen-go-vtproto:
	$(call go-install, github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto)

$(GOPATH)/bin/swagger:
	$(call go-install, github.com/go-swagger/go-swagger/cmd/swagger)

$(PBGOS): \
	$(PROTOS) \
	proto/deps
	@$(call green, "generating pb.go files...")
	$(call mkdir, $(dir $@))
	$(call proto-code-gen, $(patsubst apis/grpc/%.pb.go,apis/proto/%.proto,$@))
	find $(ROOTDIR)/apis/grpc/* -name '*.go' | xargs sed -i -E "s%google.golang.org/grpc/codes%github.com/vdaas/vald/internal/net/grpc/codes%g"
	find $(ROOTDIR)/apis/grpc/* -name '*.go' | xargs sed -i -E "s%google.golang.org/grpc/status%github.com/vdaas/vald/internal/net/grpc/status%g"
	find $(ROOTDIR)/apis/grpc/* -name '*.go' | xargs sed -i -E "s%\"io\"%\"github.com/vdaas/vald/internal/io\"%g"

$(OPENAPISPECS): \
	$(PROTOS) \
	proto/deps
	@$(call green, "generating openapi.json files...")
	$(call mkdir, $(dir $@))
	$(call protoc-gen, $(patsubst apis/openapi/%.openapi.json,apis/proto/%.proto,$@), --openapi_out=enum_type=string:$(dir $@))

$(OPENAPIJSONSCHEMAS): \
	$(PROTOS) \
	proto/deps
	@$(call green, "generating schema.json files...")
	$(call mkdir, $(dir $@))
	$(call protoc-gen, $(patsubst apis/jsonschema/%.schema.json,apis/proto/%.proto,$@), --jsonschema_out=$(dir $@))

$(PBDOCS): \
	$(PROTOS) \
	proto/deps

apis/docs/v1/docs.md: $(PROTOS_V1)
	@$(call green, "generating documents for API v1...")
	$(call mkdir, $(dir $@))
	$(call protoc-gen, $(PROTOS_V1), --plugin=protoc-gen-doc=$(GOPATH)/bin/protoc-gen-doc --doc_opt=markdown$(COMMA)docs.md --doc_out=$(dir $@))
