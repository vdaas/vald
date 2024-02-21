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

ORG                             ?= vdaas
NAME                            = vald
GOPKG                           = github.com/$(ORG)/$(NAME)
DATETIME                        = $(eval DATETIME := $(shell date -u +%Y/%m/%d_%H:%M:%S%z))$(DATETIME)
TAG                            ?= latest
CRORG                          ?= $(ORG)
GHCRORG                         = ghcr.io/$(ORG)/$(NAME)
AGENT_NGT_IMAGE                 = $(NAME)-agent-ngt
AGENT_FAISS_IMAGE               = $(NAME)-agent-faiss
AGENT_SIDECAR_IMAGE             = $(NAME)-agent-sidecar
CI_CONTAINER_IMAGE              = $(NAME)-ci-container
DEV_CONTAINER_IMAGE             = $(NAME)-dev-container
DISCOVERER_IMAGE                = $(NAME)-discoverer-k8s
FILTER_GATEWAY_IMAGE            = $(NAME)-filter-gateway
MIRROR_GATEWAY_IMAGE            = $(NAME)-mirror-gateway
HELM_OPERATOR_IMAGE             = $(NAME)-helm-operator
LB_GATEWAY_IMAGE                = $(NAME)-lb-gateway
LOADTEST_IMAGE                  = $(NAME)-loadtest
INDEX_CORRECTION_IMAGE          = $(NAME)-index-correction
INDEX_CREATION_IMAGE            = $(NAME)-index-creation
INDEX_SAVE_IMAGE                = $(NAME)-index-save
READREPLICA_ROTATE_IMAGE        = $(NAME)-readreplica-rotate
MANAGER_INDEX_IMAGE             = $(NAME)-manager-index
BENCHMARK_JOB_IMAGE             = $(NAME)-benchmark-job
BENCHMARK_OPERATOR_IMAGE        = $(NAME)-benchmark-operator
MAINTAINER                      = "$(ORG).org $(NAME) team <$(NAME)@$(ORG).org>"

VERSION ?= $(eval VERSION := $(shell cat versions/VALD_VERSION))$(VERSION)

NGT_REPO = github.com/yahoojapan/NGT

GOPRIVATE = $(GOPKG),$(GOPKG)/apis,$(GOPKG)-client-go
GOPROXY = "https://proxy.golang.org,direct"
GOPATH := $(eval GOPATH := $(shell go env GOPATH))$(GOPATH)
GO_VERSION := $(eval GO_VERSION := $(shell cat versions/GO_VERSION))$(GO_VERSION)
GOARCH := $(eval GOARCH := $(shell go env GOARCH))$(GOARCH)
GOBIN := $(eval GOBIN := $(or $(shell go env GOBIN),$(GOPATH)/bin))$(GOBIN)
GOCACHE := $(eval GOCACHE := $(shell go env GOCACHE))$(GOCACHE)
GOOS := $(eval GOOS := $(shell go env GOOS))$(GOOS)
GO_CLEAN_DEPS := true
GOTEST_TIMEOUT = 30m
CGO_ENABLED = 1

RUST_HOME = /usr/local/lib/rust
RUSTUP_HOME = $(RUST_HOME)/rustup
CARGO_HOME = $(RUST_HOME)/cargo

TEST_NOT_IMPL_PLACEHOLDER = NOT IMPLEMENTED BELOW

TEMP_DIR := $(eval TEMP_DIR := $(shell mktemp -d))$(TEMP_DIR)

