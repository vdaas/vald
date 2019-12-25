#
# Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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

FROM vdaas/vald-base:latest AS builder

ENV ORG vdaas
ENV REPO vald
ENV PKG agent/ngt
ENV APP_NAME ngt

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/internal
COPY internal .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/apis/grpc
COPY apis/grpc .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/pkg/${PKG}
COPY pkg/${PKG} .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd/${PKG}
COPY cmd/${PKG} .

WORKDIR /tmp
COPY versions/NGT_VERSION .
RUN export NGT_VERSION="$(cat NGT_VERSION)" \
    && curl -LO "https://github.com/yahoojapan/NGT/archive/v${NGT_VERSION}.tar.gz" \
    && tar zxf "v${NGT_VERSION}.tar.gz" -C /tmp \
    && cd "/tmp/NGT-${NGT_VERSION}" \
    && cmake . \
    && make -j -C "/tmp/NGT-${NGT_VERSION}" \
    && make install -C "/tmp/NGT-${NGT_VERSION}"

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
COPY versions/GO_VERSION .
COPY versions/VALD_VERSION .
COPY versions/NGT_VERSION .
COPY .git .
RUN GO_VERSION="$(cat GO_VERSION)" \
    && VALD_VERSION="$(cat VALD_VERSION)" \
    && NGT_VERSION="$(cat NGT_VERSION)" \
    && GIT_COMMIT="$(git rev-list -1 HEAD)" \
    && GOOS="$(go env GOOS)" \
    && GOARCH="$(go env GOARCH)" \
    && CGO_ENABLED=1 \
    && CGO_CXXFLAGS="-g -Ofast -march=native" \
    CGO_FFLAGS="-g -Ofast -march=native" \
    CGO_LDFLAGS="-g -Ofast -march=native" \
    GO111MODULE=on \
    go build \
    --ldflags "-s -w -linkmode 'external' \
    -extldflags '-static -fPIC -m64 -pthread -fopenmp -std=c++17 -lstdc++ -lm' \
    -X 'github.com/${ORG}/${REPO}/internal/info.Version=${VALD_VERSION}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.GitCommit=${GIT_COMMIT}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.GoVersion=${GO_VERSION}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.GoOS=${GOOS}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.GoArch=${GOARCH}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.CGOEnabled=${CGO_ENABLED}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.NGTVersion=${NGT_VERSION}'" \
    -a \
    -tags "cgo netgo" \
    -trimpath \
    -installsuffix "cgo netgo" \
    -o "${APP_NAME}" \
    "cmd/${PKG}/main.go" \
    && upx -9 -o "/usr/bin/${APP_NAME}" "${APP_NAME}"

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd/${PKG}
RUN cp sample.yaml /tmp/config.yaml

# Start From Scratch For Running Environment
FROM scratch
# Start From Alpine For Debug Environment
# FROM alpine:latest
LABEL maintainer "kpango <i.can.feel.gravity@gmail.com>"

ENV APP_NAME ngt

# Copy certificates for SSL/TLS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy permissions
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builder /usr/bin/${APP_NAME} /go/bin/${APP_NAME}
COPY --from=builder /tmp/config.yaml /etc/server/config.yaml

ENTRYPOINT ["/go/bin/ngt"]
