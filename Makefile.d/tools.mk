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

.PHONY: golangci-lint/install
## install golangci-lint
golangci-lint/install: \
	$(BINDIR)/golangci-lint

$(BINDIR)/golangci-lint:
	curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
	| sh -s -- -b $(BINDIR) $(GOLANGCILINT_VERSION)

.PHONY: goimports/install
## install goimports
goimports/install: \
	$(GOBIN)/goimports

$(GOBIN)/goimports:
	$(call go-tool-install)

.PHONY: strictgoimports/install
## install strictgoimports
strictgoimports/install: \
	$(GOBIN)/strictgoimports

$(GOBIN)/strictgoimports:
	$(call go-tool-install)

.PHONY: gofumpt/install
## install gofumpt
gofumpt/install: \
	$(GOBIN)/gofumpt

$(GOBIN)/gofumpt:
	$(call go-tool-install)

.PHONY: golines/install
## install golines
golines/install: \
	$(GOBIN)/golines

$(GOBIN)/golines:
	$(call go-tool-install)

.PHONY: crlfmt/install
## install crlfmt
crlfmt/install: \
	$(GOBIN)/crlfmt

$(GOBIN)/crlfmt:
	$(call go-tool-install)

.PHONY: actionlint/install
## install actionlint
actionlint/install: \
	$(GOBIN)/actionlint

$(GOBIN)/actionlint:
	$(call go-tool-install)

.PHONY: ghalint/install
## install ghalint
ghalint/install: \
	$(GOBIN)/ghalint

$(GOBIN)/ghalint:
	$(call go-tool-install)

.PHONY: pinact/install
## install pinact
pinact/install: \
	$(GOBIN)/pinact

$(GOBIN)/pinact:
	$(call go-tool-install)

.PHONY: ghatm/install
## install ghatm
ghatm/install: \
	$(GOBIN)/ghatm

$(GOBIN)/ghatm:
	$(call go-tool-install)

.PHONY: buf/install
## install buf
buf/install: \
	$(GOBIN)/buf

$(GOBIN)/buf:
	$(call go-tool-install)

.PHONY: k9s/install
## install k9s
k9s/install: \
	$(GOBIN)/k9s

$(GOBIN)/k9s:
	$(call go-tool-install)

.PHONY: stern/install
## install stern
stern/install: \
	$(GOBIN)/stern

$(GOBIN)/stern:
	$(call go-tool-install)

.PHONY: yamlfmt/install
## install yamlfmt
yamlfmt/install: \
	$(GOBIN)/yamlfmt

$(GOBIN)/yamlfmt:
	$(call go-tool-install)

.PHONY: gomodifytags/install
## install gomodifytags
gomodifytags/install: \
	$(GOBIN)/gomodifytags

$(GOBIN)/gomodifytags:
	$(call go-tool-install)

.PHONY: impl/install
## install impl
impl/install: \
	$(GOBIN)/impl

$(GOBIN)/impl:
	$(call go-tool-install)

.PHONY: delve/install
## install delve
delve/install: \
	$(GOBIN)/dlv

$(GOBIN)/dlv:
	$(call go-tool-install)

.PHONY: staticcheck/install
## install staticcheck
staticcheck/install: \
	$(GOBIN)/staticcheck

$(GOBIN)/staticcheck:
	$(call go-tool-install)

.PHONY: tparse/install
## install tparse
tparse/install: \
	$(GOBIN)/tparse

$(GOBIN)/tparse:
	$(call go-tool-install)

.PHONY: gotestfmt/install
## install gotestfmt
gotestfmt/install: \
	$(GOBIN)/gotestfmt

$(GOBIN)/gotestfmt:
	$(call go-tool-install)

.PHONY: gotests/install
## install gotests
gotests/install: \
	$(GOBIN)/gotests

$(GOBIN)/gotests:
	$(call go-tool-install)

.PHONY: protoc-gen-doc/install
## install protoc-gen-doc
protoc-gen-doc/install: \
	$(GOBIN)/protoc-gen-doc

$(GOBIN)/protoc-gen-doc:
	$(call go-tool-install)

