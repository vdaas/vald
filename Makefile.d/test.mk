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

.PHONY: gotests/install
gotests/install:
	go get -u github.com/cweill/gotests/...

.PHONY: gotests/gen
## generate missing go test files
gotests/gen: \
	$(GO_TEST_SOURCES) \
	$(GO_OPTION_TEST_SOURCES) \
	gotests/patch

.PHONY: gotests/patch
## apply patches to generated go test files
gotests/patch: \
	$(GO_TEST_SOURCES) \
	$(GO_OPTION_TEST_SOURCES)
	@$(call green, "apply patches to go test files...")
	find $(ROOTDIR)/internal/k8s/* -name '*_test.go' | xargs sed -i -E "s%k8s.io/apimachinery/pkg/api/errors%github.com/vdaas/vald/internal/errors%g"
	find $(ROOTDIR)/* -name '*_test.go' | xargs sed -i -E "s%cockroachdb/errors%vdaas/vald/internal/errors%g"
	find $(ROOTDIR)/* -name '*_test.go' | xargs sed -i -E "s%pkg/errors%vdaas/vald/internal/errors%g"
	find $(ROOTDIR)/* -name '*_test.go' | xargs sed -i -E "s%go-errors/errors%vdaas/vald/internal/errors%g"
	find $(ROOTDIR)/internal/errors -name '*_test.go' | xargs sed -i -E "s%\"github.com/vdaas/vald/internal/errors\"%%g"
	find $(ROOTDIR)/internal/errors -name '*_test.go' | xargs sed -i -E "s/errors\.//g"

# force to rebuild all GO_TEST_SOURCES targets
.PHONY: $(GO_TEST_SOURCES)
$(GO_TEST_SOURCES): \
	./assets/test/templates/common
	@$(call green, $(patsubst %,"generating go test file: %",$@))
	gotests -w -template_dir ./assets/test/templates/common -all $(patsubst %_test.go,%.go,$@)

# force to rebuild all GO_OPTION_TEST_SOURCES targets
.PHONY: $(GO_OPTION_TEST_SOURCES)
$(GO_OPTION_TEST_SOURCES): \
	./assets/test/templates/option
	@$(call green, $(patsubst %,"generating go test file: %",$@))
	gotests -w -template_dir ./assets/test/templates/option -all $(patsubst %_test.go,%.go,$@)
