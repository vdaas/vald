#
# Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

JAEGER_OPERATOR_WAIT_DURATION := 0
MIRROR01_NAMESPACE = vald-01
MIRROR02_NAMESPACE = vald-02
MIRROR03_NAMESPACE = vald-03
MIRROR_APP_NAME    = vald-mirror-gateway

.PHONY: k8s/manifest/clean
## clean k8s manifests
k8s/manifest/clean:
	rm -rf \
		k8s/agent \
		k8s/discoverer \
		k8s/gateway \
		k8s/manager \
		k8s/index

.PHONY: k8s/manifest/update
## update k8s manifests using helm templates
k8s/manifest/update: \
	k8s/manifest/clean
	helm template \
		--values $(HELM_VALUES) \
		$(HELM_EXTRA_OPTIONS) \
		--output-dir $(TEMP_DIR) \
		charts/vald
	mkdir -p k8s/gateway
	mkdir -p k8s/manager
	mkdir -p k8s/index/job
	mkdir -p k8s/index/job/readreplica
	mv $(TEMP_DIR)/vald/templates/agent k8s/agent
	mv $(TEMP_DIR)/vald/templates/discoverer k8s/discoverer
	mv $(TEMP_DIR)/vald/templates/gateway k8s/gateway
	mv $(TEMP_DIR)/vald/templates/manager/index k8s/manager/index
	mv $(TEMP_DIR)/vald/templates/index/job/correction k8s/index/job/correction
	mv $(TEMP_DIR)/vald/templates/index/job/creation k8s/index/job/creation
	mv $(TEMP_DIR)/vald/templates/index/job/save k8s/index/job/save
	mv $(TEMP_DIR)/vald/templates/index/job/readreplica/rotate k8s/index/job/readreplica/rotate
	rm -rf $(TEMP_DIR)

.PHONY: k8s/manifest/helm-operator/clean
## clean k8s manifests for helm-operator
k8s/manifest/helm-operator/clean:
	rm -rf \
	    k8s/operator/helm

.PHONY: k8s/manifest/helm-operator/update
## update k8s manifests for helm-operatorusing helm templates
k8s/manifest/helm-operator/update: \
	k8s/manifest/helm-operator/clean
	helm template \
	    --output-dir $(TEMP_DIR) \
	    charts/vald-helm-operator
	mkdir -p k8s/operator
	mv $(TEMP_DIR)/vald-helm-operator/templates k8s/operator/helm
	rm -rf $(TEMP_DIR)
	cp -r charts/vald-helm-operator/crds k8s/operator/helm/crds

.PHONY: k8s/manifest/benchmark-operator/clean
## clean k8s manifests for benchmark-operator
k8s/manifest/benchmark-operator/clean:
	rm -rf \
	    k8s/tools/benchmark/operator

.PHONY: k8s/manifest/benchmark-operator/update
## update k8s manifests for benchmark-operator using helm templates
k8s/manifest/benchmark-operator/update: \
	k8s/manifest/benchmark-operator/clean
	helm template \
	    --output-dir $(TEMP_DIR) \
	    charts/vald-benchmark-operator
	mkdir -p k8s/tools/benchmark
	mv $(TEMP_DIR)/vald-benchmark-operator/templates k8s/tools/benchmark/operator
	rm -rf $(TEMP_DIR)
	cp -r charts/vald-benchmark-operator/crds k8s/tools/benchmark/operator/crds

.PHONY: k8s/manifest/readreplica/clean
## clean k8s manifests for readreplica
k8s/manifest/readreplica/clean:
	rm -rf \
	    k8s/readreplica

.PHONY: k8s/manifest/readreplica/update
## update k8s manifests for readreplica using helm templates
k8s/manifest/readreplica/update: \
	k8s/manifest/readreplica/clean
	helm template \
	    --output-dir $(TEMP_DIR) \
	    charts/vald-readreplica
	mv $(TEMP_DIR)/vald-readreplica/templates k8s/readreplica
	rm -rf $(TEMP_DIR)