BUF_VERSION               := $(eval BUF_VERSION := $(shell cat versions/BUF_VERSION))$(BUF_VERSION)
NGT_VERSION 		  := $(eval NGT_VERSION := $(shell cat versions/NGT_VERSION))$(NGT_VERSION)
FAISS_VERSION             := $(eval FAISS_VERSION := $(shell cat versions/FAISS_VERSION))$(FAISS_VERSION)
GOLANGCILINT_VERSION      := $(eval GOLANGCILINT_VERSION := $(shell cat versions/GOLANGCILINT_VERSION))$(GOLANGCILINT_VERSION)
HELM_DOCS_VERSION         := $(eval HELM_DOCS_VERSION := $(shell cat versions/HELM_DOCS_VERSION))$(HELM_DOCS_VERSION)
HELM_VERSION              := $(eval HELM_VERSION := $(shell cat versions/HELM_VERSION))$(HELM_VERSION)
JAEGER_OPERATOR_VERSION   := $(eval JAEGER_OPERATOR_VERSION := $(shell cat versions/JAEGER_OPERATOR_VERSION))$(JAEGER_OPERATOR_VERSION)
KIND_VERSION              := $(eval KIND_VERSION := $(shell cat versions/KIND_VERSION))$(KIND_VERSION)
KUBECTL_VERSION           := $(eval KUBECTL_VERSION := $(shell cat versions/KUBECTL_VERSION))$(KUBECTL_VERSION)
KUBELINTER_VERSION        := $(eval KUBELINTER_VERSION := $(shell cat versions/KUBELINTER_VERSION))$(KUBELINTER_VERSION)
NGT_VERSION               := $(eval NGT_VERSION := $(shell cat versions/NGT_VERSION))$(NGT_VERSION)
OPERATOR_SDK_VERSION      := $(eval OPERATOR_SDK_VERSION := $(shell cat versions/OPERATOR_SDK_VERSION))$(OPERATOR_SDK_VERSION)
OTEL_OPERATOR_VERSION     := $(eval OTEL_OPERATOR_VERSION := $(shell cat versions/OTEL_OPERATOR_VERSION))$(OTEL_OPERATOR_VERSION)
PROMETHEUS_STACK_VERSION  := $(eval PROMETHEUS_STACK_VERSION := $(shell cat versions/PROMETHEUS_STACK_VERSION))$(PROMETHEUS_STACK_VERSION)
PROTOBUF_VERSION          := $(eval PROTOBUF_VERSION := $(shell cat versions/PROTOBUF_VERSION))$(PROTOBUF_VERSION)
REVIEWDOG_VERSION         := $(eval REVIEWDOG_VERSION := $(shell cat versions/REVIEWDOG_VERSION))$(REVIEWDOG_VERSION)
TELEPRESENCE_VERSION      := $(eval TELEPRESENCE_VERSION := $(shell cat versions/TELEPRESENCE_VERSION))$(TELEPRESENCE_VERSION)
VALDCLI_VERSION           := $(eval VALDCLI_VERSION := $(shell cat versions/VALDCLI_VERSION))$(VALDCLI_VERSION)
YQ_VERSION                := $(eval YQ_VERSION := $(shell cat versions/YQ_VERSION))$(YQ_VERSION)
BUF_VERSION               := $(eval BUF_VERSION := $(shell cat versions/BUF_VERSION))$(BUF_VERSION)
ZLIB_VERSION              := $(eval ZLIB_VERSION := $(shell cat versions/ZLIB_VERSION))$(ZLIB_VERSION)
HDF5_VERSION              := $(eval HDF5_VERSION := $(shell cat versions/HDF5_VERSION))$(HDF5_VERSION)

OTEL_OPERATOR_RELEASE_NAME ?= opentelemetry-operator
PROMETHEUS_RELEASE_NAME    ?= prometheus

SWAP_DEPLOYMENT_TYPE ?= deployment
SWAP_IMAGE           ?= ""
SWAP_TAG             ?= latest

BINDIR ?= /usr/local/bin

UNAME := $(eval UNAME := $(shell uname -s))$(UNAME)
OS := $(eval OS := $(shell echo $(UNAME) | tr '[:upper:]' '[:lower:]'))$(OS)
ARCH := $(eval ARCH := $(shell uname -m))$(ARCH)
PWD := $(eval PWD := $(shell pwd))$(PWD)

ifeq ($(UNAME),Linux)
CPU_INFO_FLAGS := $(eval CPU_INFO_FLAGS := $(shell cat /proc/cpuinfo | grep flags | cut -d " " -f 2- | head -1))$(CPU_INFO_FLAGS)
CORES := $(eval CORES := $(shell nproc 2>/dev/null || getconf _NPROCESSORS_ONLN 2>/dev/null))$(CORES)
else ifeq ($(UNAME),Darwin)
CPU_INFO_FLAGS := ""
CORES := $(eval CORES := $(shell sysctl -n hw.ncpu 2>/dev/null || getconf _NPROCESSORS_ONLN 2>/dev/null))$(CORES)
else
CPU_INFO_FLAGS := ""
CORES := 1
endif

