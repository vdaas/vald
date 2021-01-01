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
ARG BASE_TAG=latest
ARG DISTROLESS_IMAGE=gcr.io/distroless/static
ARG DISTROLESS_IMAGE_TAG=nonroot
ARG UPX_OPTIONS=-9
ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>

FROM golang:${GO_VERSION} AS golang

FROM vdaas/vald-base:${BASE_TAG} AS builder
ARG UPX_OPTIONS

ENV ORG vdaas
ENV REPO vald
ENV PKG manager/index
ENV APP_NAME index

COPY --from=golang /usr/local/go $GOROOT

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

FROM ${DISTROLESS_IMAGE}:${DISTROLESS_IMAGE_TAG}
LABEL maintainer "${MAINTAINER}"

ENV APP_NAME index

COPY --from=builder /usr/bin/${APP_NAME} /go/bin/${APP_NAME}

USER nonroot:nonroot

ENTRYPOINT ["/go/bin/index"]
