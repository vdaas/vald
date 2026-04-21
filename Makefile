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

SHELL = bash
.DEFAULT_GOAL := all
ORG ?= vdaas
NAME = vald
REPO = $(ORG)/$(NAME)
GOPKG = github.com/$(REPO)
DATETIME = $(eval DATETIME := $(shell date -u +%Y/%m/%d_%H:%M:%S%z))$(DATETIME)
TAG ?= latest
CRORG ?= $(ORG)
GHCRORG = ghcr.io/$(REPO)
AGENT_IMAGE = $(NAME)-agent
AGENT_FAISS_IMAGE = $(AGENT_IMAGE)-faiss
AGENT_NGT_IMAGE = $(AGENT_IMAGE)-ngt
AGENT_SIDECAR_IMAGE = $(AGENT_IMAGE)-sidecar
BENCHMARK_JOB_IMAGE = $(NAME)-benchmark-job
BENCHMARK_OPERATOR_IMAGE = $(NAME)-benchmark-operator
BINFMT_IMAGE = $(NAME)-binfmt
BUILDBASE_IMAGE = $(NAME)-buildbase
BUILDKIT_IMAGE = $(NAME)-buildkit
BUILDKIT_SYFT_SCANNER_IMAGE = $(BUILDKIT_IMAGE)-syft-scanner
DEV_CONTAINER_IMAGE = $(NAME)-dev-container
DISCOVERER_IMAGE = $(NAME)-discoverer-k8s
EXAMPLE_CLIENT_IMAGE = $(NAME)-example-client
FILTER_GATEWAY_IMAGE = $(NAME)-filter-gateway
HELM_OPERATOR_IMAGE = $(NAME)-helm-operator
INDEX_CORRECTION_IMAGE = $(NAME)-index-correction
INDEX_CREATION_IMAGE = $(NAME)-index-creation
INDEX_DELETION_IMAGE = $(NAME)-index-deletion
INDEX_EXPORTATION_IMAGE = $(NAME)-index-exportation
INDEX_OPERATOR_IMAGE = $(NAME)-index-operator
INDEX_SAVE_IMAGE = $(NAME)-index-save
LB_GATEWAY_IMAGE = $(NAME)-lb-gateway
MANAGER_INDEX_IMAGE = $(NAME)-manager-index
MIRROR_GATEWAY_IMAGE = $(NAME)-mirror-gateway
READREPLICA_ROTATE_IMAGE = $(NAME)-readreplica-rotate
E2E_IMAGE = $(NAME)-e2e
MAINTAINER = "$(ORG).org $(NAME) team <$(NAME)@$(ORG).org>"

XARGS_NO_RUN_IF_EMPTY := $(eval XARGS_NO_RUN_IF_EMPTY := $(shell xargs --version 2>/dev/null | head -1 | grep -qi gnu && echo -r))$(XARGS_NO_RUN_IF_EMPTY)
DEADLINK_CHECK_PATH ?= ""
DEADLINK_IGNORE_PATH ?= ""
DEADLINK_CHECK_FORMAT = html

DEFAULT_BUILDKIT_SYFT_SCANNER_IMAGE = $(GHCRORG)/$(BUILDKIT_SYFT_SCANNER_IMAGE):nightly

VERSION ?= $(eval VERSION := $(shell cat versions/VALD_VERSION))$(VERSION)

NGT_REPO = github.com/NGT-labs/NGT

NGT_EXTRA_CMAKE_FLAGS ?=

TEST_NOT_IMPL_PLACEHOLDER = NOT IMPLEMENTED BELOW

TEMP_DIR := $(eval TEMP_DIR := $(shell mktemp -d))$(TEMP_DIR)
USR_LOCAL = /usr/local
BINDIR = $(USR_LOCAL)/bin
LIB_PATH = $(USR_LOCAL)/lib
INCLUDE_PATH = $(USR_LOCAL)/include
$(BINDIR) $(LIB_PATH) $(INCLUDE_PATH):
	mkdir -p $@

