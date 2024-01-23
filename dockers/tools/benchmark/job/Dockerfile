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

ARG GO_VERSION=latest
ARG ZLIB_VERSION
ARG HDF5_VERSION
ARG DISTROLESS_IMAGE=gcr.io/distroless/static
ARG DISTROLESS_IMAGE_TAG=nonroot
ARG UPX_OPTIONS=-9
ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>"

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION} AS golang

FROM --platform=${BUILDPLATFORM} ubuntu:devel AS builder

ARG UPX_OPTIONS
ARG ZLIB_VERSION
ARG HDF5_VERSION
ARG TARGETARCH

ENV GO111MODULE on
ENV DEBIAN_FRONTEND noninteractive
ENV INITRD No
ENV LANG en_US.UTF-8
ENV GOROOT /opt/go
ENV GOPATH /go
ENV PATH ${PATH}:${GOROOT}/bin:${GOPATH}/bin
ENV ORG vdaas
ENV REPO vald
ENV APP_NAME job
ENV PKG tools/benchmark/${APP_NAME}
ENV BUILD_DIR=/usr/local
ENV LIB_DIR=/usr/local/lib
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:${LIB_DIR}:/lib

# skipcq: DOK-DL3008, DOK-DL3003
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    build-essential \
    curl \
    upx \
    g++ \
    git \
    && ldconfig \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p ${LIB_DIR} \
    && curl -sOL https://github.com/madler/zlib/releases/download/v${ZLIB_VERSION}/zlib-${ZLIB_VERSION}.tar.gz \
    && mkdir -p zlib \
    && tar -xzvf zlib-${ZLIB_VERSION}.tar.gz -C zlib --strip-components 1 \
    && cd zlib \
    && ./configure --prefix=${LIB_DIR} --static \
    && make \
    && make test \
    && make install \
    && cd / \
    && mkdir -p hdf5 \
    && curl -sOL https://github.com/HDFGroup/hdf5/releases/download/${HDF5_VERSION}/${HDF5_VERSION}.tar.gz \
    && tar -xzvf ${HDF5_VERSION}.tar.gz -C hdf5 --strip-components 2 \
    && cd hdf5 \
    && ./configure --enable-build-mode=production --enable-static-exec --disable-shared --prefix=${BUILD_DIR} --with-zlib=${BUILD_DIR}/include,${LIB_DIR} LDFLAGS="-Wl,-rpath,${LIB_DIR}" \
    && make check \
    && make install \
    && ldconfig

COPY --from=golang /usr/local/go $GOROOT
RUN mkdir -p "$GOPATH/src"

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/Makefile.d
COPY Makefile.d .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
COPY Makefile .
COPY .git .
COPY go.mod .
COPY go.sum .

RUN --mount=type=cache,target="${GOPATH}/pkg",id="go-build-${TARGETARCH}" \
    --mount=type=cache,target="${HOME}/.cache/go-build",id="go-build-${TARGETARCH}" \
    make go/download

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/internal
COPY internal .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/apis/grpc
COPY apis/grpc .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/pkg/${PKG}
COPY pkg/${PKG} .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd/${PKG}
COPY cmd/${PKG} .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/versions
COPY versions .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
RUN --mount=type=cache,target="${GOPATH}/pkg",id="go-build-${TARGETARCH}" \
    --mount=type=cache,target="${HOME}/.cache/go-build",id="go-build-${TARGETARCH}" \
    make REPO=${ORG} NAME=${REPO} cmd/${PKG}/${APP_NAME} \
    && mv "cmd/${PKG}/${APP_NAME}" "/usr/bin/${APP_NAME}"

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd/${PKG}
RUN cp sample.yaml /tmp/config.yaml

FROM --platform=${BUILDPLATFORM} ${DISTROLESS_IMAGE}:${DISTROLESS_IMAGE_TAG}
LABEL maintainer="${MAINTAINER}"

ENV APP_NAME job

COPY --from=builder /usr/bin/${APP_NAME} /go/bin/${APP_NAME}
COPY --from=builder /tmp/config.yaml /etc/server/config.yaml

USER nonroot:nonroot

ENTRYPOINT ["/go/bin/job"]