.PHONY: k8s/vald/deploy
## deploy vald sample cluster to k8s
k8s/vald/deploy:
	helm template \
	    --values $(HELM_VALUES) \
	    --set defaults.image.tag=$(VERSION) \
	    --set agent.image.repository=$(CRORG)/$(AGENT_NGT_IMAGE) \
	    --set agent.sidecar.image.repository=$(CRORG)/$(AGENT_SIDECAR_IMAGE) \
	    --set discoverer.image.repository=$(CRORG)/$(DISCOVERER_IMAGE) \
	    --set gateway.filter.image.repository=$(CRORG)/$(FILTER_GATEWAY_IMAGE) \
	    --set gateway.lb.image.repository=$(CRORG)/$(LB_GATEWAY_IMAGE) \
	    --set gateway.mirror.image.repository=$(CRORG)/$(MIRROR_GATEWAY_IMAGE) \
	    --set manager.index.image.repository=$(CRORG)/$(MANAGER_INDEX_IMAGE) \
	    --set manager.index.creator.image.repository=$(CRORG)/$(INDEX_CREATION_IMAGE) \
	    --set manager.index.saver.image.repository=$(CRORG)/$(INDEX_SAVE_IMAGE) \
	    $(HELM_EXTRA_OPTIONS) \
        --include-crds \
	    --output-dir $(TEMP_DIR) \
	    charts/vald
	@echo "Permitting error because there's some cases nothing to apply"
	kubectl apply -f $(TEMP_DIR)/vald/templates/manager/index || true
	kubectl apply -f $(TEMP_DIR)/vald/templates/agent || true
	kubectl apply -f $(TEMP_DIR)/vald/templates/agent/ngt || true
	kubectl apply -f $(TEMP_DIR)/vald/templates/agent/readreplica || true
	kubectl apply -f $(TEMP_DIR)/vald/templates/discoverer || true
	kubectl apply -f $(TEMP_DIR)/vald/templates/gateway || true
	kubectl apply -f $(TEMP_DIR)/vald/templates/gateway/lb || true
	kubectl apply -f $(TEMP_DIR)/vald/crds || true
	kubectl apply -f $(TEMP_DIR)/vald/templates/gateway/mirror || true
	kubectl apply -f $(TEMP_DIR)/vald/templates/index/job/correction || true
	kubectl apply -f $(TEMP_DIR)/vald/templates/index/job/creation || true
	kubectl apply -f $(TEMP_DIR)/vald/templates/index/job/save || true
	kubectl apply -f $(TEMP_DIR)/vald/templates/index/job/readreplica/rotate || true
	rm -rf $(TEMP_DIR)
	kubectl get pods -o jsonpath="{.items[*].spec.containers[*].image}" | tr " " "\n"

.PHONY: k8s/vald/delete
## delete vald sample cluster from k8s
k8s/vald/delete:
	helm template \
	    --values $(HELM_VALUES) \
	    --set defaults.image.tag=$(VERSION) \
	    --set agent.image.repository=$(CRORG)/$(AGENT_NGT_IMAGE) \
	    --set agent.sidecar.image.repository=$(CRORG)/$(AGENT_SIDECAR_IMAGE) \
	    --set discoverer.image.repository=$(CRORG)/$(DISCOVERER_IMAGE) \
	    --set gateway.filter.image.repository=$(CRORG)/$(FILTER_GATEWAY_IMAGE) \
	    --set gateway.lb.image.repository=$(CRORG)/$(LB_GATEWAY_IMAGE) \
	    --set gateway.mirror.image.repository=$(CRORG)/$(MIRROR_GATEWAY_IMAGE) \
	    --set manager.index.image.repository=$(CRORG)/$(MANAGER_INDEX_IMAGE) \
        --include-crds \
	    --output-dir $(TEMP_DIR) \
	    charts/vald
	kubectl delete -f $(TEMP_DIR)/vald/templates/gateway/mirror
	kubectl delete -f $(TEMP_DIR)/vald/templates/index/job/readreplica/rotate
	kubectl delete -f $(TEMP_DIR)/vald/templates/index/job/save
	kubectl delete -f $(TEMP_DIR)/vald/templates/index/job/creation
	kubectl delete -f $(TEMP_DIR)/vald/templates/index/job/correction
	kubectl delete -f $(TEMP_DIR)/vald/templates/index/job/creation
	kubectl delete -f $(TEMP_DIR)/vald/templates/index/job/save
	kubectl delete -f $(TEMP_DIR)/vald/templates/gateway
	kubectl delete -f $(TEMP_DIR)/vald/templates/gateway/lb
	kubectl delete -f $(TEMP_DIR)/vald/templates/manager/index
	kubectl delete -f $(TEMP_DIR)/vald/templates/discoverer
	kubectl delete -f $(TEMP_DIR)/vald/templates/agent/readreplica || true
	kubectl delete -f $(TEMP_DIR)/vald/templates/agent/ngt || true
	kubectl delete -f $(TEMP_DIR)/vald/templates/agent
	kubectl delete -f $(TEMP_DIR)/vald/crds
	rm -rf $(TEMP_DIR)

