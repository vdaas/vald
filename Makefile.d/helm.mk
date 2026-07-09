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

.PHONY: helm/install
## install helm
helm/install: $(BINDIR)/helm

$(BINDIR)/helm:
	mkdir -p $(BINDIR)
	$(eval DARCH := $(subst aarch64,arm64,$(ARCH)))
	TAR_NAME=helm-$(HELM_VERSION)-$(OS)-$(subst x86_64,amd64,$(shell echo $(DARCH) | tr '[:upper:]' '[:lower:]')) \
	&& cd $(TEMP_DIR) \
	&& curl -fsSL --retry 3 "https://get.helm.sh/$${TAR_NAME}.tar.gz" -o "$(TEMP_DIR)/$${TAR_NAME}" \
	&& tar -xzvf "$(TEMP_DIR)/$${TAR_NAME}" --strip=1 \
	&& mv helm $(BINDIR)/helm

.PHONY: helm-docs/install
## install helm-docs
helm-docs/install: $(BINDIR)/helm-docs

$(BINDIR)/helm-docs:
	mkdir -p $(BINDIR)
	$(eval DARCH := $(subst aarch64,arm64,$(ARCH)))
	TAR_NAME=helm-docs_$(HELM_DOCS_VERSION)_$(UNAME)_$(DARCH).tar.gz \
	&& cd $(TEMP_DIR) \
	&& curl -fsSL "https://github.com/norwoodj/helm-docs/releases/download/v$(HELM_DOCS_VERSION)/$${TAR_NAME}" -o "$(TEMP_DIR)/$${TAR_NAME}" \
	&& tar -xzvf "$(TEMP_DIR)/$${TAR_NAME}" \
	&& mv helm-docs $(BINDIR)/helm-docs

.PHONY: helm/package/vald
## packaging Helm chart for Vald
helm/package/vald:
	helm package $(ROOTDIR)/charts/vald

.PHONY: helm/package/operator/helm
## packaging Helm chart for vald-helm-operator
helm/package/operator/helm: \
	helm/schema/crd/vald \
	helm/schema/crd/operator/helm
	helm package $(ROOTDIR)/charts/operator/helm

.PHONY: helm/package/operator/vald
## packaging Helm chart for vald-operator
helm/package/operator/vald:
	helm package $(ROOTDIR)/charts/operator/vald

.PHONY: helm/package/operator/benchmark
## packaging Helm chart for vald-helm-operator
helm/package/operator/benchmark: \
	helm/schema/crd/vald-benchmark-job \
	helm/schema/crd/vald-benchmark-scenario \
	helm/schema/crd/operator/benchmark
	helm package $(ROOTDIR)/charts/operator/benchmark

.PHONY: helm/package/vald-readreplica
## packaging Helm chart for vald-readreplica
helm/package/vald-readreplica:
	helm package $(ROOTDIR)/charts/vald-readreplica

.PHONY: helm/repo/add
## add Helm chart repository
helm/repo/add:
	helm repo add vald https://vald.vdaas.org/charts

.PHONY: helm/docs/vald
helm/docs/vald: $(ROOTDIR)/charts/vald/README.md

# force to rebuild
.PHONY: $(ROOTDIR)/charts/vald/README.md
$(ROOTDIR)/charts/vald/README.md: \
	$(ROOTDIR)/charts/vald/README.md.gotmpl \
	$(ROOTDIR)/charts/vald/values.yaml
	helm-docs

.PHONY: helm/docs/operator/helm
helm/docs/operator/helm: $(ROOTDIR)/charts/operator/helm/README.md

# force to rebuild
.PHONY: $(ROOTDIR)/charts/operator/helm/README.md
$(ROOTDIR)/charts/operator/helm/README.md: \
	$(ROOTDIR)/charts/operator/helm/README.md.gotmpl \
	$(ROOTDIR)/charts/operator/helm/values.yaml
	helm-docs

.PHONY: helm/docs/operator/vald
helm/docs/operator/vald: $(ROOTDIR)/charts/operator/vald/README.md

# force to rebuild
.PHONY: $(ROOTDIR)/charts/operator/vald/README.md
$(ROOTDIR)/charts/operator/vald/README.md: \
	$(ROOTDIR)/charts/operator/vald/README.md.gotmpl \
	$(ROOTDIR)/charts/operator/vald/values.yaml
	helm-docs

.PHONY: helm/docs/vald-readreplica
helm/docs/vald-readreplica: $(ROOTDIR)/charts/vald-readreplica/README.md

.PHONY: helm/docs/operator/benchmark
helm/docs/operator/benchmark: $(ROOTDIR)/charts/operator/benchmark/README.md

.PHONY: $(ROOTDIR)/charts/operator/benchmark/README.md
$(ROOTDIR)/charts/operator/benchmark/README.md: \
	$(ROOTDIR)/charts/operator/benchmark/README.md.gotmpl \
	$(ROOTDIR)/charts/operator/benchmark/values.yaml
	helm-docs

# force to rebuild
.PHONY: $(ROOTDIR)/charts/vald-readreplica/README.md
$(ROOTDIR)/charts/vald-readreplica/README.md: \
	$(ROOTDIR)/charts/vald-readreplica/README.md.gotmpl \
	$(ROOTDIR)/charts/vald-readreplica/values.yaml
	helm-docs

