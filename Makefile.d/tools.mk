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
	# Install via the Go module proxy (GOPROXY) instead of the upstream
	# install.sh: in Docker builds the anonymous GitHub release download is
	# rate-limited and returns a non-tarball body, so the script's sha256 check
	# fails deterministically. `go install` goes through proxy.golang.org and is
	# reliable. (golangci-lint v2 has no replace directives, so this builds.)
	GOBIN=$(BINDIR) CGO_ENABLED=0 go install \
		github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCILINT_VERSION)

.PHONY: goimports/install
## install goimports
goimports/install: \
	$(GOBIN)/goimports

.PHONY: strictgoimports/install
## install strictgoimports
strictgoimports/install: \
	$(GOBIN)/strictgoimports

.PHONY: gofumpt/install
## install gofumpt
gofumpt/install: \
	$(GOBIN)/gofumpt

.PHONY: golines/install
## install golines
golines/install: \
	$(GOBIN)/golines

.PHONY: crlfmt/install
## install crlfmt
crlfmt/install: \
	$(GOBIN)/crlfmt

.PHONY: actionlint/install
## install actionlint
actionlint/install: \
	$(GOBIN)/actionlint

.PHONY: ghalint/install
## install ghalint
ghalint/install: \
	$(GOBIN)/ghalint

.PHONY: pinact/install
## install pinact
pinact/install: \
	$(GOBIN)/pinact

.PHONY: ghatm/install
## install ghatm
ghatm/install: \
	$(GOBIN)/ghatm

.PHONY: buf/install
## install buf
buf/install: \
	$(GOBIN)/buf

.PHONY: k9s/install
## install k9s
k9s/install: \
	$(GOBIN)/k9s

.PHONY: stern/install
## install stern
stern/install: \
	$(GOBIN)/stern

.PHONY: yamlfmt/install
## install yamlfmt
yamlfmt/install: \
	$(GOBIN)/yamlfmt

.PHONY: gomodifytags/install
## install gomodifytags
gomodifytags/install: \
	$(GOBIN)/gomodifytags

.PHONY: impl/install
## install impl
impl/install: \
	$(GOBIN)/impl

.PHONY: delve/install
## install delve
delve/install: \
	$(GOBIN)/dlv

.PHONY: staticcheck/install
## install staticcheck
staticcheck/install: \
	$(GOBIN)/staticcheck

.PHONY: tparse/install
## install tparse
tparse/install: \
	$(GOBIN)/tparse

.PHONY: gotestfmt/install
## install gotestfmt
gotestfmt/install: \
	$(GOBIN)/gotestfmt

.PHONY: gotests/install
## install gotests
gotests/install: \
	$(GOBIN)/gotests

.PHONY: protoc-gen-doc/install
## install protoc-gen-doc
protoc-gen-doc/install: \
	$(GOBIN)/protoc-gen-doc

GO_BINS = \
	$(GOBIN)/goimports \
	$(GOBIN)/strictgoimports \
	$(GOBIN)/gofumpt \
	$(GOBIN)/golines \
	$(GOBIN)/crlfmt \
	$(GOBIN)/actionlint \
	$(GOBIN)/ghalint \
	$(GOBIN)/pinact \
	$(GOBIN)/ghatm \
	$(GOBIN)/buf \
	$(GOBIN)/k9s \
	$(GOBIN)/stern \
	$(GOBIN)/yamlfmt \
	$(GOBIN)/gomodifytags \
	$(GOBIN)/impl \
	$(GOBIN)/dlv \
	$(GOBIN)/staticcheck \
	$(GOBIN)/tparse \
	$(GOBIN)/gotestfmt \
	$(GOBIN)/gotests \
	$(GOBIN)/protoc-gen-doc

$(GO_BINS):
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
	$(GO_ENV_VARS) go install -mod=readonly golang.org/x/tools/gopls@latest

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
	&& $(SUDO) rm -rf $(GOROOT) \
	&& $(SUDO) mv $(TEMP_DIR)/go $(GOROOT) \
	&& $(GOROOT)/bin/go version

.PHONY: rust/install
## install rust
rust/install: $(CARGO_HOME)/bin/cargo