BUN_INSTALL ?= $(USR_LOCAL)
BUN_GLOBAL_BIN ?= $(eval BUN_GLOBAL_BIN := $(shell command -v bun >/dev/null 2>&1 && bun pm bin -g 2>/dev/null || echo ""))$(BUN_GLOBAL_BIN)

GOPRIVATE := $(GOPKG),$(GOPKG)/apis,$(GOPKG)-client-go
GOPROXY ?= https://proxy.golang.org,direct
GOPATH ?= $(eval GOPATH := $(shell go env GOPATH 2>/dev/null))$(GOPATH)
GOARCH ?= $(eval GOARCH := $(shell go env GOARCH 2>/dev/null))$(GOARCH)
GOBIN ?= $(eval GOBIN := $(or $(shell go env GOBIN 2>/dev/null),$(GOPATH)/bin))$(GOBIN)
GOCACHE ?= $(eval GOCACHE := $(shell go env GOCACHE 2>/dev/null))$(GOCACHE)
GOOS ?= $(eval GOOS := $(shell go env GOOS 2>/dev/null))$(GOOS)
GOEXPERIMENT := "greenteagc,cgocheck2,newinliner,synchashtriemap,jsonv2"
GO_CLEAN_DEPS := true
GOTEST_TIMEOUT = 30m
CGO_ENABLED = 1
GODEBUG := gotestjsonbuildtext=1

RUST_HOME ?= $(LIB_PATH)/rust
RUSTUP_HOME ?= $(RUST_HOME)/rustup
CARGO_HOME ?= $(RUST_HOME)/cargo

BUF_VERSION := $(eval BUF_VERSION := $(shell cat versions/BUF_VERSION))$(BUF_VERSION)
BUSYBOX_VERSION := $(eval BUSYBOX_VERSION := $(shell cat versions/BUSYBOX_VERSION))$(BUSYBOX_VERSION)
CMAKE_VERSION := $(eval CMAKE_VERSION := $(shell cat versions/CMAKE_VERSION))$(CMAKE_VERSION)
CSI_DRIVER_HOST_PATH_VERSION := $(eval CSI_DRIVER_HOST_PATH_VERSION := $(shell cat versions/CSI_DRIVER_HOST_PATH_VERSION))$(CSI_DRIVER_HOST_PATH_VERSION)
DOCKER_VERSION := $(eval DOCKER_VERSION := $(shell cat versions/DOCKER_VERSION))$(DOCKER_VERSION)
FAISS_VERSION := $(eval FAISS_VERSION := $(shell cat versions/FAISS_VERSION))$(FAISS_VERSION)
GOLANGCILINT_VERSION := $(eval GOLANGCILINT_VERSION := $(shell cat versions/GOLANGCILINT_VERSION))$(GOLANGCILINT_VERSION)
GO_VERSION := $(eval GO_VERSION := $(shell cat versions/GO_VERSION))$(GO_VERSION)
GRAFANA_VERSION := $(eval GRAFANA_VERSION := $(shell cat versions/GRAFANA_VERSION))$(GRAFANA_VERSION)
HDF5_VERSION := $(eval HDF5_VERSION := $(shell cat versions/HDF5_VERSION))$(HDF5_VERSION)
HELM_DOCS_VERSION := $(eval HELM_DOCS_VERSION := $(shell cat versions/HELM_DOCS_VERSION))$(HELM_DOCS_VERSION)
HELM_VERSION := $(eval HELM_VERSION := $(shell cat versions/HELM_VERSION))$(HELM_VERSION)
JAEGER_OPERATOR_VERSION := $(eval JAEGER_OPERATOR_VERSION := $(shell cat versions/JAEGER_OPERATOR_VERSION))$(JAEGER_OPERATOR_VERSION)
K3D_VERSION := $(eval K3D_VERSION := $(shell cat versions/K3D_VERSION))$(K3D_VERSION)
K3S_VERSION := $(eval K3S_VERSION := $(shell cat versions/K3S_VERSION))$(K3S_VERSION)
KIND_VERSION := $(eval KIND_VERSION := $(shell cat versions/KIND_VERSION))$(KIND_VERSION)
KUBECTL_VERSION := $(eval KUBECTL_VERSION := $(shell cat versions/KUBECTL_VERSION))$(KUBECTL_VERSION)
KUBELINTER_VERSION := $(eval KUBELINTER_VERSION := $(shell cat versions/KUBELINTER_VERSION))$(KUBELINTER_VERSION)
LLVM_VERSION := $(eval LLVM_VERSION := $(shell cat versions/LLVM_VERSION))$(LLVM_VERSION)
NGT_VERSION := $(eval NGT_VERSION := $(shell cat versions/NGT_VERSION))$(NGT_VERSION)
NINJA_VERSION := $(eval NINJA_VERSION := $(shell cat versions/NINJA_VERSION))$(NINJA_VERSION)
OPERATOR_SDK_VERSION := $(eval OPERATOR_SDK_VERSION := $(shell cat versions/OPERATOR_SDK_VERSION))$(OPERATOR_SDK_VERSION)
OTEL_OPERATOR_VERSION := $(eval OTEL_OPERATOR_VERSION := $(shell cat versions/OTEL_OPERATOR_VERSION))$(OTEL_OPERATOR_VERSION)
PROMETHEUS_STACK_VERSION := $(eval PROMETHEUS_STACK_VERSION := $(shell cat versions/PROMETHEUS_STACK_VERSION))$(PROMETHEUS_STACK_VERSION)
PROTOBUF_VERSION := $(eval PROTOBUF_VERSION := $(shell cat versions/PROTOBUF_VERSION))$(PROTOBUF_VERSION)
REVIEWDOG_VERSION := $(eval REVIEWDOG_VERSION := $(shell cat versions/REVIEWDOG_VERSION))$(REVIEWDOG_VERSION)
RUST_VERSION := $(eval RUST_VERSION := $(shell cat versions/RUST_VERSION))$(RUST_VERSION)
SNAPSHOTTER_VERSION := $(eval SNAPSHOTTER_VERSION := $(shell cat versions/SNAPSHOTTER_VERSION))$(SNAPSHOTTER_VERSION)
TELEPRESENCE_VERSION := $(eval TELEPRESENCE_VERSION := $(shell cat versions/TELEPRESENCE_VERSION))$(TELEPRESENCE_VERSION)
USEARCH_VERSION := $(eval USEARCH_VERSION := $(shell cat versions/USEARCH_VERSION))$(USEARCH_VERSION)
YQ_VERSION := $(eval YQ_VERSION := $(shell cat versions/YQ_VERSION))$(YQ_VERSION)
ZLIB_VERSION := $(eval ZLIB_VERSION := $(shell cat versions/ZLIB_VERSION))$(ZLIB_VERSION)

