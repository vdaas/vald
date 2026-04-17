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

.PHONY: version
## print vald version
version: \
	version/vald

.PHONY: version/vald
## print vald version
version/vald:
	@echo $(VERSION)

.PHONY: version/go
## print go version
version/go:
	@echo $(GO_VERSION)

.PHONY: version/rust
## print rust version
version/rust:
	@echo $(RUST_VERSION)

.PHONY: version/ngt
## print NGT version
version/ngt:
	@echo $(NGT_VERSION)

.PHONY: version/faiss
## print Faiss version
version/faiss:
	@echo $(FAISS_VERSION)

.PHONY: version/usearch
## print usearch version
version/usearch:
	@echo $(USEARCH_VERSION)

.PHONY: version/docker
## print Kubernetes version
version/docker:
	@echo $(DOCKER_VERSION)

.PHONY: version/k8s
## print Kubernetes version
version/k8s:
	@echo $(KUBECTL_VERSION)

.PHONY: version/kind
## print kind version
version/kind:
	@echo $(KIND_VERSION)

.PHONY: version/helm
## print helm version
version/helm:
	@echo $(HELM_VERSION)

.PHONY: version/llvm
## print llvm version
version/llvm:
	@echo $(LLVM_VERSION)

.PHONY: version/openmp
## print openmp version
version/openmp:
	@echo $(LLVM_VERSION)

.PHONY: version/yq
## print yq version
version/yq:
	@echo $(YQ_VERSION)

.PHONY: version/telepresence
## print telepresence version
version/telepresence:
	@echo $(TELEPRESENCE_VERSION)