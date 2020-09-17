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

REPO                           ?= vdaas
NAME                            = vald
GOPKG                           = github.com/$(REPO)/$(NAME)
DATETIME                        = $(eval DATETIME := $(shell date -u +%Y/%m/%d_%H:%M:%S%z))$(DATETIME)
TAG                             = $(eval TAG := $(shell date -u +%Y%m%d-%H%M%S))$(TAG)
BASE_IMAGE                      = $(NAME)-base
AGENT_IMAGE                     = $(NAME)-agent-ngt
AGENT_SIDECAR_IMAGE             = $(NAME)-agent-sidecar
GATEWAY_IMAGE                   = $(NAME)-gateway
DISCOVERER_IMAGE                = $(NAME)-discoverer-k8s
META_REDIS_IMAGE                = $(NAME)-meta-redis
META_CASSANDRA_IMAGE            = $(NAME)-meta-cassandra
MANAGER_BACKUP_MYSQL_IMAGE      = $(NAME)-manager-backup-mysql
MANAGER_BACKUP_CASSANDRA_IMAGE  = $(NAME)-manager-backup-cassandra
MANAGER_COMPRESSOR_IMAGE        = $(NAME)-manager-compressor
MANAGER_INDEX_IMAGE             = $(NAME)-manager-index
CI_CONTAINER_IMAGE              = $(NAME)-ci-container
HELM_OPERATOR_IMAGE             = $(NAME)-helm-operator
LOADTEST_IMAGE                  = $(NAME)-loadtest

VERSION := $(eval VALD_VERSION := $(shell cat versions/VALD_VERSION))$(VALD_VERSION)

NGT_VERSION := $(eval NGT_VERSION := $(shell cat versions/NGT_VERSION))$(NGT_VERSION)
NGT_REPO = github.com/yahoojapan/NGT

GO_VERSION := $(eval GO_VERSION := $(shell cat versions/GO_VERSION))$(GO_VERSION)
GOOS := $(eval GOOS := $(shell go env GOOS))$(GOOS)
GOARCH := $(eval GOARCH := $(shell go env GOARCH))$(GOARCH)
GOPATH := $(eval GOPATH := $(shell go env GOPATH))$(GOPATH)
GOCACHE := $(eval GOCACHE := $(shell go env GOCACHE))$(GOCACHE)

TENSORFLOW_C_VERSION := $(eval TENSORFLOW_C_VERSION := $(shell cat versions/TENSORFLOW_C_VERSION))$(TENSORFLOW_C_VERSION)

OPERATOR_SDK_VERSION := $(eval OPERATOR_SDK_VERSION := $(shell cat versions/OPERATOR_SDK_VERSION))$(OPERATOR_SDK_VERSION)

KIND_VERSION         ?= v0.8.1
HELM_VERSION         ?= v3.2.4
HELM_DOCS_VERSION    ?= 0.13.0
VALDCLI_VERSION      ?= v0.0.50
TELEPRESENCE_VERSION ?= 0.105

SWAP_DEPLOYMENT_TYPE ?= deployment
SWAP_IMAGE           ?= ""
SWAP_TAG             ?= latest

BINDIR ?= /usr/local/bin

UNAME := $(eval UNAME := $(shell uname))$(UNAME)

ifeq ($(UNAME),Linux)
CPU_INFO_FLAGS := $(eval CPU_INFO_FLAGS := $(shell cat /proc/cpuinfo | grep flags | cut -d " " -f 2- | head -1))$(CPU_INFO_FLAGS)
else
CPU_INFO_FLAGS := ""
endif

GIT_COMMIT := $(eval GIT_COMMIT := $(shell git rev-list -1 HEAD))$(GIT_COMMIT)

MAKELISTS := Makefile $(shell find Makefile.d -type f -regex ".*\.mk")

ROOTDIR = $(eval ROOTDIR := $(shell git rev-parse --show-toplevel))$(ROOTDIR)
PROTODIRS := $(eval PROTODIRS := $(shell find apis/proto -type d | sed -e "s%apis/proto/%%g" | grep -v "apis/proto"))$(PROTODIRS)
BENCH_DATASET_BASE_DIR = hack/benchmark/assets
BENCH_DATASET_MD5_DIR_NAME = checksum
BENCH_DATASET_HDF5_DIR_NAME = dataset
BENCH_DATASET_MD5_DIR = $(BENCH_DATASET_BASE_DIR)/$(BENCH_DATASET_MD5_DIR_NAME)
BENCH_DATASET_HDF5_DIR = $(BENCH_DATASET_BASE_DIR)/$(BENCH_DATASET_HDF5_DIR_NAME)

