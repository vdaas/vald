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

ARG GO_VERSION=latest
ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>"

FROM golang:${GO_VERSION} AS golang

FROM ubuntu:focal AS builder

LABEL maintainer "${MAINTAINER}"

ENV GO111MODULE on
ENV DEBIAN_FRONTEND noninteractive
ENV INITRD No
ENV LANG en_US.UTF-8
ENV GOROOT /opt/go
ENV GOPATH /go
ENV PATH ${PATH}:${GOROOT}/bin:${GOPATH}/bin

COPY --from=golang /usr/local/go $GOROOT
RUN mkdir -p $GOPATH/src

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
    nodejs \
    npm \
    sed \
    unzip \
    zip \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR ${GOPATH}/src/github.com/vdaas/vald

COPY go.mod .
COPY go.sum .
COPY Makefile .
COPY Makefile.d Makefile.d
COPY apis/proto apis/proto
COPY versions versions
COPY hack/go.mod.default hack/go.mod.default

RUN make deps

RUN make ngt/install \
    && make helm/install \
    && make helm-docs/install \
    && make kind/install \
    && make valdcli/install \
    && make kubelinter/install \
    && make yq/install \
    && make tparse/install \
    && make golangci-lint/install \
    && make reviewdog/install \
    && make kubectl/install \
    && make protobuf/install

ENV PATH=$PATH:$GOPATH/bin

RUN rm -rf "${GOPATH}/src/github.com/vdaas/vald/*"
