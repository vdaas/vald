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

ARG DISTROLESS_IMAGE=gcr.io/distroless/static
ARG DISTROLESS_IMAGE_TAG=nonroot
ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>"

FROM --platform=${TARGETPLATFORM} ubuntu:devel AS builder

ARG GO_VERSION
ARG TARGETARCH
ARG TARGETOS

ENV GO111MODULE on
ENV GOPATH /go
ENV GOROOT /opt/go
ENV DEBIAN_FRONTEND noninteractive
ENV INITRD No
ENV LANG en_US.UTF-8
ENV ORG vdaas
ENV REPO vald
ENV APP_NAME sidecar
ENV PKG agent/sidecar
ENV PKG_INTERNAL agent/internal
ENV PATH ${PATH}:${GOROOT}/bin:${GOPATH}/bin

# skipcq: DOK-DL3008
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    build-essential \
    curl \
    upx \
    git \
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
    make GO_VERSION=${GO_VERSION} ARCH=${TARGETARCH} OS=${TARGETOS} go/install \
    && make GOARCH=${TARGETARCH} GOOS=${TARGETOS} go/download

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/internal
COPY internal .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/apis/grpc
COPY apis/grpc .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/pkg/${PKG}
COPY pkg/${PKG} .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/pkg/${PKG_INTERNAL}
COPY pkg/${PKG_INTERNAL} .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd/${PKG}
COPY cmd/${PKG} .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/versions
COPY versions .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
RUN --mount=type=cache,target="${GOPATH}/pkg",id="go-build-${TARGETARCH}" \
    --mount=type=cache,target="${HOME}/.cache/go-build",id="go-build-${TARGETARCH}" \
    make GOARCH=${TARGETARCH} GOOS=${TARGETOS} REPO=${ORG} NAME=${REPO} cmd/${PKG}/${APP_NAME} \
    && mv "cmd/${PKG}/${APP_NAME}" "/usr/bin/${APP_NAME}"

FROM --platform=${TARGETPLATFORM} ${DISTROLESS_IMAGE}:${DISTROLESS_IMAGE_TAG}
LABEL maintainer="${MAINTAINER}"

ENV APP_NAME sidecar

COPY --from=builder /usr/bin/${APP_NAME} /go/bin/${APP_NAME}

USER nonroot:nonroot

ENTRYPOINT ["/go/bin/sidecar"]
