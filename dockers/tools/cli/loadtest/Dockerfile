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

FROM golang:${GO_VERSION} AS builder

ENV GO111MODULE on
ENV LANG en_US.UTF-8
ENV ORG vdaas
ENV REPO vald
ENV APP_NAME loadtest
ENV PKG tools/cli/${APP_NAME}

# skipcq: DOK-DL3008
RUN apt-get update && apt-get install -y --no-install-recommends \
    upx \
    git \
    libhdf5-dev \
    && ldconfig \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p "$GOPATH/src"

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

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/hack/benchmark/assets/x1b
COPY hack/benchmark/assets/x1b .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
COPY versions/GO_VERSION .
COPY versions/VALD_VERSION .
COPY .git .
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN GO_VERSION="$(< GO_VERSION)" \
    && VALD_VERSION="$(< VALD_VERSION)" \
    && GIT_COMMIT="$(git rev-list -1 HEAD)" \
    && CPU_INFO_FLAGS="$(< /proc/cpuinfo | grep flags | cut -d " " -f 2- | head -1 || true)" \
    && GOOS="$(go env GOOS)" \
    && GOARCH="$(go env GOARCH)" \
    && CGO_ENABLED=1 \
    && CGO_CXXFLAGS="-g -Ofast -march=native" \
    CGO_FFLAGS="-g -Ofast -march=native" \
    CGO_LDFLAGS="-g -Ofast -march=native" \
    GO111MODULE=on \
    go build \
    --ldflags "-s -w \
    -X 'github.com/${ORG}/${REPO}/internal/info.Version=${VALD_VERSION}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.GitCommit=${GIT_COMMIT}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.GoVersion=${GO_VERSION}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.GoOS=${GOOS}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.GoArch=${GOARCH}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.CGOEnabled=${CGO_ENABLED}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.BuildCPUInfoFlags=${CPU_INFO_FLAGS}'" \
    -a \
    -tags "cgo netgo" \
    -trimpath \
    -installsuffix "cgo netgo" \
    -o "${APP_NAME}" \
    "cmd/${PKG}/main.go" \
    && upx -9 -o "/usr/bin/${APP_NAME}" "${APP_NAME}"

# Start From Scratch For Running Environment
FROM ubuntu:devel
# Start From Alpine For Debug Environment
# FROM alpine:latest
LABEL maintainer "${MAINTAINER}"

ENV APP_NAME loadtest

COPY --from=builder /usr/lib/x86_64-linux-gnu/libaec*  /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libhdf5*  /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libsz*  /usr/lib/x86_64-linux-gnu/

# Copy certificates for SSL/TLS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy permissions
COPY --from=builder /etc/passwd /etc/passwd
# workaround: https://github.com/moby/moby/issues/37965
RUN true
# Copy our static executable
COPY --from=builder /usr/bin/${APP_NAME} /go/bin/${APP_NAME}

ENTRYPOINT ["/go/bin/loadtest"]
