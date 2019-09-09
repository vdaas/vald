.PHONY: \
    clean \
    bench \
    init \
    deps \
    images \
    profile \
    test \
    lint \
    contributors \
    coverage \
    create-index \
    core-bench \
    core-bench-lite \
    core-bench-clean \
    e2e-bench \
    proto-all \
    pbgo \
    swagger \
    graphql \
    pbdoc \
    pbpy \
    clean-proto-artifacts \
    proto-deps \
    grpcio-tools


REPO               ?= vdaas
GOPKG               = github.com/${REPO}/vald
TAG                 = $(shell date -u +%Y%m%d-%H%M%S)
AGENT_IMAGE         = vald-agnet
PROXY_IMAGE         = vald-proxy
DISCOVERER_IMAGE    = vald-discoverer
KVS_IMAGE           = vald-metadata

NGT_VERSION = 1.7.9
NGT_REPO = github.com/yahoojapan/NGT

GO_VERSION:=$(shell go version)

PROTODIRS := $(shell ls apis/proto)
PBGODIRS = $(PROTODIRS:%=apis/grpc/%)
SWAGGERDIRS = $(PROTODIRS:%=apis/swagger/%)
GRAPHQLDIRS = $(PROTODIRS:%=apis/graphql/%)
PBDOCDIRS = $(PROTODIRS:%=apis/docs/%)

PROTOS := $(shell find apis/proto -type f -regex ".*\.proto")
PBGOS = $(PROTOS:apis/proto/%.proto=apis/grpc/%.pb.go)
PBPYS = $(PROTOS:apis/proto/%.proto=apis/grpc/%_pb2.py)
GRPCPYS = $(PROTOS:apis/proto/%.proto=apis/grpc/%_pb2_grpc.py)
SWAGGERS = $(PROTOS:apis/proto/%.proto=apis/swagger/%.swagger.json)
GRAPHQLS = $(PROTOS:apis/proto/%.proto=apis/graphql/%.pb.graphqls)
GQLCODES = $(GRAPHQLS:apis/graphql/%.pb.graphqls=apis/graphql/%.generated.go)
PBDOCS = $(PROTOS:apis/proto/%.proto=apis/docs/%.md)

red    = /bin/echo -e "\x1b[31m\#\# $1\x1b[0m"
green  = /bin/echo -e "\x1b[32m\#\# $1\x1b[0m"
yellow = /bin/echo -e "\x1b[33m\#\# $1\x1b[0m"
blue   = /bin/echo -e "\x1b[34m\#\# $1\x1b[0m"
pink   = /bin/echo -e "\x1b[35m\#\# $1\x1b[0m"
cyan   = /bin/echo -e "\x1b[36m\#\# $1\x1b[0m"

define go-get
	GO111MODULE=off go get -u $1
endef

define mkdir
	mkdir -p $1
endef