GIT_COMMIT := $(eval GIT_COMMIT := $(shell git rev-list -1 HEAD))$(GIT_COMMIT)

MAKELISTS := Makefile $(shell find Makefile.d -type f -regex ".*\.mk")

ROOTDIR = $(eval ROOTDIR := $(or $(shell git rev-parse --show-toplevel), $(PWD)))$(ROOTDIR)
PROTODIRS := $(eval PROTODIRS := $(shell find apis/proto -type d | sed -e "s%apis/proto/%%g" | grep -v "apis/proto"))$(PROTODIRS)
BENCH_DATASET_BASE_DIR = hack/benchmark/assets
BENCH_DATASET_MD5_DIR_NAME = checksum
BENCH_DATASET_HDF5_DIR_NAME = dataset
BENCH_DATASET_MD5_DIR = $(BENCH_DATASET_BASE_DIR)/$(BENCH_DATASET_MD5_DIR_NAME)
BENCH_DATASET_HDF5_DIR = $(BENCH_DATASET_BASE_DIR)/$(BENCH_DATASET_HDF5_DIR_NAME)

PROTOS := $(eval PROTOS := $(shell find apis/proto -type f -regex ".*\.proto"))$(PROTOS)
PROTOS_V1 := $(eval PROTOS_V1 := $(filter apis/proto/v1/%.proto,$(PROTOS)))$(PROTOS_V1)
PBGOS = $(PROTOS:apis/proto/%.proto=apis/grpc/%.pb.go)
SWAGGERS = $(PROTOS:apis/proto/%.proto=apis/swagger/%.swagger.json)
PBDOCS = apis/docs/v1/docs.md

ifeq ($(GOARCH),amd64)
CFLAGS ?= -mno-avx512f -mno-avx512dq -mno-avx512cd -mno-avx512bw -mno-avx512vl
CXXFLAGS ?= $(CFLAGS)
EXTLDFLAGS ?= -m64
else ifeq ($(GOARCH),arm64)
CFLAGS ?=
CXXFLAGS ?= $(CFLAGS)
EXTLDFLAGS ?= -march=armv8-a
else
CFLAGS ?=
CXXFLAGS ?= $(CFLAGS)
EXTLDFLAGS ?=
endif

BENCH_DATASET_MD5S := $(eval BENCH_DATASET_MD5S := $(shell find $(BENCH_DATASET_MD5_DIR) -type f -regex ".*\.md5"))$(BENCH_DATASET_MD5S)
BENCH_DATASETS = $(BENCH_DATASET_MD5S:$(BENCH_DATASET_MD5_DIR)/%.md5=$(BENCH_DATASET_HDF5_DIR)/%.hdf5)

BENCH_LARGE_DATASET_BASE_DIR = $(BENCH_DATASET_BASE_DIR)/large/dataset

SIFT1B_ROOT_DIR = $(BENCH_LARGE_DATASET_BASE_DIR)/sift1b

SIFT1B_BASE_FILE = $(SIFT1B_ROOT_DIR)/bigann_base.bvecs
SIFT1B_LEARN_FILE = $(SIFT1B_ROOT_DIR)/bigann_learn.bvecs
SIFT1B_QUERY_FILE = $(SIFT1B_ROOT_DIR)/bigann_query.bvecs
SIFT1B_GROUNDTRUTH_DIR = $(SIFT1B_ROOT_DIR)/gnd

SIFT1B_BASE_URL = ftp://ftp.irisa.fr/local/texmex/corpus/

DEEP1B_ROOT_DIR = $(BENCH_LARGE_DATASET_BASE_DIR)/deep1b

DEEP1B_BASE_FILE = $(DEEP1B_ROOT_DIR)/deep1B_base.fvecs
DEEP1B_LEARN_FILE = $(DEEP1B_ROOT_DIR)/deep1B_learn.fvecs
DEEP1B_QUERY_FILE = $(DEEP1B_ROOT_DIR)/deep1B_queries.fvecs
DEEP1B_GROUNDTRUTH_FILE = $(DEEP1B_ROOT_DIR)/deep1B_groundtruth.ivecs

