#
# Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
$(BENCH_DATASETS): $(BENCH_DATASET_MD5S) $(BENCH_DATASET_HDF5_DIR)
	@$(call green, "downloading datasets for benchmark...")
	curl -fsSL -o $@ http://ann-benchmarks.com/$(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,%.hdf5,$@)
	(cd $(BENCH_DATASET_BASE_DIR); \
	    md5sum -c $(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,$(BENCH_DATASET_MD5_DIR_NAME)/%.md5,$@) || \
	    (rm -f $(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,$(BENCH_DATASET_HDF5_DIR_NAME)/%.hdf5,$@) && exit 1))

$(BENCH_DATASET_HDF5_DIR):
	$(call mkdir, $@)
	$(call rm, -rf, $@/*)

%.large_dataset_dir:
	@test -f $* || mkdir -p $*

$(SIFT1B_BASE_FILE) $(SIFT1B_LEARN_FILE) $(SIFT1B_QUERY_FILE): | $(SIFT1B_ROOT_DIR).large_dataset_dir
	test -f $@ || curl -fsSL $(SIFT1B_BASE_URL)$(subst $(SIFT1B_ROOT_DIR)/,,$@).gz | gunzip -d > $@

$(SIFT1B_GROUNDTRUTH_DIR): | $(SIFT1B_ROOT_DIR).large_dataset_dir
	test -f $@ || curl -fsSL $(SIFT1B_BASE_URL)bigann_gnd.tar.gz | tar -C $(SIFT1B_ROOT_DIR) -zx

$(DEEP1B_GROUNDTRUTH_FILE) $(DEEP1B_QUERY_FILE) $(DEEP1B_BASE_CHUNK_FILES) $(DEEP1B_LEARN_CHUNK_FILES): | $(DEEP1B_ROOT_DIR).large_dataset_dir
	test -f $@ || curl -fsSL "$(shell curl -fsSL "$(DEEP1B_API_URL)$(subst $(DEEP1B_ROOT_DIR),,$@)" | sed -e 's/^{\(.*\)}$$/\1/' | tr ',' '\n' | grep href | cut -d ':' -f 2- | tr -d '"')" -o $@

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
bench/datasets/basedir/print:
	@echo $(BENCH_DATASET_BASE_DIR)

.PHONY: bench/datasets/md5dir/print
bench/datasets/md5dir/print:
	@echo $(BENCH_DATASET_MD5_DIR)

.PHONY: bench/datasets/hdf5dir/print
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

pprof/%.cpu.svg: \
	pprof/%.bin
	go tool pprof \
	    --svg \
	    $< \
	    $(patsubst %.svg,%.out,$@) \
	    > $@

pprof/%.mem.svg: \
	pprof/%.bin
	go tool pprof \
	    --svg \
	    $< \
	    $(patsubst %.svg,%.out,$@) \
	    > $@

.PHONY: bench
## run all benchmarks
bench: \
	bench/core \
	bench/agent \
	bench/ngtd \
	bench/gateway

.PHONY: bench/core
## run benchmarks for core
bench/core: \
	bench/core/ngt \
	bench/core/gongt

.PHONY: bench/core/ngt
## run benchmark for NGT core
bench/core/ngt: \
	bench/core/ngt/sequential \
	bench/core/ngt/parallel

.PHONY: bench/core/ngt/sequential
## run benchmark for NGT core sequential methods
bench/core/ngt/sequential: \
	pprof/core/ngt/sequential.cpu.svg \
	pprof/core/ngt/sequential.mem.svg
pprof/core/ngt/sequential.bin: \
	hack/benchmark/core/ngt/ngt_bench_test.go
	mkdir -p $(dir $@)
	go test \
	    -count=1 \
	    -timeout=1h \
	    -bench=NGTSequential \
	    -benchmem \
	    -o $@ \
	    -cpuprofile $(patsubst %.bin,%.cpu.out,$@) \
	    -memprofile $(patsubst %.bin,%.mem.out,$@) \
	    -trace $(patsubst %.bin,%.trace.out,$@) \
	    $< \
	    -dataset=$(DATASET_ARGS)

.PHONY: bench/core/ngt/parallel
## run benchmark for NGT core parallel methods
bench/core/ngt/parallel: \
	pprof/core/ngt/parallel.cpu.svg \
	pprof/core/ngt/parallel.mem.svg
pprof/core/ngt/parallel.bin: \
	hack/benchmark/core/ngt/ngt_bench_test.go
	mkdir -p $(dir $@)
	go test \
	    -count=1 \
	    -timeout=1h \
	    -bench=NGTParallel \
	    -benchmem \
	    -o $@ \
	    -cpuprofile $(patsubst %.bin,%.cpu.out,$@) \
	    -memprofile $(patsubst %.bin,%.mem.out,$@) \
	    -trace $(patsubst %.bin,%.trace.out,$@) \
	    $< \
	    -dataset=$(DATASET_ARGS)

.PHONY: bench/core/gongt
## run benchmark for gongt core
bench/core/gongt: \
	bench/core/gongt/sequential \
	bench/core/gongt/parallel

.PHONY: bench/core/gongt/sequential
## run benchmark for gongt core sequential methods
bench/core/gongt/sequential: \
	pprof/core/gongt/sequential.cpu.svg \
	pprof/core/gongt/sequential.mem.svg
pprof/core/gongt/sequential.bin: \
	hack/benchmark/core/gongt/gongt_bench_test.go
	mkdir -p $(dir $@)
	go test \
	    -count=1 \
	    -timeout=1h \
	    -bench=GoNGTSequential \
	    -benchmem \
	    -o $@ \
	    -cpuprofile $(patsubst %.bin,%.cpu.out,$@) \
	    -memprofile $(patsubst %.bin,%.mem.out,$@) \
	    -trace $(patsubst %.bin,%.trace.out,$@) \
	    $< \
	    -dataset=$(DATASET_ARGS)

.PHONY: bench/core/gongt/parallel
## run benchmark for gongt core parallel methods
bench/core/gongt/parallel: \
	pprof/core/gongt/parallel.cpu.svg \
	pprof/core/gongt/parallel.mem.svg
pprof/core/gongt/parallel.bin: \
	hack/benchmark/core/gongt/gongt_bench_test.go
	mkdir -p $(dir $@)
	go test \
	    -count=1 \
	    -timeout=1h \
	    -bench=GoNGTParallel \
	    -benchmem \
	    -o $@ \
	    -cpuprofile $(patsubst %.bin,%.cpu.out,$@) \
	    -memprofile $(patsubst %.bin,%.mem.out,$@) \
	    -trace $(patsubst %.bin,%.trace.out,$@) \
	    $< \
	    -dataset=$(DATASET_ARGS)

.PHONY: bench/agent
## run benchmarks for vald agent
bench/agent: \
	bench/agent/stream \
	bench/agent/sequential/grpc \
	bench/agent/sequential/rest

.PHONY: bench/agent/stream
## run benchmark for agent gRPC stream
bench/agent/stream: \
	pprof/agent/stream.cpu.svg \
	pprof/agent/stream.mem.svg
pprof/agent/stream.bin: \
	hack/benchmark/e2e/agent/core/ngt/ngt_bench_test.go \
	ngt/install
	mkdir -p $(dir $@)
	go test \
	    -count=1 \
	    -timeout=1h \
	    -bench=gRPC_Stream \
	    -benchmem \
	    -o $@ \
	    -cpuprofile $(patsubst %.bin,%.cpu.out,$@) \
	    -memprofile $(patsubst %.bin,%.mem.out,$@) \
	    -trace $(patsubst %.bin,%.trace.out,$@) \
	    $< \
	    -dataset=$(DATASET_ARGS)

.PHONY: bench/agent/sequential/grpc
## run benchmark for agent gRPC sequential
bench/agent/sequential/grpc: \
	pprof/agent/sequential/grpc.cpu.svg \
	pprof/agent/sequential/grpc.mem.svg
pprof/agent/sequential/grpc.bin: \
	hack/benchmark/e2e/agent/core/ngt/ngt_bench_test.go \
	ngt/install
	mkdir -p $(dir $@)
	go test \
	    -count=1 \
	    -timeout=1h \
	    -bench=gRPC_Sequential \
	    -benchmem \
	    -o $@ \
	    -cpuprofile $(patsubst %.bin,%.cpu.out,$@) \
	    -memprofile $(patsubst %.bin,%.mem.out,$@) \
	    -trace $(patsubst %.bin,%.trace.out,$@) \
	    $< \
	    -dataset=$(DATASET_ARGS)

.PHONY: bench/agent/sequential/rest
## run benchmark for agent REST
bench/agent/sequential/rest: \
	pprof/agent/sequential/rest.cpu.svg \
	pprof/agent/sequential/rest.mem.svg
pprof/agent/sequential/rest.bin: \
	hack/benchmark/e2e/agent/core/ngt/ngt_bench_test.go \
	ngt/install
	mkdir -p $(dir $@)
	go test \
	    -count=1 \
	    -timeout=1h \
	    -bench=REST_Sequential \
	    -benchmem \
	    -o $@ \
	    -cpuprofile $(patsubst %.bin,%.cpu.out,$@) \
	    -memprofile $(patsubst %.bin,%.mem.out,$@) \
	    -trace $(patsubst %.bin,%.trace.out,$@) \
	    $< \
	    -dataset=$(DATASET_ARGS)

.PHONY: bench/ngtd
## run benchmarks for NGTD
bench/ngtd: \
	bench/ngtd/stream \
	bench/ngtd/sequential/grpc \
	bench/ngtd/sequential/rest

.PHONY: bench/ngtd/stream
## run benchmark for NGTD gRPC stream
bench/ngtd/stream: \
	pprof/ngtd/stream.cpu.svg \
	pprof/ngtd/stream.mem.svg
pprof/ngtd/stream.bin: \
	hack/benchmark/e2e/external/ngtd/ngtd_bench_test.go \
	ngt/install
	mkdir -p $(dir $@)
	go test \
	    -count=1 \
	    -timeout=1h \
	    -bench=gRPC_Stream \
	    -benchmem \
	    -o $@ \
	    -cpuprofile $(patsubst %.bin,%.cpu.out,$@) \
	    -memprofile $(patsubst %.bin,%.mem.out,$@) \
	    -trace $(patsubst %.bin,%.trace.out,$@) \
	    $< \
	    -dataset=$(DATASET_ARGS)

.PHONY: bench/ngtd/sequential/grpc
## run benchmark for NGTD gRPC sequential
bench/ngtd/sequential/grpc: \
	pprof/ngtd/sequential/grpc.cpu.svg \
	pprof/ngtd/sequential/grpc.mem.svg
pprof/ngtd/sequential/grpc.bin: \
	hack/benchmark/e2e/external/ngtd/ngtd_bench_test.go \
	ngt/install
	mkdir -p $(dir $@)
	go test \
	    -count=1 \
	    -timeout=1h \
	    -bench=gRPC_Sequential \
	    -benchmem \
	    -o $@ \
	    -cpuprofile $(patsubst %.bin,%.cpu.out,$@) \
	    -memprofile $(patsubst %.bin,%.mem.out,$@) \
	    -trace $(patsubst %.bin,%.trace.out,$@) \
	    $< \
	    -dataset=$(DATASET_ARGS)

.PHONY: bench/ngtd/sequential/rest
## run benchmark for NGTD REST stream
bench/ngtd/sequential/rest: \
	pprof/ngtd/sequential/rest.cpu.svg \
	pprof/ngtd/sequential/rest.mem.svg
pprof/ngtd/sequential/rest.bin: \
	hack/benchmark/e2e/external/ngtd/ngtd_bench_test.go \
	ngt/install
	mkdir -p $(dir $@)
	go test \
	    -count=1 \
	    -timeout=1h \
	    -bench=REST_Sequential \
	    -benchmem \
	    -o $@ \
	    -cpuprofile $(patsubst %.bin,%.cpu.out,$@) \
	    -memprofile $(patsubst %.bin,%.mem.out,$@) \
	    -trace $(patsubst %.bin,%.trace.out,$@) \
	    $< \
	    -dataset=$(DATASET_ARGS)

.PHONY: bench/gateway
## run benchmarks for gateway
bench/gateway: \
	bench/gateway/sequential

.PHONY: bench/gateway/sequential
## run benchmark for gateway sequential
bench/gateway/sequential: \
	pprof/gateway/sequential.cpu.svg \
	pprof/gateway/sequential.mem.svg
pprof/gateway/sequential.bin: \
	hack/benchmark/e2e/gateway/vald/vald_bench_test.go \
	ngt/install
	mkdir -p $(dir $@)
	go test \
	    -count=1 \
	    -timeout=1h \
	    -bench=Sequential \
	    -benchmem \
	    -o $@ \
	    -cpuprofile $(patsubst %.bin,%.cpu.out,$@) \
	    -memprofile $(patsubst %.bin,%.mem.out,$@) \
	    -trace $(patsubst %.bin,%.trace.out,$@) \
	    $< \
	    -dataset=$(DATASET_ARGS)

.PHONY: profile
## execute profile
profile: \
	clean \
	deps \
	bench \
	profile/agent/stream \
	profile/agent/sequential/grpc \
	profile/agent/sequential/rest

.PHONY: profile/agent/stream
profile/agent/stream:
	$(call profile-web,pprof/agent/stream)

.PHONY: profile/agent/sequential/grpc
profile/agent/sequential/grpc:
	$(call profile-web,pprof/agent/sequential/grpc)

.PHONY: profile/agent/sequential/rest
profile/agent/sequential/rest:
	$(call profile-web,pprof/agent/sequential/rest)

.PHONY: profile/ngtd/stream
profile/ngtd/stream:
	$(call profile-web,pprof/ngtd/stream)

.PHONY: profile/ngtd/sequential/grpc
profile/ngtd/sequential/grpc:
	$(call profile-web,pprof/ngtd/sequential/grpc)

.PHONY: profile/ngtd/sequential/rest
profile/ngtd/sequential/rest:
	$(call profile-web,pprof/ngtd/sequential/rest)

.PHONY: metrics
## calculate all metrics
metrics: \
	metrics/agent

.PHONY: metrics/agent
## calculate agent metrics
metrics/agent: \
	metrics/agent/core/ngt

.PHONY: metrics/agent/core/ngt
## calculate agent/core/ngt metrics
metrics/agent/core/ngt: $(ROOTDIR)/metrics.gob

$(ROOTDIR)/metrics.gob:
	go test -v --timeout=1h ./hack/benchmark/e2e/agent/core/ngt/... -output=$(ROOTDIR)/metrics.gob

.PHONY: metrics/chart
## create metrics chart
metrics/chart: $(ROOTDIR)/assets/image/metrics.svg

$(ROOTDIR)/assets/image/metrics.svg: $(ROOTDIR)/metrics.gob
	go run ./hack/tools/metrics/main.go -title "Recall-QPS" -x Recall -y QPS -width 960 -height 720 -input=$(ROOTDIR)/metrics.gob -output=$(ROOTDIR)/assets/image/metrics.svg

.PHONY: bench/kill
## kill all benchmark processes
bench/kill:
	ps aux  \
		| grep go  \
		| grep -v nvim \
		| grep -v tmux \
		| grep -v gopls \
		| grep -v "rg go" \
		| grep -v "grep go" \
		| awk '{print $1}' \
		| xargs kill -9
