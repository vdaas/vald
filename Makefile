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

NGT_VERSION := $(shell cat versions/NGT_VERSION)
NGT_REPO = github.com/yahoojapan/NGT

GO_VERSION := $(shell cat versions/GO_VERSION)

MAKELISTS := Makefile $(shell find Makefile.d -type f -regex ".*\.mk")

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

PROTO_PATHS = \
	$(PROTODIRS:%=./apis/proto/%) \
	$(GOPATH)/src/github.com/protocolbuffers/protobuf/src \
	$(GOPATH)/src/github.com/gogo/protobuf/protobuf \
	$(GOPATH)/src/github.com/googleapis/googleapis \
	$(GOPATH)/src/github.com/danielvladco/go-proto-gql \
	$(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate

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
	    printf "  \x1b[32;01m%-35s\x1b[0m %s\n", helpCommand, helpMessage; \
	  } \
	} \
	{ lastLine = $$0 }' $(MAKELISTS) | sort -u
	@printf "\n"

.PHONY: clean
## clean
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

.PHONY: license
## add license to files
license:
	go run hack/license/gen/main.go ./

.PHONY: deps
## install dependencies
deps: \
	clean \
	proto/deps \
	proto/all
	go mod tidy
	go mod vendor
	rm -rf vendor

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

.PHONY: test
## run tests
test: clean init
	GO111MODULE=on go test --race -coverprofile=cover.out ./...

.PHONY: lint
## run lints
lint:
	gometalinter --enable-all . | rg -v comment

.PHONY: contributors
## show contributors
contributors:
	git log --format='%aN <%aE>' | sort -fu > CONTRIBUTORS

.PHONY: coverage
## calculate coverages
coverage:
	go test -v -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	rm -f coverage.out

include Makefile.d/bench.mk
include Makefile.d/docker.mk
include Makefile.d/proto.mk