DEEP1B_BASE_DIR = $(DEEP1B_ROOT_DIR)/base
DEEP1B_BASE_CHUNK_FILES = $(shell printf "$(DEEP1B_BASE_DIR)/base_%02d\n" {0..36})
DEEP1B_LEARN_DIR = $(DEEP1B_ROOT_DIR)/learn
DEEP1B_LEARN_CHUNK_FILES = $(shell printf "$(DEEP1B_LEARN_DIR)/learn_%02d\n" {0..13})

DEEP1B_API_URL = https://cloud-api.yandex.net/v1/disk/public/resources/download?public_key=https://yadi.sk/d/11eDCm7Dsn9GA&path=

DATASET_ARGS ?= identity-128
ADDRESS_ARGS ?= ""

HOST      ?= localhost
PORT      ?= 80
NUMBER    ?= 10
DIMENSION ?= 6
NUMPANES  ?= 4
MEAN      ?= 0.0
STDDEV    ?= 1.0

BODY = ""

PROTO_PATHS = \
	$(PWD) \
	$(GOPATH)/src \
	$(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate \
	$(GOPATH)/src/github.com/googleapis/googleapis \
	$(GOPATH)/src/github.com/planetscale/vtprotobuf \
	$(GOPATH)/src/github.com/protocolbuffers/protobuf \
	$(GOPATH)/src/google.golang.org/genproto \
	$(ROOTDIR) \
	$(ROOTDIR)/apis/proto/v1

# [Warning]
# The below packages have no original implementation.
# You should not add any features.
# - internal/copress/gob
# - internal/compress/gzip
# - internal/compress/lz4
# - internal/compress/zstd
# - internal/db/storage/blob/s3/sdk/s3
# - internal/db/rdb/mysql/dbr
# - internal/test/comparator
# - internal/test/mock
GO_SOURCES = $(eval GO_SOURCES := $(shell find \
		$(ROOTDIR)/cmd \
		$(ROOTDIR)/hack \
		$(ROOTDIR)/internal \
		$(ROOTDIR)/pkg \
		-not -path '$(ROOTDIR)/cmd/cli/*' \
		-not -path '$(ROOTDIR)/internal/core/algorithm/ngt/*' \
		-not -path '$(ROOTDIR)/internal/core/algorithm/faiss/*' \
		-not -path '$(ROOTDIR)/internal/compress/gob/*' \
		-not -path '$(ROOTDIR)/internal/compress/gzip/*' \
		-not -path '$(ROOTDIR)/internal/compress/lz4/*' \
		-not -path '$(ROOTDIR)/internal/compress/zstd/*' \
		-not -path '$(ROOTDIR)/internal/db/storage/blob/s3/sdk/s3/*' \
		-not -path '$(ROOTDIR)/internal/db/rdb/mysql/dbr/*' \
		-not -path '$(ROOTDIR)/internal/test/comparator/*' \
		-not -path '$(ROOTDIR)/internal/test/mock/*' \
		-not -path '$(ROOTDIR)/hack/benchmark/internal/client/ngtd/*' \
		-not -path '$(ROOTDIR)/hack/benchmark/internal/starter/agent/*' \
		-not -path '$(ROOTDIR)/hack/benchmark/internal/starter/external/*' \
		-not -path '$(ROOTDIR)/hack/benchmark/internal/starter/gateway/*' \
		-not -path '$(ROOTDIR)/hack/gorules/*' \
		-not -path '$(ROOTDIR)/hack/license/*' \
		-not -path '$(ROOTDIR)/hack/swagger/*' \
		-not -path '$(ROOTDIR)/hack/tools/*' \
		-not -path '$(ROOTDIR)/tests/*' \
		-type f \
		-name '*.go' \
		-not -regex '.*options?\.go' \
		-not -name '*_test.go' \
		-not -name '*_mock.go' \
		-not -name 'doc.go'))$(GO_SOURCES)
GO_OPTION_SOURCES = $(eval GO_OPTION_SOURCES := $(shell find \
		$(ROOTDIR)/cmd \
		$(ROOTDIR)/hack \
		$(ROOTDIR)/internal \
		$(ROOTDIR)/pkg \
		-not -path '$(ROOTDIR)/cmd/cli/*' \
		-not -path '$(ROOTDIR)/internal/core/algorithm/ngt/*' \
		-not -path '$(ROOTDIR)/internal/core/algorithm/faiss/*' \
		-not -path '$(ROOTDIR)/internal/compress/gob/*' \
		-not -path '$(ROOTDIR)/internal/compress/gzip/*' \
		-not -path '$(ROOTDIR)/internal/compress/lz4/*' \
		-not -path '$(ROOTDIR)/internal/compress/zstd/*' \
		-not -path '$(ROOTDIR)/internal/db/storage/blob/s3/sdk/s3/*' \
		-not -path '$(ROOTDIR)/internal/db/rdb/mysql/dbr/*' \
		-not -path '$(ROOTDIR)/internal/test/comparator/*' \
		-not -path '$(ROOTDIR)/internal/test/mock/*' \
		-not -path '$(ROOTDIR)/hack/benchmark/internal/client/ngtd/*' \
		-not -path '$(ROOTDIR)/hack/benchmark/internal/starter/agent/*' \
		-not -path '$(ROOTDIR)/hack/benchmark/internal/starter/external/*' \
		-not -path '$(ROOTDIR)/hack/benchmark/internal/starter/gateway/*' \
		-not -path '$(ROOTDIR)/hack/gorules/*' \
		-not -path '$(ROOTDIR)/hack/license/*' \
		-not -path '$(ROOTDIR)/hack/swagger/*' \
		-not -path '$(ROOTDIR)/hack/tools/*' \
		-not -path '$(ROOTDIR)/tests/*' \
		-type f \
		-regex '.*options?\.go' \
		-not -name '*_test.go' \
		-not -name '*_mock.go' \
		-not -name 'doc.go'))$(GO_OPTION_SOURCES)

GO_SOURCES_INTERNAL = $(eval GO_SOURCES_INTERNAL := $(shell find \
		$(ROOTDIR)/internal \
		-type f \
		-name '*.go' \
		-not -name '*_test.go' \
		-not -name '*_mock.go' \
		-not -name 'doc.go'))$(GO_SOURCES_INTERNAL)

GO_TEST_SOURCES = $(GO_SOURCES:%.go=%_test.go)
GO_OPTION_TEST_SOURCES = $(GO_OPTION_SOURCES:%.go=%_test.go)

GO_ALL_TEST_SOURCES = $(GO_TEST_SOURCES) $(GO_OPTION_TEST_SOURCES)

DOCKER                ?= docker
DOCKER_OPTS           ?=
BUILDKIT_INLINE_CACHE ?= 1

DISTROLESS_IMAGE      ?= gcr.io/distroless/static
DISTROLESS_IMAGE_TAG  ?= nonroot
UPX_OPTIONS           ?= -9
GOLINES_MAX_WIDTH     ?= 200

K8S_SLEEP_DURATION_FOR_WAIT_COMMAND ?= 5

ifeq ($(origin KUBECONFIG), undefined)
KUBECONFIG := $(HOME)/.kube/config
endif
K8S_KUBECTL_VERSION ?= $(eval K8S_KUBECTL_VERSION := $(shell kubectl version --short))$(K8S_KUBECTL_VERSION)
K8S_SERVER_VERSION ?= $(eval K8S_SERVER_VERSION := $(shell echo "$(K8S_KUBECTL_VERSION)" | sed -e "s/.*Server.*\(v[0-9]\.[0-9]*\)\..*/\1/g"))$(K8S_SERVER_VERSION)

# values file to use when deploying sample vald cluster with make k8s/vald/deploy
HELM_VALUES ?= $(ROOTDIR)/charts/vald/values/dev.yaml
# extra options to pass to helm when deploying sample vald cluster with make k8s/vald/deploy
HELM_EXTRA_OPTIONS ?=

COMMA := ,
SHELL = bash

E2E_BIND_HOST                      ?= 127.0.0.1
E2E_BIND_PORT                      ?= 8082
E2E_TIMEOUT                        ?= 30m
E2E_DATASET_NAME                   ?= fashion-mnist-784-euclidean.hdf5
E2E_INSERT_COUNT                   ?= 10000
E2E_SEARCH_COUNT                   ?= 1000
E2E_SEARCH_BY_ID_COUNT             ?= 100
E2E_GET_OBJECT_COUNT               ?= 10
E2E_UPDATE_COUNT                   ?= 10
E2E_UPSERT_COUNT                   ?= 10
E2E_REMOVE_COUNT                   ?= 3
E2E_WAIT_FOR_CREATE_INDEX_DURATION ?= 8m
E2E_TARGET_NAME                    ?= vald-lb-gateway
E2E_TARGET_POD_NAME                ?= $(eval E2E_TARGET_POD_NAME := $(shell kubectl get pods --selector=app=$(E2E_TARGET_NAME) -n $(E2E_TARGET_NAMESPACE) | tail -1 | cut -f1 -d " "))$(E2E_TARGET_POD_NAME)
E2E_TARGET_NAMESPACE               ?= default
E2E_TARGET_PORT                    ?= 8081
E2E_PORTFORWARD_ENABLED            ?= true

TEST_RESULT_DIR ?= /tmp

include Makefile.d/functions.mk

.PHONY: maintainer
## print maintainer
maintainer:
	@echo $(MAINTAINER)

.PHONY: help
## print all available commands
help:
	@awk '/^[a-zA-Z_0-9%:\\\/-]+:/ { \
	  helpMessage = match(lastLine, /^## (.*)/); \
	  if (helpMessage) { \
	    helpCommand = $$1; \
	    helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
      gsub("\\\\", "", helpCommand); \
      gsub(":+$$", "", helpCommand); \
	    printf "  \x1b[32;01m%-38s\x1b[0m %s\n", helpCommand, helpMessage; \
	  } \
	} \
	{ lastLine = $$0 }' $(MAKELISTS) | sort -u
	@printf "\n"