.PHONY: k8s/multi/vald/deploy
## deploy multiple vald sample clusters to k8s
k8s/multi/vald/deploy:
	-@kubectl create ns $(MIRROR01_NAMESPACE)
	-@kubectl create ns $(MIRROR02_NAMESPACE)
	-@kubectl create ns $(MIRROR03_NAMESPACE)
	helm install vald-cluster-01 charts/vald \
		-f $(ROOTDIR)/charts/vald/values/multi-vald/dev-vald-with-mirror.yaml \
		-f $(ROOTDIR)/charts/vald/values/multi-vald/dev-vald-01.yaml \
	    -n $(MIRROR01_NAMESPACE)
	helm install vald-cluster-02 charts/vald \
		-f $(ROOTDIR)/charts/vald/values/multi-vald/dev-vald-with-mirror.yaml \
		-f $(ROOTDIR)/charts/vald/values/multi-vald/dev-vald-02.yaml \
	    -n $(MIRROR02_NAMESPACE)
	helm install vald-cluster-03 charts/vald \
		-f $(ROOTDIR)/charts/vald/values/multi-vald/dev-vald-with-mirror.yaml \
		-f $(ROOTDIR)/charts/vald/values/multi-vald/dev-vald-03.yaml \
		-n $(MIRROR03_NAMESPACE)
	kubectl wait --for=condition=ready pod -l app=$(MIRROR_APP_NAME) --timeout=120s -n $(MIRROR01_NAMESPACE)
	kubectl wait --for=condition=ready pod -l app=$(MIRROR_APP_NAME) --timeout=120s -n $(MIRROR02_NAMESPACE)
	kubectl wait --for=condition=ready pod -l app=$(MIRROR_APP_NAME) --timeout=120s -n $(MIRROR03_NAMESPACE)
	kubectl apply -f $(ROOTDIR)/charts/vald/values/multi-vald/mirror-target.yaml \
		-n $(MIRROR03_NAMESPACE)

.PHONY: k8s/multi/vald/delete
## delete multiple vald sample clusters to k8s
k8s/multi/vald/delete:
	helm uninstall vald-cluster-01 -n vald-01
	helm uninstall vald-cluster-02 -n vald-02
	helm uninstall vald-cluster-03 -n vald-03
	-@kubectl delete ns vald-01 vald-02 vald-03

.PHONY: k8s/vald-helm-operator/deploy
## deploy vald-helm-operator to k8s
k8s/vald-helm-operator/deploy:
	helm template \
	    --output-dir $(TEMP_DIR) \
	    --set image.tag=$(VERSION) \
	    $(HELM_EXTRA_OPTIONS) \
	    --include-crds \
	    charts/vald-helm-operator
	kubectl create -f $(TEMP_DIR)/vald-helm-operator/crds/valdrelease.yaml
	kubectl create -f $(TEMP_DIR)/vald-helm-operator/crds/valdhelmoperatorrelease.yaml
	kubectl apply -f $(TEMP_DIR)/vald-helm-operator/templates
	sleep 2
	kubectl wait --for=condition=ready pod -l name=vald-helm-operator --timeout=600s

