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

.PHONY: valdcli/install
## install valdcli to deps dir
valdcli/install: libs/bin/valdcli

ifeq ($(UNAME),Darwin)
libs/bin/valdcli:
	mkdir -p libs/bin
	curl -LO https://github.com/rinx/vald-client-clj/releases/download/$(VALDCLI_VERSION)/valdcli-macos.zip
	unzip valdcli-macos.zip
	rm -f valdcli-macos.zip
	mv valdcli libs/bin/valdcli
else
libs/bin/valdcli:
	mkdir -p libs/bin
	curl -LO https://github.com/rinx/vald-client-clj/releases/download/$(VALDCLI_VERSION)/valdcli-linux-static.zip
	unzip valdcli-linux-static.zip
	rm -f valdcli-linux-static.zip
	mv valdcli libs/bin/valdcli
endif

.PHONY: valdcli/xpanes/insert
## insert randomized vectors using valdcli and xpanes
valdcli/xpanes/insert: libs/bin/valdcli
	xpanes -c "$(ROOTDIR)/libs/bin/valdcli rand-vecs -n $(NUMBER) -d $(DIMENSION) --with-ids | $(ROOTDIR)/libs/bin/valdcli -h $(HOST) -p $(PORT) stream-insert --elapsed-time" $$(seq 1 $(NUMPANES))

.PHONY: valdcli/xpanes/search
## search randomized vectors using valdcli and xpanes
valdcli/xpanes/search: libs/bin/valdcli
	xpanes -c "$(ROOTDIR)/libs/bin/valdcli rand-vecs -n $(NUMBER) -d $(DIMENSION) | $(ROOTDIR)/libs/bin/valdcli -h $(HOST) -p $(PORT) stream-search --elapsed-time" $$(seq 1 $(NUMPANES))
