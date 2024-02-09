# syntax = docker/dockerfile:latest
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

ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>"

# skipcq: DOK-DL3026
# we should stick to official devcontainers as a base image because it is well tested to be used with the vscode devcontainer extension.
FROM --platform=${BUILDPLATFORM} mcr.microsoft.com/vscode/devcontainers/base:debian
LABEL maintainer="${MAINTAINER}"

ARG GO_VERSION
ARG TARGETARCH
ARG TARGETOS

ENV GOROOT /opt/go
ENV GOPATH /go
ENV PATH ${PATH}:${GOROOT}/bin:${GOPATH}/bin:/usr/local/bin

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
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

WORKDIR /tmp
RUN BIN_NAME="go" \
    && TAR_NAME="${BIN_NAME}${GO_VERSION}.${TARGETOS}-${TARGETARCH}.tar.gz" \
    && curl -fsSLO "https://go.dev/dl/${TAR_NAME}" \
    && tar zxf "${TAR_NAME}" \
    && rm "${TAR_NAME}" \
    && mv ${BIN_NAME} ${GOROOT} \
    && ${BIN_NAME} version \
    && mkdir -p "${GOPATH}/src"

WORKDIR ${GOPATH}/src/github.com/vdaas/vald
COPY Makefile .
COPY Makefile.d Makefile.d
COPY versions versions
COPY rust rust
COPY hack/go.mod.default hack/go.mod.default
COPY example/client/go.mod.default example/client/go.mod.default

RUN --mount=type=cache,target="${GOPATH}/pkg",id="go-build-${TARGETARCH}" \
    --mount=type=cache,target="${HOME}/.cache/go-build",id="go-build-${TARGETARCH}" \
    make GOARCH=${TARGETARCH} GOOS=${TARGETOS} deps GO_CLEAN_DEPS=false \
    && make ngt/install \
    && make faiss/install \
    && make helm/install \
    && make helm-docs/install \
    && make valdcli/install \
    && make yq/install \
    && make GOARCH=${TARGETARCH} GOOS=${TARGETOS} golangci-lint/install \
    && make reviewdog/install \
    && make kubectl/install \
    && make k3d/install \
    && make buf/install \
    && make k9s/install \
    && make minikube/install \
    && make stern/install \
    && make telepresence/install \
    && echo "installing golang vscode extension dependencies" \
    && make gopls/install \
    && make gotests/install \
    && make gomodifytags/install \
    && make impl/install \
    && make delve/install \
    && make staticcheck/install