.PHONY: k8s/vald-helm-operator/delete
## delete vald-helm-operator from k8s
k8s/vald-helm-operator/delete:
	helm template \
	    --output-dir $(TEMP_DIR) \
	    --set image.tag=$(VERSION) \
	    --include-crds \
	    charts/vald-helm-operator
	kubectl delete -f $(TEMP_DIR)/vald-helm-operator/templates
	kubectl wait --for=delete pod -l name=vald-helm-operator --timeout=600s
	kubectl delete -f $(TEMP_DIR)/vald-helm-operator/crds
	rm -rf $(TEMP_DIR)

.PHONY: k8s/vald-readreplica/deploy
## deploy vald-readreplica to k8s
k8s/vald-readreplica/deploy:
	helm template \
	    --values $(HELM_VALUES) \
	    --set defaults.image.tag=$(VERSION) \
	    --set agent.image.repository=$(CRORG)/$(AGENT_NGT_IMAGE) \
	    --set agent.sidecar.image.repository=$(CRORG)/$(AGENT_SIDECAR_IMAGE) \
	    --set discoverer.image.repository=$(CRORG)/$(DISCOVERER_IMAGE) \
	    --set gateway.filter.image.repository=$(CRORG)/$(FILTER_GATEWAY_IMAGE) \
	    --set gateway.lb.image.repository=$(CRORG)/$(LB_GATEWAY_IMAGE) \
	    --set manager.index.image.repository=$(CRORG)/$(MANAGER_INDEX_IMAGE) \
	    --set manager.index.creator.image.repository=$(CRORG)/$(INDEX_CREATION_IMAGE) \
	    --set manager.index.saver.image.repository=$(CRORG)/$(INDEX_SAVE_IMAGE) \
	    $(HELM_EXTRA_OPTIONS) \
	    --output-dir $(TEMP_DIR) \
	    charts/vald-readreplica
	kubectl apply -f $(TEMP_DIR)/vald-readreplica/templates
	sleep 2
	kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=vald-readreplica --timeout=600s

.PHONY: k8s/vald-readreplica/delete
## delete vald-helm-operator from k8s
k8s/vald-readreplica/delete:
	helm template \
	    --values $(HELM_VALUES) \
	    --set defaults.image.tag=$(VERSION) \
	    --set agent.image.repository=$(CRORG)/$(AGENT_NGT_IMAGE) \
	    --set agent.sidecar.image.repository=$(CRORG)/$(AGENT_SIDECAR_IMAGE) \
	    --set discoverer.image.repository=$(CRORG)/$(DISCOVERER_IMAGE) \
	    --set gateway.filter.image.repository=$(CRORG)/$(FILTER_GATEWAY_IMAGE) \
	    --set gateway.lb.image.repository=$(CRORG)/$(LB_GATEWAY_IMAGE) \
	    --set manager.index.image.repository=$(CRORG)/$(MANAGER_INDEX_IMAGE) \
	    --output-dir $(TEMP_DIR) \
	    charts/vald-readreplica
	kubectl delete -f $(TEMP_DIR)/vald-readreplica/templates
	kubectl wait --for=delete pod -l app.kubernetes.io/name=vald-readreplica --timeout=600s
	rm -rf $(TEMP_DIR)

.PHONY: k8s/vr/deploy
## deploy ValdRelease resource to k8s
k8s/vr/deploy: \
	yq/install \
	k8s/metrics/metrics-server/deploy
	yq eval \
	    '{"apiVersion": "vald.vdaas.org/v1", "kind": "ValdRelease", "metadata":{"name":"vald-cluster"}, "spec": .}' \
	    $(HELM_VALUES) \
	    | kubectl apply -f -

.PHONY: k8s/vr/delete
## delete ValdRelease resource from k8s
k8s/vr/delete: \
	k8s/metrics/metrics-server/delete
	kubectl delete vr vald-cluster

