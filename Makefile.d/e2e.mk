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

.PHONY: e2e
## run e2e
e2e:
	$(call run-e2e-crud-test,-run TestE2EStandardCRUD)

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

