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

ARG GO_VERSION=latest
ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>"

FROM golang:${GO_VERSION} AS golang

FROM ubuntu:devel AS builder

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

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    cmake \
    curl \
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

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/internal
COPY internal .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/apis/grpc
COPY apis/grpc .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/pkg/${PKG}
COPY pkg/${PKG} .
# copy ngt/service/kvs and ngt/service/vqueue
WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/pkg/agent/core/ngt/service
COPY pkg/agent/core/ngt/service/kvs ./kvs
COPY pkg/agent/core/ngt/service/vqueue ./vqueue

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/pkg/${PKG_INTERNAL}
COPY pkg/${PKG_INTERNAL} .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd/${PKG}
COPY cmd/${PKG} .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/versions
COPY versions .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/Makefile.d
COPY Makefile.d .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
COPY Makefile .
RUN update-alternatives --set libblas.so-x86_64-linux-gnu /usr/lib/x86_64-linux-gnu/libmkl_rt.so \
     && make faiss/install

COPY .git .

RUN make REPO=${ORG} NAME=${REPO} cmd/${PKG}/${APP_NAME} \
    && mv "cmd/${PKG}/${APP_NAME}" "/usr/bin/${APP_NAME}"

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd/${PKG}
RUN cp sample.yaml /tmp/config.yaml

FROM ubuntu:devel
LABEL maintainer "${MAINTAINER}"

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

ENTRYPOINT ["/go/bin/faiss"]