.PHONY: all
## execute clean and deps
all: clean deps

.PHONY: clean
## clean
clean: \
	clean-generated \
	proto/all \
	deps \
	format

.PHONY: clean-generated
## clean generated files
clean-generated:
	mv $(ROOTDIR)/apis/grpc/v1/vald/vald.go $(TEMP_DIR)/vald.go
	mv $(ROOTDIR)/apis/grpc/v1/agent/core/agent.go $(TEMP_DIR)/agent.go
	mv $(ROOTDIR)/apis/grpc/v1/payload/interface.go $(TEMP_DIR)/interface.go
	mv $(ROOTDIR)/apis/grpc/v1/mirror/mirror.go $(TEMP_DIR)/mirror.go
	rm -rf \
		$(ROOTDIR)/*.log \
		$(ROOTDIR)/*.svg \
		$(ROOTDIR)/apis/docs \
		$(ROOTDIR)/apis/swagger \
		$(ROOTDIR)/apis/grpc \
		$(ROOTDIR)/bench \
		$(ROOTDIR)/pprof \
		$(ROOTDIR)/libs
	mkdir -p $(ROOTDIR)/apis/grpc/v1/vald
	mv $(TEMP_DIR)/vald.go $(ROOTDIR)/apis/grpc/v1/vald/vald.go
	mkdir -p $(ROOTDIR)/apis/grpc/v1/agent/core
	mv $(TEMP_DIR)/agent.go $(ROOTDIR)/apis/grpc/v1/agent/core/agent.go
	mkdir -p $(ROOTDIR)/apis/grpc/v1/payload
	mv $(TEMP_DIR)/interface.go $(ROOTDIR)/apis/grpc/v1/payload/interface.go
	mkdir -p $(ROOTDIR)/apis/grpc/v1/mirror
	mv $(TEMP_DIR)/mirror.go $(ROOTDIR)/apis/grpc/v1/mirror/mirror.go

.PHONY: license
## add license to files
license:
	$(call gen-license,$(ROOTDIR),$(MAINTAINER))

.PHONY: init
## initialize development environment
init: \
	git/config/init \
	git/hooks/init \
	deps \
	ngt/install

.PHONY: tools/install
## install development tools
tools/install: \
	helm/install \
	kind/install \
	valdcli/install \
	telepresence/install \
	textlint/install

.PHONY: update
## update deps, license, and run golines, gofumpt, goimports
update: \
	clean-generated \
	update/libs \
	update/actions \
	proto/all \
	deps \
	update/template \
	go/deps \
	rust/deps \
	format

.PHONY: format
## format go codes
format: \
	license \
	format/proto \
	format/go \
	format/json \
	format/md \
	format/yaml

.PHONY: format/go
## run golines, gofumpt, goimports for all go files
format/go: \
	golines/install \
	gofumpt/install \
	strictgoimports/install \
	goimports/install
	find $(ROOTDIR)/ -type d -name .git -prune -o -type f -regex '.*[^\.pb]\.go' -print | xargs -P$(CORES) $(GOBIN)/golines -w -m $(GOLINES_MAX_WIDTH)
	find $(ROOTDIR)/ -type d -name .git -prune -o -type f -regex '.*[^\.pb]\.go' -print | xargs -P$(CORES) $(GOBIN)/gofumpt -w
	find $(ROOTDIR)/ -type d -name .git -prune -o -type f -regex '.*[^\.pb]\.go' -print | xargs -P$(CORES) $(GOBIN)/strictgoimports -w
	find $(ROOTDIR)/ -type d -name .git -prune -o -type f -regex '.*\.go' -print | xargs -P$(CORES) $(GOBIN)/goimports -w

.PHONY: format/go/test
## run golines, gofumpt, goimports for go test files
format/go/test: \
	golines/install \
	gofumpt/install \
	strictgoimports/install \
	goimports/install
	find $(ROOTDIR) -name '*_test.go' | xargs -P$(CORES) $(GOBIN)/golines -w -m $(GOLINES_MAX_WIDTH)
	find $(ROOTDIR) -name '*_test.go' | xargs -P$(CORES) $(GOBIN)/gofumpt -w
	find $(ROOTDIR) -name '*_test.go' | xargs -P$(CORES) $(GOBIN)/strictgoimports -w
	find $(ROOTDIR) -name '*_test.go' | xargs -P$(CORES) $(GOBIN)/goimports -w

.PHONY: format/yaml
format/yaml: \
	prettier/install\
	yamlfmt/install
	-find $(ROOTDIR) -name "*.yaml" -type f | grep -v templates | grep -v s3 | xargs -P$(CORES) -I {} prettier --write {}
	-find $(ROOTDIR) -name "*.yml" -type f | grep -v templates | grep -v s3 | xargs -P$(CORES) -I {} prettier --write {}
	-find $(ROOTDIR) -name "*.yaml" -type f | grep -v templates | grep -v s3 | xargs -P$(CORES) -I {} yamlfmt {}
	-find $(ROOTDIR) -name "*.yml" -type f | grep -v templates | grep -v s3 | xargs -P$(CORES) -I {} yamlfmt {}

.PHONY: format/md
format/md: \
	prettier/install
	prettier --write \
	    "$(ROOTDIR)/charts/**/*.md" \
	    "$(ROOTDIR)/apis/**/*.md" \
	    "$(ROOTDIR)/tests/**/*.md" \
	    "$(ROOTDIR)/*.md"

.PHONY: format/json
format/json: \
	prettier/install
	prettier --write \
	    "$(ROOTDIR)/apis/**/*.json" \
	    "$(ROOTDIR)/charts/**/*.json" \
	    "$(ROOTDIR)/hack/**/*.json"

