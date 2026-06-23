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

# --- Color Definitions ---

red = printf "\x1b[31m\#\# %s\x1b[0m\n" $1
green = printf "\x1b[32m\#\# %s\x1b[0m\n" $1
yellow = printf "\x1b[33m\#\# %s\x1b[0m\n" $1
blue = printf "\x1b[34m\#\# %s\x1b[0m\n" $1
pink = printf "\x1b[35m\#\# %s\x1b[0m\n" $1
cyan = printf "\x1b[36m\#\# %s\x1b[0m\n" $1

# --- Environment Variables & Constants ---

LLD ?= $(shell command -v ld.lld 2>/dev/null || command -v lld 2>/dev/null)
NM ?= $(shell command -v llvm-nm 2>/dev/null || ls /usr/bin/llvm-nm-* 2>/dev/null | sort -V | tail -1 | grep . || command -v gcc-nm 2>/dev/null || ls /usr/bin/gcc-nm-* 2>/dev/null | sort -V | tail -1 | grep . || command -v nm)
RANLIB ?= $(shell command -v llvm-ranlib 2>/dev/null || ls /usr/bin/llvm-ranlib-* 2>/dev/null | sort -V | tail -1 | grep . || command -v gcc-ranlib 2>/dev/null || ls /usr/bin/gcc-ranlib-* 2>/dev/null | sort -V | tail -1 | grep . || command -v ranlib)

CC_ENV_VARS = \
	CC="$(CC)" \
	CXX="$(CXX)" \
	LLD="$(LLD)" \
	AR="$(AR)" \
	NM="$(NM)" \
	RANLIB="$(RANLIB)"

GO_ENV_VARS = \
	$(CC_ENV_VARS) \
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	GO111MODULE=on \
	GO_VERSION=$(GO_VERSION)

CGO_ENV_VARS = \
	CGO_ENABLED=$(CGO_ENABLED) \
	CGO_CFLAGS_ALLOW="$(CFLAGS_BASE)" \
	CGO_CFLAGS="$(CFLAGS) $1" \
	CGO_CXXFLAGS="$(CXXFLAGS) $1" \
	CGO_FFLAGS="$(CFLAGS) $1" \
	CGO_LDFLAGS="$2"

DOCKER_BUILD_ARGS = \
	--build-arg BUILDKIT_INLINE_CACHE=$(BUILDKIT_INLINE_CACHE) \
	--build-arg GO_VERSION=$(GO_VERSION) \
	--build-arg RUST_VERSION=$(RUST_VERSION) \
	--build-arg MAINTAINER=$(MAINTAINER) \
	--build-arg EMAIL="$(EMAIL)" \
	--build-arg CC="$(CC)" \
	--build-arg CXX="$(CXX)" \
	--build-arg LLD="$(LLD)" \
	--build-arg AR="$(AR)" \
	--build-arg NM="$(NM)" \
	--build-arg RANLIB="$(RANLIB)" \
	--build-arg CFLAGS="$(CFLAGS)" \
	--build-arg CXXFLAGS="$(CXXFLAGS)" \
	--build-arg LDFLAGS="$(LDFLAGS)" \
	--build-arg SUDO=""

GO_INFO_LDFLAGS = \
	-X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)' \
	-X '$(GOPKG)/internal/info.BuildTime=$(DATETIME)' \
	-X '$(GOPKG)/internal/info.CGOEnabled=$(if $(filter 1,$(strip $(CGO_ENABLED))),true,false)' \
	-X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
	-X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
	-X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
	-X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
	-X '$(GOPKG)/internal/info.Version=$(VERSION)'

GO_BUILD_FLAGS = -modcacherw -mod=readonly -a -trimpath
GO_TEST_BASE_FLAGS = -short -shuffle=on -race -mod=readonly -cover -timeout=$(GOTEST_TIMEOUT)
GO_TEST_FLAGS = $(GO_TEST_BASE_FLAGS) -ldflags="$(GO_TEST_LDFLAGS)"
GO_TEST_ENV = $(GO_ENV_VARS) CGO_LDFLAGS="$(TEST_LDFLAGS)"

# --- Filesystem Utilities ---

define mkdir
	mkdir -p $1
endef

# --- Go Tools & Maintenance ---