OTEL_OPERATOR_RELEASE_NAME ?= opentelemetry-operator
PROMETHEUS_RELEASE_NAME ?= prometheus

SWAP_DEPLOYMENT_TYPE ?= deployment
SWAP_IMAGE ?= ""
SWAP_TAG ?= latest

UNAME := $(eval UNAME := $(shell uname -s))$(UNAME)
OS := $(eval OS := $(shell echo $(UNAME) | tr '[:upper:]' '[:lower:]'))$(OS)
ARCH := $(eval ARCH := $(shell uname -m))$(ARCH)
PWD := $(eval PWD := $(shell pwd))$(PWD)

ifeq ($(UNAME),Linux)
CPU_INFO_FLAGS := $(eval CPU_INFO_FLAGS := $(shell cat /proc/cpuinfo | grep flags | cut -d " " -f 2- | head -1))$(CPU_INFO_FLAGS)
CORES := $(eval CORES := $(shell nproc 2>/dev/null || getconf _NPROCESSORS_ONLN 2>/dev/null))$(CORES)
else ifeq ($(UNAME),Darwin)
CPU_INFO_FLAGS := $(eval CPU_INFO_FLAGS := $(shell sysctl -n machdep.cpu.brand_string 2>/dev/null || echo "Apple Silicon"))$(CPU_INFO_FLAGS)
CORES := $(eval CORES := $(shell sysctl -n hw.ncpu 2>/dev/null || getconf _NPROCESSORS_ONLN 2>/dev/null))$(CORES)
else
CPU_INFO_FLAGS := ""
CORES := 1
endif

