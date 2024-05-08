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

ARG DISTROLESS_IMAGE=gcr.io/distroless/cc-debian12
ARG DISTROLESS_IMAGE_TAG=nonroot
ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>"

FROM --platform=${TARGETPLATFORM} ubuntu:devel AS builder

ARG TARGETARCH
ARG TARGETOS
ARG RUST_VERSION

ENV RUST_HOME /usr/local/lib/rust
ENV RUSTUP_HOME ${RUST_HOME}/rustup
ENV CARGO_HOME ${RUST_HOME}/cargo
ENV DEBIAN_FRONTEND noninteractive
ENV INITRD No
ENV LANG en_US.UTF-8
ENV PATH ${PATH}:${RUSTUP_HOME}/bin:${CARGO_HOME}/bin
ENV ORG vdaas
ENV REPO vald
ENV PKG agent/core/ngt
ENV APP_NAME ngt

# skipcq: DOK-DL3008
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    cmake \
    curl \
    g++ \
    gcc \
    git \
    liblapack-dev \
    libomp-dev \
    libopenblas-dev \
    unzip \
    upx \
    && ldconfig \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/Makefile.d
COPY Makefile.d .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
COPY Makefile .
COPY .git .

RUN make RUST_VERSION=${RUST_VERSION} rust/install

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd
COPY cmd .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/versions
COPY versions .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/rust
COPY rust .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
RUN make ngt/install

RUN make rust/target/release/agent \
    && mv ./rust/target/release/agent "/usr/bin/${APP_NAME}"

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd/${PKG}
RUN cp sample.yaml /tmp/config.yaml

FROM --platform=${TARGETPLATFORM} ${DISTROLESS_IMAGE}:${DISTROLESS_IMAGE_TAG}
LABEL maintainer="${MAINTAINER}"

ENV APP_NAME ngt

COPY --from=builder /usr/bin/${APP_NAME} /usr/bin/${APP_NAME}
COPY --from=builder /tmp/config.yaml /etc/server/config.yaml

USER nonroot:nonroot

ENTRYPOINT ["/usr/bin/ngt"]
