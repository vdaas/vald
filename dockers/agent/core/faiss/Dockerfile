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

<<<<<<< HEAD
# ARG DISTROLESS_IMAGE=gcr.io/distroless/static
# ARG DISTROLESS_IMAGE_TAG=nonroot
ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>"

FROM --platform=${TARGETPLATFORM} ubuntu:devel AS builder

ARG TARGETARCH
ARG TARGETOS
ARG GO_VERSION
=======
ARG GO_VERSION=latest
ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>"

FROM golang:${GO_VERSION} AS golang

FROM ubuntu:devel AS builder
>>>>>>> feature/agent/qbg

ENV GO111MODULE on
ENV DEBIAN_FRONTEND noninteractive
ENV INITRD No
ENV LANG en_US.UTF-8
ENV GOROOT /opt/go
ENV GOPATH /go
ENV PATH ${PATH}:${GOROOT}/bin:${GOPATH}/bin
ENV ORG vdaas
ENV REPO vald
ENV PKG agent/core/faiss
ENV PKG_INTERNAL agent/internal
ENV APP_NAME faiss

# skipcq: DOK-DL3008
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    cmake \
    curl \
<<<<<<< HEAD
    g++ \
    gcc \
    git \
    intel-mkl \
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
COPY go.mod .
COPY go.sum .

RUN --mount=type=cache,target="${GOPATH}/pkg",id="go-build-${TARGETARCH}" \
    --mount=type=cache,target="${HOME}/.cache/go-build",id="go-build-${TARGETARCH}" \
    make GO_VERSION=${GO_VERSION} go/install \
    && make go/download
=======
    gcc \
    git \
    g++ \
    intel-mkl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY --from=golang /usr/local/go $GOROOT
RUN mkdir $GOPATH

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}

COPY go.mod .
COPY go.sum .

RUN go mod download
>>>>>>> feature/agent/qbg

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/internal
COPY internal .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/apis/grpc
COPY apis/grpc .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/pkg/${PKG}
COPY pkg/${PKG} .
<<<<<<< HEAD
=======
# copy ngt/service/kvs and ngt/service/vqueue
WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/pkg/agent/core/ngt/service
COPY pkg/agent/core/ngt/service/kvs ./kvs
COPY pkg/agent/core/ngt/service/vqueue ./vqueue
>>>>>>> feature/agent/qbg

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/pkg/${PKG_INTERNAL}
COPY pkg/${PKG_INTERNAL} .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd/${PKG}
COPY cmd/${PKG} .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/versions
COPY versions .

<<<<<<< HEAD
WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
RUN update-alternatives --set libblas.so-x86_64-linux-gnu /usr/lib/x86_64-linux-gnu/libmkl_rt.so \
     && make faiss/install

RUN --mount=type=cache,target="${GOPATH}/pkg",id="go-build-${TARGETARCH}" \
    --mount=type=cache,target="${HOME}/.cache/go-build",id="go-build-${TARGETARCH}" \
    make GOARCH=${TARGETARCH} GOOS=${TARGETOS} REPO=${ORG} NAME=${REPO} cmd/${PKG}/${APP_NAME} \
=======
WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/Makefile.d
COPY Makefile.d .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
COPY Makefile .
RUN update-alternatives --set libblas.so-x86_64-linux-gnu /usr/lib/x86_64-linux-gnu/libmkl_rt.so \
     && make faiss/install

COPY .git .

RUN make REPO=${ORG} NAME=${REPO} cmd/${PKG}/${APP_NAME} \
>>>>>>> feature/agent/qbg
    && mv "cmd/${PKG}/${APP_NAME}" "/usr/bin/${APP_NAME}"

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd/${PKG}
RUN cp sample.yaml /tmp/config.yaml

<<<<<<< HEAD
FROM --platform=${TARGETPLATFORM} ubuntu:devel
LABEL maintainer="${MAINTAINER}"
=======
FROM ubuntu:devel
LABEL maintainer "${MAINTAINER}"
>>>>>>> feature/agent/qbg

ENV APP_NAME faiss

COPY --from=builder /usr/bin/${APP_NAME} /go/bin/${APP_NAME}
COPY --from=builder /tmp/config.yaml /etc/server/config.yaml

COPY --from=builder /usr/local/lib/libfaiss.so /usr/local/lib/libfaiss.so
COPY --from=builder /lib/x86_64-linux-gnu/libstdc++.so.6 /lib/x86_64-linux-gnu/libstdc++.so.6
COPY --from=builder /lib/x86_64-linux-gnu/libgcc_s.so.1 /lib/x86_64-linux-gnu/libgcc_s.so.1
COPY --from=builder /lib/x86_64-linux-gnu/libc.so.6 /lib/x86_64-linux-gnu/libc.so.6
COPY --from=builder /lib/x86_64-linux-gnu/libmkl_intel_lp64.so /lib/x86_64-linux-gnu/libmkl_intel_lp64.so
COPY --from=builder /lib/x86_64-linux-gnu/libmkl_sequential.so /lib/x86_64-linux-gnu/libmkl_sequential.so
COPY --from=builder /lib/x86_64-linux-gnu/libmkl_core.so /lib/x86_64-linux-gnu/libmkl_core.so
COPY --from=builder /lib/x86_64-linux-gnu/libgomp.so.1 /lib/x86_64-linux-gnu/libgomp.so.1
COPY --from=builder /lib/x86_64-linux-gnu/libm.so.6 /lib/x86_64-linux-gnu/libm.so.6
COPY --from=builder /lib/x86_64-linux-gnu/libdl.so.2 /lib/x86_64-linux-gnu/libdl.so.2
COPY --from=builder /lib/x86_64-linux-gnu/libpthread.so.0 /lib/x86_64-linux-gnu/libpthread.so.0
COPY --from=builder /lib/x86_64-linux-gnu/libmkl_avx2.so /lib/x86_64-linux-gnu/libmkl_avx2.so
RUN ldconfig -v

<<<<<<< HEAD
USER nonroot:nonroot

=======
>>>>>>> feature/agent/qbg
ENTRYPOINT ["/go/bin/faiss"]
