#
# Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

.PHONY: e2e
## run e2e
e2e:
	$(call run-e2e-crud-test,-run TestE2EStandardCRUD)

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
