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

cmd/agent/core/ngt/ngt: \
	ngt/install \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/agent/core/ngt -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/agent/core/ngt ./pkg/agent/internal -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	export CFLAGS="$(CFLAGS)" \
	    && export CXXFLAGS="$(CXXFLAGS)" \
	    && export CGO_ENABLED=1 \
	    && export CGO_CXXFLAGS="-g -Ofast -march=native" \
	    && export CGO_FFLAGS="-g -Ofast -march=native" \
	    && export CGO_LDFLAGS="-g -Ofast -march=native" \
	    && export GO111MODULE=on \
	    && go build \
	    --ldflags "-s -w -linkmode 'external' \
	    -extldflags '-static -fPIC -m64 -pthread -fopenmp -std=c++17 -lstdc++ -lm' \
	    -X '$(GOPKG)/internal/info.Version=$(VERSION)' \
	    -X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
	    -X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
	    -X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
	    -X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
	    -X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
	    -X '$(GOPKG)/internal/info.NGTVersion=$(NGT_VERSION)' \
	    -X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)'" \
	    -a \
	    -tags "cgo netgo" \
	    -trimpath \
	    -installsuffix "cgo netgo" \
	    -o $@ \
	    $(dir $@)main.go

cmd/agent/sidecar/sidecar: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/agent/sidecar -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/agent/sidecar ./pkg/agent/internal -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	export CGO_ENABLED=1 \
	    && export GO111MODULE=on \
	    && go build \
	    --ldflags "-s -w -linkmode 'external' \
	    -extldflags '-static' \
	    -X '$(GOPKG)/internal/info.Version=$(VERSION)' \
	    -X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
	    -X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
	    -X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
	    -X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
	    -X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
	    -X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)'" \
	    -a \
	    -tags netgo \
	    -trimpath \
	    -installsuffix netgo \
	    -o $@ \
	    $(dir $@)main.go

cmd/discoverer/k8s/discoverer: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/discoverer/k8s -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/discoverer/k8s -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	export CGO_ENABLED=1 \
	    && export GO111MODULE=on \
	    && go build \
	    --ldflags "-s -w -linkmode 'external' \
	    -extldflags '-static' \
	    -X '$(GOPKG)/internal/info.Version=$(VERSION)' \
	    -X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
	    -X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
	    -X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
	    -X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
	    -X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
	    -X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)'" \
	    -a \
	    -tags netgo \
	    -installsuffix netgo \
	    -trimpath \
	    -o $@ \
	    $(dir $@)main.go

cmd/gateway/vald/vald: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/gateway/vald -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/gateway/vald -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	export CGO_ENABLED=1 \
	    && export GO111MODULE=on \
	    && go build \
	    --ldflags "-s -w -linkmode 'external' \
	    -extldflags '-static' \
	    -X '$(GOPKG)/internal/info.Version=$(VERSION)' \
	    -X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
	    -X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
	    -X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
	    -X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
	    -X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
	    -X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)'" \
	    -a \
	    -tags netgo \
	    -installsuffix netgo \
	    -trimpath \
	    -o $@ \
	    $(dir $@)main.go

cmd/meta/redis/meta: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/meta/redis -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/meta/redis -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	export CGO_ENABLED=1 \
	    && export GO111MODULE=on \
	    && go build \
	    --ldflags "-s -w -linkmode 'external' \
	    -extldflags '-static' \
	    -X '$(GOPKG)/internal/info.Version=$(VERSION)' \
	    -X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
	    -X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
	    -X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
	    -X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
	    -X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
	    -X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)'" \
	    -a \
	    -tags netgo \
	    -installsuffix netgo \
	    -trimpath \
	    -o $@ \
	    $(dir $@)main.go

