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
.PHONY: kind/start
## start kind (kubernetes in docker) cluster
kind/start:
<<<<<<< HEAD
	kind create cluster --name $(NAME)
	export KUBECONFIG="$(kind get kubeconfig-path --name=$(NAME))"
=======
	kind create cluster --name vald
	export KUBECONFIG="$(shell kind get kubeconfig-path --name="kind")"
>>>>>>> 8ce8c6d46716604d8c644fa1cb44550255f22e49

.PHONY: kind/login
## login kind (kubernetes in docker) cluster
kind/login:
<<<<<<< HEAD
	export KUBECONFIG="$(kind get kubeconfig-path --name=$(NAME))"
=======
	export KUBECONFIG="$(shell kind get kubeconfig-path --name="kind")"
>>>>>>> 8ce8c6d46716604d8c644fa1cb44550255f22e49

.PHONY: kind/stop
## stop kind (kubernetes in docker) cluster
kind/stop:
	kind delete cluster --name $(NAME)
