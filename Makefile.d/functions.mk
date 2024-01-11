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
	golangci-lint run --config .golangci.yml
endef

define go-vet
	cat <(GOARCH=amd64 go vet $(ROOTDIR)/...) \
	  <(GOARCH=386 go vet $(ROOTDIR)/...) \
	  <(GOARCH=arm go vet $(ROOTDIR)/...) \
	  | rg -v "Mutex" | sort | uniq
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

define run-e2e-multi-crud-test
	GOPRIVATE=$(GOPRIVATE) \
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

define gen-vald-crd
	mv charts/$1/crds/$2.yaml $(TEMP_DIR)/$2.yaml
	GOPRIVATE=$(GOPRIVATE) \
	go run -mod=readonly hack/helm/schema/crd/main.go \
	charts/$1/$3.yaml > $(TEMP_DIR)/$2-spec.yaml
	$(BINDIR)/yq eval-all 'select(fileIndex==0).spec.versions[0].schema.openAPIV3Schema.properties.spec = select(fileIndex==1).spec | select(fileIndex==0)' \
	$(TEMP_DIR)/$2.yaml $(TEMP_DIR)/$2-spec.yaml > charts/$1/crds/$2.yaml
endef
