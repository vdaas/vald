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
red    = printf "\x1b[31m\#\# %s\x1b[0m\n" $1
green  = printf "\x1b[32m\#\# %s\x1b[0m\n" $1
yellow = printf "\x1b[33m\#\# %s\x1b[0m\n" $1
blue   = printf "\x1b[34m\#\# %s\x1b[0m\n" $1
pink   = printf "\x1b[35m\#\# %s\x1b[0m\n" $1
cyan   = printf "\x1b[36m\#\# %s\x1b[0m\n" $1

define go-tool-install
	go install tool
endef

define mkdir
	mkdir -p $1
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

define go-lint
	golangci-lint run --config $(ROOTDIR)/.golangci.json --fix
endef

define go-vet
	cat <(GOARCH=amd64 go vet $(ROOTDIR)/...) \
	  <(GOARCH=386 go vet $(ROOTDIR)/...) \
	  <(GOARCH=arm go vet $(ROOTDIR)/...) \
	  | grep -v "Mutex" | sort | uniq
endef

define go-build
	echo $(GO_SOURCES_INTERNAL)
	echo $(PBGOS)
	echo $(shell find $(ROOTDIR)/cmd/$1 -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	echo $(shell find $(ROOTDIR)/pkg/$1 -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	CFLAGS="$(CFLAGS)" \
	CXXFLAGS="$(CXXFLAGS)" \
	CGO_ENABLED=$(CGO_ENABLED) \
	CGO_CXXFLAGS="$3" \
	CGO_FFLAGS="$3" \
	CGO_LDFLAGS="$3" \
	GO111MODULE=on \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	GOPRIVATE=$(GOPRIVATE) \
	GO_VERSION=$(GO_VERSION) \
	go build \
		--ldflags "-w $2 \
		-extldflags '$3' \
		-X '$(GOPKG)/internal/info.AlgorithmInfo=$5' \
		-X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)' \
		-X '$(GOPKG)/internal/info.BuildTime=$(DATETIME)' \
		-X '$(GOPKG)/internal/info.CGOEnabled=$(if $(filter 1,$(strip $(CGO_ENABLED))),true,false)' \
		-X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
		-X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
		-X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
		-X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
		-X '$(GOPKG)/internal/info.Version=$(VERSION)' \
		-buildid=" \
		-modcacherw \
		-mod=readonly \
		-a \
		-tags "osusergo netgo static_build$4" \
		-trimpath \
		-o $6 \
		$(ROOTDIR)/cmd/$1/main.go
	$6 -version
endef

define go-example-build
	echo $(GO_SOURCES_INTERNAL)
	echo $(PBGOS)
	echo $(shell find $(ROOTDIR)/$1 -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	cd $(ROOTDIR)/$1 && \
	CFLAGS="$(CFLAGS)" \
	CXXFLAGS="$(CXXFLAGS)" \
	CGO_ENABLED=$(CGO_ENABLED) \
	CGO_CXXFLAGS="$3" \
	CGO_FFLAGS="$3" \
	CGO_LDFLAGS="$3" \
	GO111MODULE=on \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	GOPRIVATE=$(GOPRIVATE) \
	GO_VERSION=$(GO_VERSION) \
	go build \
		--ldflags "-w $2 \
		-extldflags '$3' \
		-buildid=" \
		-modcacherw \
		-mod=readonly \
		-a \
		-tags "osusergo netgo static_build$4" \
		-trimpath \
		-o $(ROOTDIR)/$6 \
		main.go
endef

define go-e2e-build
	echo $(GO_SOURCES_INTERNAL)
	echo $(PBGOS)
	echo $(shell find $(ROOTDIR)/$1 -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	CFLAGS="$(CFLAGS)" \
	CXXFLAGS="$(CXXFLAGS)" \
	CGO_ENABLED=$(CGO_ENABLED) \
	CGO_CXXFLAGS="$2" \
	CGO_FFLAGS="$2" \
	CGO_LDFLAGS="$2" \
	GO111MODULE=on \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	GOPRIVATE=$(GOPRIVATE) \
	GO_VERSION=$(GO_VERSION) \
	go test \
		-c \
		-v \
		-race \
		-mod=readonly \
		-tags "e2e" \
		-ldflags "-extldflags '-static' \
		-X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)' \
		-X '$(GOPKG)/internal/info.BuildTime=$(DATETIME)' \
		-X '$(GOPKG)/internal/info.CGOEnabled=$(if $(filter 1,$(strip $(CGO_ENABLED))),true,false)' \
		-X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
		-X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
		-X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
		-X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
		-X '$(GOPKG)/internal/info.Version=$(VERSION)'" \
		-o $(ROOTDIR)/$3 \
		$(ROOTDIR)/$1
endef

define telepresence
	[ -z $(SWAP_IMAGE) ] && IMAGE=$2 || IMAGE=$(SWAP_IMAGE) \
	&& echo "telepresence replaces $(SWAP_DEPLOYMENT_TYPE)/$1 with $${IMAGE}:$(SWAP_TAG)" \
	&& telepresence \
	    --swap-deployment $1 \
	    --docker-run --rm -it $${IMAGE}:$(SWAP_TAG)
	    ## will be available after merge this commit into telepresence head branch
	    ## https://github.com/telepresenceio/telepresence/commit/bb7473fbf19ed4f61796a5e32747e23de6ab03da
	    ## --deployment-type "$(SWAP_DEPLOYMENT_TYPE)"
endef

define run-v2-e2e-crud-test
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_CFLAGS="$(CGO_CFLAGS)" \
	CGO_LDFLAGS="$(CGO_LDFLAGS)" \
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
	E2E_EXPECTED_INDEX="$(E2E_EXPECTED_INDEX)" \
	go test \
	    -race \
	    -v \
	    -mod=readonly \
	    $1 \
	    $(ROOTDIR)/tests/v2/e2e/crud \
	    -tags "e2e" \
	    -timeout $(E2E_TIMEOUT) \
	    -config $(E2E_CONFIG)
endef

define run-e2e-crud-test
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test \
	    -race \
	    -mod=readonly \
	    $1 \
	    -v $(ROOTDIR)/tests/e2e/crud/crud_test.go \
	    -tags "e2e" \
	    -timeout $(E2E_TIMEOUT) \
	    -host=$(E2E_BIND_HOST) \
	    -port=$(E2E_BIND_PORT) \
	    -dataset=$(ROOTDIR)/hack/benchmark/assets/dataset/$(E2E_DATASET_NAME) \
	    -insert-num=$(E2E_INSERT_COUNT) \
	    -correction-insert-num=$(E2E_INSERT_COUNT) \
	    -insert-num=$(E2E_INSERT_COUNT) \
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
	    -kubeconfig=$(KUBECONFIG)
endef

define run-e2e-crud-faiss-test
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test \
	    -race \
	    -mod=readonly \
	    $1 \
	    -v $(ROOTDIR)/tests/e2e/crud/crud_faiss_test.go \
	    -tags "e2e" \
	    -timeout $(E2E_TIMEOUT) \
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
	    -namespace=$(E2E_TARGET_NAMESPACE)
endef

define run-e2e-multi-crud-test
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test \
	    -race \
	    -mod=readonly \
	    $1 \
            -v $(ROOTDIR)/tests/e2e/multiapis/multiapis_test.go \
	    -tags "e2e" \
	    -timeout $(E2E_TIMEOUT) \
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
	    -kubeconfig=$(KUBECONFIG)
endef

define run-e2e-max-dim-test
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test \
	    -race \
	    -mod=readonly \
            -v $(ROOTDIR)/tests/e2e/performance/max_vector_dim_test.go \
	    -tags "e2e" \
	    -timeout $(E2E_TIMEOUT) \
            -file $(E2E_MAX_DIM_RESULT_FILEPATH) \
	    -host=$(E2E_BIND_HOST) \
	    -port=$(E2E_BIND_PORT) \
	    -bit=${E2E_MAX_DIM_BIT} \
	    -portforward=$(E2E_PORTFORWARD_ENABLED) \
	    -portforward-pod-name=$(E2E_TARGET_POD_NAME) \
	    -portforward-pod-port=$(E2E_TARGET_PORT) \
	    -namespace=$(E2E_TARGET_NAMESPACE) \
	    -kubeconfig=$(KUBECONFIG)
endef

define run-e2e-sidecar-test
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	CGO_LDFLAGS="$(TEST_LDFLAGS)" \
	go test \
	    -race \
	    -mod=readonly \
	    $1 \
	    -v $(ROOTDIR)/tests/e2e/sidecar/sidecar_test.go \
	    -tags "e2e" \
	    -timeout $(E2E_TIMEOUT) \
	    -host=$(E2E_BIND_HOST) \
	    -port=$(E2E_BIND_PORT) \
	    -dataset=$(ROOTDIR)/hack/benchmark/assets/dataset/$(E2E_DATASET_NAME) \
	    -insert-num=$(E2E_INSERT_COUNT) \
	    -search-num=$(E2E_SEARCH_COUNT) \
	    -portforward=$(E2E_PORTFORWARD_ENABLED) \
	    -portforward-pod-name=$(E2E_TARGET_POD_NAME) \
	    -portforward-pod-port=$(E2E_TARGET_PORT) \
	    -namespace=$(E2E_TARGET_NAMESPACE) \
	    -kubeconfig=$(KUBECONFIG)
endef

# This function generate only implementation tests, with the following conditions:
# - Generate all go tests on `./cmd`, `./hack` and `./internal` packages with some exclusion (see $GO_SOURCES)
# - Skip generating go tests under './pkg/*/router/*' and './pkg/*/handler/test/*' package
# - Generate only 'New()' test on './pkg/*/usecase'
# - Generate only exported function tests on `./pkg` package
define gen-go-test-sources
	@$(call green, "Generating go test files in parallel (cores: $(CORES))...")
	@for f in $(GO_SOURCES); do echo "$$f"; done | \
	xargs -I {} -P$(CORES) bash -c '\
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
		gotests -w -template_dir $(ROOTDIR)/assets/test/templates/common $$GOTESTS_OPTION $(patsubst %_test.go,%.go,$$f); \
		RESULT=$$?; \
		if [ ! $$RESULT -eq 0 ]; then \
			echo "Error generating test for $$f: $$RESULT"; \
			exit 1; \
		fi; \
	'
endef

# This function generate only option tests, with the following conditions:
# - Generate all go tests on `./cmd`, `./hack` and `./internal` packages with exclusion (see $GO_SOURCES)
# - Skip generating go tests under './pkg/*/router' and './pkg/*/handler/test' and './pkg/*/usecase' package
# - Generate only exported function tests on `./pkg` package
define gen-go-option-test-sources
	@$(call green, "Generating go option test files in parallel (cores: $(CORES))...")
	@for f in $(GO_OPTION_SOURCES); do echo "$$f"; done | \
	xargs -I {} -P$(CORES) bash -c '\
		f="{}"; \
		GOTESTS_OPTION="-all"; \
		if [[ $$f =~ \.\/pkg\/.*\/router\/.* || $$f =~ \.\/pkg\/.*\/handler\/rest\/.* || $$f =~ \.\/pkg\/.*\/usecase\/.* ]]; then \
			echo "Skip generating go option test file: $$f"; \
			exit 0; \
		elif [[ $$f =~ \.\/pkg\/.* ]]; then \
			GOTESTS_OPTION=" -exported "; \
		fi; \
		echo "Generating go option test file: $$f with option $$GOTESTS_OPTION"; \
		gotests -w -template_dir $(ROOTDIR)/assets/test/templates/option $$GOTESTS_OPTION $(patsubst %_test.go,%.go,$$f); \
		RESULT=$$?; \
		if [ ! $$RESULT -eq 0 ]; then \
			echo "Error generating option test for $$f: $$RESULT"; \
			exit 1; \
		fi; \
	'
endef

define gen-license
	BIN_PATH="$(TEMP_DIR)/vald-license-gen"; \
	rm -rf $$BIN_PATH; \
	MAINTAINER=$2 \
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	go build -modcacherw \
		-mod=readonly \
		-a \
		-tags "osusergo netgo static_build" \
		-trimpath \
		-o $$BIN_PATH $(ROOTDIR)/hack/license/gen/main.go; \
	$$BIN_PATH $1; \
	rm -rf $$BIN_PATH
endef

define gen-dockerfile
	BIN_PATH="$(TEMP_DIR)/vald-dockerfile-gen"; \
	rm -rf $$BIN_PATH; \
	MAINTAINER=$2 \
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	go build -modcacherw \
		-mod=readonly \
		-a \
		-tags "osusergo netgo static_build" \
		-trimpath \
		-o $$BIN_PATH $(ROOTDIR)/hack/docker/gen/main.go; \
	$$BIN_PATH $1; \
	rm -rf $$BIN_PATH
endef

define gen-dashboard
	go run $(ROOTDIR)/hack/grafana/gen/src
endef

define gen-vald-helm-schema
	BIN_PATH="$(TEMP_DIR)/vald-helm-schema-gen"; \
	rm -rf $$BIN_PATH; \
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	go build -modcacherw \
		-mod=readonly \
		-a \
		-tags "osusergo netgo static_build" \
		-trimpath \
		-o $$BIN_PATH $(ROOTDIR)/hack/helm/schema/gen/main.go; \
	$$BIN_PATH charts/$1.yaml > charts/$1.schema.json; \
	rm -rf $$BIN_PATH
endef

define gen-vald-crd
	if [[ -f $(ROOTDIR)/charts/$1/crds/$2.yaml ]]; then \
		mv $(ROOTDIR)/charts/$1/crds/$2.yaml $(TEMP_DIR)/$2.yaml; \
	fi;
	BIN_PATH="$(TEMP_DIR)/vald-helm-crd-schema-gen"; \
	rm -rf $$BIN_PATH; \
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	go build -modcacherw \
		-mod=readonly \
		-a \
		-tags "osusergo netgo static_build" \
		-trimpath \
		-o $$BIN_PATH $(ROOTDIR)/hack/helm/schema/crd/main.go; \
	$$BIN_PATH $(ROOTDIR)/charts/$3.yaml > $(TEMP_DIR)/$2-spec.yaml; \
	rm -rf $$BIN_PATH; \
	$(BINDIR)/yq eval-all 'select(fileIndex==0).spec.versions[0].schema.openAPIV3Schema.properties.spec = select(fileIndex==1).spec | select(fileIndex==0)' \
	$(TEMP_DIR)/$2.yaml $(TEMP_DIR)/$2-spec.yaml > $(ROOTDIR)/charts/$1/crds/$2.yaml
endef

define update-github-actions
	@set -e; for ACTION_NAME in $1; do \
		if [ -n "$$ACTION_NAME" ] && [ "$$ACTION_NAME" != "security-and-quality" ]; then \
			FILE_NAME=`echo $$ACTION_NAME | tr '/' '_' | tr '-' '_' | tr '[:lower:]' '[:upper:]'`; \
			if [ -n "$$FILE_NAME" ]; then \
				case "$$ACTION_NAME" in \
					"aquasecurity/trivy-action" | "machine-learning-apps/actions-chatops" ) VERSION="master";; \
					* ) \
						REPO_NAME=`echo $$ACTION_NAME | cut -d'/' -f1-2`; \
						echo "$$ACTION_NAME to $$REPO_NAME"; \
						VERSION=`curl -fsSL https://api.github.com/repos/$$REPO_NAME/tags?per_page=1 | grep -Po '"name": "\K.*?(?=")' | head -n1 | sed 's/v//g' | sed -E 's/[^0-9.]+//g'`; \
						;; \
				esac; \
				if [ -n "$$VERSION" ]; then \
					OLD_VERSION=`cat $(ROOTDIR)/versions/actions/$$FILE_NAME`; \
					echo "updating $$ACTION_NAME version file $$FILE_NAME from $$OLD_VERSION to $$VERSION"; \
					echo $$VERSION > $(ROOTDIR)/versions/actions/$$FILE_NAME; \
				else \
					VERSION=`cat $(ROOTDIR)/versions/actions/$$FILE_NAME`; \
					echo "No version found for $$ACTION_NAME version file $$FILE_NAME=$$VERSION"; \
				fi; \
				if [ "$$ACTION_NAME" = "cirrus-actions/rebase" ]; then \
					VERSION_PREFIX=$$VERSION; \
					find $(ROOTDIR)/.github -type f -exec sed -i "s%$$ACTION_NAME@.*%$$ACTION_NAME@$$VERSION_PREFIX%g" {} +; \
				elif echo $$VERSION | grep -qE '^[0-9]'; then \
					VERSION_PREFIX=`echo $$VERSION | cut -c 1`; \
					find $(ROOTDIR)/.github -type f -exec sed -i "s%$$ACTION_NAME@.*%$$ACTION_NAME@v$$VERSION_PREFIX%g" {} +; \
				else \
					VERSION_PREFIX=$$VERSION; \
					find $(ROOTDIR)/.github -type f -exec sed -i "s%$$ACTION_NAME@.*%$$ACTION_NAME@$$VERSION_PREFIX%g" {} +; \
				fi; \
			else \
				echo "No action version file found for $$ACTION_NAME version file $$FILE_NAME" >&2; \
			fi \
		else \
			echo "No action found for $$ACTION_NAME" >&2; \
		fi \
	done
endef

define gen-deadlink-checker
	BIN_PATH="$(TEMP_DIR)/vald-deadlink-checker-gen"; \
	rm -rf $$BIN_PATH; \
	MAINTAINER=$2 \
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
	go build -modcacherw \
		-mod=readonly \
		-a \
		-tags "osusergo netgo static_build" \
		-trimpath \
		-o $$BIN_PATH $(ROOTDIR)/hack/tools/deadlink/main.go; \
	$$BIN_PATH -path $3 -ignore-path $4 -format $5 $1; \
	rm -rf $$BIN_PATH
endef

define gen-api-document
	buf generate --template=apis/docs/buf.gen.tmpl.yaml --path $2
	cat apis/docs/v1/payload.md.tmpl apis/docs/v1/_doc.md.tmpl > apis/docs/v1/doc.md.tmpl; \
	buf generate --template=apis/docs/buf.gen.doc.yaml --path $2; \
	mv $(ROOTDIR)/apis/docs/v1/doc.md $1; \
	rm apis/docs/v1/*doc.md.tmpl
endef