GIT_COMMIT := $(eval GIT_COMMIT := $(shell git rev-list -1 HEAD))$(GIT_COMMIT)

ROOTDIR = $(eval ROOTDIR := $(or $(shell git rev-parse --show-toplevel), $(PWD)))$(ROOTDIR)
MAKELISTS := Makefile $(shell cat $(ROOTDIR)/.gitfiles | grep '^Makefile\.d/.*\.mk$$' | sed -e 's%^%$(ROOTDIR)/%')
PROTODIRS := $(eval PROTODIRS := $(shell cat $(ROOTDIR)/.gitfiles | grep '^apis/proto/.*\.proto$$' | xargs -n1 dirname | sort -u | sed -e "s%^apis/proto/%%g" | grep -v '^$$'))$(PROTODIRS)
BENCH_DATASET_BASE_DIR = hack/benchmark/assets
BENCH_DATASET_MD5_DIR_NAME = checksum
BENCH_DATASET_HDF5_DIR_NAME = dataset
BENCH_DATASET_MD5_DIR = $(BENCH_DATASET_BASE_DIR)/$(BENCH_DATASET_MD5_DIR_NAME)
BENCH_DATASET_HDF5_DIR = $(BENCH_DATASET_BASE_DIR)/$(BENCH_DATASET_HDF5_DIR_NAME)

PROTOS := $(eval PROTOS := $(shell cat $(ROOTDIR)/.gitfiles | grep '^apis/proto/.*\.proto$$' | sed -e 's%^%$(ROOTDIR)/%'))$(PROTOS)
PROTOS_V1 := $(eval PROTOS_V1 := $(filter apis/proto/v1/%.proto,$(PROTOS)))$(PROTOS_V1)
PBGOS = $(PROTOS:apis/proto/%.proto=apis/grpc/%.pb.go)
SWAGGERS = $(PROTOS:apis/proto/%.proto=apis/swagger/%.swagger.json)
PBDOCS = $(ROOTDIR)/apis/docs/v1/docs.md
PROTO_VALD_APIS := $(eval PROTO_VALD_APIS := $(filter $(ROOTDIR)/apis/proto/v1/vald/%.proto,$(PROTOS)))$(PROTO_VALD_APIS)
PROTO_VALD_API_DOCS := $(PROTO_VALD_APIS:$(ROOTDIR)/apis/proto/v1/vald/%.proto=$(ROOTDIR)/apis/docs/v1/%.md)
PROTO_MIRROR_APIS := $(eval PROTO_MIRROR_APIS := $(filter $(ROOTDIR)/apis/proto/v1/mirror/%.proto,$(PROTOS)))$(PROTO_MIRROR_APIS)
PROTO_MIRROR_API_DOCS := $(PROTO_MIRROR_APIS:$(ROOTDIR)/apis/proto/v1/mirror/%.proto=$(ROOTDIR)/apis/docs/v1/%.md)

select-binary = $(firstword $(wildcard $(BINDIR)/$(1)) $(wildcard /usr/bin/$(1)) $(shell which $(1) 2>/dev/null) $(1))

CC = $(call select-binary,clang)
CXX = $(call select-binary,clang++)
AR = $(call select-binary,llvm-ar)
NM = $(call select-binary,llvm-nm)
RANLIB = $(call select-binary,llvm-ranlib)
LLD = $(call select-binary,ld.lld)
NINJA = $(call select-binary,ninja)

export CC CXX AR NM RANLIB LDFLAGS CFLAGS CXXFLAGS SUDO

#
# Common flags
#
CFLAGS_BASE = -fPIC

OPT_CFLAGS = \
	-O3 \
	-fno-plt \
	-ffast-math \
	-ffp-contract=fast \
	-fmerge-all-constants \
	-funroll-loops \
	-falign-functions=32 \
	-ffunction-sections \
	-fdata-sections

ARCH_CFLAGS_amd64 = \
	-march=native \
	-mtune=native \
	-mno-avx512f \
	-mno-avx512dq \
	-mno-avx512cd \
	-mno-avx512bw \
	-mno-avx512vl