.PHONY: go/tools/install
## install go tools
go/tools/install:
	$(call go-tool-install)

.PHONY: gopls/install
## install gopls
gopls/install: \
	$(GOBIN)/gopls

$(GOBIN)/gopls:
	GO111MODULE=on go install -mod=readonly golang.org/x/tools/gopls@latest

.PHONY: prettier/install
## Install prettier via Bun (global)
prettier/install: $(BUN_GLOBAL_BIN)/prettier
$(BUN_GLOBAL_BIN)/prettier: bun/install
	command -v prettier >/dev/null 2>&1 || bun add --global prettier

.PHONY: reviewdog/install
## install reviewdog
reviewdog/install: $(BINDIR)/reviewdog

$(BINDIR)/reviewdog:
	curl -fsSL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh \
	| sh -s -- -b $(BINDIR) $(REVIEWDOG_VERSION)

.PHONY: mbake/install
## install mbake
mbake/install: $(BINDIR)/mbake

$(BINDIR)/mbake:
	pip install mbake --break-system-packages --prefix /usr

.PHONY: kubectl/install
## install kubectl
kubectl/install: $(BINDIR)/kubectl

$(BINDIR)/kubectl:
	$(eval DARCH := $(subst aarch64,arm64,$(ARCH)))
	curl -fsSL "https://dl.k8s.io/release/$(KUBECTL_VERSION)/bin/$(OS)/$(subst x86_64,amd64,$(shell echo $(DARCH) | tr '[:upper:]' '[:lower:]'))/kubectl" -o $(BINDIR)/kubectl
	chmod a+x $(BINDIR)/kubectl

.PHONY: textlint/install
## Install textlint & rules via Bun (global)
textlint/install: $(BUN_GLOBAL_BIN)/textlint
$(BUN_GLOBAL_BIN)/textlint: bun/install
	bun add --global \
	textlint \
	textlint-rule-en-spell \
	textlint-rule-prh \
	textlint-rule-write-good

.PHONY: textlint/ci/install
## Install textlint & rules for CI via Bun (local devDependencies)
textlint/ci/install: bun/install
	[ -f package.json ] || (bun init -y >/dev/null 2>&1 || echo '{}' > package.json)
	bun add --dev \
	textlint \
	textlint-rule-en-spell \
	textlint-rule-prh \
	textlint-rule-write-good

.PHONY: cspell/install
## Install cspell & dictionaries via Bun (global)
cspell/install: $(BUN_GLOBAL_BIN)/cspell
$(BUN_GLOBAL_BIN)/cspell: bun/install
	bun add --global \
	cspell@latest \
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

.PHONY: bun/install
## Install Bun runtime into $(BUN_INSTALL) if not already installed
bun/install: $(BINDIR)/bun

$(BINDIR)/bun:
	curl -fsSL https://bun.sh/install | BUN_INSTALL=$(BUN_INSTALL) bash

.PHONY: go/install
## install go
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
## install rust
rust/install: $(CARGO_HOME)/bin/cargo

$(CARGO_HOME)/bin/cargo:
	curl --proto '=https' --tlsv1.2 -fsSL https://sh.rustup.rs | CARGO_HOME=${CARGO_HOME} RUSTUP_HOME=${RUSTUP_HOME} sh -s -- --default-toolchain $(RUST_VERSION) -y
	rustup toolchain install $(RUST_VERSION)
	rustup default $(RUST_VERSION)
	source "${CARGO_HOME}/env"

.PHONY: rustfmt/install
## install rustfmt
$(CARGO_HOME)/bin/rustfmt: $(CARGO_HOME)/bin/cargo
	CARGO_HOME=${CARGO_HOME} RUSTUP_HOME=${RUSTUP_HOME} \
	$(CARGO_HOME)/bin/rustup component add rustfmt

rustfmt/install:
	$(MAKE) rust/install
	$(MAKE) $(CARGO_HOME)/bin/rustfmt

.PHONY: zlib/install
## install zlib
zlib/install: $(LIB_PATH)/libz.a

