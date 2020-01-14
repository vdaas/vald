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
$(BENCH_DATASETS): $(BENCH_DATASET_MD5S) $(BENCH_DATASET_HDF5_DIR)
	@$(call green, "downloading datasets for benchmark...")
	curl -fsSL -o $@ http://ann-benchmarks.com/$(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,%.hdf5,$@)
	(cd $(BENCH_DATASET_BASE_DIR); \
	    md5sum -c $(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,$(BENCH_DATASET_MD5_DIR_NAME)/%.md5,$@) || \
	    (rm -f $(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,$(BENCH_DATASET_HDF5_DIR_NAME)/%.hdf5,$@) && exit 1))

$(BENCH_DATASET_HDF5_DIR):
	$(call mkdir, $@)
	$(call rm, -rf, $@/*)

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
bench/core/ngt/sequential:
	$(call bench-pprof,pprof/core/ngt,core,NGTSequential,sequential,\
    		./hack/benchmark/core/ngt/ngt_bench_test.go \
    		-dataset=$(DATASET_ARGS))

.PHONY: bench/core/ngt/parallel
## run benchmark for NGT core parallel methods
bench/core/ngt/parallel:
	$(call bench-pprof,pprof/core/ngt,core,NGTParallel,parallel,\
    		./hack/benchmark/core/ngt/ngt_bench_test.go \
    		-dataset=$(DATASET_ARGS))

.PHONY: bench/core/gongt
## run benchmark for gongt core
bench/core/gongt: \
	bench/core/gongt/sequential \
    bench/core/gongt/parallel

.PHONY: bench/core/gongt/sequential
## run benchmark for gongt core sequential methods
bench/core/gongt/sequential:
	$(call bench-pprof,pprof/core/gongt,core,GoNGTSequential,sequential,\
    		./hack/benchmark/core/ngt/gongt_bench_test.go \
    		-dataset=$(DATASET_ARGS))

.PHONY: bench/core/gongt/parallel
## run benchmark for gongt core parallel methods
bench/core/gongt/parallel:
	$(call bench-pprof,pprof/core/gongt,core,GoNGTParallel,parallel,\
    		./hack/benchmark/core/ngt/gongt_bench_test.go \
    		-dataset=$(DATASET_ARGS))

.PHONY: bench/agent
## run benchmarks for vald agent
bench/agent: \
	bench/agent/stream \
	bench/agent/sequential/grpc \
	bench/agent/sequential/rest

.PHONY: bench/agent/stream
## run benchmark for agent gRPC stream
bench/agent/stream: \
	ngt/install
	$(call bench-pprof,pprof/agent/ngt,agent,gRPCStream,stream,\
		./hack/benchmark/e2e/agent/ngt/ngt_bench_test.go \
		 -dataset=$(DATASET_ARGS) -address=$(ADDRESS_ARGS))

.PHONY: bench/agent/sequential/grpc
## run benchmark for agent gRPC sequential
bench/agent/sequential/grpc: \
	ngt/install
	$(call bench-pprof,pprof/agent/ngt,agent,gRPCSequential,sequential-grpc,\
		./hack/benchmark/e2e/agent/ngt/ngt_bench_test.go \
		 -dataset=$(DATASET_ARGS) -address=$(ADDRESS_ARGS))

.PHONY: bench/agent/sequential/rest
## run benchmark for agent REST
bench/agent/sequential/rest: \
	ngt/install
	$(call bench-pprof,pprof/agent/ngt,agent,RESTSequential,sequential-rest,\
		./hack/benchmark/e2e/agent/ngt/ngt_bench_test.go \
		 -dataset=$(DATASET_ARGS) -address=$(ADDRESS_ARGS))

.PHONY: bench/ngtd
## run benchmarks for NGTD
bench/ngtd: \
	bench/ngtd/stream \
	bench/ngtd/sequential/grpc \
	bench/ngtd/sequential/rest

.PHONY: bench/ngtd/stream
## run benchmark for NGTD gRPC stream
bench/ngtd/stream: \
	ngt/install
	$(call bench-pprof,pprof/external/ngtd,ngtd,gRPCStream,stream,\
		./hack/benchmark/e2e/external/ngtd/ngtd_bench_test.go \
		 -dataset=$(DATASET_ARGS) -address=$(ADDRESS_ARGS))

.PHONY: bench/ngtd/sequential/grpc
## run benchmark for NGTD gRPC sequential
bench/ngtd/sequential/grpc: \
	ngt/install
	$(call bench-pprof,pprof/external/ngtd,ngtd,gRPCSequential,sequential-grpc,\
		./hack/benchmark/e2e/external/ngtd/ngtd_bench_test.go \
		 -dataset=$(DATASET_ARGS) -address=$(ADDRESS_ARGS))

.PHONY: bench/ngtd/sequential/rest
## run benchmark for NGTD REST stream
bench/ngtd/sequential/rest: \
	ngt/install
	$(call bench-pprof,pprof/external/ngtd,ngtd,RESTSequential,sequential-rest,\
		./hack/benchmark/e2e/external/ngtd/ngtd_bench_test.go \
		 -dataset=$(DATASET_ARGS) -address=$(ADDRESS_ARGS))

.PHONY: bench/gateway
## run benchmarks for gateway
bench/gateway: \
	bench/gateway/sequential

.PHONY: bench/gateway/sequential
## run benchmark for gateway sequential
bench/gateway/sequential: \
	ngt/install
	$(call bench-pprof,pprof/gateway/vald,vald,Sequential,sequential,\
		./hack/benchmark/e2e/gateway/vald/vald_bench_test.go \
		 -dataset=$(DATASET_ARGS) -address=$(ADDRESS_ARGS))

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
	$(call profile-web,pprof/agent/ngt,agent,stream,":6061",":6062",":6063")

.PHONY: profile/agent/sequential/grpc
profile/agent/sequential/grpc:
	$(call profile-web,pprof/agent/ngt,agent,sequential-grpc,":6061",":6062",":6063")

.PHONY: profile/agent/sequential/rest
profile/agent/sequential/rest:
	$(call profile-web,pprof/agent/ngt,agent,sequential-rest,":6061",":6062",":6063")

.PHONY: profile/ngtd/stream
profile/ngtd/stream:
	$(call profile-web,pprof/external/ngtd,ngtd,stream,":6061",":6062",":6063")

.PHONY: profile/ngtd/sequential/grpc
profile/ngtd/sequential/grpc:
	$(call profile-web,pprof/external/ngtd,ngtd,sequential-grpc,":6061",":6062",":6063")

.PHONY: profile/ngtd/sequential/rest
profile/ngtd/sequential/rest:
	$(call profile-web,pprof/external/ngtd,ngtd,sequential-rest,":6061",":6062",":6063")

.PHONY: metrics
## calculate all metrics
metrics: \
	metrics/agent

.PHONY: metrics/agent
## calculate agent metrics
metrics/agent: \
	metrics/agent/ngt

.PHONY: metrics/agent/ngt
## calculate agent/ngt metrics
metrics/agent/ngt: $(ROOTDIR)/metrics.gob

$(ROOTDIR)/metrics.gob:
	go test -v --timeout=1h ./hack/benchmark/e2e/agent/ngt/... -output=$(ROOTDIR)/metrics.gob

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