ARCH_CFLAGS_arm64 = \
	-march=armv8-a

#
# Linker flags only
#
LDFLAGS_BASE = \
	-L$(LIB_PATH) \
	-fuse-ld=$(LLD) \
	-pthread \
	-Wl,-z,relro \
	-Wl,-z,now \
	-Wl,--gc-sections

# Optional/common libs
LDLIBS_BASE = \
	-lm

LDFLAGS = $(LDFLAGS_BASE) $(LDLIBS_BASE)
NGT_LDFLAGS = -L$(LIB_PATH) -fopenmp -lopenblas -llapack -lgfortran
FAISS_LDFLAGS = $(NGT_LDFLAGS)
HDF5_LDFLAGS = $(LDFLAGS) -lhdf5 -lhdf5_hl -lz -ldl
CGO_LDFLAGS = $(FAISS_LDFLAGS) $(HDF5_LDFLAGS)
TEST_LDFLAGS = $(LDFLAGS) $(CGO_LDFLAGS)

GO_LDFLAGS = -s -w -extld '$(CXX)'
GO_STATIC_LDFLAGS = $(GO_LDFLAGS) -linkmode=external -extldflags '-static'

ifeq ($(GOARCH),amd64)
CFLAGS = $(CFLAGS_BASE) $(OPT_CFLAGS) $(ARCH_CFLAGS_amd64)
CXXFLAGS = $(CFLAGS) -std=gnu++23

ifeq ($(GOOS),darwin)
EXTLDFLAGS ?= -m64 -L$(LIB_PATH) -stdlib=libc++
else
EXTLDFLAGS ?= -m64 -L$(LIB_PATH) -stdlib=libc++ -Wl,--no-keep-memory
endif

else ifeq ($(GOARCH),arm64)

ifeq ($(GOOS),darwin)
CFLAGS = -I$(shell brew --prefix hdf5)/include
CGO_CFLAGS ?= $(CFLAGS)
CGO_LDFLAGS = -L$(shell brew --prefix hdf5)/lib -L$(shell brew --prefix zlib)/lib $(HDF5_LDFLAGS)
EXTLDFLAGS ?= -L$(LIB_PATH) -stdlib=libc++
CXXFLAGS ?= $(CFLAGS)
CXXFLAGS += -std=gnu++23
else
CFLAGS = $(CFLAGS_BASE) $(OPT_CFLAGS) $(ARCH_CFLAGS_arm64)
CXXFLAGS = $(CFLAGS) -std=gnu++23
EXTLDFLAGS ?= -march=armv8-a -L$(LIB_PATH) -Wl,--no-keep-memory
endif

else
CFLAGS = $(CFLAGS_BASE) $(OPT_CFLAGS)
CXXFLAGS = $(CFLAGS) -std=gnu++23

ifeq ($(GOOS),darwin)
EXTLDFLAGS ?= -L$(LIB_PATH) -stdlib=libc++
else
EXTLDFLAGS ?= -L$(LIB_PATH) -stdlib=libc++ -Wl,--no-keep-memory
endif
endif

BENCH_DATASET_MD5S := $(eval BENCH_DATASET_MD5S := $(shell cat $(ROOTDIR)/.gitfiles | grep '^$(BENCH_DATASET_MD5_DIR)/.*\.md5$$' | sed -e 's%^%$(ROOTDIR)/%'))$(BENCH_DATASET_MD5S)
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

HOST ?= localhost
PORT ?= 80
NUMBER ?= 10
DIMENSION ?= 6
NUMPANES ?= 4
MEAN ?= 0.0
STDDEV ?= 1.0

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
space := $() $()
comma := ,

# Directories to exclude from Go sources
GO_EXCLUDE_DIRS = \
	cmd/cli \
	internal/core/algorithm/ngt \
	internal/core/algorithm/faiss \
	internal/compress/gob \
	internal/compress/gzip \
	internal/compress/lz4 \
	internal/compress/zstd \
	internal/db/storage/blob/s3/sdk/s3 \
	internal/db/rdb/mysql/dbr \
	internal/test/comparator \
	internal/test/mock \
	hack/benchmark/internal/client/ngtd \
	hack/benchmark/internal/starter/agent \
	hack/benchmark/internal/starter/external \
	hack/benchmark/internal/starter/gateway \
	hack/gorules \
	hack/license \
	hack/docker \
	hack/swagger \
	hack/tools \
	tests