define go-tool-install
	cat $(ROOTDIR)/hack/go.tools | \
	xargs $(XARGS_NO_RUN_IF_EMPTY) -I {} -P $(CORES) \
	sh -c 'if ! out=$$( $(GO_ENV_VARS) go install -ldflags="$(GO_STATIC_LDFLAGS)" {} 2>&1 >/dev/null); then echo "--- Failed to install {} ---" >&2; echo "$$out" >&2; exit 255; fi';
endef

define go-lint
	golangci-lint run --config $(ROOTDIR)/.golangci.json --fix
endef

define go-vet
	cat <(GOARCH=amd64 go vet $(ROOTDIR)/...) \
	<(GOARCH=386 go vet $(ROOTDIR)/...) \
	<(GOARCH=arm go vet $(ROOTDIR)/...) \
	| grep -v "Mutex" | sort | uniq
endef

define update-template
	sed -i -e "s/^- $1 Version: .*$$/- $1 Version: $2/" $(ROOTDIR)/.github/ISSUE_TEMPLATE/bug_report.md
	sed -i -e "s/^- $1 Version: .*$$/- $1 Version: $2/" $(ROOTDIR)/.github/ISSUE_TEMPLATE/security_issue_report.md
	sed -i -e "s/^- $1 Version: .*$$/- $1 Version: $2/" $(ROOTDIR)/.github/PULL_REQUEST_TEMPLATE.md
endef

export UPDATE_GITHUB_ACTION_SCRIPT = \
	ACTION="$$1"; \
	FILE_NAME=$$(echo $$ACTION | tr '/' '_' | tr '-' '_' | tr '[:lower:]' '[:upper:]'); \
	if [ -n "$$FILE_NAME" ]; then \
		case "$$ACTION" in \
			"aquasecurity/trivy-action" | "machine-learning-apps/actions-chatops" ) VERSION="master";; \
			* ) \
			REPO_NAME=$$(echo $$ACTION | cut -d'/' -f1-2); \
			echo "$$ACTION to $$REPO_NAME"; \
			VERSION=$$(curl -fsSL https://api.github.com/repos/$$REPO_NAME/tags?per_page=1 | \
			grep -Po '"name": "\K.*?(?=")' | \
			head -n1 | sed 's/v//g' | sed -E 's/[^0-9.]+//g'); \
			;; \
		esac; \
	if [ -n "$$VERSION" ]; then \
		OLD_VERSION=$$(cat $(ROOTDIR)/versions/actions/$$FILE_NAME); \
		echo "updating $$ACTION version file $$FILE_NAME from $$OLD_VERSION to $$VERSION"; \
		echo $$VERSION > $(ROOTDIR)/versions/actions/$$FILE_NAME; \
	else \
		VERSION=$$(cat $(ROOTDIR)/versions/actions/$$FILE_NAME); \
		echo "No version found for $$ACTION version file $$FILE_NAME=$$VERSION"; \
	fi; \
	if [ "$$ACTION" = "cirrus-actions/rebase" ]; then \
		VERSION_PREFIX=$$VERSION; \
		cat $(ROOTDIR)/.gitfiles | grep '^\.github/' | sed -e 's%^%$(ROOTDIR)/%' | \
		xargs -I __FILE__ -P $(CORES) sed -i "s%$$ACTION@.*%$$ACTION@$$VERSION_PREFIX%g" "__FILE__"; \
	elif echo $$VERSION | grep -qE '^[0-9]'; then \
		VERSION_PREFIX=$$(echo $$VERSION | cut -c 1); \
		cat $(ROOTDIR)/.gitfiles | grep '^\.github/' | sed -e 's%^%$(ROOTDIR)/%' | \
		xargs -I __FILE__ -P $(CORES) sed -i "s%$$ACTION@.*%$$ACTION@v$$VERSION_PREFIX%g" "__FILE__"; \
	else \
		VERSION_PREFIX=$$VERSION; \
		cat $(ROOTDIR)/.gitfiles | grep '^\.github/' | sed -e 's%^%$(ROOTDIR)/%' | \
		xargs -I __FILE__ -P $(CORES) sed -i "s%$$ACTION@.*%$$ACTION@$$VERSION_PREFIX%g" "__FILE__"; \
	fi; \
	fi