.PHONY: format/proto
format/proto: \
	buf/install
	buf format -w

.PHONY: deps
## resolve dependencies
deps: \
	proto/deps \
	deps/install

.PHONY: deps/install
## install dependencies
deps/install: \
	golines/install \
	gofumpt/install \
	strictgoimports/install \
	goimports/install \
	prettier/install \
	go/deps \
	go/example/deps \
	rust/deps

.PHONY: version
## print vald version
version: \
	version/vald

.PHONY: version/vald
## print vald version
version/vald:
	@echo $(VERSION)

.PHONY: version/go
## print go version
version/go:
	@echo $(GO_VERSION)

.PHONY: version/ngt
## print NGT version
version/ngt:
	@echo $(NGT_VERSION)

.PHONY: version/k8s
## print Kubernetes version
version/k8s:
	@echo $(KUBECTL_VERSION)

.PHONY: version/kind
version/kind:
	@echo $(KIND_VERSION)

.PHONY: version/helm
version/helm:
	@echo $(HELM_VERSION)

.PHONY: version/yq
version/yq:
	@echo $(YQ_VERSION)

.PHONY: version/valdcli
version/valdcli:
	@echo $(VALDCLI_VERSION)

.PHONY: version/telepresence
version/telepresence:
	@echo $(TELEPRESENCE_VERSION)

