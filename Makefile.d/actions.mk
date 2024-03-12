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

ACTIONS_LIST := $(eval ACTIONS_LIST := $(shell grep --include=\*.{yml,yaml} -rohE 'uses: .+' .github/ | awk '!/\.github/ && !seen[$$0]++ {sub(/@.*/, "", $$0); print substr($$0, 7)}'))$(ACTIONS_LIST)

.PHONY: list/actions
## show variation of external actions
list/actions:
	@echo $(ACTIONS_LIST)

.PHONY: update/actions
# update github actions version
update/actions:
	$(call update-github-actions, $(ACTIONS_LIST))
