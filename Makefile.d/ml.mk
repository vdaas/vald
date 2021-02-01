#
# Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
.PHONY: ml/models/clean
ml/models/clean:
	rm -rf hack/ml/models

.PHONY: ml/models/tensorflow/init
ml/models/tensorflow/init:
	mkdir -p hack/ml/models/tensorflow

.PHONY: ml/models/tensorflow/download
## download tensorflow model
ml/models/tensorflow/download: \
	ml/models/clean \
	ml/models/tensorflow/init \
	ml/models/tensorflow/download/bert \
	ml/models/tensorflow/download/insightface

.PHONY: ml/models/tensorflow/download/bert
ml/models/tensorflow/download/bert:
	curl -LO https://github.com/vdaas/ml/raw/master/tensorflow/bert.tar.gz
	tar -xvf bert.tar.gz -C hack/ml/models/tensorflow
	rm bert.tar.gz

.PHONY: ml/models/tensorflow/download/insightface
ml/models/tensorflow/download/insightface:
	curl -LO https://github.com/vdaas/ml/raw/master/tensorflow/insightface.tar.gz
	tar -xvf insightface.tar.gz -C hack/ml/models/tensorflow
	rm insightface.tar.gz