define update-github-actions
	@echo "$1" | tr ' ' '\n' | grep -v '^$$' | grep -v '^security-and-quality$$' | \
	xargs -I {} -P $(CORES) bash -c 'eval "$$UPDATE_GITHUB_ACTION_SCRIPT"' _ {}
endef

# --- Go Build & Example ---

# go-base implements base go command helper
# Arguments:
# $1: Command (build, test, install, run)
# $2: Extra Environment Variables
# $3: Command Flags
# $4: ldflags
# $5: Tags
# $6: Output path (-o flag)
# $7: Target path (cmd path or package path)
# $8: Additional arguments
define go-base
	$(GO_ENV_VARS) $2 go $(strip $1) \
		$(if $(filter build,$(strip $1)),$(GO_BUILD_FLAGS)) \
		$(if $(filter test,$(strip $1)),$(GO_TEST_BASE_FLAGS)) \
		-ldflags "$(strip -buildid= $(GO_INFO_LDFLAGS) $4)" \
		$(strip $3) \
		$(if $(strip $5),-tags "$(strip $5)") \
		$(if $(strip $6),-o $(strip $6)) \
		$(strip $7) \
		$8
endef

define go-build
	$(call go-base,build, \
		$(call CGO_ENV_VARS,-DNGT_LARGE_DATASET,$3), \
		-mod=readonly, \
		$(GO_LDFLAGS) $2 -extldflags '-static $3' \
		-X '$(GOPKG)/internal/info.AlgorithmInfo=$5', \
		osusergo netgo static_build $4, \
		$6, \
		$(ROOTDIR)/cmd/$1/main.go)
	OMP_TOOL=disabled $6 -version
endef

define go-example-build
	cd $(ROOTDIR)/$1 && \
	$(call go-base,build, \
	$(call CGO_ENV_VARS,-DNGT_LARGE_DATASET,$3), \
	, \
	$(GO_LDFLAGS) $2 -extldflags '-static $3', \
	osusergo netgo static_build$4, \
	$(ROOTDIR)/$6, \
	$(ROOTDIR)/example/client/main.go)
endef

define go-e2e-build
	$(call go-base,test, \
		$(call CGO_ENV_VARS,,$2), \
		-c -v -race -mod=readonly, \
		$(GO_LDFLAGS) -linkmode=external \
		-extldflags '-static $2', \
		e2e, \
		$(ROOTDIR)/$3, \
		$(ROOTDIR)/$1)
endef

# --- Go Test & Benchmark ---

define go-test
	$(call go-base,test, \
		CGO_LDFLAGS="$(TEST_LDFLAGS)", \
		, \
		$(GO_TEST_LDFLAGS), \
		, \
		, \
		$1)
endef

define go-test-tparse
	set -euo pipefail; \
	rm -rf "$(TEST_RESULT_DIR)/$$(echo $2 | sed -e 's%/%-%g')-result.json"; \
	$(call go-base,test, \
	CGO_LDFLAGS="$(TEST_LDFLAGS)", \
	-json, \
	$(GO_TEST_LDFLAGS), \
	, \
	, \
	$1, \
	| tee "$(TEST_RESULT_DIR)/$$(echo $2 | sed -e 's%/%-%g')-result.json" \
	| tparse -pass -notests)
endef

