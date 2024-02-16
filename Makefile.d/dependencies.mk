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

.PHONY: update/libs
## update vald libraries including tools
update/libs: \
	update/chaos-mesh \
	update/faiss \
	update/go \
	update/golangci-lint \
	update/helm \
	update/helm-docs \
	update/helm-operator \
	update/jaeger-operator \
	update/k3s \
	update/kind \
	update/kube-linter \
	update/kubectl \
	update/ngt \
	update/prometheus-stack \
	update/protobuf \
	update/reviewdog \
	update/telepresence \
	update/vald \
	update/valdcli \
	update/yq \
	update/zlib \
	update/hdf5

.PHONY: go/download
## download Go package dependencies
go/download:
	GOPRIVATE=$(GOPRIVATE) go mod download

.PHONY: go/deps
## install Go package dependencies
go/deps:
	sed -i "3s/go [0-9]\+\.[0-9]\+\(\.[0-9]\+\)\?/go $(GO_VERSION)/g" $(ROOTDIR)/hack/go.mod.default
	if $(GO_CLEAN_DEPS); then \
        	rm -rf $(ROOTDIR)/vendor \
        		/go/pkg \
        		$(GOCACHE) \
        		$(ROOTDIR)/go.sum \
        		$(ROOTDIR)/go.mod ; \
        	cp $(ROOTDIR)/hack/go.mod.default $(ROOTDIR)/go.mod ; \
        	GOPRIVATE=$(GOPRIVATE) go mod tidy ; \
        	go clean -cache -modcache -testcache -i -r ; \
        	rm -rf $(ROOTDIR)/vendor \
        		/go/pkg \
        		$(GOCACHE) \
        		$(ROOTDIR)/go.sum \
        		$(ROOTDIR)/go.mod ; \
        	cp $(ROOTDIR)/hack/go.mod.default $(ROOTDIR)/go.mod ; \
	fi
	cp $(ROOTDIR)/hack/go.mod.default $(ROOTDIR)/go.mod
	GOPRIVATE=$(GOPRIVATE) go mod tidy
	go get -u all 2>/dev/null || true

.PHONY: go/example/deps
## install Go package dependencies
go/example/deps:
	rm -rf $(ROOTDIR)/vendor \
		$(GOCACHE) \
	        $(ROOTDIR)/example/client/vendor \
	        $(ROOTDIR)/example/client/go.mod \
	        $(ROOTDIR)/example/client/go.sum
	sed -i "3s/go [0-9]\+\.[0-9]\+\(\.[0-9]\+\)\?/go $(GO_VERSION)/g" $(ROOTDIR)/example/client/go.mod.default
	cp $(ROOTDIR)/example/client/go.mod.default $(ROOTDIR)/example/client/go.mod
	cd $(ROOTDIR)/example/client && GOPRIVATE=$(GOPRIVATE) go mod tidy && cd -

.PHONY: rust/deps
## install Rust package dependencies
rust/deps: \
	rust/install
	if [ -x "$(CARGO_HOME)/bin/cargo" ]; then \
		cd $(ROOTDIR)/rust \
			&& $(CARGO_HOME)/bin/rustup default stable \
			&& $(CARGO_HOME)/bin/cargo update \
			&& cd -;\
	else \
		echo "Cargo not found. Please install Cargo or add it to your PATH."; \
		exit 1; \
	fi

.PHONY: update/chaos-mesh
## update chaos-mesh version
update/chaos-mesh:
	curl --silent https://api.github.com/repos/chaos-mesh/chaos-mesh/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' | sed 's/v//g' > $(ROOTDIR)/versions/CHAOS_MESH_VERSION