.PHONY: ngt/install
## install NGT
ngt/install: /usr/local/include/NGT/Capi.h
/usr/local/include/NGT/Capi.h:
	curl -LO https://github.com/yahoojapan/NGT/archive/v$(NGT_VERSION).tar.gz
	tar zxf v$(NGT_VERSION).tar.gz -C $(TEMP_DIR)/
	cd $(TEMP_DIR)/NGT-$(NGT_VERSION) && \
		cmake -DCMAKE_C_FLAGS="$(CFLAGS)" -DCMAKE_CXX_FLAGS="$(CXXFLAGS)" .
	make -j -C $(TEMP_DIR)/NGT-$(NGT_VERSION)
	make install -C $(TEMP_DIR)/NGT-$(NGT_VERSION)
	rm -rf v$(NGT_VERSION).tar.gz
	rm -rf $(TEMP_DIR)/NGT-$(NGT_VERSION)
	ldconfig

.PHONY: faiss/install
## install Faiss
faiss/install: /usr/local/lib/libfaiss.so
/usr/local/lib/libfaiss.so:
	curl -LO https://github.com/facebookresearch/faiss/archive/v$(FAISS_VERSION).tar.gz
	tar zxf v$(FAISS_VERSION).tar.gz -C $(TEMP_DIR)/
	cd $(TEMP_DIR)/faiss-$(FAISS_VERSION) && \
		cmake -DFAISS_ENABLE_GPU=OFF -DFAISS_ENABLE_PYTHON=OFF -DBUILD_TESTING=OFF -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED_LIBS=ON -B build . && \
		make -C build -j faiss && \
		make -C build install
	rm -rf v$(FAISS_VERSION).tar.gz
	rm -rf $(TEMP_DIR)/faiss-$(FAISS_VERSION)
	ldconfig

