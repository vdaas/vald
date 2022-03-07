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
ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>"

FROM golang:${GO_VERSION} AS builder

ENV GO111MODULE on
ENV LANG en_US.UTF-8
ENV ORG vdaas
ENV REPO vald
ENV PKG discoverer/k8s
ENV APP_NAME discoverer

RUN apt-get update && apt-get install -y --no-install-recommends \
    upx \
    git \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p $GOPATH/src

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
    && mv "cmd/${PKG}/${APP_NAME}" "/usr/bin/${APP_NAME}"

FROM ${DISTROLESS_IMAGE}:${DISTROLESS_IMAGE_TAG}
LABEL maintainer "${MAINTAINER}"

ENV APP_NAME discoverer

COPY --from=builder /usr/bin/${APP_NAME} /go/bin/${APP_NAME}

USER nonroot:nonroot

ENTRYPOINT ["/go/bin/discoverer"]
