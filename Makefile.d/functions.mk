#
# Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
red    = printf "\x1b[31m\#\# %s\x1b[0m\n" $1
green  = printf "\x1b[32m\#\# %s\x1b[0m\n" $1
yellow = printf "\x1b[33m\#\# %s\x1b[0m\n" $1
blue   = printf "\x1b[34m\#\# %s\x1b[0m\n" $1
pink   = printf "\x1b[35m\#\# %s\x1b[0m\n" $1
cyan   = printf "\x1b[36m\#\# %s\x1b[0m\n" $1

define go-install
	GO111MODULE=on go install $1@latest
endef

define mkdir
	mkdir -p $1
endef

define proto-code-gen
	protoc \
		$(PROTO_PATHS:%=-I %) \
                --go_out=$(GOPATH)/src --plugin protoc-gen-go="$(GOPATH)/bin/protoc-gen-go" \
                --go-vtproto_out=$(GOPATH)/src --plugin protoc-gen-go-vtproto="$(GOPATH)/bin/protoc-gen-go-vtproto" \
                --go-vtproto_opt=features=grpc+marshal+unmarshal+size+equal+clone \
		$1
endef
                # --go-vtproto_opt=pool=$(ROOTDIR)/apis/proto/v1/payload.Object.Vector \
                # --go-vtproto_opt=pool=$(ROOTDIR)/apis/proto/v1/payload.Insert.MultiRequest \
                # --go-vtproto_opt=pool=$(ROOTDIR)/apis/proto/v1/payload.Insert.Request \
                # --go-vtproto_opt=pool=$(ROOTDIR)/apis/proto/v1/payload.Object.Vector \
                # --go-vtproto_opt=pool=$(ROOTDIR)/apis/proto/v1/payload.Object.Vectors \
                # --go-vtproto_opt=pool=$(ROOTDIR)/apis/proto/v1/payload.Search.ObjectRequest \
                # --go-vtproto_opt=pool=$(ROOTDIR)/apis/proto/v1/payload.Search.Request \
                # --go-vtproto_opt=pool=$(ROOTDIR)/apis/proto/v1/payload.Update.MultiRequest \
                # --go-vtproto_opt=pool=$(ROOTDIR)/apis/proto/v1/payload.Update.Request \
                # --go-vtproto_opt=pool=$(ROOTDIR)/apis/proto/v1/payload.Upsert.MultiRequest \
                # --go-vtproto_opt=pool=$(ROOTDIR)/apis/proto/v1/payload.Upsert.Request \

define protoc-gen
	protoc \
		$(PROTO_PATHS:%=-I %) \
		$2 \
		$1
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

# Generate all go tests on cmd, hack and internal packages with exclusion (see GO_SOURCES),
# and generate only exported function tests on pkg package with exclusion (see GO_SOURCES_PKG),
# and generate only New() test on pkg/*/usecase.
define gen-go-test-sources
	@for f in $(GO_SOURCES); do \
		echo "Generating go test file: $$f"; \
		gotests -w -template_dir $(ROOTDIR)/assets/test/templates/common -all $(patsubst %_test.go,%.go,$$f); \
		RESULT=$$?; \
		if [ ! $$RESULT -eq 0 ]; then \
			echo $$RESULT; \
			exit 1; \
		fi; \
	done
	@for f in $(GO_SOURCES_PKG); do \
		echo "Generating go test file: $$f"; \
		gotests -w -exported -template_dir $(ROOTDIR)/assets/test/templates/common -all $(patsubst %_test.go,%.go,$$f); \
		RESULT=$$?; \
		if [ ! $$RESULT -eq 0 ]; then \
			echo $$RESULT; \
			exit 1; \
		fi; \
	done
	@for f in $(GO_SOURCES_PKG_USECASE); do \
		echo "Generating go test file: $$f"; \
		gotests -w -exported -template_dir $(ROOTDIR)/assets/test/templates/common -only New $(patsubst %_test.go,%.go,$$f); \
		RESULT=$$?; \
		if [ ! $$RESULT -eq 0 ]; then \
			echo $$RESULT; \
			exit 1; \
		fi; \
	done
endef

define gen-go-option-test-sources
	@for f in $(GO_OPTION_SOURCES); do \
		echo "Generating go option test file: $$f"; \
		gotests -w -template_dir $(ROOTDIR)/assets/test/templates/common -all $(patsubst %_test.go,%.go,$$f); \
		RESULT=$$?; \
		if [ ! $$RESULT -eq 0 ]; then \
			echo $$RESULT; \
			exit 1; \
		fi; \
	done
	@for f in $(GO_OPTION_SOURCES_PKG); do \
		echo "Generating go option test file: $$f"; \
		gotests -w -exported -template_dir $(ROOTDIR)/assets/test/templates/common -all $(patsubst %_test.go,%.go,$$f); \
		RESULT=$$?; \
		if [ ! $$RESULT -eq 0 ]; then \
			echo $$RESULT; \
			exit 1; \
		fi; \
	done
endef