$(CARGO_HOME)/bin/cargo:
	curl --proto '=https' --tlsv1.2 -fsSL https://sh.rustup.rs | CARGO_HOME=$(CARGO_HOME) RUSTUP_HOME=$(RUSTUP_HOME) sh -s -- --default-toolchain $(RUST_VERSION) -y && \
	$(CARGO_HOME)/bin/rustup toolchain install $(RUST_VERSION) && \
	$(CARGO_HOME)/bin/rustup default $(RUST_VERSION) && \
	$(CARGO_HOME)/bin/rustup component add rust-src rust-analyzer rust-docs rustfmt clippy

.PHONY: rustfmt/install
## install rustfmt
$(CARGO_HOME)/bin/rustfmt: $(CARGO_HOME)/bin/cargo
	CARGO_HOME=${CARGO_HOME} RUSTUP_HOME=${RUSTUP_HOME} \
	$(CARGO_HOME)/bin/rustup component add rustfmt

rustfmt/install:
	$(MAKE) rust/install
	$(MAKE) $(CARGO_HOME)/bin/rustfmt

.PHONY: clang/install
## install clang, lld, and llvm from source
clang/install: $(USR_LOCAL)/bin/clang

.PHONY: openmp/install
## install openmp from source
openmp/install: $(USR_LOCAL)/bin/clang

$(USR_LOCAL)/bin/clang:
	$(call cmake-install,https://github.com/llvm/llvm-project/releases/download/llvmorg-$(LLVM_VERSION)/llvm-project-$(LLVM_VERSION).src.tar.xz,llvm, \
		-DLLVM_ENABLE_PROJECTS="clang;lld" \
		-DLLVM_ENABLE_RUNTIMES="openmp;libcxx;libcxxabi;libunwind" \
		-DLLVM_TARGETS_TO_BUILD="X86;AArch64" \
		-DLLVM_INCLUDE_TESTS=OFF \
		-DLLVM_INCLUDE_EXAMPLES=OFF \
		-DLLVM_ENABLE_RTTI=ON \
		-DLLVM_ENABLE_ZLIB=ON \
		-DLLVM_ENABLE_ZSTD=ON \
		-DLLVM_ENABLE_LIBXML2=OFF \
		-DLLVM_ENABLE_TERMINFO=OFF \
		-DLLVM_ENABLE_BINDINGS=OFF \
		-DLLVM_ENABLE_OCAMLDOC=OFF \
		-DLLVM_ENABLE_DOXYGEN=OFF \
		-DLLVM_ENABLE_SPHINX=OFF \
		-DLIBOMP_INSTALL_ALIASES=ON \
		-DLLVM_ENABLE_LTO=Thin \
		-DCMAKE_INTERPROCEDURAL_OPTIMIZATION=OFF \
		-DCMAKE_C_FLAGS="$(CFLAGS_BASE)" \
		-DCMAKE_CXX_FLAGS="$(CFLAGS_BASE)",,, \
		llvm)

.PHONY: zlib/install
## install zlib
zlib/install: $(LIB_PATH)/libz.a

$(LIB_PATH)/libz.a: | $(LIB_PATH) ninja/install
	$(call cmake-install,https://github.com/madler/zlib/releases/download/v$(ZLIB_VERSION)/zlib-$(ZLIB_VERSION).tar.gz,zlib, \
		-DZLIB_BUILD_SHARED=OFF \
		-DZLIB_BUILD_STATIC=ON \
		-DZLIB_COMPAT=ON \
		-DZLIB_USE_STATIC_LIBS=ON, \
		$(SUDO) rm -f $(USR_LOCAL)/include/zlib.h $(USR_LOCAL)/include/zconf.h $(LIB_PATH)/libz.a $(USR_LOCAL)/share/man/man3/zlib.3,,,zlibstatic)

.PHONY: hdf5/install
## install hdf5
hdf5/install: $(LIB_PATH)/libhdf5.a

