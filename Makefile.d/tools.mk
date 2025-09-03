#
# Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

.PHONY: golangci-lint/install
## install golangci-lint
golangci-lint/install: $(BINDIR)/golangci-lint

$(BINDIR)/golangci-lint:
	curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $(BINDIR) $(GOLANGCILINT_VERSION)

.PHONY: goimports/install
goimports/install:
	@echo "goimports is now managed via go.mod tool directive. Use: goimports ..."



.PHONY: strictgoimports/install
strictgoimports/install:
	@echo "strictgoimports is now managed via go.mod tool directive. Use: strictgoimports ..."



.PHONY: gofumpt/install
gofumpt/install:
	@echo "gofumpt is now managed via go.mod tool directive. Use: gofumpt ..."



.PHONY: golines/install
golines/install:
	@echo "golines is now managed via go.mod tool directive. Use: golines ..."



.PHONY: crlfmt/install
crlfmt/install:
	@echo "crlfmt is now managed via go.mod tool directive. Use: crlfmt ..."



.PHONY: prettier/install
prettier/install: $(NPM_GLOBAL_PREFIX)/bin/prettier
$(NPM_GLOBAL_PREFIX)/bin/prettier:
	npm config -g set registry http://registry.npmjs.org/
	npm cache clean --force
	type prettier || npm install -g prettier

.PHONY: reviewdog/install
## install reviewdog
reviewdog/install: $(BINDIR)/reviewdog

$(BINDIR)/reviewdog:
	curl -fsSL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh \
		| sh -s -- -b $(BINDIR) $(REVIEWDOG_VERSION)

.PHONY: kubectl/install
kubectl/install: $(BINDIR)/kubectl

$(BINDIR)/kubectl:
	$(eval DARCH := $(subst aarch64,arm64,$(ARCH)))
	curl -fsSL "https://dl.k8s.io/release/$(KUBECTL_VERSION)/bin/$(OS)/$(subst x86_64,amd64,$(shell echo $(DARCH) | tr '[:upper:]' '[:lower:]'))/kubectl" -o $(BINDIR)/kubectl
	chmod a+x $(BINDIR)/kubectl

.PHONY: textlint/install
textlint/install: $(NPM_GLOBAL_PREFIX)/bin/textlint

$(NPM_GLOBAL_PREFIX)/bin/textlint:
	npm install -g textlint textlint-rule-en-spell textlint-rule-prh textlint-rule-write-good

.PHONY: textlint/ci/install
textlint/ci/install:
	npm init -y
	npm install --save-dev textlint textlint-rule-en-spell textlint-rule-prh textlint-rule-write-good

.PHONY: cspell/install
cspell/install: $(NPM_GLOBAL_PREFIX)/bin/cspell

$(NPM_GLOBAL_PREFIX)/bin/cspell:
	npm install -g cspell@latest \
		@cspell/dict-cpp \
		@cspell/dict-docker \
		@cspell/dict-en_us \
		@cspell/dict-fullstack \
		@cspell/dict-git \
		@cspell/dict-golang \
		@cspell/dict-k8s \
		@cspell/dict-makefile \
		@cspell/dict-markdown \
		@cspell/dict-npm \
		@cspell/dict-public-licenses \
		@cspell/dict-rust \
		@cspell/dict-shell
	cspell link add @cspell/dict-cpp
	cspell link add @cspell/dict-docker
	cspell link add @cspell/dict-en_us
	cspell link add @cspell/dict-fullstack
	cspell link add @cspell/dict-git
	cspell link add @cspell/dict-golang
	cspell link add @cspell/dict-k8s
	cspell link add @cspell/dict-makefile
	cspell link add @cspell/dict-markdown
	cspell link add @cspell/dict-npm
	cspell link add @cspell/dict-public-licenses
	cspell link add @cspell/dict-rust
	cspell link add @cspell/dict-shell

.PHONY: buf/install
buf/install: $(BINDIR)/buf

$(BINDIR)/buf:
	curl -fsSL \
	"https://github.com/bufbuild/buf/releases/download/$(BUF_VERSION)/buf-$(shell uname -s)-$(shell uname -m)" \
	-o "${BINDIR}/buf" && \
	chmod +x "${BINDIR}/buf"

