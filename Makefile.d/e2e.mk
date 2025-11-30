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

.PHONY: e2e
## run e2e
e2e:
	$(call run-e2e-crud-test,-run TestE2EStandardCRUD)

.PHONY: e2e/v2
## run e2e
e2e/v2:
	$(call run-v2-e2e-crud-test,-run TestE2EStrategy)

.PHONY: e2e/faiss
## run e2e/faiss
e2e/faiss:
	#$(call run-e2e-crud-faiss-test,-run TestE2EInsertOnly)
	#$(call run-e2e-crud-faiss-test,-run TestE2ESearchOnly)
	#$(call run-e2e-crud-faiss-test,-run TestE2EUpdateOnly)
	#$(call run-e2e-crud-faiss-test,-run TestE2ERemoveOnly)
	#$(call run-e2e-crud-faiss-test,-run TestE2EInsertAndSearch)
	$(call run-e2e-crud-faiss-test,-run TestE2EStandardCRUD)

.PHONY: e2e/skip
## run e2e with skip exists operation
e2e/skip:
	$(call run-e2e-crud-test,-run TestE2ECRUDWithSkipStrictExistCheck)

.PHONY: e2e/multi
## run e2e multiple apis
e2e/multi:
	$(call run-e2e-multi-crud-test,-run TestE2EMultiAPIs)

.PHONY: e2e/insert
## run insert e2e
e2e/insert:
	$(call run-e2e-crud-test,-run TestE2EInsertOnly)

.PHONY: e2e/update
## run update e2e
e2e/update:
	$(call run-e2e-crud-test,-run TestE2EUpdateOnly)

.PHONY: e2e/search
## run search e2e
e2e/search:
	$(call run-e2e-crud-test,-run TestE2ESearchOnly)

.PHONY: e2e/linearsearch
## run linearsearch e2e
e2e/linearsearch:
	$(call run-e2e-crud-test,-run TestE2ELinearSearchOnly)

.PHONY: e2e/upsert
## run upsert e2e
e2e/upsert:
	$(call run-e2e-crud-test,-run TestE2EUpsertOnly)

.PHONY: e2e/remove
## run remove e2e
e2e/remove:
	$(call run-e2e-crud-test,-run TestE2ERemoveOnly)

.PHONY: e2e/remove/timestamp
## run removeByTimestamp e2e
e2e/remove/timestamp:
	$(call run-e2e-crud-test,-run TestE2ERemoveByTimestampOnly)

.PHONY: e2e/insert/search
## run insert and search e2e
e2e/insert/search:
	$(call run-e2e-crud-test,-run TestE2EInsertAndSearch)

.PHONY: e2e/index/job/correction
## run index correction job e2e
e2e/index/job/correction:
	$(call run-e2e-crud-test,-run TestE2EIndexJobCorrection)

.PHONY: e2e/readreplica
## run readreplica e2e
e2e/readreplica:
	$(call run-e2e-crud-test,-run TestE2EReadReplica)

.PHONY: e2e/maxdim
## run e2e/maxdim
e2e/maxdim:
	$(call run-e2e-max-dim-test)

.PHONY: e2e/sidecar
## run e2e with sidecar operation
e2e/sidecar:
	$(call run-e2e-sidecar-test,-run TestE2EForSidecar)

.PHONY: e2e/actions/run/stream/crud
## run GitHub Actions E2E test (Stream CRUD)
e2e/actions/run/stream/crud: \
	hack/benchmark/assets/dataset/$(E2E_DATASET_NAME) \
	k3d/restart
	kubectl wait -n kube-system --for=condition=Available deployment/metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	sleep 2
	kubectl wait -n kube-system --for=condition=Ready pod -l k8s-app=metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl wait -n kube-system --for=condition=ContainersReady pod -l k8s-app=metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	$(MAKE) k8s/vald/deploy \
	HELM_VALUES=$(ROOTDIR)/.github/helm/values/values-lb.yaml
	sleep 3
	kubectl wait --for=condition=Ready pod -l "app=$(LB_GATEWAY_IMAGE)" --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl wait --for=condition=ContainersReady pod -l "app=$(LB_GATEWAY_IMAGE)" --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl get pods
	pod_name=$$(kubectl get pods --selector="app=$(LB_GATEWAY_IMAGE)" | tail -1 | awk '{print $$1}'); \
	echo $$pod_name; \
	$(MAKE) E2E_TARGET_POD_NAME=$$pod_name e2e
	$(MAKE) k8s/vald/delete
	$(MAKE) k3d/delete