PROTO_PATHS = \
	-I ./apis/proto/payload \
	-I ./apis/proto/agent \
	-I ./apis/proto/vald \
	-I ./apis/proto/discoverer \
	-I ./apis/proto/meta_manager \
	-I ./apis/proto/egress_filter \
	-I ./apis/proto/ingress_filter \
	-I ./apis/proto/backup_manager \
	-I ./apis/proto/traffic_manager \
	-I ./apis/proto/replication_manager \
	-I $(GOPATH)/src/github.com/protocolbuffers/protobuf/src \
	-I $(GOPATH)/src/github.com/gogo/protobuf/protobuf \
	-I $(GOPATH)/src/github.com/googleapis/googleapis \
	-I $(GOPATH)/src/github.com/danielvladco/go-proto-gql \
	-I $(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate

define protoc-gen
	protoc \
		$(PROTO_PATHS) \
		$2 \
		$1
endef

all:

clean:
	go clean -cache ./...
	# go clean -modcache
	rm -rf ./*.log
	rm -rf ./*.svg
	rm -rf ./go.mod
	rm -rf ./go.sum
	rm -rf bench
	rm -rf pprof
	rm -rf vendor
	rm -rf apis/docs
	rm -rf apis/grpc
	rm -rf apis/swagger
	rm -rf apis/graphql

bench:
	go test -count=5 -run=NONE -bench . -benchmem

init:
	GO111MODULE=on go mod init
	GO111MODULE=on go mod vendor

# deps: clean init
deps:
	go get github.com/envoyproxy/protoc-gen-validate \
		github.com/gogo/protobuf/gogoproto \
		github.com/gogo/protobuf/jsonpb \
		github.com/gogo/protobuf/proto \
		github.com/gogo/protobuf/protoc-gen-gogo \
		github.com/danielvladco/go-proto-gql \
		github.com/googleapis/googleapis
	rm -rf vendor
	curl -LO https://github.com/yahoojapan/NGT/archive/v${NGT_VERSION}.tar.gz
	tar zxf v${NGT_VERSION}.tar.gz -C /tmp
	# cd /tmp/NGT-${NGT_VERSION}&& cmake -DNGT_AVX_DISABLED=1 .
	cd /tmp/NGT-${NGT_VERSION}&& cmake .
	make -j -C /tmp/NGT-${NGT_VERSION}
	make install -C /tmp/NGT-${NGT_VERSION}
	rm -rf v${NGT_VERSION}.tar.gz
	rm -rf /tmp/NGT-${NGT_VERSION}

images:
	docker build -f dockers/agent/ngt/Dockerfile -t $(REPO)/$(AGENT_IMAGE) .
	# docker build -f dockers/proxy/gateway/vald/Dockerfile -t $(REPO)/$(PROXY_IMAGE) .

profile: clean
	mkdir pprof
	mkdir bench
	go test -count=10 -run=NONE -bench . -benchmem -o pprof/test.bin -cpuprofile pprof/cpu.out -memprofile pprof/mem.out
	go tool pprof --svg pprof/test.bin pprof/mem.out > bench/mem.svg
	go tool pprof --svg pprof/test.bin pprof/cpu.out > bench/cpu.svg

test: clean init
	GO111MODULE=on go test --race -coverprofile=cover.out ./...

lint:
	gometalinter --enable-all . | rg -v comment

contributors:
	git log --format='%aN <%aE>' | sort -fu > CONTRIBUTORS

coverage:
	go test -v -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	rm -f coverage.out

create-index:
	$(MAKE) -C ./hack/core/ngt create

core-bench: create-index
	$(MAKE) -C ./hach/core/ngt bench

core-bench-lite: create-index
	$(MAKE) -C ./hach/core/ngt bench-lite

core-bench-clean:
	$(MAKE) -C ./hack/core/ngt clean

e2e-bench: pbpy
	$(MAKE) -C ./hack/e2e/benchmark bench

proto-all: \
    pbgo \
    pbpy \
    pbdoc \
    swagger
    # swagger \
    # graphql

pbgo: $(PBGOS)
swagger: $(SWAGGERS)
graphql: $(GRAPHQLS) $(GQLCODES)
pbdoc: $(PBDOCS)
pbpy: $(PBPYS) $(GRPCPYS)

clean-proto-artifacts:
	rm -rf apis/grpc apis/swagger apis/graphql

proto-deps: \
    $(GOPATH)/bin/gqlgen \
    $(GOPATH)/bin/protoc-gen-doc \
    $(GOPATH)/bin/protoc-gen-go \
    $(GOPATH)/bin/protoc-gen-gofast \
    $(GOPATH)/bin/protoc-gen-gogofast \
    $(GOPATH)/bin/protoc-gen-gogqlgen \
    $(GOPATH)/bin/protoc-gen-gql \
    $(GOPATH)/bin/protoc-gen-gqlgencfg \
    $(GOPATH)/bin/protoc-gen-grpc-gateway \
    $(GOPATH)/bin/protoc-gen-swagger \
    $(GOPATH)/bin/protoc-gen-validate \
    $(GOPATH)/bin/prototool \
    $(GOPATH)/bin/swagger \
    $(GOPATH)/src/github.com/googleapis/googleapis \
    $(GOPATH)/src/github.com/protocolbuffers/protobuf \
    grpcio-tools

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

$(GOPATH)/bin/protoc-gen-go:
	$(call go-get, github.com/golang/protobuf/protoc-gen-go)

$(GOPATH)/bin/protoc-gen-gofast:
	$(call go-get, github.com/gogo/protobuf/protoc-gen-gofast)

$(GOPATH)/bin/protoc-gen-gogofast:
	$(call go-get, github.com/gogo/protobuf/protoc-gen-gogofast)

$(GOPATH)/bin/protoc-gen-grpc-gateway:
	$(call go-get, github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway)

$(GOPATH)/bin/protoc-gen-swagger:
	$(call go-get, github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger)

$(GOPATH)/bin/protoc-gen-gql:
	$(call go-get, github.com/danielvladco/go-proto-gql/protoc-gen-gql)

$(GOPATH)/bin/protoc-gen-gogqlgen:
	$(call go-get, github.com/danielvladco/go-proto-gql/protoc-gen-gogqlgen)

$(GOPATH)/bin/protoc-gen-gqlgencfg:
	$(call go-get, github.com/danielvladco/go-proto-gql/protoc-gen-gqlgencfg)

$(GOPATH)/bin/protoc-gen-validate:
	$(call go-get, github.com/envoyproxy/protoc-gen-validate)

$(GOPATH)/bin/prototool:
	$(call go-get, github.com/uber/prototool/cmd/prototool)

$(GOPATH)/bin/protoc-gen-doc:
	$(call go-get, github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc)

$(GOPATH)/bin/swagger:
	$(call go-get, github.com/go-swagger/go-swagger/cmd/swagger)

$(GOPATH)/bin/gqlgen:
	$(call go-get, github.com/99designs/gqlgen)

grpcio-tools:
	pip3 install -U grpcio-tools

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

$(PBGOS): proto-deps $(PBGODIRS)
	@$(call green, "generating pb.go files...")
	$(call protoc-gen, $(patsubst apis/grpc/%.pb.go,apis/proto/%.proto,$@), --gogofast_out=plugins=grpc:$(GOPATH)/src)
	# we have to enable validate after https://github.com/envoyproxy/protoc-gen-validate/pull/257 is merged
	# $(call protoc-gen, $(patsubst apis/grpc/%.pb.go,apis/proto/%.proto,$@), --gogofast_out=plugins=grpc:$(GOPATH)/src --validate_out=lang=gogo:$(GOPATH)/src)

$(PBPYS): proto-deps $(PBGODIRS)
	@$(call green, "generating pb2.py files...")
	python3 -m grpc_tools.protoc $(PROTO_PATHS) --python_out=$(dir $@) $(patsubst apis/grpc/%_pb2.py,apis/proto/%.proto,$@)

$(GRPCPYS): proto-deps $(PBGODIRS)
	@$(call green, "generating pb2_grpc.py files...")
	python3 -m grpc_tools.protoc $(PROTO_PATHS) --grpc_python_out=$(dir $@) $(patsubst apis/grpc/%_pb2_grpc.py,apis/proto/%.proto,$@)

$(SWAGGERS): proto-deps $(SWAGGERDIRS)
	@$(call green, "generating swagger.json files...")
	$(call protoc-gen, $(patsubst apis/swagger/%.swagger.json,apis/proto/%.proto,$@), --swagger_out=json_names_for_fields=true:$(dir $@))

$(GRAPHQLS): proto-deps $(GRAPHQLDIRS)
	@$(call green, "generating pb.graphqls files...")
	$(call protoc-gen, $(patsubst apis/graphql/%.pb.graphqls,apis/proto/%.proto,$@), --gql_out=paths=source_relative:$(dir $@))

$(GQLCODES): proto-deps $(GRAPHQLS)
	@$(call green, "generating graphql generated.go files...")
	sh hack/graphql/gqlgen.sh $(dir $@) $(patsubst apis/graphql/%.generated.go,apis/graphql/%.pb.graphqls,$@) $@

$(PBDOCS): proto-deps $(PBDOCDIRS)
	@$(call green, "generating documents files...")
	$(call protoc-gen, $(patsubst apis/docs/%.md,apis/proto/%.proto,$@), --plugin=protoc-gen-doc=$(GOPATH)/bin/protoc-gen-doc --doc_out=$(dir $@))
