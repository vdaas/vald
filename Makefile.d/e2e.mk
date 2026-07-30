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

.PHONY: e2e/v2
## run e2e
e2e/v2:
	$(call run-v2-e2e-crud-test,-run TestE2EStrategy)

# One benchmark iteration already performs an execution's full configured
# request load (num/parallelism/qps), so a fixed small iteration count is the
# meaningful default; raise it (e.g. 3x) to let benchstat compute variance.
# Each execution is looped independently at >1x, so state-dependent steps
# (e.g. create_index right after insert) can legitimately fail on later
# iterations unless the scenario's expect also lists the status codes such a
# re-execution returns (e.g. failedprecondition for an empty uncommitted queue).
E2E_BENCH_TIME ?= 1x

.PHONY: e2e/v2/bench
## run e2e scenario as Go benchmarks (BenchmarkE2EStrategy, one execution pass per iteration)
e2e/v2/bench:
	$(call run-v2-e2e-crud-test,-run '^$$' -bench BenchmarkE2EStrategy -benchtime $(E2E_BENCH_TIME))

.PHONY: e2e/v2/operator
## run e2e/v2 for vald-operator
e2e/v2/operator:
	$(MAKE) E2E_CONFIG=$(E2E_CONFIG_DIR)/operator.yaml \
		E2E_MANIFEST_PATH=$(ROOTDIR)/k8s/operator/vald/samples/controller_v1_valdoperatorrelease.yaml \
		e2e/v2

# E2E_MAX_DIM_BIT must be 1 ~ 32 (bit < 1 || maxBit < bit is rejected), the
# same range the removed v1 tests/e2e/performance/max_vector_dim_test.go used
# to enforce.
ifneq ($(shell test $(E2E_MAX_DIM_BIT) -ge 1 -a $(E2E_MAX_DIM_BIT) -le 32 2>/dev/null && echo ok),ok)
$(error E2E_MAX_DIM_BIT must be between 1 and 32, got $(E2E_MAX_DIM_BIT))
endif

# E2E_MAX_DIM is 2^E2E_MAX_DIM_BIT, except at the uint32 boundary (bit 32)
# where it is capped to math.MaxUint32, mirroring .github/workflows/e2e-max-dim.yaml.
ifeq ($(E2E_MAX_DIM_BIT),32)
E2E_MAX_DIM := $(shell echo $$(( (1 << 32) - 1 )))
else
E2E_MAX_DIM := $(shell echo $$(( 1 << $(E2E_MAX_DIM_BIT) )))
endif

.PHONY: e2e/v2/maxdim
## run e2e/v2/maxdim (config-based: insert & search a single vector of dimension 2^E2E_MAX_DIM_BIT)
e2e/v2/maxdim:
	$(call run-v2-e2e-max-dim-test,-run TestE2EStrategy)

.PHONY: e2e/v2/actions/run/unary/crud
## run GitHub Actions E2E/V2 test (Unary CRUD)
e2e/v2/actions/run/unary/crud: \
	hack/benchmark/assets/dataset/$(E2E_DATASET_NAME) \
	k3d/restart
	kubectl wait -n kube-system --for=condition=Available deployment/metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	sleep 2
	kubectl wait -n kube-system --for=condition=Ready pod -l app.kubernetes.io/name=metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl wait -n kube-system --for=condition=ContainersReady pod -l app.kubernetes.io/name=metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	$(MAKE) k8s/vald/deploy \
	VERSION=$(VERSION) \
	HELM_VALUES=$(ROOTDIR)/.github/helm/values/values-lb.yaml
	sleep 3
	kubectl wait --for=condition=Ready pod -l "app=$(LB_GATEWAY_IMAGE)" --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl wait --for=condition=ContainersReady pod -l "app=$(LB_GATEWAY_IMAGE)" --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl get pods
	$(MAKE) E2E_CONFIG="$(E2E_CONFIG_DIR)/unary_crud.yaml" \
		E2E_TIMEOUT=30m \
		E2E_PARALLELISM="4" \
		E2E_INSERT_COUNT="10000" \
		E2E_EXPECTED_INDEX="30000" \
		E2E_QPS="30" \
		E2E_SEARCH_COUNT="10" \
		E2E_UPDATE_COUNT="100" \
		E2E_BULK_SIZE="10" \
		e2e/v2
	$(MAKE) k8s/vald/delete
	$(MAKE) k3d/delete

.PHONY: e2e/v2/actions/run/operator
## run GitHub Actions E2E/V2 test (vald-operator)
e2e/v2/actions/run/operator: \
	k3d/restart
	@docker buildx version > /dev/null 2>&1 || $(MAKE) docker-buildx/install
	$(MAKE) docker/build/operator
	$(K3D_COMMAND) image import $(CRORG)/$(OPERATOR_IMAGE):$(TAG) -c $(K3D_CLUSTER_NAME)
	$(MAKE) k8s/operator/vald/deploy
	$(MAKE) E2E_TIMEOUT=15m \
		E2E_TARGET_NAMESPACE=default \
		E2E_TARGET_NAME=vald-operator \
		e2e/v2/operator
	$(MAKE) k8s/operator/vald/delete
	$(MAKE) k3d/delete

.PHONY: e2e/v2/actions/run/faiss
## run GitHub Actions E2E/V2 test (FAISS agent backend, CRUD)
e2e/v2/actions/run/faiss: \
	hack/benchmark/assets/dataset/$(E2E_DATASET_NAME) \
	k3d/restart
	kubectl wait -n kube-system --for=condition=Available deployment/metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	sleep 2
	kubectl wait -n kube-system --for=condition=Ready pod -l app.kubernetes.io/name=metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl wait -n kube-system --for=condition=ContainersReady pod -l app.kubernetes.io/name=metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	$(MAKE) k8s/vald/deploy \
	VERSION=$(VERSION) \
	HELM_VALUES=$(ROOTDIR)/.github/helm/values/values-faiss.yaml \
	HELM_EXTRA_OPTIONS="--set agent.image.repository=$(CRORG)/$(AGENT_FAISS_IMAGE)"
	sleep 3
	kubectl wait --for=condition=Ready pod -l "app=$(LB_GATEWAY_IMAGE)" --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl wait --for=condition=ContainersReady pod -l "app=$(LB_GATEWAY_IMAGE)" --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl get pods
	$(MAKE) E2E_CONFIG="$(E2E_CONFIG_DIR)/faiss_crud.yaml" \
		E2E_TIMEOUT=30m \
		E2E_PARALLELISM="4" \
		E2E_INSERT_COUNT="10000" \
		E2E_EXPECTED_INDEX="30000" \
		E2E_QPS="30" \
		E2E_SEARCH_COUNT="10" \
		E2E_UPDATE_COUNT="100" \
		E2E_BULK_SIZE="10" \
		e2e/v2
	$(MAKE) k8s/vald/delete
	$(MAKE) k3d/delete
