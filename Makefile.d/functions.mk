#
# Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
                --go-vtproto_opt=features=grpc+marshal+unmarshal+size+pool \
                --go-vtproto_opt=pool=$(GOPKG)/apis/proto/v1/payload.Search.Request \
                --go-vtproto_opt=pool=$(GOPKG)/apis/proto/v1/payload.Object.Vector \
		$1
endef

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
	cat <(GOARCH=amd64 go vet ./...) \
	  <(GOARCH=386 go vet ./...) \
	  <(GOARCH=arm go vet ./...) \
	  | rg -v "Mutex" | sort | uniq
endef


define telepresence
	[ -z $(SWAP_IMAGE) ] && IMAGE=$2 || IMAGE=$(SWAP_IMAGE) \
	&& echo "telepresence replaces $(SWAP_DEPLOYMENT_TYPE)/$1 with $${IMAGE}:$(SWAP_TAG)" \
	&& telepresence \
	    --swap-deployment $1 \
	    --docker-run --rm -it $${IMAGE}:$(SWAP_TAG)
	    ## will be available after merge this commit into telepresence master branch
	    ## https://github.com/telepresenceio/telepresence/commit/bb7473fbf19ed4f61796a5e32747e23de6ab03da
	    ## --deployment-type "$(SWAP_DEPLOYMENT_TYPE)"
endef


define run-e2e-crud-test
	go test \
	    -race \
	    -mod=readonly \
	    $1 \
	    -v $(ROOTDIR)/tests/e2e/crud/crud_test.go \
	    -tags "e2e" \
	    -timeout $(E2E_TIMEOUT) \
	    -host=$(E2E_BIND_HOST) \
	    -port=$(E2E_BIND_PORT) \
	    -dataset=$(ROOTDIR)/hack/benchmark/assets/dataset/$(E2E_DATASET_NAME).hdf5 \
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
	    -namespace=$(E2E_TARGET_NAMESPACE)
endef