.PHONY: lint
## run lints
lint: vet
	$(call go-lint)

.PHONY: vet
## run go vet
vet:
	$(call go-vet)

.PHONY: docs/lint
## run lint for document
docs/lint:
	textlint docs/**/*.md

.PHONY: changelog/update
## update changelog
changelog/update:
	echo "# CHANGELOG" > $(TEMP_DIR)/CHANGELOG.md
	echo "" >> $(TEMP_DIR)/CHANGELOG.md
	$(MAKE) -s changelog/next/print >> $(TEMP_DIR)/CHANGELOG.md
	echo "" >> $(TEMP_DIR)/CHANGELOG.md
	tail -n +2 CHANGELOG.md >> $(TEMP_DIR)/CHANGELOG.md
	mv -f $(TEMP_DIR)/CHANGELOG.md CHANGELOG.md

.PHONY: changelog/next/print
## print next changelog entry
changelog/next/print:
	@cat hack/CHANGELOG.template.md | \
	    sed -e 's/{{ version }}/$(VERSION)/g'
	@echo "$$BODY"

include Makefile.d/actions.mk
include Makefile.d/bench.mk
include Makefile.d/build.mk
include Makefile.d/client.mk
include Makefile.d/dependencies.mk
include Makefile.d/docker.mk
include Makefile.d/e2e.mk
include Makefile.d/git.mk
include Makefile.d/helm.mk
include Makefile.d/k3d.mk
include Makefile.d/k8s.mk
include Makefile.d/kind.mk
include Makefile.d/proto.mk
include Makefile.d/test.mk
include Makefile.d/tools.mk
include Makefile.d/minikube.mk
