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

.PHONY: golangci-lint/install
## install golangci-lint
golangci-lint/install: $(BINDIR)/golangci-lint

$(BINDIR)/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $(BINDIR) $(GOLANGCILINT_VERSION)

.PHONY: goimports/install
goimports/install: $(GOPATH)/bin/goimports

$(GOPATH)/bin/goimports:
	$(call go-install, golang.org/x/tools/cmd/goimports)

.PHONY: strictgoimports/install
strictgoimports/install: $(GOPATH)/bin/strictgoimports

$(GOPATH)/bin/strictgoimports:
	$(call go-install, github.com/momotaro98/strictgoimports/cmd/strictgoimports)

.PHONY: gofumpt/install
gofumpt/install: $(GOPATH)/bin/gofumpt

$(GOPATH)/bin/gofumpt:
	$(call go-install, mvdan.cc/gofumpt)

.PHONY: golines/install
golines/install: $(GOPATH)/bin/golines

$(GOPATH)/bin/golines:
	$(call go-install, github.com/segmentio/golines)

.PHONY: prettier/install
prettier/install: $(BINDIR)/prettier
$(BINDIR)/prettier:
	type prettier || npm install -g prettier

.PHONY: reviewdog/install
## install reviewdog
reviewdog/install: $(BINDIR)/reviewdog

$(BINDIR)/reviewdog:
	curl -sSfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh \
		| sh -s -- -b $(BINDIR) $(REVIEWDOG_VERSION)

.PHONY: kubectl/install
kubectl/install: $(BINDIR)/kubectl

$(BINDIR)/kubectl:
	curl -L "https://dl.k8s.io/release/$(KUBECTL_VERSION)/bin/$(shell echo $(UNAME) | tr '[:upper:]' '[:lower:]')/$(subst x86_64,amd64,$(shell echo $(ARCH) | tr '[:upper:]' '[:lower:]'))/kubectl" -o $(BINDIR)/kubectl
	chmod a+x $(BINDIR)/kubectl

.PHONY: protobuf/install
protobuf/install: $(BINDIR)/protoc

$(BINDIR)/protoc:
	curl -L "https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOBUF_VERSION)/protoc-$(PROTOBUF_VERSION)-$(subst darwin,osx,$(shell echo $(UNAME) | tr '[:upper:]' '[:lower:]'))-$(ARCH).zip" -o /tmp/protoc.zip
	unzip -o /tmp/protoc.zip -d /usr/local bin/protoc
	unzip -o /tmp/protoc.zip -d /usr/local 'include/*'
	rm -f /tmp/protoc.zip

.PHONY: textlint/install
textlint/install:
	npm install -g textlint textlint-rule-en-spell textlint-rule-prh textlint-rule-write-good

.PHONY: textlint/ci/install
textlint/ci/install:
	npm init -y
	npm install --save-dev textlint textlint-rule-en-spell textlint-rule-prh textlint-rule-write-good

.PHONY: buf/install
buf/install: $(BINDIR)/buf

$(BINDIR)/buf:
	curl -sSL \
	"https://github.com/bufbuild/buf/releases/download/$(BUF_VERSION)/buf-$(shell uname -s)-$(shell uname -m)" \
	-o "${BINDIR}/buf" && \
	chmod +x "${BINDIR}/buf"
