#
# Copyright (C) 2019 kpango (Yusuke Kato)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

.PHONY: \
	all \
	clean \
	license \
	bench \
	init \
	deps \
	go-version \
	ngt-version \
	ngt \
	images \
	dockers-base-image-name \
	dockers-base-image \
	dockers-agent-ngt-image-name \
	dockers-agent-ngt-image \
	dockers-discoverer-k8s-image-name \
	dockers-discoverer-k8s-image \
	dockers-gateway-vald-image-name \
	dockers-gateway-vald-image \
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
	pbdoc \
	swagger \
	graphql \
	bench-datasets \
	clean-bench-datasets \
	clean-proto-artifacts \
	proto-deps \
	bench-agent-stream \
	bench-agent-sequential-grpc \
	bench-agent-sequential-rest \
	bench-ngtd-stream \
	bench-ngtd-sequential-grpc \
	bench-ngtd-sequential-rest \
	profile-agent-stream \
	profile-agent-sequential-grpc \
	profile-agent-sequential-rest \
	profile-ngtd-stream \
	profile-ngtd-sequential-grpc \
	profile-ngtd-sequential-rest \
	kill-bench

REPO                   ?= vdaas
NAME                    = vald
GOPKG                   = github.com/$(REPO)/$(NAME)
TAG                     = $(shell date -u +%Y%m%d-%H%M%S)
BASE_IMAGE              = $(NAME)-base
AGENT_IMAGE             = $(NAME)-agent-ngt
GATEWAY_IMAGE           = $(NAME)-gateway
DISCOVERER_IMAGE        = $(NAME)-discoverer-k8s
KVS_IMAGE               = $(NAME)-meta-redis
BACKUP_MANAGER_IMAGE    = $(NAME)-manager-backup-mysql

NGT_VERSION := $(shell cat resources/NGT_VERSION)
NGT_REPO = github.com/yahoojapan/NGT

GO_VERSION := $(shell cat resources/GO_VERSION)

PROTODIRS := $(shell find apis/proto -type d | sed -e "s%apis/proto/%%g" | grep -v "apis/proto")
PBGODIRS = $(PROTODIRS:%=apis/grpc/%)
SWAGGERDIRS = $(PROTODIRS:%=apis/swagger/%)
GRAPHQLDIRS = $(PROTODIRS:%=apis/graphql/%)
PBDOCDIRS = $(PROTODIRS:%=apis/docs/%)

BENCH_DATASET_BASE_DIR = hack/e2e/benchmark/assets
BENCH_DATASET_MD5_DIR_NAME = checksum
BENCH_DATASET_HDF5_DIR_NAME = dataset
BENCH_DATASET_MD5_DIR = $(BENCH_DATASET_BASE_DIR)/$(BENCH_DATASET_MD5_DIR_NAME)
BENCH_DATASET_HDF5_DIR = $(BENCH_DATASET_BASE_DIR)/$(BENCH_DATASET_HDF5_DIR_NAME)

PROTOS := $(shell find apis/proto -type f -regex ".*\.proto")
PBGOS = $(PROTOS:apis/proto/%.proto=apis/grpc/%.pb.go)
SWAGGERS = $(PROTOS:apis/proto/%.proto=apis/swagger/%.swagger.json)
GRAPHQLS = $(PROTOS:apis/proto/%.proto=apis/graphql/%.pb.graphqls)
GQLCODES = $(GRAPHQLS:apis/graphql/%.pb.graphqls=apis/graphql/%.generated.go)
PBDOCS = $(PROTOS:apis/proto/%.proto=apis/docs/%.md)

BENCH_DATASET_MD5S := $(shell find $(BENCH_DATASET_MD5_DIR) -type f -regex ".*\.md5")
BENCH_DATASETS = $(BENCH_DATASET_MD5S:$(BENCH_DATASET_MD5_DIR)/%.md5=$(BENCH_DATASET_HDF5_DIR)/%.hdf5)

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
	$(PROTODIRS:%=-I ./apis/proto/%) \
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

define bench-pprof
	rm -rf $1
	mkdir -p $1
	@$(call green, "starting $4 $2 benchmark")
	go test -count=1 \
		-timeout=1h \
		-bench=$3 \
		-benchmem \
		-o $1/$2.bin \
		-cpuprofile $1/cpu-$4.out \
		-memprofile $1/mem-$4.out \
		-trace $1/trace-$4.out \
		$5 \
		| tee $1/result-$4.out
	go tool pprof --svg \
		$1/agent.bin \
		$1/cpu-$4.out \
		> $1/cpu-$4.svg
	go tool pprof --svg \
		$1/agent.bin \
		$1/mem-$4.out \
		> $1/mem-$4.svg
endef

