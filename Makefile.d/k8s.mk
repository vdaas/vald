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
.PHONY: k8s/vald/deploy
## deploy vald sample cluster to k8s
k8s/vald/deploy: \
	k8s/external/mysql/deploy \
	k8s/external/redis/deploy
	kubectl apply -f k8s/metrics/metrics-server
	kubectl apply -f k8s/manager/backup/mysql
	kubectl apply -f k8s/manager/compressor
	kubectl apply -f k8s/agent/ngt
	kubectl apply -f k8s/discoverer/k8s
	kubectl apply -f k8s/meta/redis
	kubectl apply -f k8s/gateway/vald

.PHONY: k8s/vald/remove
## remove vald sample cluster from k8s
k8s/vald/remove: \
	k8s/external/mysql/remove \
	k8s/external/redis/remove
	kubectl delete -f k8s/gateway/vald
	kubectl delete -f k8s/manager/backup/mysql
	kubectl delete -f k8s/manager/compressor
	kubectl delete -f k8s/meta/redis
	kubectl delete -f k8s/discoverer/k8s
	kubectl delete -f k8s/agent/ngt
	kubectl delete -f k8s/metrics/metrics-server

.PHONY: k8s/external/mysql/deploy
## deploy mysql to k8s
k8s/external/mysql/deploy:
	kubectl create configmap mysql-config --from-file=$(ROOTDIR)/assets/ddl/mysql/ddl.sql
	kubectl apply -f k8s/external/mysql

.PHONY: k8s/external/mysql/remove
## remove mysql from k8s
k8s/external/mysql/remove:
	kubectl delete -f k8s/external/mysql
	kubectl delete configmap mysql-config

.PHONY: k8s/external/mysql/initialize
## initialize mysql on k8s
k8s/external/mysql/initialize:
	-kubectl delete -f k8s/jobs/db/initialize/mysql
	-kubectl delete configmap mysql-config
	kubectl create configmap mysql-config --from-file=$(ROOTDIR)/assets/ddl/mysql/ddl.sql
	kubectl apply -f k8s/external/mysql/secret.yaml
	kubectl apply -f k8s/jobs/db/initialize/mysql

.PHONY: k8s/external/redis/deploy
## deploy redis to k8s
k8s/external/redis/deploy:
	kubectl apply -f k8s/external/redis

.PHONY: k8s/external/redis/remove
## remove redis from k8s
k8s/external/redis/remove:
	kubectl delete -f k8s/external/redis

.PHONY: k8s/external/redis/initialize
## initialize redis on k8s
k8s/external/redis/initialize:
	-kubectl delete -f k8s/jobs/db/initialize/redis
	kubectl apply -f k8s/external/redis/secret.yaml
	kubectl apply -f k8s/jobs/db/initialize/redis

.PHONY: k8s/external/cassandra/deploy
## deploy cassandra to k8s
k8s/external/cassandra/deploy:
	kubectl create configmap cassandra-initdb --from-file=$(ROOTDIR)/assets/ddl/cassandra/init.cql
	kubectl apply -f k8s/external/cassandra

.PHONY: k8s/external/cassandra/remove
## remove cassandra from k8s
k8s/external/cassandra/remove:
	kubectl delete -f k8s/external/cassandra
	kubectl delete configmap cassandra-initdb

.PHONY: k8s/external/cassandra/initialize
## initialize cassandra on k8s
k8s/external/cassandra/initialize:
	-kubectl delete -f k8s/jobs/db/initialize/cassandra
	-kubectl delete configmap cassandra-initdb
	kubectl create configmap cassandra-initdb --from-file=$(ROOTDIR)/assets/ddl/cassandra/init.cql
	kubectl apply -f k8s/jobs/db/initialize/cassandra

.PHONY: k8s/linkerd/deploy
## deploy linkerd to k8s
k8s/linkerd/deploy:
	linkerd install | kubectl apply -f -
	kubectl annotate namespace \
		kubectl config get-contexts --no-headers \
			"$(kubectl config current-context)"  \
			| awk "{print \$5}" | sed "s/^$/default/" \
		linkerd.io/inject=enabled
