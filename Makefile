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

SHELL = /usr/bin/env bash
.SHELLFLAGS := -eu -o pipefail -c
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
OPERATOR_IMAGE = $(NAME)-operator
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
$(LIB_PATH):
	mkdir -p $(LIB_PATH)

BUN_INSTALL ?= $(USR_LOCAL)
BUN_GLOBAL_BIN ?= $(eval BUN_GLOBAL_BIN := $(or $(shell bun pm bin -g 2>/dev/null),$(BINDIR)))$(BUN_GLOBAL_BIN)

GOPRIVATE := $(GOPKG),$(GOPKG)/apis,$(GOPKG)-client-go
GOPROXY := "https://proxy.golang.org,direct"
GOPATH ?= $(eval GOPATH := $(shell go env GOPATH))$(GOPATH)
GOARCH ?= $(eval GOARCH := $(shell go env GOARCH))$(GOARCH)
GOBIN ?= $(eval GOBIN := $(or $(shell go env GOBIN),$(GOPATH)/bin))$(GOBIN)
GOCACHE ?= $(eval GOCACHE := $(shell go env GOCACHE))$(GOCACHE)
GOOS ?= $(eval GOOS := $(shell go env GOOS))$(GOOS)
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
DOCKER_VERSION := $(eval DOCKER_VERSION := $(shell cat versions/DOCKER_VERSION))$(DOCKER_VERSION)
DOCKER_BUILDX_VERSION := $(eval DOCKER_BUILDX_VERSION := $(shell cat versions/DOCKER_BUILDX_VERSION))$(DOCKER_BUILDX_VERSION)
FAISS_VERSION := $(eval FAISS_VERSION := $(shell cat versions/FAISS_VERSION))$(FAISS_VERSION)
USEARCH_VERSION := $(eval USEARCH_VERSION := $(shell cat versions/USEARCH_VERSION))$(USEARCH_VERSION)
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
TELEPRESENCE_VERSION := $(eval TELEPRESENCE_VERSION := $(shell cat versions/TELEPRESENCE_VERSION))$(TELEPRESENCE_VERSION)
YQ_VERSION := $(eval YQ_VERSION := $(shell cat versions/YQ_VERSION))$(YQ_VERSION)
ZLIB_VERSION := $(eval ZLIB_VERSION := $(shell cat versions/ZLIB_VERSION))$(ZLIB_VERSION)
SNAPSHOTTER_VERSION := $(eval SNAPSHOTTER_VERSION := $(shell cat versions/SNAPSHOTTER_VERSION))$(SNAPSHOTTER_VERSION)
CSI_DRIVER_HOST_PATH_VERSION := $(eval CSI_DRIVER_HOST_PATH_VERSION := $(shell cat versions/CSI_DRIVER_HOST_PATH_VERSION))$(CSI_DRIVER_HOST_PATH_VERSION)

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
MAKELISTS := Makefile $(shell find $(ROOTDIR)/Makefile.d -type f -regex ".*\.mk")
PROTODIRS := $(eval PROTODIRS := $(shell find $(ROOTDIR)/apis/proto -type d | sed -e "s%apis/proto/%%g" | grep -v "apis/proto"))$(PROTODIRS)
BENCH_DATASET_BASE_DIR = hack/benchmark/assets
BENCH_DATASET_MD5_DIR_NAME = checksum
BENCH_DATASET_HDF5_DIR_NAME = dataset
BENCH_DATASET_MD5_DIR = $(BENCH_DATASET_BASE_DIR)/$(BENCH_DATASET_MD5_DIR_NAME)
BENCH_DATASET_HDF5_DIR = $(BENCH_DATASET_BASE_DIR)/$(BENCH_DATASET_HDF5_DIR_NAME)