$(LIB_PATH)/libz.a: $(LIB_PATH)
	curl -fsSL https://github.com/madler/zlib/releases/download/v$(ZLIB_VERSION)/zlib-$(ZLIB_VERSION).tar.gz -o $(TEMP_DIR)/zlib-$(ZLIB_VERSION).tar.gz \
	&& mkdir -p $(TEMP_DIR)/zlib \
	&& tar -xzvf $(TEMP_DIR)/zlib-$(ZLIB_VERSION).tar.gz -C $(TEMP_DIR)/zlib --strip-components 1 \
	&& cd $(TEMP_DIR)/zlib \
	&& mkdir -p build \
	&& cd build \
	&& cmake	-DCMAKE_BUILD_TYPE=Release \
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
## install hdf5
hdf5/install: $(LIB_PATH)/libhdf5.a

<<<<<<< HEAD
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
=======
$(LIB_PATH)/libhdf5.a: | $(LIB_PATH) zlib/install
	$(call cmake-install,https://github.com/HDFGroup/hdf5/archive/refs/tags/$(HDF5_VERSION).tar.gz,hdf5, \
		-DHDF5_BUILD_CPP_LIB=OFF \
		-DHDF5_BUILD_HL_LIB=ON \
		-DHDF5_BUILD_STATIC_EXECS=OFF \
		-DHDF5_BUILD_TOOLS=OFF \
		-DHDF5_BUILD_EXAMPLES=OFF \
		-DHDF5_ENABLE_Z_LIB_SUPPORT=ON \
		-DH5_ZLIB_INCLUDE_DIR=$(USR_LOCAL)/include \
		-DH5_ZLIB_LIBRARY=$(LIB_PATH)/libz.a \
		-DCMAKE_INTERPROCEDURAL_OPTIMIZATION=OFF \
		-DCMAKE_EXE_LINKER_FLAGS="" \
		-DCMAKE_SHARED_LINKER_FLAGS="" \
		-DCMAKE_MODULE_LINKER_FLAGS="" \
		-DCMAKE_C_FLAGS="$(CFLAGS) -fno-lto" \
		-DCMAKE_CXX_FLAGS="$(CXXFLAGS) -fno-lto", \
		$(SUDO) rm -f $(USR_LOCAL)/include/H5*.h $(USR_LOCAL)/include/hdf5*.h $(LIB_PATH)/libhdf5*)

.PHONY: libomp/install
## install libomp static library from LLVM openmp source
libomp/install: $(LIB_PATH)/libomp.a