GO_EXCLUDE_PATTERN = ^($(subst $(space),|,$(GO_EXCLUDE_DIRS)))/

GO_SOURCES = $(eval GO_SOURCES := $(shell cat $(ROOTDIR)/.gitfiles | \
	grep -E '^(cmd|hack|internal|pkg)/' | \
	grep -v -E '$(GO_EXCLUDE_PATTERN)' | \
	grep '\.go$$' | \
	grep -v -E 'options?\.go$$|_test\.go$$|_mock\.go$$|doc\.go$$' \
	))$(GO_SOURCES)

GO_OPTION_SOURCES = $(eval GO_OPTION_SOURCES := $(shell cat $(ROOTDIR)/.gitfiles | \
	grep -E '^(cmd|hack|internal|pkg)/' | \
	grep -v -E '$(GO_EXCLUDE_PATTERN)' | \
	grep -E 'options?\.go$$' | \
	grep -v -E '_test\.go$$|_mock\.go$$|doc\.go$$' \
	))$(GO_OPTION_SOURCES)

GO_SOURCES_INTERNAL = $(eval GO_SOURCES_INTERNAL := $(shell cat $(ROOTDIR)/.gitfiles | \
	grep '^internal/' | \
	grep '\.go$$' | \
	grep -v -E '_test\.go$$|_mock\.go$$|doc\.go$$' \
	))$(GO_SOURCES_INTERNAL)

GO_TEST_SOURCES = $(GO_SOURCES:%.go=%_test.go)
GO_OPTION_TEST_SOURCES = $(GO_OPTION_SOURCES:%.go=%_test.go)

GO_ALL_TEST_SOURCES = $(GO_TEST_SOURCES) $(GO_OPTION_TEST_SOURCES)

DOCKER ?= docker
DOCKER_OPTS ?=
BUILDKIT_INLINE_CACHE ?= 1

DISTROLESS_IMAGE ?= gcr.io/distroless/static
DISTROLESS_IMAGE_TAG ?= nonroot
UPX_OPTIONS ?= -9
GOLINES_MAX_WIDTH ?= 200

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

# extra options to pass to textlint
# extra options to pass to cspell
E2E_CONFIG_DIR ?= $(ROOTDIR)/tests/v2/e2e/assets
E2E_CONFIG_NAME ?= unary_crud.yaml
E2E_CONFIG ?= $(E2E_CONFIG_DIR)/$(E2E_CONFIG_NAME)
E2E_ADDR ?= $(E2E_BIND_HOST):$(E2E_BIND_PORT)
E2E_BIND_HOST ?= 127.0.0.1
E2E_BIND_PORT ?= 8082
E2E_DATASET_NAME ?= fashion-mnist-784-euclidean.hdf5
E2E_GET_OBJECT_COUNT ?= 10
E2E_INSERT_COUNT ?= 60000
E2E_EXPECTED_INDEX ?= 180000
E2E_PARALLELISM ?= 10
E2E_QPS ?= 3000
E2E_SEARCH_COUNT ?= 1000
E2E_BULK_SIZE ?= 100
E2E_PORTFORWARD_ENABLED ?= true
E2E_REMOVE_COUNT ?= 3
E2E_SEARCH_BY_ID_COUNT ?= 100
E2E_TARGET_NAME ?= vald-lb-gateway
E2E_TARGET_NAMESPACE ?= default
E2E_TARGET_POD_NAME ?= $(eval E2E_TARGET_POD_NAME := $(shell kubectl get pods --selector=app=$(E2E_TARGET_NAME) -n $(E2E_TARGET_NAMESPACE) | tail -1 | cut -f1 -d " "))$(E2E_TARGET_POD_NAME)
E2E_TARGET_PORT ?= 8081
E2E_TIMEOUT ?= 60m
E2E_UPDATE_COUNT ?= 10
E2E_UPSERT_COUNT ?= 10
E2E_WAIT_FOR_CREATE_INDEX_DURATION ?= 8m
E2E_WAIT_FOR_START_TIMEOUT ?= 10m
E2E_SEARCH_FROM ?= 0
E2E_SEARCH_BY_ID_FROM ?= 0
E2E_INSERT_FROM ?= 0
E2E_UPDATE_FROM ?= 0
E2E_UPSERT_FROM ?= 0
E2E_GET_OBJECT_FROM ?= 0
E2E_REMOVE_FROM ?= 0

