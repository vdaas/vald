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
ARG DISTROLESS_IMAGE=gcr.io/distroless/static
ARG DISTROLESS_IMAGE_TAG=nonroot
ARG OPERATOR_SDK_VERSION
ARG UPX_OPTIONS=-9
ARG VERSION="v1"
ARG GROUP="vald.vdaas.org"
ARG VALD_KIND="ValdRelease"
ARG VALD_HELM_OPERATOR_KIND="ValdHelmOperatorRelease"
ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>"

# skipcq: DOK-DL3026
FROM quay.io/operator-framework/helm-operator:${OPERATOR_SDK_VERSION} AS operator
FROM golang:${GO_VERSION} AS golang
FROM ubuntu:devel AS builder

ARG OPERATOR_SDK_VERSION
ARG VERSION
ARG GROUP
ARG VALD_KIND
ARG VALD_HELM_OPERATOR_KIND
ARG UPX_OPTIONS

ENV GO111MODULE on
ENV DEBIAN_FRONTEND noninteractive
ENV INITRD No
ENV LANG en_US.UTF-8
ENV GOROOT /opt/go
ENV GOPATH /go
ENV PATH ${PATH}:${GOROOT}/bin:${GOPATH}/bin
ENV ORG vdaas
ENV REPO vald
ENV APP_NAME helm-operator

# skipcq: DOK-DL3008
RUN apt-get clean \
    && rm -rf \
        /var/lib/apt/lists/* \
        /var/cache/* \
    && apt-get update -y \
    && apt-get upgrade -y \
    && apt-get install -y --no-install-recommends --fix-missing \
    ca-certificates \
    build-essential \
    curl \
    upx \
    git \
    && apt-get clean \
    && rm -rf \
        /var/lib/apt/lists/* \
        /var/cache/* \
    && apt-get autoremove

COPY --from=golang /usr/local/go $GOROOT
RUN mkdir -p "$GOPATH/src"

RUN { \
        echo "---"; \
        echo "- version: ${VERSION}"; \
        echo "  group: ${GROUP}"; \
        echo "  kind: ${VALD_KIND}"; \
        echo "  chart: helm-charts/vald"; \
        echo "- version: ${VERSION}"; \
        echo "  group: ${GROUP}"; \
        echo "  kind: ${VALD_HELM_OPERATOR_KIND}"; \
        echo "  chart: helm-charts/vald-helm-operator"; \
    } > /tmp/watches.yaml

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/Makefile.d
COPY Makefile.d .
WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
COPY Makefile .
COPY .git .
COPY go.mod .
COPY go.sum .

RUN make go/download

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/internal/errors
COPY internal/errors .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/internal/log
COPY internal/log .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/internal/conv
COPY internal/conv .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/internal/strings
COPY internal/strings .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/versions
COPY versions .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/charts
COPY charts .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/hack/helm/schema/gen
COPY hack/helm/schema/gen .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
RUN make helm/schema/vald helm/schema/vald-helm-operator \
    && cp -r charts /charts

COPY --from=operator /usr/local/bin/${APP_NAME} /usr/bin/${APP_NAME}

RUN upx ${UPX_OPTIONS} "/usr/bin/${APP_NAME}"

FROM ${DISTROLESS_IMAGE}:${DISTROLESS_IMAGE_TAG}
LABEL maintainer="${MAINTAINER}"

ENV APP_NAME helm-operator

ENV HOME=/opt/helm
WORKDIR ${HOME}

COPY --from=builder /tmp/watches.yaml watches.yaml
COPY --from=builder /charts/vald helm-charts/vald
COPY --from=builder /charts/vald-helm-operator helm-charts/vald-helm-operator
COPY --from=builder /usr/bin/${APP_NAME} /usr/bin/${APP_NAME}

USER nonroot:nonroot

ENTRYPOINT ["/usr/bin/helm-operator", "run", "--watches-file=./watches.yaml"]