.PHONY: k8s/vald-benchmark-operator/deploy
## deploy vald-benchmark-operator to k8s
k8s/vald-benchmark-operator/deploy:
	helm template \
	    --output-dir $(TEMP_DIR) \
	    --set image.tag=${VERSION} \
	    --include-crds \
	    charts/vald-benchmark-operator
	kubectl create -f $(TEMP_DIR)/vald-benchmark-operator/crds/valdbenchmarkjob.yaml
	kubectl create -f $(TEMP_DIR)/vald-benchmark-operator/crds/valdbenchmarkscenario.yaml
	kubectl create -f $(TEMP_DIR)/vald-benchmark-operator/crds/valdbenchmarkoperatorrelease.yaml
	kubectl apply -f $(TEMP_DIR)/vald-benchmark-operator/templates
	sleep 2
	kubectl wait --for=condition=ready pod -l name=vald-benchmark-operator --timeout=600s

.PHONY: k8s/vald-benchmark-operator/delete
## delete vald-benchmark-operator from k8s
k8s/vald-benchmark-operator/delete:
	helm template \
	    --output-dir $(TEMP_DIR) \
	    --set image.tag=${VERSION} \
	    --include-crds \
	    charts/vald-benchmark-operator
	kubectl delete -f $(TEMP_DIR)/vald-benchmark-operator/templates
	kubectl wait --for=delete pod -l name=vald-benchmark-operator --timeout=600s
	kubectl delete -f $(TEMP_DIR)/vald-benchmark-operator/crds
	rm -rf $(TEMP_DIR)

.PHONY: k8s/external/cert-manager/deploy
## deploy cert-manager
k8s/external/cert-manager/deploy:
	kubectl apply -f https://github.com/jetstack/cert-manager/releases/latest/download/cert-manager.yaml
	sleep $(K8S_SLEEP_DURATION_FOR_WAIT_COMMAND)
	kubectl wait -n cert-manager --for=condition=ready pod -l app=cert-manager --timeout=60s
	kubectl wait -n cert-manager --for=condition=ready pod -l app=cainjector --timeout=60s
	kubectl wait -n cert-manager --for=condition=ready pod -l app=webhook --timeout=60s
	kubectl wait -n cert-manager --for=condition=Available deployment --timeout=60s --all
	sleep 20

.PHONY: k8s/external/cert-manager/delete
## delete cert-manager
k8s/external/cert-manager/delete:
	kubectl delete -f https://github.com/jetstack/cert-manager/releases/latest/download/cert-manager.yaml

.PHONY: k8s/external/minio/deploy
## deploy minio
k8s/external/minio/deploy:
	kubectl apply -f k8s/external/minio/deployment.yaml
	kubectl apply -f k8s/external/minio/svc.yaml
	sleep $(K8S_SLEEP_DURATION_FOR_WAIT_COMMAND)
	kubectl wait --for=condition=ready pod -l app=minio --timeout=600s
	kubectl apply -f k8s/external/minio/mb-job.yaml
	sleep $(K8S_SLEEP_DURATION_FOR_WAIT_COMMAND)
	kubectl wait --for=condition=complete job/minio-make-bucket --timeout=600s

.PHONY: k8s/external/minio/delete
## delete minio
k8s/external/minio/delete:
	kubectl delete -f k8s/external/minio

.PHONY: k8s/metrics/metrics-server/deploy
## deploy metrics-serrver
k8s/metrics/metrics-server/deploy:
	kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
	kubectl patch deployment metrics-server -n kube-system -p '{"spec":{"template":{"spec":{"containers":[{"name":"metrics-server","args":["--cert-dir=/tmp", "--secure-port=4443", "--kubelet-insecure-tls","--kubelet-preferred-address-types=InternalIP"]}]}}}}'
	sleep $(K8S_SLEEP_DURATION_FOR_WAIT_COMMAND)
	# kubectl wait -n kube-system --for=condition=ready pod -l k8s-app=metrics-server --timeout=600s

.PHONY: k8s/metrics/metrics-server/delete
## delete metrics-serrver
k8s/metrics/metrics-server/delete:
	kubectl delete -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

.PHONY: k8s/metrics/prometheus/deploy
## deploy prometheus
k8s/metrics/prometheus/deploy:
	kubectl apply -f k8s/metrics/prometheus

.PHONY: k8s/metrics/prometheus/delete
## delete prometheus
k8s/metrics/prometheus/delete:
	kubectl delete -f k8s/metrics/prometheus