JAEGER_OPERATOR_WAIT_DURATION := 0

TEST_RESULT_DIR ?= /tmp

CERT_DIR ?= $(ROOTDIR)/internal/test/data/tls
DOMAIN ?= $(NAME).$(ORG).org
EMAIL ?= $(NAME)@$(ORG).org
STREET_ADDR ?= 1-3 Kioicho, Tokyo Garden Terrace Kioicho Tower, Chiyoda-ku, Tokyo 102-8282, Japan
CA_CN ?= $(DOMAIN) Root CA
CA_DAYS ?= 3650
END_ENTITY_DAYS ?= 825
CA_KEY := $(CERT_DIR)/ca.key
CA_CRT := $(CERT_DIR)/ca.crt
CA_PEM := $(CERT_DIR)/ca.pem
CA_SRL := $(CERT_DIR)/ca.srl

SERVER_KEY := $(CERT_DIR)/server.key
SERVER_CSR := $(CERT_DIR)/server.csr
SERVER_CRT := $(CERT_DIR)/server.crt

CLIENT_KEY := $(CERT_DIR)/client.key
CLIENT_CSR := $(CERT_DIR)/client.csr
CLIENT_CRT := $(CERT_DIR)/client.crt

INVALID_CA := $(CERT_DIR)/invalid-ca.pem
INVALID_SERVER_CRT := $(CERT_DIR)/invalid-server.crt

TEXTLINT_EXTRA_OPTIONS ?=

CSPELL_EXTRA_OPTIONS ?=

include Makefile.d/functions.mk

.PHONY: maintainer
## print maintainer
maintainer:
	@echo $(MAINTAINER)

.PHONY: help
## print all available commands
help:
	@awk 'BEGIN { \
		while (getline < ".gitfiles" > 0) if ($$0 !~ /^#/) files[$$0] = 1 \
			} \
	/^[a-zA-Z_0-9%:\\/-]+:/ { \
		if (match(lastLine, /^## (.*)/)) { \
			cmd = $$1; sub(/:+$$/, "", cmd); gsub(/\\\\/, "", cmd); \
			msg = substr(lastLine, RSTART + 3, RLENGTH); \
			if (index(cmd, "%") > 0 && match(msg, /\(([^)]+%[^)]*)\)/)) { \
				pat = substr(msg, RSTART + 1, RLENGTH - 2); \
				msg = substr(msg, 1, RSTART - 1) substr(msg, RSTART + RLENGTH); \
				sub(/^[ \t]+|[ \t]+$$/, "", msg); \
				pre = substr(pat, 1, index(pat, "%") - 1); \
				suf = substr(pat, index(pat, "%") + 1); \
				split("", o); \
				for (f in files) { \
					if (index(f, pre) == 1) { \
						r = substr(f, length(pre) + 1); \
						if (suf == "") { sub(/\/.*/, "", r); if (r != "") o[r] = 1 } \
						else if (index(r, suf) > 0) { sub(suf ".*", "", r); if (index(r, "/") == 0 && r != "") o[r] = 1 } \
							} \
				} \
				n = 0; split("", a); for (v in o) a[++n] = v; \
				for (i=1; i<=n; i++) for (j=i+1; j<=n; j++) if (a[i]>a[j]) { t=a[i]; a[i]=a[j]; a[j]=t } \
					s = ""; \
					for (i=1; i<=n; i++) s = s (s==""?"":", ") a[i]; \
						if (s != "") msg = msg " (Available: " s ")"; \
							} \
			printf "        \x1b[32;01m%-38s\x1b[0m %s\n", cmd, msg; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKELISTS) | sort -u
	@printf "\n"

