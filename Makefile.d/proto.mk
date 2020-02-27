#
# Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	pbgo \
	pbdoc \
	swagger
	# mv apis/docs/payload/docs.md ./docs.md
	# sed -i -E "s%#payload%apis/docs/payload/docs.md#payload%g" apis/docs
	# mv ./docs.md apis/docs/payload/docs.md
	# swagger \
	# graphql

.PHONY: pbgo
pbgo: $(PBGOS)

.PHONY: swagger
swagger: $(SWAGGERS)

.PHONY: graphql
graphql: $(GRAPHQLS) $(GQLCODES)

.PHONY: pbdoc
pbdoc: $(PBDOCS)

.PHONY: proto/clean
## clean proto artifacts
proto/clean:
	rm -rf apis/grpc apis/swagger apis/graphql apis/docs

.PHONY: proto/paths/print
## print proto paths
proto/paths/print:
	@echo $(PROTO_PATHS)

.PHONY: proto/deps
## install protobuf dependencies
proto/deps: \
	$(GOPATH)/bin/gqlgen \
	$(GOPATH)/bin/protoc-gen-doc \
	$(GOPATH)/bin/protoc-gen-go \
	$(GOPATH)/bin/protoc-gen-gogo \
	$(GOPATH)/bin/protoc-gen-gofast \
	$(GOPATH)/bin/protoc-gen-gogofast \
	$(GOPATH)/bin/protoc-gen-gogofaster \
	$(GOPATH)/bin/protoc-gen-gogoslick \
	$(GOPATH)/bin/protoc-gen-gogqlgen \
	$(GOPATH)/bin/protoc-gen-gql \
	$(GOPATH)/bin/protoc-gen-gqlgencfg \
	$(GOPATH)/bin/protoc-gen-grpc-gateway \
	$(GOPATH)/bin/protoc-gen-swagger \
	$(GOPATH)/bin/protoc-gen-validate \
	$(GOPATH)/bin/prototool \
	$(GOPATH)/bin/swagger \
	$(GOPATH)/src/google.golang.org/genproto \
	$(GOPATH)/src/github.com/protocolbuffers/protobuf \
	$(GOPATH)/src/github.com/googleapis/googleapis

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

$(GOPATH)/src/google.golang.org/genproto:
	$(call go-get-no-mod, google.golang.org/genproto/...)

$(GOPATH)/bin/protoc-gen-go:
	$(call go-get-no-mod, github.com/golang/protobuf/protoc-gen-go)

$(GOPATH)/bin/protoc-gen-gogo:
	$(call go-get-no-mod, github.com/gogo/protobuf/protoc-gen-gogo)

$(GOPATH)/bin/protoc-gen-gofast:
	$(call go-get-no-mod, github.com/gogo/protobuf/protoc-gen-gofast)

$(GOPATH)/bin/protoc-gen-gogofast:
	$(call go-get-no-mod, github.com/gogo/protobuf/protoc-gen-gogofast)

$(GOPATH)/bin/protoc-gen-gogofaster:
	$(call go-get-no-mod, github.com/gogo/protobuf/protoc-gen-gogofaster)

$(GOPATH)/bin/protoc-gen-gogoslick:
	$(call go-get-no-mod, github.com/gogo/protobuf/protoc-gen-gogoslick)

$(GOPATH)/bin/protoc-gen-grpc-gateway:
	$(call go-get-no-mod, github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway)

$(GOPATH)/bin/protoc-gen-swagger:
	$(call go-get-no-mod, github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger)

$(GOPATH)/bin/protoc-gen-gql:
	$(call go-get-no-mod, github.com/danielvladco/go-proto-gql/protoc-gen-gql)

$(GOPATH)/bin/protoc-gen-gogqlgen:
	$(call go-get-no-mod, github.com/danielvladco/go-proto-gql/protoc-gen-gogqlgen)

$(GOPATH)/bin/protoc-gen-gqlgencfg:
	$(call go-get-no-mod, github.com/danielvladco/go-proto-gql/protoc-gen-gqlgencfg)

$(GOPATH)/bin/protoc-gen-validate:
	$(call go-get-no-mod, github.com/envoyproxy/protoc-gen-validate)

$(GOPATH)/bin/prototool:
	$(call go-get-no-mod, github.com/uber/prototool/cmd/prototool)

$(GOPATH)/bin/protoc-gen-doc:
	$(call go-get-no-mod, github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc)

$(GOPATH)/bin/swagger:
	$(call go-get-no-mod, github.com/go-swagger/go-swagger/cmd/swagger)

$(GOPATH)/bin/gqlgen:
	$(call go-get, github.com/99designs/gqlgen)

$(PBGODIRS):
	$(call mkdir, $@)
	$(call rm, -rf, $@/*)

$(SWAGGERDIRS):
	$(call mkdir, $@)
	$(call rm, -rf, $@/*)

$(GRAPHQLDIRS):
	$(call mkdir, $@)
	$(call rm, -rf, $@/*)

$(PBDOCDIRS):
	$(call mkdir, $@)
	$(call rm, -rf, $@/*)

$(PBPYDIRS):
	$(call mkdir, $@)
	$(call rm, -rf, $@/*)

$(PBGOS): proto/deps $(PBGODIRS)
	@$(call green, "generating pb.go files...")
	$(call protoc-gen, $(patsubst apis/grpc/%.pb.go,apis/proto/%.proto,$@), --gogofast_out=plugins=grpc:$(GOPATH)/src)
	# we have to enable validate after https://github.com/envoyproxy/protoc-gen-validate/pull/257 is merged
	# $(call protoc-gen, $(patsubst apis/grpc/%.pb.go,apis/proto/%.proto,$@), --gogofast_out=plugins=grpc:$(GOPATH)/src --validate_out=lang=gogo:$(GOPATH)/src)

$(SWAGGERS): proto/deps $(SWAGGERDIRS)
	@$(call green, "generating swagger.json files...")
	$(call protoc-gen, $(patsubst apis/swagger/%.swagger.json,apis/proto/%.proto,$@), --swagger_out=json_names_for_fields=true:$(dir $@))

$(GRAPHQLS): proto/deps $(GRAPHQLDIRS)
	@$(call green, "generating pb.graphqls files...")
	$(call protoc-gen, $(patsubst apis/graphql/%.pb.graphqls,apis/proto/%.proto,$@), --gql_out=paths=source_relative:$(dir $@))

$(GQLCODES): proto/deps $(GRAPHQLS)
	@$(call green, "generating graphql generated.go files...")
	sh hack/graphql/gqlgen.sh $(dir $@) $(patsubst apis/graphql/%.generated.go,apis/graphql/%.pb.graphqls,$@) $@

$(PBDOCS): proto/deps $(PBDOCDIRS)
	@$(call green, "generating documents...")
	$(call protoc-gen, $(patsubst apis/docs/%.md,apis/proto/%.proto,$@), --plugin=protoc-gen-doc=$(GOPATH)/bin/protoc-gen-doc --doc_opt=markdown$(COMMA)docs.md --doc_out=$(dir $@))