PROTOS := $(eval PROTOS := $(shell find $(ROOTDIR)/apis/proto -type f -regex ".*\.proto"))$(PROTOS)
PROTOS_V1 := $(eval PROTOS_V1 := $(filter apis/proto/v1/%.proto,$(PROTOS)))$(PROTOS_V1)
PBGOS = $(PROTOS:apis/proto/%.proto=apis/grpc/%.pb.go)
SWAGGERS = $(PROTOS:apis/proto/%.proto=apis/swagger/%.swagger.json)
PBDOCS = $(ROOTDIR)/apis/docs/v1/docs.md
PROTO_VALD_APIS := $(eval PROTO_VALD_APIS := $(filter $(ROOTDIR)/apis/proto/v1/vald/%.proto,$(PROTOS)))$(PROTO_VALD_APIS)
PROTO_VALD_API_DOCS := $(PROTO_VALD_APIS:$(ROOTDIR)/apis/proto/v1/vald/%.proto=$(ROOTDIR)/apis/docs/v1/%.md)
PROTO_MIRROR_APIS := $(eval PROTO_MIRROR_APIS := $(filter $(ROOTDIR)/apis/proto/v1/mirror/%.proto,$(PROTOS)))$(PROTO_MIRROR_APIS)
PROTO_MIRROR_API_DOCS := $(PROTO_MIRROR_APIS:$(ROOTDIR)/apis/proto/v1/mirror/%.proto=$(ROOTDIR)/apis/docs/v1/%.md)

# Prefer the LLVM toolchain when present (containers install clang/lld/llvm
# via clangLTOBuildDeps); fall back to GCC elsewhere. Mirrors the LLD/NM/RANLIB
# detection in Makefile.d/functions.mk.
CC ?= $(shell command -v clang 2>/dev/null || command -v gcc 2>/dev/null || command -v cc)
CXX ?= $(shell command -v clang++ 2>/dev/null || command -v g++ 2>/dev/null || command -v c++)

# CI-built artifacts must not bake the builder's ISA: with LTO the LINK-stage
# -march decides the final codegen, so a native build on an AVX-512 builder
# silently defeats the -mno-avx512* compile guards below and crashes on
# runners without (or with quirky) AVX-512. The portable per-arch defaults
# (x86-64-v3 / armv8-a) are set in the GOARCH block below so that a literal
# x86 target name never leaks onto an arm64 build; -march/-mtune below are
# emitted only when MARCH/MTUNE are set. Override MARCH=native explicitly
# for local single-machine performance builds.

# ThinLTO for clang, classic fat-object LTO for GCC (-ffat-lto-objects is a
# GCC concept; clang's ThinLTO archives rely on llvm-ar instead).
LTO_FLAGS ?= $(if $(findstring clang,$(notdir $(CC))),-flto=thin,-flto=auto -ffat-lto-objects)
# lld understands ThinLTO archives natively; only wire it up for clang.
LLD_FLAGS ?= $(if $(and $(findstring clang,$(notdir $(CC))),$(LLD)),-fuse-ld=lld)

# Native Darwin builds use Homebrew's native dependencies. These variables are
# intentionally overridable: callers can resolve them as an unprivileged user
# and pass the values to a build that uses SUDO for installation.
ifeq ($(GOOS),darwin)
OPENMP_PREFIX ?= $(shell brew --prefix libomp 2>/dev/null)
HDF5_PREFIX ?= $(shell brew --prefix hdf5 2>/dev/null)
ZLIB_PREFIX ?= $(shell brew --prefix zlib 2>/dev/null)
OPENMP_CFLAGS ?= -I$(OPENMP_PREFIX)/include
NATIVE_LTO_FLAGS ?= -flto=thin

LDFLAGS = -fPIC -pthread -std=c++17 -lc++ -lm -L$(OPENMP_PREFIX)/lib -Wl,-rpath,$(USR_LOCAL)/lib -Wl,-rpath,$(OPENMP_PREFIX)/lib -lomp -framework Accelerate -lpthread
NGT_LDFLAGS =
HDF5_LDFLAGS = -lhdf5 -lhdf5_hl -lz -lm
CGO_LDFLAGS = -L$(HDF5_PREFIX)/lib -L$(ZLIB_PREFIX)/lib $(HDF5_LDFLAGS)
FAISS_LDFLAGS =
FAISS_CMAKE_C_FLAGS = $(CFLAGS)
FAISS_CMAKE_CXX_FLAGS = $(CXXFLAGS) $(NATIVE_LTO_FLAGS) $(OPENMP_CFLAGS)
FAISS_CMAKE_EXTRA_FLAGS = -DOpenMP_ROOT=$(OPENMP_PREFIX)
else
OPENMP_CFLAGS =
NATIVE_LTO_FLAGS ?= $(LTO_FLAGS)