.PHONY: perm
## set correct permissions for dirs and files
perm:
	@if [ -f "$(ROOTDIR)/.gitfiles" ]; then \
		grep -vE '^\s*#' "$(ROOTDIR)/.gitfiles" \
		| xargs -n1 dirname | sort -u \
		| sed -e 's%^%$(ROOTDIR)/%' \
		| xargs $(XARGS_NO_RUN_IF_EMPTY) -I {} -P"$(CORES)" chmod 755 "{}"; \
	fi
	@if [ -f "$(ROOTDIR)/.gitfiles" ]; then \
		grep -vE '^\s*#' "$(ROOTDIR)/.gitfiles" | grep -v gitignore \
		| xargs $(XARGS_NO_RUN_IF_EMPTY) -I {} -P"$(CORES)" chmod 644 "{}"; \
	fi
	if [ -d "$(ROOTDIR)/.git" ]; then \
		chmod 750 "$(ROOTDIR)/.git"; \
		if [ -f "$(ROOTDIR)/.git/config" ]; then \
			chmod 644 "$(ROOTDIR)/.git/config"; \
		fi; \
	if [ -d "$(ROOTDIR)/.git/hooks" ]; then \
	find "$(ROOTDIR)/.git/hooks" -type f -exec chmod 755 {} \;; \
	fi; \
	fi
	if [ -f "$(ROOTDIR)/.gitignore" ]; then \
		chmod 644 "$(ROOTDIR)/.gitignore"; \
	fi
	if [ -f "$(ROOTDIR)/.gitattributes" ]; then \
		chmod 644 "$(ROOTDIR)/.gitattributes"; \
	fi

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
	mv $(ROOTDIR)/apis/docs/buf.gen.*.yaml $(TEMP_DIR)/
	mv $(ROOTDIR)/apis/docs/v1/*.tmpl $(TEMP_DIR)/
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
	mkdir -p $(ROOTDIR)/apis/docs/v1
	mv $(TEMP_DIR)/buf.gen.*.yaml $(ROOTDIR)/apis/docs
	mv $(TEMP_DIR)/*.tmpl $(ROOTDIR)/apis/docs/v1

.PHONY: files
## add current repository file list to .gitfiles
files:
	@if [ ! -f $(ROOTDIR)/.gitfiles ]; then \
		printf '\n%.0s' {1..15} > $(ROOTDIR)/.gitfiles; \
	else \
		head -n 15 $(ROOTDIR)/.gitfiles > $(TEMP_DIR)/.gitfiles.tmp; \
		git ls-files --cached --others --exclude-standard | uniq >> $(TEMP_DIR)/.gitfiles.tmp; \
		mv $(TEMP_DIR)/.gitfiles.tmp $(ROOTDIR)/.gitfiles; \
	fi

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
	telepresence/install \
	textlint/install

.PHONY: update
## update deps, license, and run golines, gofumpt, goimports
update:
	- @$(MAKE) clean-generated
	- @$(MAKE) update/libs
	- @$(MAKE) update/actions
	- @$(MAKE) proto/all
	- @$(MAKE) deps
	- @$(MAKE) update/template
	- @$(MAKE) go/deps
	- @$(MAKE) go/example/deps
	- @$(MAKE) rust/deps
	- @$(MAKE) format

include Makefile.d/actions.mk
include Makefile.d/bench.mk
include Makefile.d/build.mk
include Makefile.d/dependencies.mk
include Makefile.d/docker.mk
include Makefile.d/e2e.mk
include Makefile.d/format.mk
include Makefile.d/generate.mk
include Makefile.d/git.mk
include Makefile.d/helm.mk
include Makefile.d/k3d.mk
include Makefile.d/k0s.mk
include Makefile.d/k8s.mk
include Makefile.d/kind.mk
include Makefile.d/lint.mk
include Makefile.d/minikube.mk
include Makefile.d/proto.mk
include Makefile.d/test.mk
include Makefile.d/tools.mk
include Makefile.d/tls.mk
include Makefile.d/version.mk