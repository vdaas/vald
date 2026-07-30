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

$(BENCH_DATASETS): $(BENCH_DATASET_MD5S) $(BENCH_DATASET_HDF5_DIR)
	@$(call green, "downloading datasets for benchmark...")
	curl -fsSL -o $@ https://ann-benchmarks.com/$(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,%.hdf5,$@)
	(cd $(BENCH_DATASET_BASE_DIR); \
	md5sum -c $(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,$(BENCH_DATASET_MD5_DIR_NAME)/%.md5,$@) || \
	(rm -f $(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,$(BENCH_DATASET_HDF5_DIR_NAME)/%.hdf5,$@) && exit 1))

$(BENCH_DATASET_HDF5_DIR):
	$(call mkdir, $@)
	$(call rm, -rf, $@/*)

%.large_dataset_dir:
	@test -f $* || mkdir -p $*

$(SIFT1B_BASE_FILE) $(SIFT1B_LEARN_FILE) $(SIFT1B_QUERY_FILE): | $(SIFT1B_ROOT_DIR).large_dataset_dir
	test -f $@ || \
	curl -fsSL $(SIFT1B_BASE_URL)$(subst $(SIFT1B_ROOT_DIR)/,,$@).gz | \
	gunzip -d > $@

$(SIFT1B_GROUNDTRUTH_DIR): | $(SIFT1B_ROOT_DIR).large_dataset_dir
	test -f $@ || \
	curl -fsSL $(SIFT1B_BASE_URL)bigann_gnd.tar.gz | \
	tar -C $(SIFT1B_ROOT_DIR) -zx

$(DEEP1B_GROUNDTRUTH_FILE) $(DEEP1B_QUERY_FILE) $(DEEP1B_BASE_CHUNK_FILES) $(DEEP1B_LEARN_CHUNK_FILES): | $(DEEP1B_ROOT_DIR).large_dataset_dir
	test -f $@ || \
	curl -fsSL "$(shell curl -fsSL "$(DEEP1B_API_URL)$(subst $(DEEP1B_ROOT_DIR),,$@)" | \
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
