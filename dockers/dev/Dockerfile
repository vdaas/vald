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

ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>"

FROM ubuntu:devel AS builder

# skipcq: DOK-DL3008
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    unzip \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /tmp

# skipcq: DOK-DL3026
FROM mcr.microsoft.com/vscode/devcontainers/go:1 AS base
LABEL maintainer "${MAINTAINER}"

# skipcq: DOK-DL3008
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

RUN make ngt/install \
    && make helm/install \
    && make helm-docs/install \
    && make kind/install \
    && make valdcli/install \
    && make yq/install \
    && make golangci-lint/install \
    && make reviewdog/install \
    && make kubectl/install \
    && make protobuf/install

WORKDIR ${GOPATH}
