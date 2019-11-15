#
# Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
$(BENCH_DATASET_ARGS): $(BENCH_DATASET_MD5S)
	@$(call green, "downloading datasets for benchmark...")
	curl -fsSL -o $@ http://ann-benchmarks.com/$(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,%.hdf5,$@)
	(cd $(BENCH_DATASET_BASE_DIR); \
	    md5sum -c $(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,$(BENCH_DATASET_MD5_DIR_NAME)/%.md5,$@) || \
	    (rm -f $(patsubst $(BENCH_DATASET_HDF5_DIR)/%.hdf5,$(BENCH_DATASET_HDF5_DIR_NAME)/%.hdf5,$@) && exit 1))

.PHONY: bench/datasets
## fetch datasets for benchmark
bench/datasets: $(BENCH_DATASET_ARGS)

.PHONY: bench/datasets/clean
## clean datasets for benchmark
bench/datasets/clean:
	rm -rf $(BENCH_DATASETS)

.PHONY: bench/create-index
bench/create-index:
	$(MAKE) -C ./hack/core/ngt create

.PHONY: bench/core
## run benchmark for NGT core
bench/core: bench/create-index
	$(MAKE) -C ./hach/core/ngt bench

.PHONY: bench/core/lite
## run lite benchmark for NGT core
bench/core/lite: bench/create-index
	$(MAKE) -C ./hach/core/ngt bench-lite

.PHONY: bench/core/clean
## clean results for NGT core benchmark
bench/core/clean:
	$(MAKE) -C ./hack/core/ngt clean

.PHONY: bench/e2e
## run e2e benchmark
bench/e2e:
	$(MAKE) -C ./hack/e2e/benchmark bench

.PHONY: bench
## run benchmarks
bench: \
	bench/agent/stream \
	bench/agent/sequential/grpc \
	bench/agent/sequential/rest

.PHONY: bench/agent
## run benchmarks for agent
bench/agent: \
	bench/agent/stream \
	bench/agent/sequential/grpc \
	bench/agent/sequential/rest

.PHONY: bench/agent/stream
bench/agent/stream: \
	ngt/install
	$(call bench-pprof,pprof/agent/ngt,agent,gRPCStream,stream,\
		./hack/e2e/benchmark/agent/ngt/ngt_bench_test.go \
		 -dataset=$(DATASET_ARGS) -address=$(ADDRESS_ARGS))

.PHONY: bench/agent/sequential/grpc
bench/agent/sequential/grpc: \
	ngt/install
	$(call bench-pprof,pprof/agent/ngt,agent,gRPCSequential,sequential-grpc,\
		./hack/e2e/benchmark/agent/ngt/ngt_bench_test.go \
		 -dataset=$(DATASET_ARGS) -address=$(ADDRESS_ARGS))

.PHONY: bench/agent/sequential/rest
bench/agent/sequential/rest: \
	ngt/install
	$(call bench-pprof,pprof/agent/ngt,agent,RESTSequential,sequential-rest,\
		./hack/e2e/benchmark/agent/ngt/ngt_bench_test.go \
		 -dataset=$(DATASET_ARGS) -address=$(ADDRESS_ARGS))

.PHONY: bench/ngtd
## run benchmarks for NGTD
bench/ngtd: \
	bench/ngtd/stream \
	bench/ngtd/sequential/grpc \
	bench/ngtd/sequential/rest

.PHONY: bench/ngtd/stream
bench/ngtd/stream: \
	ngt/install
	$(call bench-pprof,pprof/external/ngtd,ngtd,gRPCStream,stream,\
		./hack/e2e/benchmark/external/ngtd/ngtd_bench_test.go \
		 -dataset=$(DATASET_ARGS) -address=$(ADDRESS_ARGS))

.PHONY: bench/ngtd/sequential/grpc
bench/ngtd/sequential/grpc: \
	ngt/install
	$(call bench-pprof,pprof/external/ngtd,ngtd,gRPCSequential,sequential-grpc,\
		./hack/e2e/benchmark/external/ngtd/ngtd_bench_test.go \
		 -dataset=$(DATASET_ARGS) -address=$(ADDRESS_ARGS))

.PHONY: bench/ngtd/sequential/rest
bench/ngtd/sequential/rest: \
	ngt/install
	$(call bench-pprof,pprof/external/ngtd,ngtd,RESTSequential,sequential-rest,\
		./hack/e2e/benchmark/external/ngtd/ngtd_bench_test.go \
		 -dataset=$(DATASET_ARGS) -address=$(ADDRESS_ARGS))

.PHONY: bench/gateway/sequential
bench/gateway/sequential: \
	ngt/install
	$(call bench-pprof,pprof/gateway/vald,vald,Sequential,sequential,\
		./hack/e2e/benchmark/gateway/vald/vald_bench_test.go \
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