.PHONY: k8s/metrics/prometheus/operator/deploy
## deploy prometheus operator
k8s/metrics/prometheus/operator/deploy:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm install ${PROMETHEUS_RELEASE_NAME} prometheus-community/kube-prometheus-stack --version ${PROMETHEUS_STACK_VERSION} --set grafana.enabled=false

.PHONY: k8s/metrics/prometheus/operator/delete
## delete prometheus operator
k8s/metrics/prometheus/operator/delete:
	helm uninstall ${PROMETHEUS_RELEASE_NAME}

.PHONY: k8s/metrics/grafana/deploy
## deploy grafana
k8s/metrics/grafana/deploy:
	kubectl apply -f k8s/metrics/grafana/dashboards
	kubectl apply -f k8s/metrics/grafana

.PHONY: k8s/metrics/grafana/delete
## delete grafana
k8s/metrics/grafana/delete:
	kubectl delete -f k8s/metrics/grafana/dashboards
	kubectl delete -f k8s/metrics/grafana

.PHONY: k8s/metrics/jaeger/deploy
## deploy jaeger
k8s/metrics/jaeger/deploy:
	helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
	helm install jaeger jaegertracing/jaeger-operator --version $(JAEGER_OPERATOR_VERSION)
	kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=jaeger-operator --timeout=60s
	kubectl wait --for=condition=available deployment/jaeger-jaeger-operator --timeout=60s
	sleep $(JAEGER_OPERATOR_WAIT_DURATION)
	kubectl apply -f k8s/metrics/jaeger/jaeger.yaml

.PHONY: k8s/metrics/jaeger/delete
## delete jaeger
k8s/metrics/jaeger/delete:
	kubectl delete -f k8s/metrics/jaeger
	helm uninstall jaeger

.PHONY: k8s/metrics/loki/deploy
## deploy loki and promtail
k8s/metrics/loki/deploy:
	kubectl apply -f k8s/metrics/loki

.PHONY: k8s/metrics/loki/delete
## delete loki and promtail
k8s/metrics/loki/delete:
	kubectl delete -f k8s/metrics/loki

.PHONY: k8s/metrics/tempo/deploy
## deploy tempo and jaeger-agent
k8s/metrics/tempo/deploy:
	kubectl apply -f k8s/metrics/tempo

.PHONY: k8s/metrics/tempo/delete
## delete tempo and jaeger-agent
k8s/metrics/tempo/delete:
	kubectl delete -f k8s/metrics/tempo

.PHONY: k8s/metrics/profefe/deploy
## deploy profefe
k8s/metrics/profefe/deploy:
	kubectl apply -f k8s/metrics/profefe

.PHONY: k8s/metrics/profefe/delete
## delete profefe
k8s/metrics/profefe/delete:
	kubectl delete -f k8s/metrics/profefe

.PHONY: k8s/metrics/pyroscope/deploy
## deploy pyroscope
k8s/metrics/pyroscope/deploy:
	kubectl apply -k k8s/metrics/pyroscope/base

.PHONY: k8s/metrics/pyroscope/delete
## delete pyroscope
k8s/metrics/pyroscope/delete:
	kubectl delete -k k8s/metrics/pyroscope/base

.PHONY: k8s/metrics/pyroscope/pv/deploy
## deploy pyroscope on persistent volume
k8s/metrics/pyroscope/pv/deploy:
	kubectl apply -k k8s/metrics/pyroscope/overlay

.PHONY: k8s/metrics/pyroscope/pv/delete
## delete pyroscope on persistent volume
k8s/metrics/pyroscope/pv/delete:
	kubectl delete -k k8s/metrics/pyroscope/overlay

.PHONY: k8s/linkerd/deploy
## deploy linkerd to k8s
k8s/linkerd/deploy:
	linkerd check --pre
	linkerd install | kubectl apply -f -
	linkerd check
	kubectl annotate namespace \
		default \
		linkerd.io/inject=enabled

.PHONY: k8s/linkerd/delete
## delete linkerd from k8s
k8s/linkerd/delete:
	linkerd install --ignore-cluster | kubectl delete -f -