.PHONY: helm/schema/all
helm/schema/all: \
	helm/schema/vald \
	helm/schema/operator/helm \
	helm/schema/operator/vald \
	helm/schema/vald-benchmark-job \
	helm/schema/vald-benchmark-scenario \
	helm/schema/operator/benchmark

.PHONY: helm/schema/vald
## generate json schema for Vald Helm Chart
helm/schema/vald: $(ROOTDIR)/charts/vald/values.schema.json

$(ROOTDIR)/charts/vald/values.schema.json: \
	$(ROOTDIR)/charts/vald/values.yaml
	$(call gen-vald-helm-schema,vald/values)

.PHONY: helm/schema/operator/helm
## generate json schema for Vald Helm Operator Chart
helm/schema/operator/helm: $(ROOTDIR)/charts/operator/helm/values.schema.json

$(ROOTDIR)/charts/operator/helm/values.schema.json: \
	$(ROOTDIR)/charts/operator/helm/values.yaml
	$(call gen-vald-helm-schema,operator/helm/values)

.PHONY: helm/schema/operator/vald
## generate json schema for Vald Operator Chart
helm/schema/operator/vald: $(ROOTDIR)/charts/operator/vald/values.schema.json

$(ROOTDIR)/charts/operator/vald/values.schema.json: \
	$(ROOTDIR)/charts/operator/vald/values.yaml
	$(call gen-vald-helm-schema,operator/vald/values)

.PHONY: helm/schema/vald-benchmark-job
## generate json schema for Vald Benchmark Job Chart
helm/schema/vald-benchmark-job: $(ROOTDIR)/charts/operator/benchmark/job-values.schema.json

$(ROOTDIR)/charts/operator/benchmark/job-values.schema.json: \
	$(ROOTDIR)/charts/operator/benchmark/schemas/job-values.yaml
	$(call gen-vald-helm-schema,operator/benchmark/schemas/job-values)

.PHONY: helm/schema/vald-benchmark-scenario
## generate json schema for Vald Benchmark Job Chart
helm/schema/vald-benchmark-scenario: $(ROOTDIR)/charts/operator/benchmark/scenario-values.schema.json

$(ROOTDIR)/charts/operator/benchmark/scenario-values.schema.json: \
	$(ROOTDIR)/charts/operator/benchmark/schemas/scenario-values.yaml
	$(call gen-vald-helm-schema,operator/benchmark/schemas/scenario-values)

.PHONY: helm/schema/operator/benchmark
## generate json schema for Vald Benchmark Operator Chart
helm/schema/operator/benchmark: $(ROOTDIR)/charts/operator/benchmark/values.schema.json

$(ROOTDIR)/charts/operator/benchmark/values.schema.json: \
	$(ROOTDIR)/charts/operator/benchmark/values.yaml
	$(call gen-vald-helm-schema,operator/benchmark/values)

.PHONY: helm/schema/crd/all
helm/schema/crd/all: \
	helm/schema/crd/vald \
	helm/schema/crd/operator/helm \
	helm/schema/crd/vald/mirror-target \
	helm/schema/crd/vald-benchmark-job \
	helm/schema/crd/vald-benchmark-scenario \
	helm/schema/crd/operator/benchmark

.PHONY: helm/schema/crd/vald
## generate OpenAPI v3 schema for ValdRelease
helm/schema/crd/vald: \
	yq/install
	$(call gen-vald-crd,operator/helm,valdrelease,vald/values)

.PHONY: helm/schema/crd/operator/helm
## generate OpenAPI v3 schema for ValdHelmOperatorRelease
helm/schema/crd/operator/helm: \
	yq/install
	$(call gen-vald-crd,operator/helm,valdhelmoperatorrelease,operator/helm/values)

.PHONY: helm/schema/crd/vald/mirror-target
## generate OpenAPI v3 schema for ValdMirrorTarget
helm/schema/crd/vald/mirror-target: \
	yq/install
	$(call gen-vald-crd,vald,valdmirrortarget,vald/schemas/mirror-target-values)

.PHONY: helm/schema/crd/vald-benchmark-job
## generate OpenAPI v3 schema for ValdBenchmarkJobRelease
helm/schema/crd/vald-benchmark-job: \
	yq/install
	$(call gen-vald-crd,operator/benchmark,valdbenchmarkjob,operator/benchmark/schemas/job-values)

.PHONY: helm/schema/crd/vald-benchmark-scenario
## generate OpenAPI v3 schema for ValdBenchmarkScenarioRelease
helm/schema/crd/vald-benchmark-scenario: \
	yq/install
	$(call gen-vald-crd,operator/benchmark,valdbenchmarkscenario,operator/benchmark/schemas/scenario-values)

.PHONY: helm/schema/crd/operator/benchmark
## generate OpenAPI v3 schema for ValdBenchmarkOperatorRelease
helm/schema/crd/operator/benchmark: \
	yq/install
	$(call gen-vald-crd,operator/benchmark,valdbenchmarkoperatorrelease,operator/benchmark/values)