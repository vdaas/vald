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

.PHONY: golangci-lint/install
## install golangci-lint
golangci-lint/install: $(BINDIR)/golangci-lint

$(BINDIR)/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $(BINDIR) $(GOLANGCILINT_VERSION)

.PHONY: reviewdog/install
## install reviewdog
reviewdog/install: $(BINDIR)/reviewdog

$(BINDIR)/reviewdog:
	curl -sSfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh \
		| sh -s -- -b $(BINDIR) $(REVIEWDOG_VERSION)

.PHONY: kubectl/install
kubectl/install: $(BINDIR)/kubectl

ifeq ($(UNAME),Darwin)
$(BINDIR)/kubectl:
	curl -L "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/darwin/amd64/kubectl" -o $(BINDIR)/kubectl
	chmod a+x $(BINDIR)/kubectl
else
$(BINDIR)/kubectl:
	curl -L "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" -o $(BINDIR)/kubectl
	chmod a+x $(BINDIR)/kubectl
endif

.PHONY: protobuf/install
protobuf/install: /usr/local/bin/protoc

ifeq ($(UNAME),Darwin)
/usr/local/bin/protoc:
	curl -L "https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOBUF_VERSION)/protoc-$(PROTOBUF_VERSION)-osx-x86_64.zip" -o /tmp/protoc.zip
	sudo unzip -o /tmp/protoc.zip -d /usr/local bin/protoc
	sudo unzip -o /tmp/protoc.zip -d /usr/local 'include/*'
	rm -f /tmp/protoc.zip
else
/usr/local/bin/protoc:
	curl -L "https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOBUF_VERSION)/protoc-$(PROTOBUF_VERSION)-linux-x86_64.zip" -o /tmp/protoc.zip
	unzip -o /tmp/protoc.zip -d /usr/local bin/protoc
	unzip -o /tmp/protoc.zip -d /usr/local 'include/*'
	rm -f /tmp/protoc.zip
endif
