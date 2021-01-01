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

ARG BASE_TAG=latest
ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>

FROM vdaas/vald-ci-container:${BASE_TAG} AS vald

FROM mcr.microsoft.com/vscode/devcontainers/go:1 AS base
LABEL maintainer "${MAINTAINER}"

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    cmake \
    curl \
    libhdf5-dev \
    g++ \
    gawk \
    gcc \
    git \
    graphviz \
    nodejs \
    npm \
    jq \
    sed \
    zip \
    unzip \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR ${GOPATH}/src/github.com/vdaas/vald

COPY Makefile .
COPY Makefile.d Makefile.d
COPY versions versions

RUN make ngt/install
RUN make tensorflow/install

RUN make helm/install
RUN make helm-docs/install
RUN make kind/install
RUN make valdcli/install

COPY --from=vald /usr/local/bin/protoc     /usr/local/bin/protoc
COPY --from=vald /usr/local/include/google /usr/local/include/google

COPY --from=vald /usr/local/bin/kubectl      /usr/local/bin/kubectl
COPY --from=vald /usr/local/bin/yq           /usr/local/bin/yq
COPY --from=vald ${GOPATH}/bin/golangci-lint ${GOPATH}/bin/golangci-lint
COPY --from=vald ${GOPATH}/bin/reviewdog     ${GOPATH}/bin/reviewdog

WORKDIR ${GOPATH}