$(LIB_PATH)/libomp.a: | ninja/install
	@$(call green, "Installing libomp...")
	$(SUDO) mkdir -p $(LIB_PATH) $(INCLUDE_PATH)
	# Fast path: reuse a static libomp.a already shipped by the system LLVM/OpenMP
	# packages when one exists. Debian/Ubuntu's `llvm` package (clangLTOBuildDeps)
	# ships only libomp.so, so this branch typically misses there and the build
	# falls through to the from-source compile below; it hits on distros/LLVM
	# installs that do package a static libomp.a, avoiding the multi-minute build.
	@SYSTEM_LIBOMP="$$(ls /usr/lib/llvm-*/lib/libomp.a /usr/lib/x86_64-linux-gnu/libomp.a /usr/lib/aarch64-linux-gnu/libomp.a 2>/dev/null | head -1 || true)"; \
	if [ -n "$$SYSTEM_LIBOMP" ]; then \
		$(SUDO) cp "$$SYSTEM_LIBOMP" "$(LIB_PATH)/libomp.a" && $(SUDO) ldconfig; \
	else \
		command -v python3 >/dev/null 2>&1 || { command -v apt-get >/dev/null 2>&1 \
		&& $(SUDO) apt-get update -qq \
		&& $(SUDO) apt-get install -y --no-install-recommends python3; } || true; \
		rm -rf $(TEMP_DIR)/libomp $(TEMP_DIR)/libomp-archive; \
		curl -fsSL "https://github.com/llvm/llvm-project/archive/refs/tags/llvmorg-$(LLVM_VERSION).tar.gz" \
		-o $(TEMP_DIR)/libomp-archive; \
		mkdir -p $(TEMP_DIR)/libomp; \
		tar -xzf $(TEMP_DIR)/libomp-archive -C $(TEMP_DIR)/libomp --strip-components 1 \
		llvm-project-llvmorg-$(LLVM_VERSION)/openmp \
		llvm-project-llvmorg-$(LLVM_VERSION)/cmake; \
		cd $(TEMP_DIR)/libomp/openmp \
		&& _AR=$$(if [ -x '$(AR)' ]; then echo '$(AR)'; else command -v llvm-ar 2>/dev/null || command -v ar; fi) \
		&& _NM=$$(if [ -x '$(NM)' ]; then echo '$(NM)'; else command -v llvm-nm 2>/dev/null || command -v nm; fi) \
		&& _RANLIB=$$(if [ -x '$(RANLIB)' ]; then echo '$(RANLIB)'; else command -v llvm-ranlib 2>/dev/null || ls /usr/bin/llvm-ranlib-* 2>/dev/null | sort -V | tail -1 | grep . || command -v gcc-ranlib 2>/dev/null || ls /usr/bin/gcc-ranlib-* 2>/dev/null | sort -V | tail -1 | grep . || command -v ranlib; fi) \
		&& env LDFLAGS="" cmake -G Ninja \
		-DCMAKE_BUILD_TYPE=Release \
		-DCMAKE_POLICY_VERSION_MINIMUM=$(CMAKE_VERSION) \
		-DCMAKE_C_COMPILER="$(CC)" \
		-DCMAKE_CXX_COMPILER="$(CXX)" \
		-DCMAKE_AR="$${_AR}" \
		-DCMAKE_NM="$${_NM}" \
		-DCMAKE_RANLIB="$${_RANLIB}" \
		-DCMAKE_MAKE_PROGRAM="$(USR_LOCAL)/bin/ninja" \
		-DCMAKE_INSTALL_PREFIX="$(USR_LOCAL)" \
		-DCMAKE_INSTALL_LIBDIR="lib" \
		-DCMAKE_INSTALL_INCLUDEDIR="include" \
		-DBUILD_SHARED_LIBS=OFF \
		-DLIBOMP_ENABLE_SHARED=OFF \
		-DLIBOMP_OMPT_SUPPORT=OFF \
		-DLIBOMP_INSTALL_ALIASES=ON \
		-DPython3_EXECUTABLE="$$(command -v python3)" \
		-DCMAKE_C_FLAGS="$(CFLAGS_BASE)" \
		-DCMAKE_CXX_FLAGS="$(CFLAGS_BASE)" \
		-B $(TEMP_DIR)/libomp/openmp/build $(TEMP_DIR)/libomp/openmp \
		&& cmake --build $(TEMP_DIR)/libomp/openmp/build --parallel $(CORES) \
		&& $(SUDO) cmake --install $(TEMP_DIR)/libomp/openmp/build \
		&& rm -rf $(TEMP_DIR)/libomp $(TEMP_DIR)/libomp-archive \
		&& $(SUDO) ldconfig; \
	fi

.PHONY: ngt/install
## install NGT
ngt/install: $(USR_LOCAL)/include/NGT/Capi.h