define go-test-gotestfmt
	set -euo pipefail; \
	rm -rf "$(TEST_RESULT_DIR)/$$(echo $2 | sed -e 's%/%-%g')-result.json"; \
	$(call go-base,test, \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" GODEBUG=$(GODEBUG), \
	-json, \
	$(GO_TEST_LDFLAGS), \
	, \
	, \
	$1, \
	| tee "$(TEST_RESULT_DIR)/$$(echo $2 | sed -e 's%/%-%g')-result.json" \
	| gotestfmt $3)
endef

define go-bench
	mkdir -p $(dir $2)
	$(call go-base,test, \
		CGO_LDFLAGS="$(TEST_LDFLAGS)", \
		-mod=readonly -count=1 -timeout=1h -bench=$1 -benchmem \
		-cpuprofile $(patsubst %.bin,%.cpu.out,$2) \
		-memprofile $(patsubst %.bin,%.mem.out,$2) \
		-trace $(patsubst %.bin,%.trace.out,$2), \
		, \
		, \
		$2, \
		$3)
endef

GOTESTS_GEN_SCRIPT = \
	f="{}"; \
	GOTESTS_OPTION="-all"; \
	if [[ $$f =~ \.\/pkg\/.*\/router\/.* || $$f =~ \.\/pkg\/.*\/handler\/rest\/.* ]]; then \
		echo "Skip generating go test file: $$f"; \
		exit 0; \
	elif [[ $$f =~ \.\/pkg\/.*\/usecase\/.* ]]; then \
		GOTESTS_OPTION=" -only New "; \
	elif [[ $$f =~ \.\/pkg\/.* ]]; then \
		GOTESTS_OPTION=" -exported "; \
	fi; \
	echo "Generating go test file: $$f with option $$GOTESTS_OPTION"; \
	gotests -w -template_dir $(ROOTDIR)/assets/test/templates/common $$GOTESTS_OPTION $$f;

GOTESTS_OPTION_GEN_SCRIPT = \
	f="{}"; \
	GOTESTS_OPTION="-all"; \
	if [[ $$f =~ \.\/pkg\/.*\/router\/.* || $$f =~ \.\/pkg\/.*\/handler\/rest\/.* || $$f =~ \.\/pkg\/.*\/usecase\/.* ]]; then \
		echo "Skip generating go option test file: $$f"; \
		exit 0; \
	elif [[ $$f =~ \.\/pkg\/.* ]]; then \
		GOTESTS_OPTION=" -exported "; \
	fi; \
	echo "Generating go option test file: $$f with option $$GOTESTS_OPTION"; \
	gotests -w -template_dir $(ROOTDIR)/assets/test/templates/option $$GOTESTS_OPTION $$f;

# This function generate only implementation tests
define gen-go-test-sources
	@$(call green, "Generating go test files in parallel (cores: $(CORES))...")
	@echo "$(GO_SOURCES)" | tr ' ' '\n' | \
		xargs -I {} -P$(CORES) bash -c '$(GOTESTS_GEN_SCRIPT)'
endef

# This function generate only option tests
define gen-go-option-test-sources
	@$(call green, "Generating go option test files in parallel (cores: $(CORES))...")
	@echo "$(GO_OPTION_SOURCES)" | tr ' ' '\n' | \
		xargs -I {} -P$(CORES) bash -c '$(GOTESTS_OPTION_GEN_SCRIPT)'
endef

# --- E2E Test ---

define run-v2-e2e-crud-test
	$(call go-base,test, \
		CGO_LDFLAGS="$(TEST_LDFLAGS)" \
		CGO_CFLAGS="$(CGO_CFLAGS)" \
		E2E_ADDR="$(E2E_BIND_HOST):$(E2E_BIND_PORT)" \
		E2E_BIND_HOST="$(E2E_BIND_HOST)" \
		E2E_BIND_PORT="$(E2E_BIND_PORT)" \
		E2E_TARGET_NAMESPACE="$(E2E_TARGET_NAMESPACE)" \
		E2E_TARGET_NAME="$(E2E_TARGET_NAME)" \
		E2E_DATASET_PATH="$(ROOTDIR)/hack/benchmark/assets/dataset/$(E2E_DATASET_NAME)" \
		E2E_PARALLELISM="$(E2E_PARALLELISM)" \
		E2E_INSERT_COUNT="$(E2E_INSERT_COUNT)" \
		E2E_QPS="$(E2E_QPS)" \
		E2E_SEARCH_COUNT="$(E2E_SEARCH_COUNT)" \
		E2E_UPDATE_COUNT="$(E2E_UPDATE_COUNT)" \
		E2E_BULK_SIZE="$(E2E_BULK_SIZE)" \
		E2E_EXPECTED_INDEX="$(E2E_EXPECTED_INDEX)", \
		-race -v -mod=readonly -timeout $(E2E_TIMEOUT) -config $(E2E_CONFIG), \
		$(GO_STATIC_LDFLAGS), \
		e2e, \
		, \
		$1 $(ROOTDIR)/tests/v2/e2e/crud)
endef

define run-e2e-crud-test
	$(call go-base,test, \
		CGO_LDFLAGS="$(TEST_LDFLAGS)", \
		-race -mod=readonly -v -timeout $(E2E_TIMEOUT) \
		-host=$(E2E_BIND_HOST) \
		-port=$(E2E_BIND_PORT) \
		-dataset=$(ROOTDIR)/hack/benchmark/assets/dataset/$(E2E_DATASET_NAME) \
		-insert-num=$(E2E_INSERT_COUNT) \
		-correction-insert-num=$(E2E_INSERT_COUNT) \
		-search-num=$(E2E_SEARCH_COUNT) \
		-search-by-id-num=$(E2E_SEARCH_BY_ID_COUNT) \
		-get-object-num=$(E2E_GET_OBJECT_COUNT) \
		-update-num=$(E2E_UPDATE_COUNT) \
		-upsert-num=$(E2E_UPSERT_COUNT) \
		-remove-num=$(E2E_REMOVE_COUNT) \
		-insert-from=$(E2E_INSERT_FROM) \
		-update-from=$(E2E_UPDATE_FROM) \
		-upsert-from=$(E2E_UPSERT_FROM) \
		-remove-from=$(E2E_REMOVE_FROM) \
		-search-from=$(E2E_SEARCH_FROM) \
		-search-by-id-from=$(E2E_SEARCH_BY_ID_FROM) \
		-get-object-from=$(E2E_GET_OBJECT_FROM) \
		-wait-after-insert=$(E2E_WAIT_FOR_CREATE_INDEX_DURATION) \
		-portforward=$(E2E_PORTFORWARD_ENABLED) \
		-portforward-pod-name=$(E2E_TARGET_POD_NAME) \
		-portforward-pod-port=$(E2E_TARGET_PORT) \
		-namespace=$(E2E_TARGET_NAMESPACE) \
		-kubeconfig=$(KUBECONFIG), \
		$(GO_STATIC_LDFLAGS), \
		e2e, \
		, \
		$1 $(ROOTDIR)/tests/e2e/crud/crud_test.go)
endef

define run-e2e-crud-faiss-test
	$(call go-base,test, \
		CGO_LDFLAGS="$(TEST_LDFLAGS)", \
		-race -mod=readonly -v -timeout $(E2E_TIMEOUT) \
		-host=$(E2E_BIND_HOST) \
		-port=$(E2E_BIND_PORT) \
		-dataset=$(ROOTDIR)/hack/benchmark/assets/dataset/$(E2E_DATASET_NAME).hdf5 \
		-insert-num=$(E2E_INSERT_COUNT) \
		-search-num=$(E2E_SEARCH_COUNT) \
		-update-num=$(E2E_UPDATE_COUNT) \
		-remove-num=$(E2E_REMOVE_COUNT) \
		-wait-after-insert=$(E2E_WAIT_FOR_CREATE_INDEX_DURATION) \
		-portforward=$(E2E_PORTFORWARD_ENABLED) \
		-portforward-pod-name=$(E2E_TARGET_POD_NAME) \
		-portforward-pod-port=$(E2E_TARGET_PORT) \
		-namespace=$(E2E_TARGET_NAMESPACE), \
		$(GO_STATIC_LDFLAGS), \
		e2e, \
		, \
		$1 $(ROOTDIR)/tests/e2e/crud/crud_faiss_test.go)
endef

define run-e2e-multi-crud-test
	$(call go-base,test, \
		CGO_LDFLAGS="$(TEST_LDFLAGS)", \
		-race -mod=readonly -v -timeout $(E2E_TIMEOUT) \
		-host=$(E2E_BIND_HOST) \
		-port=$(E2E_BIND_PORT) \
		-dataset=$(ROOTDIR)/hack/benchmark/assets/dataset/$(E2E_DATASET_NAME) \
		-insert-num=$(E2E_INSERT_COUNT) \
		-search-num=$(E2E_SEARCH_COUNT) \
		-search-by-id-num=$(E2E_SEARCH_BY_ID_COUNT) \
		-get-object-num=$(E2E_GET_OBJECT_COUNT) \
		-update-num=$(E2E_UPDATE_COUNT) \
		-upsert-num=$(E2E_UPSERT_COUNT) \
		-remove-num=$(E2E_REMOVE_COUNT) \
		-wait-after-insert=$(E2E_WAIT_FOR_CREATE_INDEX_DURATION) \
		-portforward=$(E2E_PORTFORWARD_ENABLED) \
		-portforward-pod-name=$(E2E_TARGET_POD_NAME) \
		-portforward-pod-port=$(E2E_TARGET_PORT) \
		-namespace=$(E2E_TARGET_NAMESPACE) \
		-kubeconfig=$(KUBECONFIG), \
		$(GO_STATIC_LDFLAGS), \
		e2e, \
		, \
		$1 $(ROOTDIR)/tests/e2e/multiapis/multiapis_test.go)
endef

define run-e2e-max-dim-test
	$(call go-base,test, \
		CGO_LDFLAGS="$(TEST_LDFLAGS)", \
		-race -mod=readonly -v -timeout $(E2E_TIMEOUT) \
		-file $(E2E_MAX_DIM_RESULT_FILEPATH) \
		-host=$(E2E_BIND_HOST) \
		-port=$(E2E_BIND_PORT) \
		-bit=${E2E_MAX_DIM_BIT} \
		-portforward=$(E2E_PORTFORWARD_ENABLED) \
		-portforward-pod-name=$(E2E_TARGET_POD_NAME) \
		-portforward-pod-port=$(E2E_TARGET_PORT) \
		-namespace=$(E2E_TARGET_NAMESPACE) \
		-kubeconfig=$(KUBECONFIG), \
		$(GO_STATIC_LDFLAGS), \
		e2e, \
		, \
		$(ROOTDIR)/tests/e2e/performance/max_vector_dim_test.go)
endef

define run-e2e-sidecar-test
	$(call go-base,test, \
		CGO_LDFLAGS="$(TEST_LDFLAGS)", \
		-race -mod=readonly -v -timeout $(E2E_TIMEOUT) \
		-host=$(E2E_BIND_HOST) \
		-port=$(E2E_BIND_PORT) \
		-dataset=$(ROOTDIR)/hack/benchmark/assets/dataset/$(E2E_DATASET_NAME) \
		-insert-num=$(E2E_INSERT_COUNT) \
		-search-num=$(E2E_SEARCH_COUNT) \
		-portforward=$(E2E_PORTFORWARD_ENABLED) \
		-portforward-pod-name=$(E2E_TARGET_POD_NAME) \
		-portforward-pod-port=$(E2E_TARGET_PORT) \
		-namespace=$(E2E_TARGET_NAMESPACE) \
		-kubeconfig=$(KUBECONFIG), \
		$(GO_STATIC_LDFLAGS), \
		e2e, \
		, \
		$1 $(ROOTDIR)/tests/e2e/sidecar/sidecar_test.go)
endef

define run-github-actions-e2e
	$(MAKE) $1
	kubectl wait -n kube-system --for=condition=Available deployment/metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	sleep 2
	kubectl wait -n kube-system --for=condition=Ready pod -l $2 --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl wait -n kube-system --for=condition=ContainersReady pod -l $2 --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	$(MAKE) $3 HELM_VALUES=$4 $5
	sleep 3
	kubectl wait --for=condition=Ready pod -l "$6" --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl wait --for=condition=ContainersReady pod -l "$6" --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl get pods
	pod_name=$$(kubectl get pods --selector="$6" | tail -1 | awk '{print $$1}'); \
	echo $$pod_name; \
	$7
	$8
endef

# --- Code/Doc Generation ---

define gen-license
	BIN_PATH="$(TEMP_DIR)/vald-license-gen"; \
	rm -rf $$BIN_PATH; \
	MAINTAINER=$2 \
	$(call go-base,build, \
	, \
	, \
	$(GO_STATIC_LDFLAGS), \
	osusergo netgo static_build, \
	$$BIN_PATH, \
	$(ROOTDIR)/hack/license/gen/main.go) \
	&& $$BIN_PATH $1; \
	rm -rf $$BIN_PATH
endef

define gen-dockerfile
	BIN_PATH="$(TEMP_DIR)/vald-dockerfile-gen"; \
	rm -rf $$BIN_PATH; \
	MAINTAINER=$2 \
	$(call go-base,build, \
	, \
	, \
	$(GO_STATIC_LDFLAGS), \
	osusergo netgo static_build, \
	$$BIN_PATH, \
	$(ROOTDIR)/hack/docker/gen/main.go) \
	&& $$BIN_PATH $1; \
	rm -rf $$BIN_PATH
endef

define gen-dashboard
	go run $(ROOTDIR)/hack/grafana/gen/src
endef

define gen-vald-helm-schema
	BIN_PATH="$(TEMP_DIR)/vald-helm-schema-gen"; \
	rm -rf $$BIN_PATH; \
	$(call go-base,build, \
	, \
	, \
	$(GO_STATIC_LDFLAGS), \
	osusergo netgo static_build, \
	$$BIN_PATH, \
	$(ROOTDIR)/hack/helm/schema/gen/main.go) \
	&& $$BIN_PATH charts/$1.yaml > charts/$1.schema.json; \
	rm -rf $$BIN_PATH
endef

define gen-vald-crd
	if [[ -f $(ROOTDIR)/charts/$1/crds/$2.yaml ]]; then \
		mv $(ROOTDIR)/charts/$1/crds/$2.yaml $(TEMP_DIR)/$2.yaml; \
	fi;
	BIN_PATH="$(TEMP_DIR)/vald-helm-crd-schema-gen"; \
	rm -rf $$BIN_PATH; \
	$(call go-base,build, \
	, \
	, \
	$(GO_STATIC_LDFLAGS), \
	osusergo netgo static_build, \
	$$BIN_PATH, \
	$(ROOTDIR)/hack/helm/schema/crd/main.go) \
	&& $$BIN_PATH $(ROOTDIR)/charts/$3.yaml > $(TEMP_DIR)/$2-spec.yaml; \
	rm -rf $$BIN_PATH; \
	$(BINDIR)/yq eval-all \
	'select(fileIndex==0).spec.versions[0].schema.openAPIV3Schema.properties.spec = select(fileIndex==1).spec | select(fileIndex==0)' \
	$(TEMP_DIR)/$2.yaml $(TEMP_DIR)/$2-spec.yaml > $(ROOTDIR)/charts/$1/crds/$2.yaml
endef

define gen-deadlink-checker
	BIN_PATH="$(TEMP_DIR)/vald-deadlink-checker-gen"; \
	rm -rf $$BIN_PATH; \
	MAINTAINER=$2 \
	$(call go-base,build, \
	, \
	, \
	$(GO_STATIC_LDFLAGS), \
	osusergo netgo static_build, \
	$$BIN_PATH, \
	$(ROOTDIR)/hack/tools/deadlink/main.go) \
	&& $$BIN_PATH -path $3 -ignore-path $4 -format $5 $1; \
	rm -rf $$BIN_PATH
endef

define gen-api-document
	buf generate --template=apis/docs/buf.gen.tmpl.yaml --path $2
	cat apis/docs/v1/payload.md.tmpl apis/docs/v1/_doc.md.tmpl > apis/docs/v1/doc.md.tmpl; \
	buf generate --template=apis/docs/buf.gen.doc.yaml --path $2; \
	mv $(ROOTDIR)/apis/docs/v1/doc.md $1; \
	rm apis/docs/v1/*doc.md.tmpl
endef

# --- C/C++ & CMake Build/Install ---

# cmake-install merges c-lib-install and curl-cmake-install.
# Arguments:
# $1: Source (URL or Git Repo)
# $2: Name
# $3: Extra CMake Options
# $4: Pre-install/Cleanup Command
# $5: Branch (optional, for git)
# $6: Subdir (optional)
# $7: Make target (optional)
# $8: Extra install options (optional)
define cmake-install
	@$(call green, "Installing $2...")
	$(SUDO) mkdir -p $(BINDIR) $(LIB_PATH) $(INCLUDE_PATH)
	rm -rf $(TEMP_DIR)/$2 $(TEMP_DIR)/$2-archive
	if echo "$1" | grep -qE "\.git$$" || [ -n "$(strip $5)" ]; then \
		git clone --depth 1 --recurse-submodules --shallow-submodules $$(if [ -n "$(strip $5)" ]; then echo "--branch $(strip $5)"; fi) $1 $(TEMP_DIR)/$2; \
	else \
		curl -fsSL $1 -o $(TEMP_DIR)/$2-archive; \
		mkdir -p $(TEMP_DIR)/$2; \
		if echo "$1" | grep -qE "\.zip$$"; then \
			unzip -q $(TEMP_DIR)/$2-archive -d $(TEMP_DIR)/$2; \
			MV_DIR=$$(ls -1 $(TEMP_DIR)/$2 | head -n 1); \
			if [ -d "$(TEMP_DIR)/$2/$$MV_DIR" ] && [ $$(ls -1 $(TEMP_DIR)/$2 | wc -l) -eq 1 ]; then \
				mv $(TEMP_DIR)/$2/$$MV_DIR/* $(TEMP_DIR)/$2/ 2>/dev/null || true; \
				rm -rf $(TEMP_DIR)/$2/$$MV_DIR; \
			fi; \
		else \
			tar -xf $(TEMP_DIR)/$2-archive -C $(TEMP_DIR)/$2 --strip-components 1; \
		fi; \
	fi
	cd $(TEMP_DIR)/$2 && \
	cmake \
	-G Ninja \
	-DCMAKE_BUILD_TYPE=Release \
	$(if $(filter-out cmake,$2),-DCMAKE_POLICY_VERSION_MINIMUM=$(CMAKE_VERSION),) \
	-DBUILD_SHARED_LIBS=OFF \
	-DBUILD_STATIC_EXECS=ON \
	-DBUILD_TESTING=OFF \
	-DCMAKE_C_COMPILER="$(CC)" \
	-DCMAKE_C_COMPILER_AR="$(AR)" \
	-DCMAKE_C_COMPILER_RANLIB="$(RANLIB)" \
	-DCMAKE_CXX_COMPILER_AR="$(AR)" \
	-DCMAKE_CXX_COMPILER_RANLIB="$(RANLIB)" \
	-DCMAKE_CXX_COMPILER="$(CXX)" \
	-DCMAKE_ASM_COMPILER="$(CC)" \
	-DCMAKE_ASM_COMPILER_AR="$(AR)" \
	-DCMAKE_ASM_COMPILER_RANLIB="$(RANLIB)" \
	-DCMAKE_AR="$(AR)" \
	-DCMAKE_NM="$(NM)" \
	-DCMAKE_RANLIB="$(RANLIB)" \
	-DCMAKE_MAKE_PROGRAM="$(USR_LOCAL)/bin/ninja" \
	-DCMAKE_INSTALL_PREFIX="$(USR_LOCAL)" \
	-DCMAKE_INSTALL_BINDIR="bin" \
	-DCMAKE_INSTALL_LIBDIR="lib" \
	-DCMAKE_INSTALL_INCLUDEDIR="include" \
	$(if $(filter-out cmake,$2),-DCMAKE_INTERPROCEDURAL_OPTIMIZATION=ON,) \
	-DCMAKE_C_FLAGS="$(CFLAGS)" \
	-DCMAKE_CXX_FLAGS="$(CXXFLAGS)" \
	-DCMAKE_EXE_LINKER_FLAGS="$(LDFLAGS)" \
	-DCMAKE_SHARED_LINKER_FLAGS="$(LDFLAGS)" \
	-DCMAKE_MODULE_LINKER_FLAGS="$(LDFLAGS)" \
	$3 \
	-B $(TEMP_DIR)/$2/build $(TEMP_DIR)/$2$6
	cmake --build $(TEMP_DIR)/$2/build --parallel $(CORES) $(if $7,--target $7,)
	$4
	$(SUDO) cmake --install $(TEMP_DIR)/$2/build $8
	cd $(ROOTDIR)
	rm -rf $(TEMP_DIR)/$2 $(TEMP_DIR)/$2-archive
	$(SUDO) ldconfig
endef

# --- Others ---

define telepresence
	[ -z $(SWAP_IMAGE) ] && IMAGE=$2 || IMAGE=$(SWAP_IMAGE) \
	&& echo "telepresence replaces $(SWAP_DEPLOYMENT_TYPE)/$1 with $${IMAGE}:$(SWAP_TAG)" \
	&& telepresence \
	--swap-deployment $1 \
	--docker-run --rm -it $${IMAGE}:$(SWAP_TAG)
endef

define profile-web
	go tool pprof -http=":6061" \
	$1.bin \
	$1.cpu.out &
	go tool pprof -http=":6062" \
	$1.bin \
	$1.mem.out &
	go tool trace -http=":6063" \
	$1.trace.out
endef