.PHONY: k9s/install
k9s/install:
	@echo "k9s is now managed via go.mod tool directive. Use: k9s ..."



.PHONY: stern/install
stern/install:
	@echo "stern is now managed via go.mod tool directive. Use: stern ..."



.PHONY: yamlfmt/install
yamlfmt/install:
	@echo "yamlfmt is now managed via go.mod tool directive. Use: yamlfmt ..."



.PHONY: gopls/install
gopls/install:
	@echo "gopls is now managed via go.mod tool directive. Use: gopls ..."



.PHONY: gomodifytags/install
gomodifytags/install:
	@echo "gomodifytags is now managed via go.mod tool directive. Use: gomodifytags ..."



.PHONY: impl/install
impl/install:
	@echo "impl is now managed via go.mod tool directive. Use: impl ..."



.PHONY: goplay/install
goplay/install:
	@echo "goplay is now managed via go.mod tool directive. Use: goplay ..."



.PHONY: delve/install
delve/install:
	@echo "delve is now managed via go.mod tool directive. Use: dlv ..."



.PHONY: staticcheck/install
staticcheck/install:
	@echo "staticcheck is now managed via go.mod tool directive. Use: staticcheck ..."



.PHONY: go/install
go/install: $(GOROOT)/bin/go

$(GOROOT)/bin/go:
	TAR_NAME=go$(GO_VERSION).$(OS)-$(subst x86_64,amd64,$(subst aarch64,arm64,$(ARCH))).tar.gz \
	&& curl -fsSL "https://go.dev/dl/$${TAR_NAME}" -o "$(TEMP_DIR)/$${TAR_NAME}" \
	&& mkdir -p $(TEMP_DIR)/go \
	&& tar -xzvf "$(TEMP_DIR)/$${TAR_NAME}" -C $(TEMP_DIR)/go --strip-components 1 \
	&& rm -rf "$(TEMP_DIR)/$${TAR_NAME}" \
	&& mv $(TEMP_DIR)/go $(GOROOT) \
	&& $(GOROOT)/bin/go version

.PHONY: rust/install
rust/install: $(CARGO_HOME)/bin/cargo

$(CARGO_HOME)/bin/cargo:
	curl --proto '=https' --tlsv1.2 -fsSL https://sh.rustup.rs | CARGO_HOME=${CARGO_HOME} RUSTUP_HOME=${RUSTUP_HOME} sh -s -- --default-toolchain $(RUST_VERSION) -y
	rustup toolchain install $(RUST_VERSION)
	rustup default $(RUST_VERSION)
	source "${CARGO_HOME}/env"

.PHONY: zlib/install
zlib/install: $(LIB_PATH)/libz.a

$(LIB_PATH)/libz.a: $(LIB_PATH)
	curl -fsSL https://github.com/madler/zlib/releases/download/v$(ZLIB_VERSION)/zlib-$(ZLIB_VERSION).tar.gz -o $(TEMP_DIR)/zlib-$(ZLIB_VERSION).tar.gz \
	&& mkdir -p $(TEMP_DIR)/zlib \
	&& tar -xzvf $(TEMP_DIR)/zlib-$(ZLIB_VERSION).tar.gz -C $(TEMP_DIR)/zlib --strip-components 1 \
	&& cd $(TEMP_DIR)/zlib \
	&& mkdir -p build \
	&& cd build \
	&& cmake  -DCMAKE_BUILD_TYPE=Release \
		-DCMAKE_POLICY_VERSION_MINIMUM=$(CMAKE_VERSION) \
		-DBUILD_SHARED_LIBS=OFF \
		-DBUILD_STATIC_EXECS=ON \
		-DBUILD_TESTING=OFF \
		-DZLIB_BUILD_SHARED=OFF \
		-DZLIB_BUILD_STATIC=ON \
		-DZLIB_COMPAT=ON \
		-DZLIB_USE_STATIC_LIBS=ON \
		-DCMAKE_CXX_FLAGS="$(CXXFLAGS)" \
		-DCMAKE_C_FLAGS="$(CFLAGS)" \
		-DCMAKE_INSTALL_LIBDIR=$(LIB_PATH) \
		-DCMAKE_INSTALL_PREFIX=$(USR_LOCAL) \
		-B $(TEMP_DIR)/zlib/build $(TEMP_DIR)/zlib \
	&& make -j$(CORES) \
	&& make install \
	&& cd $(ROOTDIR) \
	&& rm -rf $(TEMP_DIR)/zlib-$(ZLIB_VERSION).tar.gz $(TEMP_DIR)/zlib $(LIB_PATH)/libz.s*