# LLVM libomp is only linked when building with clang: a fully static glibc
# binary carrying static libomp segfaults during process startup before any
# output (e2e runs 29947389761 attempt1 / 29950003093 — CrashLoopBackOff with
# empty --previous logs on both the GCC and clang images). With GCC, cmake's
# FindOpenMP resolves -fopenmp to libgomp, the runtime this agent has always
# shipped with.
$(USR_LOCAL)/include/NGT/Capi.h: | ninja/install $(if $(findstring clang,$(notdir $(CC))),$(LIB_PATH)/libomp.a)
	$(call cmake-install,https://github.com/NGT-labs/NGT.git,ngt, \
		-DNGT_LARGE_DATASET=ON \
		$(if $(filter arm64,$(subst aarch64,arm64,$(ARCH))),,-DNGT_AVX2=ON) \
		-DBUILD_STATIC_EXECS=OFF \
		-DCMAKE_INTERPROCEDURAL_OPTIMIZATION=OFF \
		-DCMAKE_AR=$$(command -v llvm-ar 2>/dev/null || ls /usr/bin/llvm-ar-* 2>/dev/null | sort -V | tail -1 | grep . || command -v gcc-ar 2>/dev/null || ls /usr/bin/gcc-ar-* 2>/dev/null | sort -V | tail -1 | grep . || command -v ar) \
		-DCMAKE_CXX_COMPILER_AR=$$(command -v llvm-ar 2>/dev/null || ls /usr/bin/llvm-ar-* 2>/dev/null | sort -V | tail -1 | grep . || command -v gcc-ar 2>/dev/null || ls /usr/bin/gcc-ar-* 2>/dev/null | sort -V | tail -1 | grep . || command -v ar) \
		-DCMAKE_RANLIB=$$(command -v llvm-ranlib 2>/dev/null || ls /usr/bin/llvm-ranlib-* 2>/dev/null | sort -V | tail -1 | grep . || command -v gcc-ranlib 2>/dev/null || ls /usr/bin/gcc-ranlib-* 2>/dev/null | sort -V | tail -1 | grep . || command -v ranlib) \
		-DCMAKE_C_FLAGS="$(CFLAGS) $(LTO_FLAGS) $(if $(filter Linux,$(UNAME)),-fopenmp)" \
		-DCMAKE_CXX_FLAGS="$(CXXFLAGS) $(LTO_FLAGS) $(if $(filter Linux,$(UNAME)),-fopenmp)" \
		$(if $(and $(findstring clang,$(notdir $(CC))),$(filter Linux,$(UNAME))), \
		-DOpenMP_CXX_FLAGS="-fopenmp" \
		-DOpenMP_C_FLAGS="-fopenmp" \
		-DOpenMP_CXX_LIB_NAMES="omp" \
		-DOpenMP_C_LIB_NAMES="omp" \
		$(if $(LIBOMP),-DOpenMP_omp_LIBRARY="$(LIBOMP)")) \
		-DCMAKE_THREAD_LIBS_INIT="-lpthread" \
		-DCMAKE_HAVE_THREADS_LIBRARY=1 \
		-DCMAKE_USE_PTHREADS_INIT=1 \
		-DTHREADS_PREFER_PTHREAD_FLAG=OFF \
		$(if $(OPENBLAS_PATH),-DBLAS_LIBRARIES="$(OPENBLAS_PATH)" -DLAPACK_LIBRARIES="$(OPENBLAS_PATH)",) \
		-DCMAKE_EXE_LINKER_FLAGS="$(NGT_LDFLAGS)$(if $(filter ld.lld lld,$(notdir $(LLD))), -fuse-ld=lld)" \
		-DCMAKE_SHARED_LINKER_FLAGS="$(NGT_LDFLAGS)$(if $(filter ld.lld lld,$(notdir $(LLD))), -fuse-ld=lld)" \
		-DCMAKE_MODULE_LINKER_FLAGS="$(NGT_LDFLAGS)$(if $(filter ld.lld lld,$(notdir $(LLD))), -fuse-ld=lld)" \
		$(NGT_EXTRA_CMAKE_FLAGS), \
		mkdir -p $(TEMP_DIR)/ngt/build/bin/ngt $(TEMP_DIR)/ngt/build/bin/qbg && touch $(TEMP_DIR)/ngt/build/bin/ngt/ngt $(TEMP_DIR)/ngt/build/bin/qbg/qbg, \
		v$(NGT_VERSION), \
		, \
		ngt)

.PHONY: faiss/install
## install Faiss
faiss/install: $(LIB_PATH)/libfaiss.a

# Resolve the OpenBLAS library path. faiss 1.14.2 only auto-detects the
# RHEL/Fedora-style threaded names (libopenblaso/libopenblasp); on Debian/Ubuntu
# the threaded library is plain libopenblas, which faiss's
# find_library(BLAS_PREFER_THREADED ...) misses, falling through to a
# find_package(BLAS REQUIRED) that fails under cmake 4.3.3. Pre-seeding the
# BLAS_PREFER_THREADED cache entry with the real path takes faiss's working
# branch (it sets BLAS_LIBRARIES/LAPACK_LIBRARIES from it and skips find_package).
OPENBLAS_PATH = $(shell ldconfig -p 2>/dev/null | awk '/libopenblas.*\.so.*=>/{print $$NF; exit}')