define profile-web
	@$(call green, "starting $3 $2 profiler")
	go tool pprof -http=$4 \
		$1/$2.bin \
		$1/cpu-$3.out &
	go tool pprof -http=$5 \
		$1/$2.bin \
		$1/mem-$3.out &
	go tool trace -http=$6 \
		$1/trace-$3.out 
endef

all: clean deps

clean:
	# go clean -cache -modcache ./...
	rm -rf \
		/go/pkg \
		./*.log \
		./*.svg \
		./apis/docs \
		./apis/graphql \
		./apis/grpc \
		./apis/swagger \
		./bench \
		./pprof \
		./vendor \
		./go.sum \
		./go.mod
	cp ./hack/go.mod.default ./go.mod
		# ./go.mod \

license:
	go run hack/license/gen/main.go ./
	chmod -R 0644 ./*
	chmod -R 0644 ./.*

bench: \
	bench-agent-stream \
	bench-agent-sequential-grpc \
	bench-agent-sequential-rest

init:
	GO111MODULE=on go mod vendor

deps: \
	clean \
	proto-deps \
	proto-all
	go mod tidy
	go mod vendor
	rm -rf vendor

go-version:
	@echo $(GO_VERSION)

ngt-version:
	@echo $(NGT_VERSION)

ngt: /usr/local/include/NGT/Capi.h
/usr/local/include/NGT/Capi.h:
	curl -LO https://github.com/yahoojapan/NGT/archive/v$(NGT_VERSION).tar.gz
	tar zxf v$(NGT_VERSION).tar.gz -C /tmp
	cd /tmp/NGT-$(NGT_VERSION)&& cmake .
	make -j -C /tmp/NGT-$(NGT_VERSION)
	make install -C /tmp/NGT-$(NGT_VERSION)
	rm -rf v$(NGT_VERSION).tar.gz
	rm -rf /tmp/NGT-$(NGT_VERSION)

images: \
	dockers-base-image \
	dockers-agent-ngt-image \
	dockers-discoverer-k8s-image \
	dockers-gateway-vald-image

dockers-base-image-name:
	@echo "$(REPO)/$(BASE_IMAGE)"

dockers-base-image:
	docker build -f dockers/base/Dockerfile -t $(REPO)/$(BASE_IMAGE) .

dockers-agent-ngt-image-name:
	@echo "$(REPO)/$(AGENT_IMAGE)"

dockers-agent-ngt-image: dockers-base-image
	docker build -f dockers/agent/ngt/Dockerfile -t $(REPO)/$(AGENT_IMAGE) .

dockers-discoverer-k8s-image-name:
	@echo "$(REPO)/$(DISCOVERER_IMAGE)"

dockers-discoverer-k8s-image: dockers-base-image
	docker build -f dockers/discoverer/k8s/Dockerfile -t $(REPO)/$(DISCOVERER_IMAGE) .

dockers-gateway-vald-image-name:
	@echo "$(REPO)/$(GATEWAY_IMAGE)"

dockers-gateway-vald-image: dockers-base-image
	docker build -f dockers/gateway/vald/Dockerfile -t $(REPO)/$(GATEWAY_IMAGE) .

profile: \
	clean \
	deps \
	bench \
	profile-agent-stream \
	profile-agent-sequential-grpc \
	profile-agent-sequential-rest

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

e2e-bench:
	$(MAKE) -C ./hack/e2e/benchmark bench

proto-all: \
	pbgo \
	pbdoc \
	swagger
	# swagger \
	# graphql

pbgo: $(PBGOS)
swagger: $(SWAGGERS)
graphql: $(GRAPHQLS) $(GQLCODES)
pbdoc: $(PBDOCS)

bench-datasets: $(BENCH_DATASETS)

clean-bench-datasets:
	rm -rf $(BENCH_DATASETS)

clean-proto-artifacts:
	rm -rf apis/grpc apis/swagger apis/graphql

proto-deps: \
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
	$(call go-get, google.golang.org/genproto/...)

$(GOPATH)/bin/protoc-gen-go:
	$(call go-get, github.com/golang/protobuf/protoc-gen-go)

$(GOPATH)/bin/protoc-gen-gogo:
	$(call go-get, github.com/gogo/protobuf/protoc-gen-gogo)

$(GOPATH)/bin/protoc-gen-gofast:
	$(call go-get, github.com/gogo/protobuf/protoc-gen-gofast)

$(GOPATH)/bin/protoc-gen-gogofast:
	$(call go-get, github.com/gogo/protobuf/protoc-gen-gogofast)

$(GOPATH)/bin/protoc-gen-gogofaster:
	$(call go-get, github.com/gogo/protobuf/protoc-gen-gogofaster)

$(GOPATH)/bin/protoc-gen-gogoslick:
	$(call go-get, github.com/gogo/protobuf/protoc-gen-gogoslick)

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

$(BENCH_DATASETS): $(BENCH_DATASET_MD5S)
	@$(call green, "downloading datasets for benchmark...")
	curl -fsSL -o $@ http://vectors.erikbern.com/$(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,%.hdf5,$@)
	(cd $(BENCH_DATASET_BASE_DIR); \
	    md5sum -c $(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,$(BENCH_DATASET_MD5_DIR_NAME)/%.md5,$@) || \
	    (rm -f $(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,$(BENCH_DATASET_HDF5_DIR_NAME)/%.hdf5,$@) && exit 1))

bench-agent: \
	bench-agent-stream \
	bench-agent-sequential-grpc \
	bench-agent-sequential-rest

bench-agent-stream: \
	ngt \
	$(BENCH_DATASET_HDF5_DIR)/fashion-mnist-784-euclidean.hdf5 \
	$(BENCH_DATASET_HDF5_DIR)/mnist-784-euclidean.hdf5
	$(call bench-pprof,pprof/agent/ngt,agent,gRPCStream,stream,\
		./hack/e2e/benchmark/agent/ngt/ngt_bench_test.go \
		 -dataset=fashion-mnist)

bench-agent-sequential-grpc: \
	ngt \
	$(BENCH_DATASET_HDF5_DIR)/fashion-mnist-784-euclidean.hdf5 \
	$(BENCH_DATASET_HDF5_DIR)/mnist-784-euclidean.hdf5
	$(call bench-pprof,pprof/agent/ngt,agent,gRPCSequential,sequential-grpc,\
		./hack/e2e/benchmark/agent/ngt/ngt_bench_test.go \
		 -dataset=fashion-mnist)

bench-agent-sequential-rest: \
	ngt \
	$(BENCH_DATASET_HDF5_DIR)/fashion-mnist-784-euclidean.hdf5 \
	$(BENCH_DATASET_HDF5_DIR)/mnist-784-euclidean.hdf5
	$(call bench-pprof,pprof/agent/ngt,agent,RESTSequential,sequential-rest,\
		./hack/e2e/benchmark/agent/ngt/ngt_bench_test.go \
		 -dataset=fashion-mnist)

bench-ngtd: \
	bench-ngtd-stream \
	bench-ngtd-sequential-grpc \
	bench-ngtd-sequential-rest

bench-ngtd-stream: \
	ngt \
	$(BENCH_DATASET_HDF5_DIR)/fashion-mnist-784-euclidean.hdf5 \
	$(BENCH_DATASET_HDF5_DIR)/mnist-784-euclidean.hdf5
	$(call bench-pprof,pprof/external/ngtd,ngtd,gRPCStream,stream,\
		./hack/e2e/benchmark/external/ngtd/ngtd_bench_test.go \
		 -dataset=fashion-mnist)

bench-ngtd-sequential-grpc: \
	ngt \
	$(BENCH_DATASET_HDF5_DIR)/fashion-mnist-784-euclidean.hdf5 \
	$(BENCH_DATASET_HDF5_DIR)/mnist-784-euclidean.hdf5
	$(call bench-pprof,pprof/external/ngtd,ngtd,gRPCSequential,sequential-grpc,\
		./hack/e2e/benchmark/external/ngtd/ngtd_bench_test.go \
		 -dataset=fashion-mnist)

bench-ngtd-sequential-rest: \
	ngt \
	$(BENCH_DATASET_HDF5_DIR)/fashion-mnist-784-euclidean.hdf5 \
	$(BENCH_DATASET_HDF5_DIR)/mnist-784-euclidean.hdf5
	$(call bench-pprof,pprof/external/ngtd,ngtd,RESTSequential,sequential-rest,\
		./hack/e2e/benchmark/extenal/ngtd/ngtd_bench_test.go \
		 -dataset=fashion-mnist)

profile-agent-stream:
	$(call profile-web,pprof/agent/ngt,agent,stream,":6061",":6062",":6063")

profile-agent-sequential-grpc:
	$(call profile-web,pprof/agent/ngt,agent,sequential-grpc,":6061",":6062",":6063")

profile-agent-sequential-rest:
	$(call profile-web,pprof/agent/ngt,agent,sequential-rest,":6061",":6062",":6063")

profile-ngtd-stream:
	$(call profile-web,pprof/external/ngtd,ngtd,stream,":6061",":6062",":6063")

profile-ngtd-sequential-grpc:
	$(call profile-web,pprof/external/ngtd,ngtd,sequential-grpc,":6061",":6062",":6063")

profile-ngtd-sequential-rest:
	$(call profile-web,pprof/external/ngtd,ngtd,sequential-rest,":6061",":6062",":6063")


kill-bench:
	ps aux  \
		| grep go  \
		| grep -v nvim \
		| grep -v tmux \
		| grep -v gopls \
		| grep -v "rg go" \
		| grep -v "grep go" \
		| awk '{print $1}' \
		| xargs kill -9