.PHONY: hdf5/install
hdf5/install: $(LIB_PATH)/libhdf5.a

$(LIB_PATH)/libhdf5.a: $(LIB_PATH) \
	zlib/install
	mkdir -p $(TEMP_DIR)/hdf5 \
	&& curl -fsSL https://github.com/HDFGroup/hdf5/archive/refs/tags/$(HDF5_VERSION).tar.gz -o $(TEMP_DIR)/hdf5.tar.gz \
	&& tar -xzvf $(TEMP_DIR)/hdf5.tar.gz -C $(TEMP_DIR)/hdf5 --strip-components 1 \
	&& mkdir -p $(TEMP_DIR)/hdf5/build \
	&& cd $(TEMP_DIR)/hdf5/build \
	&& cmake -DCMAKE_BUILD_TYPE=Release \
		-DCMAKE_POLICY_VERSION_MINIMUM=$(CMAKE_VERSION) \
		-DBUILD_SHARED_LIBS=OFF \
		-DBUILD_STATIC_EXECS=ON \
		-DBUILD_TESTING=OFF \
		-DHDF5_BUILD_CPP_LIB=OFF \
		-DHDF5_BUILD_HL_LIB=ON \
		-DHDF5_BUILD_STATIC_EXECS=ON \
		-DHDF5_BUILD_TOOLS=OFF \
		-DHDF5_ENABLE_Z_LIB_SUPPORT=ON \
		-DH5_ZLIB_INCLUDE_DIR=$(USR_LOCAL)/include \
		-DH5_ZLIB_LIBRARY=$(LIB_PATH)/libz.a \
		-DCMAKE_CXX_FLAGS="$(CXXFLAGS)" \
		-DCMAKE_C_FLAGS="$(CFLAGS)" \
		-DCMAKE_INSTALL_LIBDIR=$(LIB_PATH) \
		-DCMAKE_INSTALL_PREFIX=$(USR_LOCAL) \
		-B $(TEMP_DIR)/hdf5/build $(TEMP_DIR)/hdf5 \
	&& make -j$(CORES) \
	&& make install \
	&& cd $(ROOTDIR) \
	&& rm -rf $(TEMP_DIR)/hdf5.tar.gz $(TEMP_DIR)/hdf5 \
	&& ldconfig

.PHONY: yq/install
## install yq
yq/install: $(BINDIR)/yq

$(BINDIR)/yq:
	mkdir -p $(BINDIR)
	$(eval DARCH := $(subst aarch64,arm64,$(ARCH)))
	cd $(TEMP_DIR) \
	    && curl -fsSL https://github.com/mikefarah/yq/releases/download/$(YQ_VERSION)/yq_$(OS)_$(subst x86_64,amd64,$(shell echo $(DARCH) | tr '[:upper:]' '[:lower:]')) -o $(BINDIR)/yq \
	    && chmod a+x $(BINDIR)/yq

.PHONY: docker-cli/install
docker-cli/install: $(BINDIR)/docker

$(BINDIR)/docker: $(BINDIR)
	curl -fsSL https://download.docker.com/linux/static/stable/$(shell uname -m)/docker-$(shell echo $(DOCKER_VERSION) | cut -c2-).tgz -o $(TEMP_DIR)/docker.tgz \
		&& tar -xzvf $(TEMP_DIR)/docker.tgz -C $(TEMP_DIR) \
		&& mv $(TEMP_DIR)/docker/docker $(BINDIR) \
		&& rm -rf $(TEMP_DIR)/docker{.tgz,}