$(LIB_PATH)/libfaiss.a: | ninja/install $(if $(findstring clang,$(notdir $(CC))),$(LIB_PATH)/libomp.a)
	# Faiss needs the BLAS/LAPACK (OpenBLAS) -dev packages (headers + .so symlinks)
	# and a Fortran runtime. The dev-container Dockerfile ships them, but the
	# published image used by the runtime jobs can lag. A present runtime lib
	# (libopenblas.so.0) does NOT satisfy `-lopenblas` / find_package(BLAS) without
	# the -dev package, so install unconditionally on apt-based systems (idempotent).
	@command -v apt-get >/dev/null 2>&1 \
		&& { $(SUDO) apt-get update -qq \
		&& $(SUDO) apt-get install -y --no-install-recommends libopenblas-dev liblapack-dev gfortran; } || true
	$(call cmake-install,https://github.com/facebookresearch/faiss/archive/v$(FAISS_VERSION).tar.gz,faiss, \
		-DFAISS_ENABLE_PYTHON=OFF \
		-DFAISS_ENABLE_GPU=OFF \
		$(if $(OPENBLAS_PATH),-DBLAS_PREFER_THREADED="$(OPENBLAS_PATH)" -DBLAS_LIBRARIES="$(OPENBLAS_PATH)" -DLAPACK_LIBRARIES="$(OPENBLAS_PATH)",) \
		-DCMAKE_CXX_SCAN_FOR_MODULES=OFF \
		-DCMAKE_INTERPROCEDURAL_OPTIMIZATION=OFF \
		-DCMAKE_C_FLAGS="$(CFLAGS) $(LTO_FLAGS) $(if $(MARCH),-march=$(MARCH)) $(if $(MTUNE),-mtune=$(MTUNE)) $(if $(filter Linux,$(UNAME)),-fopenmp)" \
		-DCMAKE_CXX_FLAGS="$(CXXFLAGS) $(LTO_FLAGS) $(if $(MARCH),-march=$(MARCH)) $(if $(MTUNE),-mtune=$(MTUNE)) $(if $(filter Linux,$(UNAME)),-fopenmp)" \
		$(if $(and $(findstring clang,$(notdir $(CC))),$(filter Linux,$(UNAME))), \
		-DOpenMP_CXX_FLAGS="-fopenmp" \
		-DOpenMP_C_FLAGS="-fopenmp" \
		-DOpenMP_CXX_LIB_NAMES="omp" \
		-DOpenMP_C_LIB_NAMES="omp" \
		$(if $(LIBOMP),-DOpenMP_omp_LIBRARY="$(LIBOMP)")) \
		-DCMAKE_THREAD_LIBS_INIT="-lpthread" \
		-DCMAKE_HAVE_THREADS_LIBRARY=1 \
		-DCMAKE_USE_PTHREADS_INIT=1 \
		-DTHREADS_PREFER_PTHREAD_FLAG=OFF \
		-DCMAKE_EXE_LINKER_FLAGS="$(FAISS_LDFLAGS)$(if $(filter ld.lld lld,$(notdir $(LLD))), -fuse-ld=lld)" \
		-DCMAKE_SHARED_LINKER_FLAGS="$(FAISS_LDFLAGS)$(if $(filter ld.lld lld,$(notdir $(LLD))), -fuse-ld=lld)" \
		-DCMAKE_MODULE_LINKER_FLAGS="$(FAISS_LDFLAGS)$(if $(filter ld.lld lld,$(notdir $(LLD))), -fuse-ld=lld)", \
		cd $(TEMP_DIR)/faiss && $(SUDO) find faiss -name '*.h' -exec install -D -m 0644 {} $(USR_LOCAL)/include/{} \;, \
		, \
		, \
		faiss)

.PHONY: usearch/install
## install usearch
usearch/install: $(USR_LOCAL)/include/usearch.h

