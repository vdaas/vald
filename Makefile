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
TAG                             = $(shell date -u +%Y%m%d-%H%M%S)
BASE_IMAGE                      = $(NAME)-base
AGENT_IMAGE                     = $(NAME)-agent-ngt
GATEWAY_IMAGE                   = $(NAME)-gateway
DISCOVERER_IMAGE                = $(NAME)-discoverer-k8s
KVS_IMAGE                       = $(NAME)-meta-redis
NOSQL_IMAGE                     = $(NAME)-meta-cassandra
BACKUP_MANAGER_MYSQL_IMAGE      = $(NAME)-manager-backup-mysql
BACKUP_MANAGER_CASSANDRA_IMAGE  = $(NAME)-manager-backup-cassandra
MANAGER_COMPRESSOR_IMAGE        = $(NAME)-manager-compressor
MANAGER_INDEX_IMAGE             = $(NAME)-manager-index
CI_CONTAINER_IMAGE              = $(NAME)-ci-container

NGT_VERSION := $(shell cat versions/NGT_VERSION)
NGT_REPO = github.com/yahoojapan/NGT

GO_VERSION := $(shell cat versions/GO_VERSION)
GOPATH := $(shell go env GOPATH)
GOCACHE := $(shell go env GOCACHE)

TENSORFLOW_C_VERSION := $(shell cat versions/TENSORFLOW_C_VERSION)

KIND_VERSION         ?= v0.7.0
VALDCLI_VERSION      ?= v0.0.1
TELEPRESENCE_VERSION ?= 0.104

BINDIR ?= /usr/local/bin

UNAME := $(shell uname)

MAKELISTS := Makefile $(shell find Makefile.d -type f -regex ".*\.mk")

ROOTDIR = $(shell git rev-parse --show-toplevel)
PROTODIRS := $(shell find apis/proto -type d | sed -e "s%apis/proto/%%g" | grep -v "apis/proto")
PBGODIRS = $(PROTODIRS:%=apis/grpc/%)
SWAGGERDIRS = $(PROTODIRS:%=apis/swagger/%)
GRAPHQLDIRS = $(PROTODIRS:%=apis/graphql/%)
PBDOCDIRS = $(PROTODIRS:%=apis/docs/%)

BENCH_DATASET_BASE_DIR = hack/benchmark/assets
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

DATASET_ARGS ?= identity-128
ADDRESS_ARGS ?= ""

HOST      ?= localhost
PORT      ?= 80
NUMBER    ?= 10
DIMENSION ?= 6
NUMPANES  ?= 4

PROTO_PATHS = \
	$(PROTODIRS:%=./apis/proto/%) \
	$(GOPATH)/src/github.com/protocolbuffers/protobuf/src \
	$(GOPATH)/src/github.com/gogo/protobuf/protobuf \
	$(GOPATH)/src/github.com/googleapis/googleapis \
	$(GOPATH)/src/github.com/danielvladco/go-proto-gql \
	$(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate

COMMA := ,
SHELL = bash

include Makefile.d/functions.mk

.PHONY: all
## execute clean and deps
all: clean deps

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

.PHONY: clean
## clean
clean:
	go clean -cache -modcache -testcache -i -r
	rm -rf \
		/go/pkg \
		./*.log \
		./*.svg \
		./apis/docs \
		./apis/graphql \
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
	license \
	update/goimports


.PHONY: format
## format go codes
format: \
	license \
	update/goimports

.PHONY: update/goimports
## run goimports for all go files
update/goimports:
	find ./ -type f -regex ".*\.go" | xargs goimports -w

.PHONY: deps
## install dependencies
deps: \
	proto/deps
	go mod tidy

.PHONY: version/go
## print go version
version/go:
	@echo $(GO_VERSION)

.PHONY: version/ngt
## print NGT version
version/ngt:
	@echo $(NGT_VERSION)

.PHONY: ngt/install
## install NGT
ngt/install: /usr/local/include/NGT/Capi.h
/usr/local/include/NGT/Capi.h:
	curl -LO https://github.com/yahoojapan/NGT/archive/v$(NGT_VERSION).tar.gz
	tar zxf v$(NGT_VERSION).tar.gz -C /tmp
	cd /tmp/NGT-$(NGT_VERSION)&& cmake .
	make -j -C /tmp/NGT-$(NGT_VERSION)
	make install -C /tmp/NGT-$(NGT_VERSION)
	rm -rf v$(NGT_VERSION).tar.gz
	rm -rf /tmp/NGT-$(NGT_VERSION)

.PHONY: tensorflow/install
## install TensorFlow for C
tensorflow/install: /usr/local/lib/libtensorflow.so
/usr/local/lib/libtensorflow.so:
	curl -LO https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-$(TENSORFLOW_C_VERSION).tar.gz
	tar -C /usr/local -xzf libtensorflow-cpu-linux-x86_64-$(TENSORFLOW_C_VERSION).tar.gz
	rm -f libtensorflow-cpu-linux-x86_64-$(TENSORFLOW_C_VERSION).tar.gz
	ldconfig

.PHONY: test
## run tests
test:
	GO111MODULE=on go test --race -coverprofile=cover.out ./...

.PHONY: lint
## run lints
lint:
	$(call go-lint)


.PHONY: coverage
## calculate coverages
coverage:
	go test -v -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: readme/update/authors
## update authors in README
readme/update/authors:
	$(eval AUTHORS = $(shell awk '{printf "- [%s]\\(https:\\/\\/github.com\\/%s\\)\\n", $$1, $$1}' AUTHORS))
	sed -i -e ':lbl1;N;s/## Author.*## Contributor/## Author\n\n$(AUTHORS)\n## Contributor/;b lbl1;' README.md

.PHONY: readme/update/contributors
## update contributors in README
readme/update/contributors:
	$(eval CONTRIBUTORS = $(shell awk '{printf "- [%s]\\(https:\\/\\/github.com\\/%s\\)\\n", $$1, $$1}' CONTRIBUTORS))
	sed -i -e ':lbl1;N;s/## Contributor.*## LICENSE/## Contributor\n\n$(CONTRIBUTORS)\n## LICENSE/;b lbl1;' README.md

include Makefile.d/bench.mk
include Makefile.d/docker.mk
include Makefile.d/git.mk
include Makefile.d/proto.mk
include Makefile.d/k8s.mk
include Makefile.d/kind.mk
include Makefile.d/client.mk
