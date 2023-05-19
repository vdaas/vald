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

ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>"

# skipcq: DOK-DL3026
FROM mcr.microsoft.com/vscode/devcontainers/go:1 AS base
LABEL maintainer "${MAINTAINER}"

# skipcq: DOK-DL3008
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    cmake \
    curl \
    g++ \
    gawk \
    gcc \
    git \
    graphviz \
    jq \
    libhdf5-dev \
    liblapack-dev \
    libomp-dev \
    libopenblas-dev \
    nodejs \
    npm \
    sed \
    unzip \
    zip \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR ${GOPATH}/src/github.com/vdaas/vald

COPY Makefile .
COPY Makefile.d Makefile.d
COPY versions versions

RUN make ngt/install \
    && make helm/install \
    && make helm-docs/install \
    && make valdcli/install \
    && make yq/install \
    && make golangci-lint/install \
    && make reviewdog/install \
    && make protobuf/install

# k3d
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN curl -sS https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
# kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
# k9s
# k9s installs the binary to the current user which is root for devcontainer
# so change the current user to vscode to user afterwards
USER vscode
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN curl -sS https://webinstall.dev/k9s | bash

WORKDIR ${GOPATH}
