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

.PHONY: license
## add license to files
license:
	$(call gen-license,$(ROOTDIR),$(MAINTAINER))

.PHONY: dockerfile
## generate dockerfiles
dockerfile:
	$(call gen-dockerfile,$(ROOTDIR),$(MAINTAINER))

.PHONY: dashboard
## generate dashboards
dashboard: k8s/metrics/grafana/dashboards/00-vald-cluster-overview.yaml

# To cache the generated dashboards, making a generated file target
k8s/metrics/grafana/dashboards/00-vald-cluster-overview.yaml: $(shell cat $(ROOTDIR)/.gitfiles | grep '^hack/grafana/gen/' | sed -e 's%^%$(ROOTDIR)/%') versions/GRAFANA_VERSION
	$(call gen-dashboard,$(ROOTDIR),$(MAINTAINER))

.PHONY: workflow
## generate workflows
workflow:
	$(call gen-dockerfile,$(ROOTDIR),$(MAINTAINER))

.PHONY: deadlink-checker
## generate deadlink-checker
deadlink-checker:
	$(call gen-deadlink-checker,$(ROOTDIR),$(MAINTAINER),$(DEADLINK_CHECK_PATH),$(DEADLINK_IGNORE_PATH),$(DEADLINK_CHECK_FORMAT))

.PHONY: changelog/update
## update changelog
changelog/update:
	echo "# CHANGELOG" > $(TEMP_DIR)/CHANGELOG.md
	echo "" >> $(TEMP_DIR)/CHANGELOG.md
	$(MAKE) -s changelog/next/print >> $(TEMP_DIR)/CHANGELOG.md
	echo "" >> $(TEMP_DIR)/CHANGELOG.md
	tail -n +2 $(ROOTDIR)/CHANGELOG.md >> $(TEMP_DIR)/CHANGELOG.md
	mv -f $(TEMP_DIR)/CHANGELOG.md $(ROOTDIR)/CHANGELOG.md

.PHONY: changelog/next/print
## print next changelog entry
changelog/next/print:
	@cat $(ROOTDIR)/hack/CHANGELOG.template.md | \
	sed -e 's/{{ version }}/$(VERSION)/g'
	@echo "$$BODY"