# NOTE: -ffast-math must NOT appear here. On the link line the compiler driver
# pulls in crtfastmath.o, whose global constructor (set_fast_math, sets the
# MXCSR FTZ/DAZ bits) runs before main and corrupts libgcc's static C++
# exception-frame registration in this large statically-linked cgo binary
# (frame_dummy -> __register_frame_info -> classify_object_over_fdes SIGSEGV,
# reproduced locally). It buys no optimization at link time; per-TU fast-math
# for the hot C++ code already comes from NGT/faiss's own -Ofast.
LDFLAGS = -static -fPIC -pthread -std=gnu++23 -lstdc++ -lm -z relro -z now $(LTO_FLAGS) $(LLD_FLAGS) $(if $(MARCH),-march=$(MARCH)) $(if $(MTUNE),-mtune=$(MTUNE)) -fno-plt -O3 -fvisibility=hidden -ffp-contract=fast -fomit-frame-pointer -fmerge-all-constants -funroll-loops -falign-functions=32 -ffunction-sections -fdata-sections -Wl,--whole-archive -lpthread -Wl,--no-whole-archive
NGT_LDFLAGS = -fopenmp -lopenblas -llapack -lgfortran
# Resolves a shared libomp.so path for CMake's Linux OpenMP probe.
LIBOMP ?= $(shell ldconfig -p 2>/dev/null | awk '/libomp\.so[^.].*=>/{print $$NF; exit}' | grep -v '^$$' || ls /usr/lib/llvm-*/lib/libomp.so 2>/dev/null | sort -V | tail -1)
HDF5_LDFLAGS = -lhdf5 -lhdf5_hl -lsz -laec -lz -ldl -lm
CGO_LDFLAGS = $(FAISS_LDFLAGS) $(HDF5_LDFLAGS)
FAISS_LDFLAGS = $(NGT_LDFLAGS)
FAISS_CMAKE_C_FLAGS = $(CFLAGS) $(LTO_FLAGS) $(if $(MARCH),-march=$(MARCH)) $(if $(MTUNE),-mtune=$(MTUNE)) -fopenmp
FAISS_CMAKE_CXX_FLAGS = $(CXXFLAGS) $(LTO_FLAGS) $(if $(MARCH),-march=$(MARCH)) $(if $(MTUNE),-mtune=$(MTUNE)) -fopenmp
FAISS_CMAKE_EXTRA_FLAGS = -DBLA_VENDOR=OpenBLAS
endif

# TEST_LDFLAGS without -static to avoid conflicts with CGO and glibc dynamic linking requirements
ifeq ($(GOOS),darwin)
TEST_LDFLAGS_BASE = $(LDFLAGS)
else
TEST_LDFLAGS_BASE = -fPIC -pthread -std=gnu++23 -lstdc++ -lm -z relro -z now $(LTO_FLAGS) $(LLD_FLAGS) $(if $(MARCH),-march=$(MARCH)) $(if $(MTUNE),-mtune=$(MTUNE)) -fno-plt -O3 -fvisibility=hidden -ffp-contract=fast -fomit-frame-pointer -fmerge-all-constants -funroll-loops -falign-functions=32 -ffunction-sections -fdata-sections
endif
TEST_LDFLAGS = $(TEST_LDFLAGS_BASE) $(CGO_LDFLAGS)

ifeq ($(GOARCH),amd64)
MARCH ?= x86-64-v3
MTUNE ?= generic
CFLAGS ?= -mno-avx512f -mno-avx512dq -mno-avx512cd -mno-avx512bw -mno-avx512vl
CXXFLAGS ?= $(CFLAGS)
ifeq ($(GOOS),darwin)
EXTLDFLAGS ?= -m64
else
EXTLDFLAGS ?= -m64 -Wl,--no-keep-memory
endif
else ifeq ($(GOARCH),arm64)
# armv8-a is the portable arm64 baseline (analogous to x86-64-v3 on amd64);
# never emit the x86-only x86-64-v3 target name here.
MARCH ?= armv8-a
MTUNE ?= generic
CFLAGS ?=
ifeq ($(GOOS),darwin)
EXTLDFLAGS ?= -march=armv8-a
else
EXTLDFLAGS ?= -march=armv8-a -Wl,--no-keep-memory
endif
CXXFLAGS ?= $(CFLAGS)
else
CFLAGS ?=
CXXFLAGS ?= $(CFLAGS)
ifeq ($(GOOS),darwin)
EXTLDFLAGS ?=
else
EXTLDFLAGS ?= -Wl,--no-keep-memory
endif
endif