PROTOS := $(eval PROTOS := $(shell find apis/proto -type f -regex ".*\.proto"))$(PROTOS)
PBGOS = $(PROTOS:apis/proto/%.proto=apis/grpc/%.pb.go)
SWAGGERS = $(PROTOS:apis/proto/%.proto=apis/swagger/%.swagger.json)
PBDOCS = apis/docs/docs.md

CFLAGS ?= -mno-avx512f -mno-avx512dq -mno-avx512cd -mno-avx512bw -mno-avx512vl
CXXFLAGS ?= $(CFLAGS)

BENCH_DATASET_MD5S := $(eval BENCH_DATASET_MD5S := $(shell find $(BENCH_DATASET_MD5_DIR) -type f -regex ".*\.md5"))$(BENCH_DATASET_MD5S)
BENCH_DATASETS = $(BENCH_DATASET_MD5S:$(BENCH_DATASET_MD5_DIR)/%.md5=$(BENCH_DATASET_HDF5_DIR)/%.hdf5)

DATASET_ARGS ?= identity-128
ADDRESS_ARGS ?= ""

HOST      ?= localhost
PORT      ?= 80
NUMBER    ?= 10
DIMENSION ?= 6
NUMPANES  ?= 4

BODY = ""

PROTO_PATHS = \
	$(PROTODIRS:%=./apis/proto/%) \
	$(GOPATH)/src/github.com/protocolbuffers/protobuf/src \
	$(GOPATH)/src/github.com/gogo/protobuf/protobuf \
	$(GOPATH)/src/github.com/googleapis/googleapis \
	$(GOPATH)/src/github.com/danielvladco/go-proto-gql \
	$(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate

GO_SOURCES = $(eval GO_SOURCES := $(shell find \
		./cmd \
		./hack \
		./internal \
		./pkg \
		-not -path './cmd/cli/*' \
		-not -path './internal/core/ngt/*' \
		-not -path './internal/test/comparator/*' \
		-not -path './hack/benchmark/internal/client/ngtd/*' \
		-not -path './hack/benchmark/internal/starter/agent/*' \
		-not -path './hack/benchmark/internal/starter/external/*' \
		-not -path './hack/benchmark/internal/starter/gateway/*' \
		-not -path './hack/license/*' \
		-not -path './hack/swagger/*' \
		-not -path './hack/tools/*' \
		-type f \
		-name '*.go' \
		-not -regex '.*options?\.go' \
		-not -name '*_test.go' \
		-not -name 'doc.go'))$(GO_SOURCES)
GO_OPTION_SOURCES = $(eval GO_OPTION_SOURCES := $(shell find \
		./cmd \
		./hack \
		./internal \
		./pkg \
		-not -path './cmd/cli/*' \
		-not -path './internal/core/ngt/*' \
		-not -path './internal/test/comparator/*' \
		-not -path './hack/benchmark/internal/client/ngtd/*' \
		-not -path './hack/benchmark/internal/starter/agent/*' \
		-not -path './hack/benchmark/internal/starter/external/*' \
		-not -path './hack/benchmark/internal/starter/gateway/*' \
		-not -path './hack/license/*' \
		-not -path './hack/swagger/*' \
		-not -path './hack/tools/*' \
		-type f \
		-regex '.*options?\.go' \
		-not -name '*_test.go' \
		-not -name 'doc.go'))$(GO_OPTION_SOURCES)

GO_SOURCES_INTERNAL = $(eval GO_SOURCES_INTERNAL := $(shell find \
		./internal \
		-type f \
		-name '*.go' \
		-not -name '*_test.go' \
		-not -name 'doc.go'))$(GO_SOURCES_INTERNAL)

GO_TEST_SOURCES = $(GO_SOURCES:%.go=%_test.go)
GO_OPTION_TEST_SOURCES = $(GO_OPTION_SOURCES:%.go=%_test.go)

DISTROLESS_IMAGE      ?= gcr.io/distroless/static
DISTROLESS_IMAGE_TAG  ?= nonroot
UPX_OPTIONS           ?= -9

COMMA := ,
SHELL = bash

include Makefile.d/functions.mk

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
clean:
	go clean -cache -modcache -testcache -i -r
	rm -rf \
		/go/pkg \
		./*.log \
		./*.svg \
		./apis/docs \
		./apis/swagger \
		./bench \
		./pprof \
		./libs \
		$(GOCACHE) \
		./go.sum \
		./go.mod
	cp ./hack/go.mod.default ./go.mod

.PHONY: license
## add license to files
license:
	go run hack/license/gen/main.go ./

.PHONY: init
## initialize development environment
init: \
	git/config/init \
	git/hooks/init \
	deps \
	ngt/install \
	tensorflow/install

.PHONY: tools/install
## install development tools
tools/install: \
	helm/install \
	kind/install \
	valdcli/install \
	telepresence/install

.PHONY: update
## update deps, license, and run goimports
update: \
	clean \
	deps \
	proto/all \
	format \
	go/deps

.PHONY: format
## format go codes
format: \
	license \
	update/goimports \
	format/yaml

.PHONY: update/goimports
## run goimports for all go files
update/goimports:
	find ./ -type d -name .git -prune -o -type f -regex '.*\.go' -print | xargs goimports -w

.PHONY: format/yaml
format/yaml:
	prettier --write \
	    ".github/**/*.yaml" \
	    "cmd/**/*.yaml" \
	    "hack/**/*.yaml" \
	    "k8s/**/*.yaml"

.PHONY: deps
## resolve dependencies
deps: \
	proto/deps \
	deps/install

.PHONY: deps/install
## install dependencies
deps/install: \
	goimports/install \
	prettier/install \
	go/deps

.PHONY: go/deps
## install Go package dependencies
go/deps:
	go clean -cache -modcache -testcache -i -r
	rm -rf \
		/go/pkg \
		$(GOCACHE) \
		./go.sum \
		./go.mod
	cp ./hack/go.mod.default ./go.mod
	go mod tidy


.PHONY: goimports/install
goimports/install:
	go get -u golang.org/x/tools/cmd/goimports
	# GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

.PHONY: prettier/install
prettier/install:
	npm install -g npm prettier

.PHONY: version/vald
## print vald version
version/vald:
	@echo $(VALD_VERSION)

.PHONY: version/go
## print go version
version/go:
	@echo $(GO_VERSION)

.PHONY: version/ngt
## print NGT version
version/ngt:
	@echo $(NGT_VERSION)

.PHONY: version/kind
version/kind:
	@echo $(KIND_VERSION)

.PHONY: version/helm
version/helm:
	@echo $(HELM_VERSION)

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
	tar zxf v$(NGT_VERSION).tar.gz -C /tmp
	cd /tmp/NGT-$(NGT_VERSION) && \
	    cmake -DCMAKE_C_FLAGS="$(CFLAGS)" -DCMAKE_CXX_FLAGS="$(CXXFLAGS)" .
	make -j -C /tmp/NGT-$(NGT_VERSION)
	make install -C /tmp/NGT-$(NGT_VERSION)
	rm -rf v$(NGT_VERSION).tar.gz
	rm -rf /tmp/NGT-$(NGT_VERSION)
	ldconfig

.PHONY: tensorflow/install
## install TensorFlow for C
tensorflow/install: /usr/local/lib/libtensorflow.so
ifeq ($(UNAME),Darwin)
/usr/local/lib/libtensorflow.so:
	brew install libtensorflow@1
else
/usr/local/lib/libtensorflow.so:
	curl -LO https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-$(TENSORFLOW_C_VERSION).tar.gz
	tar -C /usr/local -xzf libtensorflow-cpu-linux-x86_64-$(TENSORFLOW_C_VERSION).tar.gz
	rm -f libtensorflow-cpu-linux-x86_64-$(TENSORFLOW_C_VERSION).tar.gz
	ldconfig
endif

.PHONY: lint
## run lints
lint:
	$(call go-lint)

.PHONY: changelog/update
## update changelog
changelog/update:
	echo "# CHANGELOG" > /tmp/CHANGELOG.md
	echo "" >> /tmp/CHANGELOG.md
	$(MAKE) -s changelog/next/print >> /tmp/CHANGELOG.md
	echo "" >> /tmp/CHANGELOG.md
	tail -n +2 CHANGELOG.md >> /tmp/CHANGELOG.md
	mv -f /tmp/CHANGELOG.md CHANGELOG.md

.PHONY: changelog/next/print
## print next changelog entry
changelog/next/print:
	@cat hack/CHANGELOG.template.md | \
	    sed -e 's/{{ version }}/$(VALD_VERSION)/g'
	@echo "$$BODY"

include Makefile.d/bench.mk
include Makefile.d/build.mk
include Makefile.d/docker.mk
include Makefile.d/git.mk
include Makefile.d/helm.mk
include Makefile.d/proto.mk
include Makefile.d/k8s.mk
include Makefile.d/kind.mk
include Makefile.d/client.mk
include Makefile.d/ml.mk
include Makefile.d/test.mk
