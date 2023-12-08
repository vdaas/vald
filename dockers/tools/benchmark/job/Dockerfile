#
# Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
ARG DISTROLESS_IMAGE=gcr.io/distroless/static
ARG DISTROLESS_IMAGE_TAG=nonroot
ARG UPX_OPTIONS=-9
ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>"
ARG TARGETARCH

FROM golang:${GO_VERSION} AS golang

FROM ubuntu:devel AS builder

ARG UPX_OPTIONS
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
ENV LIB_DIR=/usr/local/lib
ENV LD_LIBRARY_PATH $LD_LIBRARY_PATH:${LIB_DIR}
ENV ARCH=${TARGETARCH}
ENV XARCH x86_64
ENV AARCH aarch64

# skipcq: DOK-DL3008
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
    && curl -sOL https://github.com/madler/zlib/releases/download/v1.3/zlib-1.3.tar.gz \
    && tar -xzvf zlib-1.3.tar.gz \
    && cd zlib-1.3 \
    && ./configure --prefix="${LIB_DIR}" --static --shared \
    && make \
    && make test \
    && make install \
    && cd / \
    && curl -sOL https://github.com/HDFGroup/hdf5/releases/download/hdf5-1_14_3/hdf5-1_14_3.tar.gz \
    && tar -xzvf hdf5-1_14_3.tar.gz \
    && cd hdfsrc \
    && ./configure --enable-build-mode=production --enable-static-exec --disable-shared --prefix="/usr/local" --with-zlib=/usr/local/include,/usr/local/lib LDFLAGS="-Wl,-rpath,${LIB_DIR}" \
    && make check \
    && make install \
    && ldconfig
 #    && mkdir -p /usr/local/lib/linux-gnu \
 #    && if [ "${ARCH}" = "amd64" ]; then \
	# ARCH=${XARCH} ; \
	# # mv /lib64/ld-linux-x86-64.so.2 /usr/local/lib/linux-gnu ; \
 #    else \
	# ARCH=${AARCH}; \
 #    fi \
 #    && mv /usr/lib/${ARCH}-linux-gnu/libc.so.6 /usr/local/lib/linux-gnu/ \
 #    && /usr/sbin/ldconfig


COPY --from=golang /usr/local/go $GOROOT
RUN mkdir -p "$GOPATH/src"

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/Makefile.d
COPY Makefile.d .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
COPY Makefile .
COPY .git .
COPY go.mod .
COPY go.sum .

RUN make go/download

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
RUN make REPO=${ORG} NAME=${REPO} cmd/${PKG}/${APP_NAME} \
    && mv "cmd/${PKG}/${APP_NAME}" "/usr/bin/${APP_NAME}"

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd/${PKG}
RUN cp sample.yaml /tmp/config.yaml

# CMD ["/usr/bin/bash"]

FROM ${DISTROLESS_IMAGE}:${DISTROLESS_IMAGE_TAG}
LABEL maintainer="${MAINTAINER}"

ENV APP_NAME job

ENV LD_LIBRARY_PATH $LD_LIBRARY_PATH:/usr/lib

COPY --from=builder /usr/local/lib/* /usr/lib/
COPY --from=builder /usr/local/include/* /usr/include
COPY --from=builder /lib/x86_64-linux-gnu/libc.so.6 /lib/x86_64-linux-gnu/
# COPY --from=builder /lib64/* /lib64/
COPY --from=builder /usr/bin/${APP_NAME} /usr/bin/${APP_NAME}
COPY --from=builder /tmp/config.yaml /etc/server/config.yaml

USER nonroot:nonroot

ENTRYPOINT ["/usr/bin/job"]