.PHONY: k8s/otel/operator/deploy
## deploy opentelemetry operator
k8s/otel/operator/deploy:
	helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
	helm install ${OTEL_OPERATOR_RELEASE_NAME} open-telemetry/opentelemetry-operator --set installCRDs=true --version ${OTEL_OPERATOR_VERSION}
	kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=opentelemetry-operator --timeout=60s
	sleep 10

.PHONY: k8s/otel/operator/delete
## delete opentelemetry operator
k8s/otel/operator/delete:
	helm uninstall ${OTEL_OPERATOR_RELEASE_NAME}

.PHONY: k8s/otel/collector/deploy
## deploy opentelemetry collector
k8s/otel/collector/deploy:
	kubectl apply -f $(ROOTDIR)/k8s/metrics/otel/collector.yaml
	kubectl apply -f $(ROOTDIR)/k8s/metrics/otel/pod-monitor.yaml

.PHONY: k8s/otel/collector/delete
## delete opentelemetry collector
k8s/otel/collector/delete:
	kubectl delete -f $(ROOTDIR)/k8s/metrics/otel/collector.yaml
	kubectl delete -f $(ROOTDIR)/k8s/metrics/otel/pod-monitor.yaml

.PHONY: k8s/monitoring/deploy
## deploy monitoring stack
k8s/monitoring/deploy: \
	k8s/metrics/jaeger/deploy \
	k8s/metrics/prometheus/operator/deploy \
	k8s/metrics/grafana/deploy \
	k8s/otel/operator/deploy \
	k8s/otel/collector/deploy

.PHONY: k8s/monitoring/delete
## delete monitoring stack
k8s/monitoring/delete: \
	k8s/otel/collector/delete \
	k8s/otel/operator/delete \
	k8s/metrics/grafana/delete \
	k8s/metrics/jaeger/delete \
	k8s/metrics/prometheus/operator/delete \

.PHONY: telepresence/install
## install telepresence
telepresence/install: $(BINDIR)/telepresence

$(BINDIR)/telepresence:
	mkdir -p $(BINDIR)
	cd $(TEMP_DIR) \
	    && curl -fsSL "https://app.getambassador.io/download/tel2oss/releases/download/v$(TELEPRESENCE_VERSION)/telepresence-$(OS)-$(subst x86_64,amd64,$(shell echo $(ARCH) | tr '[:upper:]' '[:lower:]'))" -o $(BINDIR)/telepresence \
	    && chmod a+x $(BINDIR)/telepresence

.PHONY: telepresence/swap/agent-ngt
## swap agent-ngt deployment using telepresence
telepresence/swap/agent-ngt:
	$(call telepresence,vald-agent-ngt,vdaas/vald-agent-ngt)

.PHONY: telepresence/swap/agent-faiss
## swap agent-faiss deployment using telepresence
telepresence/swap/agent-faiss:
	$(call telepresence,vald-agent-faiss,vdaas/vald-agent-faiss)

.PHONY: telepresence/swap/discoverer
## swap discoverer deployment using telepresence
telepresence/swap/discoverer:
	$(call telepresence,vald-discoverer,vdaas/vald-discoverer-k8s)

.PHONY: telepresence/swap/manager-index
## swap manager-index deployment using telepresence
telepresence/swap/manager-index:
	$(call telepresence,vald-manager-index,vdaas/vald-manager-index)

.PHONY: telepresence/swap/lb-gateway
## swap lb-gateway deployment using telepresence
telepresence/swap/lb-gateway:
	$(call telepresence,vald-lb-gateway,vdaas/vald-lb-gateway)

.PHONY: kubelinter/install
## install kubelinter
kubelinter/install: $(BINDIR)/kube-linter

$(BINDIR)/kube-linter:
	mkdir -p $(BINDIR)
	cd $(TEMP_DIR) \
	    && curl -L https://github.com/stackrox/kube-linter/releases/download/$(KUBELINTER_VERSION)/kube-linter-$(OS) -o $(BINDIR)/kube-linter \
	    && chmod a+x $(BINDIR)/kube-linter
