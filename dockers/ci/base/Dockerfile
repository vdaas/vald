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

FROM --platform=${TARGETPLATFORM} ubuntu:devel AS builder
LABEL maintainer="${MAINTAINER}"

ARG GO_VERSION
ARG TARGETARCH
ARG TARGETOS

ENV GO111MODULE on
ENV DEBIAN_FRONTEND noninteractive
ENV INITRD No
ENV LANG en_US.UTF-8
ENV GOROOT /opt/go
ENV GOPATH /go
ENV PATH ${PATH}:${GOROOT}/bin:${GOPATH}/bin:/usr/local/bin

RUN sysctl -w net.ipv6.conf.all.disable_ipv6=1 \
    && sysctl -w net.ipv6.conf.default.disable_ipv6=1 \
    && sysctl -w net.ipv6.conf.lo.disable_ipv6=1

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
# skipcq: DOK-DL3008
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    cmake \
    curl \
    g++ \
    gawk \
    gcc \
    git \
    gnupg2 \
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

ENV ROOTDIR .

COPY Makefile .
COPY Makefile.d Makefile.d
COPY apis/proto apis/proto
COPY example/client/go.mod.default example/client/go.mod.default
COPY go.mod .
COPY go.sum .
COPY hack/go.mod.default hack/go.mod.default
COPY rust rust
COPY versions versions

RUN --mount=type=cache,target="${GOPATH}/pkg",id="go-build-${TARGETARCH}" \
    --mount=type=cache,target="${HOME}/.cache/go-build",id="go-build-${TARGETARCH}" \
    make GO_VERSION=${GO_VERSION} go/install \
    && make GOARCH=${TARGETARCH} GOOS=${TARGETOS} deps ROOTDIR=${ROOTDIR} GO_CLEAN_DEPS=false \
    && make ngt/install \
    && make faiss/install \
    && make helm/install \
    && make helm-docs/install \
    && make kind/install \
    && make valdcli/install \
    && make kubelinter/install \
    && make yq/install \
    && make tparse/install \
    && make GOARCH=${TARGETARCH} GOOS=${TARGETOS} gotestfmt/install \
    && make GOARCH=${TARGETARCH} GOOS=${TARGETOS} golangci-lint/install \
    && make reviewdog/install \
    && make kubectl/install \
    && make k3d/install \
    && make buf/install \
    && make k9s/install \
    && rm -rf "${GOPATH}/src/github.com/vdaas/vald/*"