.PHONY: e2e/actions/run/job
## run GitHub Actions E2E test (jobs)
e2e/actions/run/job: \
	hack/benchmark/assets/dataset/$(E2E_DATASET_NAME) \
	k3d/restart
	kubectl wait -n kube-system --for=condition=Available deployment/metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	sleep 2
	kubectl wait -n kube-system --for=condition=Ready pod -l k8s-app=metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl wait -n kube-system --for=condition=ContainersReady pod -l k8s-app=metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	$(MAKE) k8s/vald/deploy \
	HELM_VALUES=$(ROOTDIR)/.github/helm/values/values-correction.yaml
	sleep 3
	kubectl wait --for=condition=Ready pod -l "app=$(LB_GATEWAY_IMAGE)" --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl wait --for=condition=ContainersReady pod -l "app=$(LB_GATEWAY_IMAGE)" --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl get pods
	pod_name=$$(kubectl get pods --selector="app=$(LB_GATEWAY_IMAGE)" | tail -1 | awk '{print $$1}'); \
	echo $$pod_name; \
	$(MAKE) E2E_TARGET_POD_NAME=$$pod_name e2e/index/job/correction
	$(MAKE) k8s/vald/delete
	$(MAKE) k3d/delete

.PHONY: e2e/actions/run/readreplica
## run GitHub Actions E2E test (Stream CRUD with read replica )
e2e/actions/run/readreplica: \
	hack/benchmark/assets/dataset/$(E2E_DATASET_NAME) \
	kind/vs/restart
	kubectl wait -n kube-system --for=condition=Available deployment/metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	sleep 2
	kubectl wait -n kube-system --for=condition=Ready pod -l app.kubernetes.io/name=metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl wait -n kube-system --for=condition=ContainersReady pod -l app.kubernetes.io/name=metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)

	$(MAKE) k8s/vald-readreplica/deploy \
	VERSION=$(VERSION) \
	HELM_VALUES=$(ROOTDIR)/.github/helm/values/values-readreplica.yaml
	sleep 3
	kubectl wait --for=condition=Ready pod -l "app=$(LB_GATEWAY_IMAGE)" --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl wait --for=condition=ContainersReady pod -l "app=$(LB_GATEWAY_IMAGE)" --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl get pods
	pod_name=$$(kubectl get pods --selector="app=$(LB_GATEWAY_IMAGE)" | tail -1 | awk '{print $$1}'); \
	echo $$pod_name; \
	$(MAKE) E2E_TIMEOUT=30m E2E_CONFIG=$(ROOTDIR)/.github/e2e/readreplica.yaml e2e/v2
	# $(MAKE) E2E_TARGET_POD_NAME=$$pod_name e2e/readreplica
	# $(MAKE) k8s/vald/delete
	# $(MAKE) kind/vs/stop

.PHONY: e2e/actions/run/stream/crud/skip
## run GitHub Actions E2E test (Stream CRUD with SkipExistsCheck = true)
e2e/actions/run/stream/crud/skip: \
	hack/benchmark/assets/dataset/$(E2E_DATASET_NAME) \
	k3d/restart
	kubectl wait -n kube-system --for=condition=Available deployment/metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	sleep 2
	kubectl wait -n kube-system --for=condition=Ready pod -l k8s-app=metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl wait -n kube-system --for=condition=ContainersReady pod -l k8s-app=metrics-server --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	$(MAKE) k8s/vald/deploy \
	HELM_VALUES=$(ROOTDIR)/.github/helm/values/values-lb.yaml
	sleep 3
	kubectl wait --for=condition=Ready pod -l "app=$(LB_GATEWAY_IMAGE)" --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl wait --for=condition=ContainersReady pod -l "app=$(LB_GATEWAY_IMAGE)" --timeout=$(E2E_WAIT_FOR_START_TIMEOUT)
	kubectl get pods
	pod_name=$$(kubectl get pods --selector="app=$(LB_GATEWAY_IMAGE)" | tail -1 | awk '{print $$1}'); \
	echo $$pod_name; \
	$(MAKE) E2E_TARGET_POD_NAME=$$pod_name e2e/skip
	$(MAKE) k8s/vald/delete
	$(MAKE) k3d/delete

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
	$(MAKE) E2E_CONFIG_NAME=unary_crud.yaml \
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