cmd/meta/cassandra/meta: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/meta/cassandra -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/meta/cassandra -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	export CGO_ENABLED=1 \
	    && export GO111MODULE=on \
	    && go build \
	    --ldflags "-s -w -linkmode 'external' \
	    -extldflags '-static' \
	    -X '$(GOPKG)/internal/info.Version=$(VERSION)' \
	    -X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
	    -X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
	    -X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
	    -X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
	    -X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
	    -X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)'" \
	    -a \
	    -tags netgo \
	    -installsuffix netgo \
	    -trimpath \
	    -o $@ \
	    $(dir $@)main.go

cmd/manager/backup/mysql/backup: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/manager/backup/mysql -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/manager/backup/mysql -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	export CGO_ENABLED=1 \
	    && export GO111MODULE=on \
	    && go build \
	    --ldflags "-s -w -linkmode 'external' \
	    -extldflags '-static' \
	    -X '$(GOPKG)/internal/info.Version=$(VERSION)' \
	    -X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
	    -X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
	    -X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
	    -X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
	    -X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
	    -X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)'" \
	    -a \
	    -tags netgo \
	    -installsuffix netgo \
	    -trimpath \
	    -o $@ \
	    $(dir $@)main.go

cmd/manager/backup/cassandra/backup: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/manager/backup/cassandra -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/manager/backup/cassandra -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	export CGO_ENABLED=1 \
	    && export GO111MODULE=on \
	    && go build \
	    --ldflags "-s -w -linkmode 'external' \
	    -extldflags '-static' \
	    -X '$(GOPKG)/internal/info.Version=$(VERSION)' \
	    -X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
	    -X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
	    -X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
	    -X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
	    -X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
	    -X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)'" \
	    -a \
	    -tags netgo \
	    -installsuffix netgo \
	    -trimpath \
	    -o $@ \
	    $(dir $@)main.go

cmd/manager/compressor/compressor: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/manager/compressor -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/manager/compressor -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	export CGO_ENABLED=1 \
	    && export GO111MODULE=on \
	    && go build \
	    --ldflags "-s -w -linkmode 'external' \
	    -extldflags '-static' \
	    -X '$(GOPKG)/internal/info.Version=$(VERSION)' \
	    -X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
	    -X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
	    -X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
	    -X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
	    -X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
	    -X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)'" \
	    -a \
	    -tags netgo \
	    -trimpath \
	    -installsuffix netgo \
	    -o $@ \
	    $(dir $@)main.go

cmd/manager/index/index: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/manager/index -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/manager/index -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	export CGO_ENABLED=1 \
	    && export GO111MODULE=on \
	    && go build \
	    --ldflags "-s -w -linkmode 'external' \
	    -extldflags '-static' \
	    -X '$(GOPKG)/internal/info.Version=$(VERSION)' \
	    -X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
	    -X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
	    -X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
	    -X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
	    -X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
	    -X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)'" \
	    -a \
	    -tags netgo \
	    -trimpath \
	    -installsuffix netgo \
	    -o $@ \
	    $(dir $@)main.go

cmd/manager/replication/agent/agent: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/manager/replication/agent -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/manager/replication/agent -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	export CGO_ENABLED=1 \
	    && export GO111MODULE=on \
	    && go build \
	    --ldflags "-s -w -linkmode 'external' \
	    -extldflags '-static' \
	    -X '$(GOPKG)/internal/info.Version=$(VERSION)' \
	    -X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
	    -X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
	    -X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
	    -X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
	    -X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
	    -X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)'" \
	    -a \
	    -tags netgo \
	    -trimpath \
	    -installsuffix netgo \
	    -o $@ \
	    $(dir $@)main.go

cmd/manager/replication/controller/controller: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/manager/replication/controller -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/manager/replication/controller -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	export CGO_ENABLED=1 \
	    && export GO111MODULE=on \
	    && go build \
	    --ldflags "-s -w -linkmode 'external' \
	    -extldflags '-static' \
	    -X '$(GOPKG)/internal/info.Version=$(VERSION)' \
	    -X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
	    -X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
	    -X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
	    -X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
	    -X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
	    -X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)'" \
	    -a \
	    -tags netgo \
	    -trimpath \
	    -installsuffix netgo \
	    -o $@ \
	    $(dir $@)main.go