$(LIB_PATH)/libhdf5.a: | $(LIB_PATH) zlib/install
	$(call cmake-install,https://github.com/HDFGroup/hdf5/archive/refs/tags/$(HDF5_VERSION).tar.gz,hdf5, \
		-DHDF5_BUILD_CPP_LIB=OFF \
		-DHDF5_BUILD_HL_LIB=ON \
		-DHDF5_BUILD_STATIC_EXECS=OFF \
		-DHDF5_BUILD_TOOLS=OFF \
		-DHDF5_BUILD_EXAMPLES=OFF \
		-DHDF5_ENABLE_Z_LIB_SUPPORT=ON \
		-DH5_ZLIB_INCLUDE_DIR=$(USR_LOCAL)/include \
		-DH5_ZLIB_LIBRARY=$(LIB_PATH)/libz.a, \
		$(SUDO) rm -f $(USR_LOCAL)/include/H5*.h $(USR_LOCAL)/include/hdf5*.h $(LIB_PATH)/libhdf5*)

.PHONY: libomp/install
## install libomp static library from LLVM openmp source
libomp/install: $(LIB_PATH)/libomp.a

$(LIB_PATH)/libomp.a: | ninja/install
	@$(call green, "Installing libomp...")
	$(SUDO) mkdir -p $(LIB_PATH) $(INCLUDE_PATH)
	@SYSTEM_LIBOMP="$$(ls /usr/lib/llvm-*/lib/libomp.a /usr/lib/x86_64-linux-gnu/libomp.a /usr/lib/aarch64-linux-gnu/libomp.a 2>/dev/null | head -1)"; \
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

$(USR_LOCAL)/include/NGT/Capi.h: | ninja/install $(LIB_PATH)/libomp.a
	$(call cmake-install,https://github.com/NGT-labs/NGT.git,ngt, \
		-DNGT_LARGE_DATASET=ON \
		-DBUILD_STATIC_EXECS=OFF \
		-DCMAKE_INTERPROCEDURAL_OPTIMIZATION=OFF \
		-DCMAKE_AR=$$(command -v llvm-ar 2>/dev/null || ls /usr/bin/llvm-ar-* 2>/dev/null | sort -V | tail -1 | grep . || command -v gcc-ar 2>/dev/null || ls /usr/bin/gcc-ar-* 2>/dev/null | sort -V | tail -1 | grep . || command -v ar) \
		-DCMAKE_CXX_COMPILER_AR=$$(command -v llvm-ar 2>/dev/null || ls /usr/bin/llvm-ar-* 2>/dev/null | sort -V | tail -1 | grep . || command -v gcc-ar 2>/dev/null || ls /usr/bin/gcc-ar-* 2>/dev/null | sort -V | tail -1 | grep . || command -v ar) \
		-DCMAKE_RANLIB=$$(command -v llvm-ranlib 2>/dev/null || ls /usr/bin/llvm-ranlib-* 2>/dev/null | sort -V | tail -1 | grep . || command -v gcc-ranlib 2>/dev/null || ls /usr/bin/gcc-ranlib-* 2>/dev/null | sort -V | tail -1 | grep . || command -v ranlib) \
		-DCMAKE_C_FLAGS="$(CFLAGS) $(if $(filter Linux,$(UNAME)),-fopenmp)" \
		-DCMAKE_CXX_FLAGS="$(CXXFLAGS) $(if $(filter Linux,$(UNAME)),-fopenmp)" \
		$(if $(filter Linux,$(UNAME)), \
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

$(LIB_PATH)/libfaiss.a: | ninja/install $(LIB_PATH)/libomp.a
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
		-DCMAKE_C_FLAGS="$(CFLAGS) $(if $(filter Linux,$(UNAME)),-fopenmp)" \
		-DCMAKE_CXX_FLAGS="$(CXXFLAGS) $(if $(filter Linux,$(UNAME)),-fopenmp)" \
		$(if $(filter Linux,$(UNAME)), \
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
	&& rm -rf $(TEMP_DIR)/docker $(TEMP_DIR)/docker.tgz

.PHONY: fossa/install
## install fossa
fossa/install: $(BINDIR)/fossa

$(BINDIR)/fossa:
	curl -H 'Cache-Control: no-cache' https://raw.githubusercontent.com/fossas/fossa-cli/master/install.sh | sh -s -- -b $(BINDIR)

.PHONY: replace/busybox
## replace busybox version
replace/busybox:
	cat $(ROOTDIR)/.gitfiles | grep -E '\.(yaml|md)$$' | sed -e 's%^%$(ROOTDIR)/%' | xargs -I {} -P $(CORES) sed -i -E 's/busybox:([0-9]+\.[0-9]+\.[0-9]+|latest)/busybox:$(BUSYBOX_VERSION)/g' "{}"
