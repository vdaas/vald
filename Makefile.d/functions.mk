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

define bench-pprof
	rm -rf $1
	mkdir -p $1
	@$(call green, "starting $4 $2 benchmark")
	go test -count=1 \
		-timeout=1h \
		-bench=$3 \
		-benchmem \
		-o $1/$2.bin \
		-cpuprofile $1/cpu-$4.out \
		-memprofile $1/mem-$4.out \
		-trace $1/trace-$4.out \
		$5 \
		| tee $1/result-$4.out
	go tool pprof --svg \
		$1/$2.bin \
		$1/cpu-$4.out \
		> $1/cpu-$4.svg
	go tool pprof --svg \
		$1/$2.bin \
		$1/mem-$4.out \
		> $1/mem-$4.svg
endef

define profile-web
	@$(call green, "starting $3 $2 profiler")
	go tool pprof -http=$4 \
		$1/$2.bin \
		$1/cpu-$3.out &
	go tool pprof -http=$5 \
		$1/$2.bin \
		$1/mem-$3.out &
	go tool trace -http=$6 \
		$1/trace-$3.out
endef

define go-lint
	find ./ -type f -regex ".*\.go" | xargs goimports -w
	golangci-lint run --enable-all --disable=gochecknoglobals --fix --color always -j 16 --skip-dirs apis/grpc --exclude-use-default=false ./...
endef
