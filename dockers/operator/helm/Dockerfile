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
ARG OPERATOR_SDK_VERSION
ARG VERSION="v1"
ARG GROUP="vald.vdaas.org"
ARG VALD_KIND="ValdRelease"
ARG VALD_HELM_OPERATOR_KIND="ValdHelmOperatorRelease"
ARG BASE_TAG=latest
ARG MAINTAINER="vdaas.org vald team <vald@vdaas.org>"

FROM vdaas/vald-base:${BASE_TAG} AS builder
ARG OPERATOR_SDK_VERSION
ARG VERSION
ARG GROUP
ARG VALD_KIND
ARG VALD_HELM_OPERATOR_KIND

ENV ORG vdaas
ENV REPO vald

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

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/apis/grpc/v1/errors
COPY apis/grpc/v1/errors .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/internal/errors
COPY internal/errors .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/internal/log
COPY internal/log .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/versions
COPY versions .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/charts
COPY charts .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/hack/helm/schema/gen
COPY hack/helm/schema/gen .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/Makefile.d
COPY Makefile.d .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
COPY Makefile .
COPY .git .

RUN make helm/schema/vald
RUN make helm/schema/vald-helm-operator
RUN cp -r charts /charts

FROM quay.io/operator-framework/helm-operator:${OPERATOR_SDK_VERSION}
LABEL maintainer "${MAINTAINER}"

COPY --from=builder /tmp/watches.yaml ${HOME}/watches.yaml

COPY --from=builder /charts/vald ${HOME}/helm-charts/vald
COPY --from=builder /charts/vald-helm-operator ${HOME}/helm-charts/vald-helm-operator