.PHONY: update/k3s
## update k3s version
update/k3s:
	@{ \
		RESULT=$$(curl --silent https://hub.docker.com/v2/repositories/rancher/k3s/tags?page_size=1000 | jq -r '.results[].name' | grep -E '.*-k3s[0-9]+$$' | grep -v rc | sort -Vr | head -n 1); \
		if [ -n "$$RESULT" ]; then \
			echo $$RESULT > $(ROOTDIR)/versions/K3S_VERSION; \
		else \
			echo "No version found" >&2; \
		fi \
	}

.PHONY: update/go
## update go version
update/go:
	curl --silent https://go.dev/VERSION?m=text | head -n 1 | sed -e 's/go//g' > $(ROOTDIR)/versions/GO_VERSION

.PHONY: update/golangci-lint
## update golangci-lint version
update/golangci-lint:
	curl --silent https://api.github.com/repos/golangci/golangci-lint/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' > $(ROOTDIR)/versions/GOLANGCILINT_VERSION

.PHONY: update/helm
## update helm version
update/helm:
	curl --silent https://api.github.com/repos/helm/helm/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' > $(ROOTDIR)/versions/HELM_VERSION

.PHONY: update/helm-operator
## update helm-operator version
update/helm-operator:
	curl --silent https://quay.io/api/v1/repository/operator-framework/helm-operator | jq -r '.tags'| grep name | grep -v master |grep -v latest | grep -v rc | head -1 | sed -e 's/.*\"name\":\ \"\(.*\)\",/\1/g' > $(ROOTDIR)/versions/OPERATOR_SDK_VERSION

.PHONY: update/helm-docs
## update helm-docs version
update/helm-docs:
	curl --silent https://api.github.com/repos/norwoodj/helm-docs/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' | sed 's/v//g' > $(ROOTDIR)/versions/HELM_DOCS_VERSION

.PHONY: update/protobuf
## update protobuf version
update/protobuf:
	curl --silent https://api.github.com/repos/protocolbuffers/protobuf/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' | sed 's/v//g' > $(ROOTDIR)/versions/PROTOBUF_VERSION

.PHONY: update/kind
## update kind (kubernetes in docker) version
update/kind:
	curl --silent https://api.github.com/repos/kubernetes-sigs/kind/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' | sed 's/v//g' > $(ROOTDIR)/versions/KIND_VERSION

.PHONY: update/kubectl
## update kubectl (kubernetes cli) version
update/kubectl:
	curl -L -s https://dl.k8s.io/release/stable.txt > $(ROOTDIR)/versions/KUBECTL_VERSION

.PHONY: update/prometheus-stack
## update prometheus version
update/prometheus-stack:
	curl --silent https://artifacthub.io/api/v1/packages/helm/prometheus-community/kube-prometheus-stack | jq .version | sed 's/"//g' > $(ROOTDIR)/versions/PROMETHEUS_STACK_VERSION

.PHONY: update/jaeger-operator
## update jaeger-operator version
update/jaeger-operator:
	curl --silent https://artifacthub.io/api/v1/packages/helm/jaegertracing/jaeger-operator | jq .version | sed 's/"//g' > $(ROOTDIR)/versions/JAEGER_OPERATOR_VERSION

.PHONY: update/kube-linter
## update kube-linter version
update/kube-linter:
	curl --silent https://api.github.com/repos/stackrox/kube-linter/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' > $(ROOTDIR)/versions/KUBELINTER_VERSION

# .PHONY: update/otel-operator
# ## update otel-operator version
# update/otel-operator:
# 	curl --silent https://api.github.com/repos/open-telemetry/opentelemetry-operator/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' > $(ROOTDIR)/versions/OTEL_OPERATOR_VERSION

.PHONY: update/ngt
## update yahoojapan/NGT version
update/ngt:
	curl --silent https://api.github.com/repos/yahoojapan/NGT/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' | sed 's/v//g' > $(ROOTDIR)/versions/NGT_VERSION

.PHONY: update/faiss
## update facebookresearch/faiss version
update/faiss:
	curl --silent https://api.github.com/repos/facebookresearch/faiss/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' | sed 's/v//g' > $(ROOTDIR)/versions/FAISS_VERSION

.PHONY: update/reviewdog
## update reviewdog version
update/reviewdog:
	curl --silent https://api.github.com/repos/reviewdog/reviewdog/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' > $(ROOTDIR)/versions/REVIEWDOG_VERSION

.PHONY: update/telepresence
## update telepresence version
update/telepresence:
	curl --silent https://api.github.com/repos/telepresenceio/telepresence/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' | sed 's/v//g' > $(ROOTDIR)/versions/TELEPRESENCE_VERSION

.PHONY: update/yq
## update YQ version
update/yq:
	curl --silent https://api.github.com/repos/mikefarah/yq/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' > $(ROOTDIR)/versions/YQ_VERSION

.PHONY: update/zlib
## update zlib version
update/zlib:
	curl --silent https://api.github.com/repos/madler/zlib/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' | sed 's/v//g' > $(ROOTDIR)/versions/ZLIB_VERSION

.PHONY: update/hdf5
## update hdf5 version
update/hdf5:
	curl --silent https://api.github.com/repos/HDFGroup/hdf5/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' | sed 's/v//g' > $(ROOTDIR)/versions/HDF5_VERSION

.PHONY: update/vald
## update vald it's self version
update/vald:
	curl --silent https://api.github.com/repos/vdaas/vald/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' > $(ROOTDIR)/versions/VALD_VERSION

.PHONY: update/valdcli
## update vald client library made by clojure self version
update/valdcli:
	curl --silent https://api.github.com/repos/vdaas/vald-client-clj/releases/latest | grep -Po '"tag_name": "\K.*?(?=")' > $(ROOTDIR)/versions/VALDCLI_VERSION

.PHONY: update/template
## update PULL_REQUEST_TEMPLATE and ISSUE_TEMPLATE
update/template:
	$(eval GO_VERSION      := $(shell $(MAKE) -s version/go))
	$(eval NGT_VERSION     := $(shell $(MAKE) -s version/ngt))
	$(eval KUBECTL_VERSION := $(shell $(MAKE) -s version/k8s))
	sed -i -e "s/^- Go Version: .*$$/- Go Version: $(GO_VERSION)/" $(ROOTDIR)/.github/ISSUE_TEMPLATE/bug_report.md
	sed -i -e "s/^- Go Version: .*$$/- Go Version: $(GO_VERSION)/" $(ROOTDIR)/.github/ISSUE_TEMPLATE/security_issue_report.md
	sed -i -e "s/^- Go Version: .*$$/- Go Version: $(GO_VERSION)/" $(ROOTDIR)/.github/PULL_REQUEST_TEMPLATE.md

	sed -i -e "s/^- NGT Version: .*$$/- NGT Version: $(NGT_VERSION)/" $(ROOTDIR)/.github/ISSUE_TEMPLATE/bug_report.md
	sed -i -e "s/^- NGT Version: .*$$/- NGT Version: $(NGT_VERSION)/" $(ROOTDIR)/.github/ISSUE_TEMPLATE/security_issue_report.md
	sed -i -e "s/^- NGT Version: .*$$/- NGT Version: $(NGT_VERSION)/" $(ROOTDIR)/.github/PULL_REQUEST_TEMPLATE.md

	sed -i -e "s/^- Kubernetes Version: .*$$/- Kubernetes Version: $(KUBECTL_VERSION)/" $(ROOTDIR)/.github/ISSUE_TEMPLATE/bug_report.md
	sed -i -e "s/^- Kubernetes Version: .*$$/- Kubernetes Version: $(KUBECTL_VERSION)/" $(ROOTDIR)/.github/ISSUE_TEMPLATE/security_issue_report.md
	sed -i -e "s/^- Kubernetes Version: .*$$/- Kubernetes Version: $(KUBECTL_VERSION)/" $(ROOTDIR)/.github/PULL_REQUEST_TEMPLATE.md