ifeq ($(GOOS),darwin)
# Keep the C++ standard explicit in cgo invocations; Faiss requires it when
# compiling through Go rather than through its CMake-generated flags.
CFLAGS := $(CFLAGS) -I$(HDF5_PREFIX)/include -I$(ZLIB_PREFIX)/include
CXXFLAGS := $(CXXFLAGS) -std=c++17
CGO_CFLAGS ?= $(CFLAGS) $(OPENMP_CFLAGS)
CGO_CXXFLAGS ?= $(CXXFLAGS) $(OPENMP_CFLAGS)
export CGO_CFLAGS CGO_CXXFLAGS
endif

# Base compile flags (the arch-gated guards above, without the per-build
# extras appended via CGO_ENV_VARS' $1). Consumed by libomp/install in
# Makefile.d/tools.mk; previously referenced but never defined, so that
# recipe built with empty flags.
CFLAGS_BASE ?= $(CFLAGS)

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
	-not -path '$(ROOTDIR)/hack/gorules/*' \
	-not -path '$(ROOTDIR)/hack/license/*' \
	-not -path '$(ROOTDIR)/hack/docker/*' \
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
	-not -path '$(ROOTDIR)/hack/gorules/*' \
	-not -path '$(ROOTDIR)/hack/license/*' \
	-not -path '$(ROOTDIR)/hack/docker/*' \
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

DOCKER ?= docker
DOCKER_OPTS ?=
DOCKER_BUILDKIT ?= 1
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
TEXTLINT_EXTRA_OPTIONS ?=
# extra options to pass to cspell
CSPELL_EXTRA_OPTIONS ?=

COMMA := ,
SHELL = bash

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
E2E_MAX_DIM_BIT ?= 1
E2E_MAX_DIM_CONFIG_NAME ?= max_vector_dim.yaml
E2E_MAX_DIM_CONFIG ?= $(E2E_CONFIG_DIR)/$(E2E_MAX_DIM_CONFIG_NAME)
E2E_MAX_DIM_WAIT ?= 30s
E2E_MAX_DIM_RETRY_TIMEOUT ?= 5m
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
	printf "	\x1b[32;01m%-38s\x1b[0m %s\n", helpCommand, helpMessage; \
	} \
	} \
	{ lastLine = $$0 }' $(MAKELISTS) | sort -u
	@printf "\n"

.PHONY: perm
## set correct permissions for dirs and files
perm:
	find $(ROOTDIR) -type d -not -path "$(ROOTDIR)/.git*" -exec chmod 755 {} \;
	@if [ -f "$(ROOTDIR)/.gitfiles" ]; then \
		grep -vE '^\s*#' "$(ROOTDIR)/.gitfiles" | grep -v gitignore \
		| xargs $(XARGS_NO_RUN_IF_EMPTY) -I {} -P"$(CORES)" bash -c '[ -f "{}" ] || exit 0; chmod 644 "{}"'; \
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

.PHONY: deps
## resolve dependencies
deps: \
	proto/deps \
	deps/install

.PHONY: deps/install
## install dependencies
deps/install: \
	crlfmt/install \
	golines/install \
	gofumpt/install \
	strictgoimports/install \
	goimports/install \
	prettier/install \
	go/deps \
	go/example/deps \
	rust/deps

include Makefile.d/actions.mk
include Makefile.d/bench.mk
include Makefile.d/build.mk
include Makefile.d/dependencies.mk
include Makefile.d/docker.mk
include Makefile.d/e2e.mk
include Makefile.d/git.mk
include Makefile.d/helm.mk
include Makefile.d/k3d.mk
include Makefile.d/k0s.mk
include Makefile.d/k8s.mk
include Makefile.d/kind.mk
include Makefile.d/minikube.mk
include Makefile.d/proto.mk
include Makefile.d/test.mk
include Makefile.d/tools.mk
include Makefile.d/tls.mk
include Makefile.d/format.mk
include Makefile.d/generate.mk
include Makefile.d/lint.mk
include Makefile.d/version.mk
