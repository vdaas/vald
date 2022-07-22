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
ARG DISTROLESS_IMAGE=gcr.io/distroless/static
ARG DISTROLESS_IMAGE_TAG=nonroot
ARG UPX_OPTIONS=-9
ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>"

FROM golang:${GO_VERSION} AS builder

ARG UPX_OPTIONS

ENV GO111MODULE on
ENV LANG en_US.UTF-8
ENV ORG vdaas
ENV REPO vald
ENV PKG benchmark/job/search
ENV APP_NAME search

RUN apt-get update && apt-get install -y --no-install-recommends \
    upx \
    git \
    libhdf5-dev \
    && ldconfig \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p ${GOPATH}/src

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

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd/${PKG}
COPY cmd/${PKG} .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/versions
COPY versions .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/Makefile.d
COPY Makefile.d .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
COPY Makefile .

COPY .git .

RUN make REPO=${ORG} NAME=${REPO} cmd/${PKG}/${APP_NAME} \
    && upx ${UPX_OPTIONS} -o "/usr/bin/${APP_NAME}" "cmd/${PKG}/${APP_NAME}"
    # && mv cmd/${PKG}/${APP_NAME} /usr/bin/${APP_NAME}

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd/${PKG}
RUN cp sample.yaml /tmp/config.yaml

FROM ${DISTROLESS_IMAGE}:${DISTROLESS_IMAGE_TAG}
LABEL maintainer "Vald team <vald@vdaas.org>"
ENV APP_NAME search

COPY --from=builder /usr/lib/x86_64-linux-gnu/libaec*  /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libhdf5*  /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libsz*  /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libhdf5_serial.so /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libhdf5_serial_hl.so /usr/lib/x86_64-linux-gnu/
COPY --from=builder /lib/x86_64-linux-gnu/libpthread.so.0 /lib/x86_64-linux-gnu/
COPY --from=builder /lib/x86_64-linux-gnu/libc.so.6 /lib/x86_64-linux-gnu/
# COPY --from=builder /usr/lib/x86_64-linux-gnu/libsz.so.2 /usr/lib/x86_64-linux-gnu/
COPY --from=builder /lib/x86_64-linux-gnu/libz.so.1 /lib/x86_64-linux-gnu/
COPY --from=builder /lib/x86_64-linux-gnu/libdl.so.2 /lib/x86_64-linux-gnu/
COPY --from=builder /lib/x86_64-linux-gnu/libm.so.6 /lib/x86_64-linux-gnu/
# COPY --from=builder /usr/lib/x86_64-linux-gnu/libaec.so.0 /usr/lib/x86_64-linux-gnu/
COPY --from=builder /lib64/ld-linux-x86-64.so.2 /lib64/

COPY --from=builder /usr/bin/${APP_NAME} /go/bin/${APP_NAME}
COPY --from=builder /tmp/config.yaml /etc/server/config.yaml

USER nonroot:nonroot

ENTRYPOINT ["/go/bin/search"]
