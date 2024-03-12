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
red    = printf "\x1b[31m\#\# %s\x1b[0m\n" $1
green  = printf "\x1b[32m\#\# %s\x1b[0m\n" $1
yellow = printf "\x1b[33m\#\# %s\x1b[0m\n" $1
blue   = printf "\x1b[34m\#\# %s\x1b[0m\n" $1
pink   = printf "\x1b[35m\#\# %s\x1b[0m\n" $1
cyan   = printf "\x1b[36m\#\# %s\x1b[0m\n" $1

define go-install
	GO111MODULE=on go install \
	    -mod=readonly \
	    $1@latest
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
	golangci-lint run --config .golangci.yml --fix
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
	CGO_CXXFLAGS="-g -Ofast -march=native" \
	CGO_FFLAGS="-g -Ofast -march=native" \
	CGO_LDFLAGS="-g -Ofast -march=native" \
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
		-X '$(GOPKG)/internal/info.CGOEnabled=$(CGO_ENABLED)' \
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

define run-e2e-crud-test
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
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

define run-e2e-crud-faiss-test
	GOPRIVATE=$(GOPRIVATE) \
	GOARCH=$(GOARCH) \
	GOOS=$(GOOS) \
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
	@for f in $(GO_SOURCES); do \
		GOTESTS_OPTION=" -all "; \
		if [[ $$f =~ \.\/pkg\/.*\/router\/.* || $$f =~ \.\/pkg\/.*\/handler\/rest\/.* ]]; then \
			echo "Skip generating go test file: $$f"; \
			continue; \
		elif [[ $$f =~ \.\/pkg\/.*\/usecase\/.* ]]; then \
			GOTESTS_OPTION=" -only New "; \
		elif [[ $$f =~ \.\/pkg\/.* ]]; then \
			GOTESTS_OPTION=" -exported "; \
		fi; \
		echo "Generating go test file: $$f" with option  $$GOTESTS_OPTION; \
		gotests -w -template_dir $(ROOTDIR)/assets/test/templates/common $$GOTESTS_OPTION $(patsubst %_test.go,%.go,$$f); \
		RESULT=$$?; \
		if [ ! $$RESULT -eq 0 ]; then \
			echo $$RESULT; \
			exit 1; \
		fi; \
	done
endef

# This function generate only option tests, with the following conditions:
# - Generate all go tests on `./cmd`, `./hack` and `./internal` packages with exclusion (see $GO_SOURCES)
# - Skip generating go tests under './pkg/*/router' and './pkg/*/handler/test' and './pkg/*/usecase' package
# - Generate only exported function tests on './pkg` package
define gen-go-option-test-sources
	@for f in $(GO_OPTION_SOURCES); do \
		GOTESTS_OPTION=" -all "; \
		if [[ $$f =~ \.\/pkg\/.*\/router\/.* || $$f =~ \.\/pkg\/.*\/handler\/rest\/.* || $$f =~ \.\/pkg\/.*\/usecase\/.* ]]; then \
			echo "Skip generating go option test file: $$f"; \
			continue; \
		elif [[ $$f =~ \.\/pkg\/.* ]]; then \
			GOTESTS_OPTION=" -exported "; \
		fi; \
		echo "Generating go option test file: $$f"; \
		gotests -w -template_dir $(ROOTDIR)/assets/test/templates/common -all $(patsubst %_test.go,%.go,$$f); \
		RESULT=$$?; \
		if [ ! $$RESULT -eq 0 ]; then \
			echo $$RESULT; \
			exit 1; \
		fi; \
	done
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
	@for ACTION_NAME in $1; do \
		if [ -n "$$ACTION_NAME" ]; then \
			FILE_NAME=`echo $$ACTION_NAME | tr '/' '_' | tr '-' '_' | tr '[:lower:]' '[:upper:]'`; \
			if [ -n "$$FILE_NAME" ]; then \
				if [ "$$ACTION_NAME" = "aquasecurity/trivy-action" ] || [ "$$ACTION_NAME" = "machine-learning-apps/actions-chatops" ]; then \
					VERSION="master"; \
				elif [ "$$ACTION_NAME" = "softprops/action-gh-release" ]; then \
					VERSION="1.0.0"; \
				else \
					REPO_NAME=`echo $$ACTION_NAME | cut -d'/' -f1-2`; \
					VERSION=`curl --silent https://api.github.com/repos/$$REPO_NAME/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' | sed 's/v//g' | sed -E 's/[^0-9.]+//g'`;\
				fi; \
				if [ -n "$$VERSION" ]; then \
					echo "updating $$ACTION_NAME version file $$FILE_NAME to $$VERSION"; \
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

