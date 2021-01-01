#
# Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

define go-get
	GO111MODULE=on go get -u $1
endef

define go-get-no-mod
	GO111MODULE=off go get -u $1
endef

define mkdir
	mkdir -p $1
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
	find ./ -type d -name .git -prune -o -type f -regex '.*\.go' -print | xargs goimports -w
	golangci-lint run --enable-all --disable=gochecknoglobals --fix --color always -j 16 --skip-dirs apis/grpc --exclude-use-default=false ./...
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
