#
# Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
.PHONY: git/config/init
## add git configs required for development
git/config/init:
	git config commit.template ".commit_template"
	git config core.fileMode false

.PHONY: git/hooks/init
## add configs for registering pre-defined git hooks
git/hooks/init:
	ln -sf $(ROOTDIR)/hack/git/hooks/pre-commit $(ROOTDIR)/.git/hooks/pre-commit
	chmod a+x $(ROOTDIR)/.git/hooks/pre-commit