$(USR_LOCAL)/include/usearch.h: | ninja/install
	$(call cmake-install,https://github.com/unum-cloud/usearch.git,usearch, \
		-DUSEARCH_BUILD_LIB_C=ON \
		-DUSEARCH_BUILD_TEST_CPP=OFF \
		-DUSEARCH_BUILD_BENCH_CPP=OFF \
		-DUSEARCH_BUILD_TEST_C=OFF \
		-DUSEARCH_USE_FP16LIB=ON \
		-DUSEARCH_USE_OPENMP=OFF \
		-DUSEARCH_USE_SIMSIMD=ON \
		-DUSEARCH_USE_JEMALLOC=OFF \
		-DCMAKE_POSITION_INDEPENDENT_CODE=ON \
		-DCMAKE_EXE_LINKER_FLAGS="" \
		-DCMAKE_SHARED_LINKER_FLAGS="" \
		-DCMAKE_MODULE_LINKER_FLAGS="", \
		cp $(TEMP_DIR)/usearch/build/libusearch_static_c.a $(LIB_PATH)/libusearch_c.a && cp $(TEMP_DIR)/usearch/build/libusearch_static_c.a $(LIB_PATH)/libusearch_static_c.a && cp $(TEMP_DIR)/usearch/build/libusearch_c.so $(LIB_PATH)/libusearch_c.so && cp $(TEMP_DIR)/usearch/c/usearch.h $(USR_LOCAL)/include/usearch.h, \
		v$(USEARCH_VERSION))

.PHONY: cmake/install
## install CMAKE
cmake/install:
	CMAKE_ARCH=$$(if [ "$(ARCH)" = "aarch64" ] || [ "$(ARCH)" = "arm64" ]; then echo "aarch64"; else echo "x86_64"; fi); \
	TAR_NAME="cmake-$(CMAKE_VERSION)-linux-$${CMAKE_ARCH}.tar.gz" \
	&& curl -fsSL "https://github.com/Kitware/CMake/releases/download/v$(CMAKE_VERSION)/$${TAR_NAME}" -o "$(TEMP_DIR)/$${TAR_NAME}" \
	&& $(SUDO) tar -xzf "$(TEMP_DIR)/$${TAR_NAME}" -C $(USR_LOCAL) --strip-components 1 \
	&& rm -rf "$(TEMP_DIR)/$${TAR_NAME}" \
	&& cmake --version

.PHONY: ninja/install
## install ninja-build
ninja/install:
	NINJA_ARCH=$$(if [ "$(ARCH)" = "aarch64" ] || [ "$(ARCH)" = "arm64" ]; then echo "-aarch64"; else echo ""; fi); \
	TAR_NAME="ninja-linux$${NINJA_ARCH}.zip" \
	&& curl -fsSL "https://github.com/ninja-build/ninja/releases/download/v$(NINJA_VERSION)/$${TAR_NAME}" -o "$(TEMP_DIR)/$${TAR_NAME}" \
	&& $(SUDO) unzip -q -o "$(TEMP_DIR)/$${TAR_NAME}" -d $(USR_LOCAL)/bin \
	&& rm -rf "$(TEMP_DIR)/$${TAR_NAME}" \
	&& $(SUDO) chmod +x $(USR_LOCAL)/bin/ninja \
	&& ninja --version
>>>>>>> e303744f0 (fix: skip NGT AVX2 on arm64 (#3580))

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
## install docker-cli
docker-cli/install: $(BINDIR)/docker

$(BINDIR)/docker: $(BINDIR)
	curl -fsSL https://download.docker.com/linux/static/stable/$(shell uname -m)/docker-$(shell echo $(DOCKER_VERSION) | cut -c2-).tgz -o $(TEMP_DIR)/docker.tgz \
	&& tar -xzvf $(TEMP_DIR)/docker.tgz -C $(TEMP_DIR) \
	&& mv $(TEMP_DIR)/docker/docker $(BINDIR) \
	&& rm -rf $(TEMP_DIR)/docker{.tgz,}

.PHONY: replace/busybox
## replace busybox version
replace/busybox:
	find . -type f \( -name "*.yaml" -o -name "*.md" \) -exec sed -i -E 's/busybox:([0-9]+\.[0-9]+\.[0-9]+|latest)/busybox:$(BUSYBOX_VERSION)/g' {} +
