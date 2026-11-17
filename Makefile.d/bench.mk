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

# ann-benchmarks.com resets connections now and then ("curl: (56) Recv failure: Connection reset by
# peer" failed an E2E job in CI), and a truncated file must not stay behind as an up-to-date target.
# curl retries every transient error itself (CURL_RETRY_ALL); the loop around it also covers a
# download that finished but fails its checksum, removes the file between attempts and waits a
# little longer each time. The Run E2E CRUD steps of the workflows call this target.
BENCH_DATASET_DOWNLOAD_ATTEMPTS ?= 5
BENCH_DATASET_DOWNLOAD_DELAY ?= 15

$(BENCH_DATASETS): $(BENCH_DATASET_MD5S) $(BENCH_DATASET_HDF5_DIR)
	@$(call green, "downloading datasets for benchmark...")
	@name=$(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,%.hdf5,$@); \
	md5=$(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,$(BENCH_DATASET_MD5_DIR_NAME)/%.md5,$@); \
	attempt=1; \
	while :; do \
		if $(CURL) $(CURL_RETRY_ALL) -o $@ https://ann-benchmarks.com/$$name \
			&& (cd $(BENCH_DATASET_BASE_DIR) && md5sum -c $$md5); then \
			break; \
		fi; \
		rm -f $@; \
		if [ $$attempt -ge $(BENCH_DATASET_DOWNLOAD_ATTEMPTS) ]; then \
			echo "ERROR: could not download $$name in $$attempt attempts" >&2; \
			exit 1; \
		fi; \
		delay=$$(( attempt * $(BENCH_DATASET_DOWNLOAD_DELAY) )); \
		echo "WARN: downloading $$name failed (attempt $$attempt of $(BENCH_DATASET_DOWNLOAD_ATTEMPTS)); retrying in $${delay}s" >&2; \
		sleep $$delay; \
		attempt=$$(( attempt + 1 )); \
	done

