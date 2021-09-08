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

ARG GO_VERSION=latest
ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>"

FROM golang:${GO_VERSION} AS golang

FROM ubuntu:devel AS builder
ARG PROTOBUF_VERSION=3.14.0
ARG GOLANGCI_LINT_VERSION=v1.42.1
ARG REVIEWDOG_VERSION=v0.11.0

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

WORKDIR /tmp
RUN curl -OL "https://github.com/google/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip" \
    && unzip protoc-${PROTOBUF_VERSION}-linux-x86_64.zip -d protoc3 \
    && mv protoc3/bin/* /usr/local/bin/ \
    && mv protoc3/include/* /usr/local/include/ \
    && rm -rf protoc-${PROTOBUF_VERSION}-linux-x86_64.zip protoc3

WORKDIR ${GOPATH}/src/github.com/vdaas/vald

COPY go.mod .
COPY go.sum .
COPY Makefile .
COPY Makefile.d Makefile.d
COPY apis/proto apis/proto
COPY versions versions
COPY hack/go.mod.default hack/go.mod.default

RUN make deps

RUN make ngt/install
RUN make tensorflow/install

RUN make helm/install
RUN make helm-docs/install
RUN make kind/install
RUN make valdcli/install
RUN make kubelinter/install
RUN make yq/install

RUN make tparse/install

RUN curl -L "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" -o /usr/local/bin/kubectl \
    && chmod a+x /usr/local/bin/kubectl

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}
RUN curl -sSfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh \
    | sh -s -- -b $(go env GOPATH)/bin ${REVIEWDOG_VERSION}

ENV PATH=$PATH:$GOPATH/bin

RUN rm -rf "${GOPATH}/src/github.com/vdaas/vald/*"