$(BENCH_DATASET_HDF5_DIR):
	$(call mkdir, $@)
	$(call rm, -rf, $@/*)

%.large_dataset_dir:
	@test -f $* || mkdir -p $*

$(SIFT1B_BASE_FILE) $(SIFT1B_LEARN_FILE) $(SIFT1B_QUERY_FILE): | $(SIFT1B_ROOT_DIR).large_dataset_dir
	test -f $@ || \
	$(CURL) $(SIFT1B_BASE_URL)$(subst $(SIFT1B_ROOT_DIR)/,,$@).gz | \
	gunzip -d > $@

$(SIFT1B_GROUNDTRUTH_DIR): | $(SIFT1B_ROOT_DIR).large_dataset_dir
	test -f $@ || \
	$(CURL) $(SIFT1B_BASE_URL)bigann_gnd.tar.gz | \
	tar -C $(SIFT1B_ROOT_DIR) -zx

$(DEEP1B_GROUNDTRUTH_FILE) $(DEEP1B_QUERY_FILE) $(DEEP1B_BASE_CHUNK_FILES) $(DEEP1B_LEARN_CHUNK_FILES): | $(DEEP1B_ROOT_DIR).large_dataset_dir
	test -f $@ || \
	$(CURL) "$(shell $(CURL) "$(DEEP1B_API_URL)$(subst $(DEEP1B_ROOT_DIR),,$@)" | \
	sed -e 's/^{\(.*\)}$$/\1/' | \
	tr ',' '\n' | \
	grep href | \
	cut -d ':' -f 2- | \
	tr -d '"')" -o $@

$(DEEP1B_BASE_FILE): | $(DEEP1B_BASE_DIR).large_dataset_dir $(DEEP1B_BASE_CHUNK_FILES)
	cat $(DEEP1B_BASE_CHUNK_FILES) > $@

$(DEEP1B_LEARN_FILE): | $(DEEP1B_LEARN_DIR).large_dataset_dir $(DEEP1B_LEARN_CHUNK_FILES)
	cat $(DEEP1B_LEARN_CHUNK_FILES) > $@

.PHONY: bench/datasets
## fetch datasets for benchmark
bench/datasets: $(BENCH_DATASETS)

.PHONY: bench/datasets/clean
## clean datasets for benchmark
bench/datasets/clean:
	rm -rf $(BENCH_DATASETS)

.PHONY: bench/datasets/basedir/print
## print benchmark dataset base directory path
bench/datasets/basedir/print:
	@echo $(BENCH_DATASET_BASE_DIR)

.PHONY: bench/datasets/md5dir/print
## print benchmark dataset md5 directory path
bench/datasets/md5dir/print:
	@echo $(BENCH_DATASET_MD5_DIR)

.PHONY: bench/datasets/hdf5dir/print
## print benchmark dataset hdf5 directory path
bench/datasets/hdf5dir/print:
	@echo $(BENCH_DATASET_HDF5_DIR)

.PHONY: bench/datasets/large
## fetch large datasets for benchmark
bench/datasets/large: \
	bench/datasets/large/sift1b \
	bench/datasets/large/deep1b

.PHONY: bench/datasets/large/sift1b
## fetch sift1b dataset for benchmark
bench/datasets/large/sift1b: \
	$(SIFT1B_BASE_FILE) \
	$(SIFT1B_LEARN_FILE) \
	$(SIFT1B_QUERY_FILE) \
	$(SIFT1B_GROUNDTRUTH_DIR)

.PHONY: bench/datasets/large/deep1b
## fetch deep1b dataset for benchmark
bench/datasets/large/deep1b: \
	$(DEEP1B_BASE_FILE) \
	$(DEEP1B_LEARN_FILE) \
	$(DEEP1B_QUERY_FILE) \
	$(DEEP1B_GROUNDTRUTH_FILE)

.PHONY: bench
## run all benchmarks
## NOTE: the agent gRPC (stream/sequential) and gateway sequential gRPC
## benchmarks that used to run here were removed together with
## hack/benchmark/e2e; their recall/QPS measurement is now covered by the
## tests/v2/e2e/assets/agent_recall_qps.yaml scenario (see e2e/v2 targets in
## Makefile.d/e2e.mk). This target is intentionally left as a no-op
## placeholder for future benchmark categories.
bench:

BENCH_PKG ?= ./...
BENCH_RUN ?= .
BENCH_COUNT ?= 10
BENCH_TIME ?= 1s
BENCH_OUT ?= $(TEMP_DIR)/bench.txt

.PHONY: bench/go
## run Go micro-benchmarks without -race/-cover (unlike `test`), e.g.
## make bench/go BENCH_PKG=./internal/core/algorithm/ngt/... BENCH_RUN=Search BENCH_OUT=/tmp/new.txt
## then compare two runs with `make bench/compare BENCH_OLD=/tmp/old.txt BENCH_NEW=/tmp/new.txt`
bench/go: \
	ngt/install
	$(GO_ENV_VARS) CGO_LDFLAGS="$(TEST_LDFLAGS)" go test -mod=readonly -run '^$$' \
		-bench '$(BENCH_RUN)' -benchmem -count $(BENCH_COUNT) -benchtime $(BENCH_TIME) \
		-timeout=$(GOTEST_TIMEOUT) -ldflags "$(GO_TEST_LDFLAGS)" $(BENCH_PKG) | tee $(BENCH_OUT)

.PHONY: bench/compare
## compare two bench/go outputs with benchstat
bench/compare: \
	benchstat/install
	$(GOBIN)/benchstat $(BENCH_OLD) $(BENCH_NEW)

.PHONY: bench/kill
## kill all benchmark processes
bench/kill:
	ps aux \
	| grep go \
	| grep -v nvim \
	| grep -v tmux \
	| grep -v gopls \
	| grep -v "rg go" \
	| grep -v "grep go" \
	| awk '{print $1}' \
	| xargs -P$(CORES) kill -9
