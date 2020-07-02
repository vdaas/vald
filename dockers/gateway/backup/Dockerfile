#
# Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
ENV PKG gateway/backup
ENV APP_NAME backup

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/internal
COPY internal .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/apis/grpc
COPY apis/grpc .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/pkg/${PKG}
COPY pkg/${PKG} .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd/${PKG}
COPY cmd/${PKG} .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
COPY versions/GO_VERSION .
COPY versions/VALD_VERSION .
COPY .git .
RUN GO_VERSION="$(cat GO_VERSION)" \
    && VALD_VERSION="$(cat VALD_VERSION)" \
    && GIT_COMMIT="$(git rev-list -1 HEAD)" \
    && CPU_INFO_FLAGS="$(cat /proc/cpuinfo | grep flags | cut -d " " -f 2- | head -1)" \
    && GOOS="$(go env GOOS)" \
    && GOARCH="$(go env GOARCH)" \
    && CGO_ENABLED=0 \
    && GO111MODULE=on \
    go build \
    --ldflags "-s -w -linkmode 'external' \
    -extldflags '-static' \
    -X 'github.com/${ORG}/${REPO}/internal/info.Version=${VALD_VERSION}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.GitCommit=${GIT_COMMIT}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.GoVersion=${GO_VERSION}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.GoOS=${GOOS}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.GoArch=${GOARCH}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.CGOEnabled=${CGO_ENABLED}' \
    -X 'github.com/${ORG}/${REPO}/internal/info.BuildCPUInfoFlags=${CPU_INFO_FLAGS}'" \
    -a \
    -tags netgo \
    -installsuffix netgo \
    -trimpath \
    -o "${APP_NAME}" \
    "cmd/${PKG}/main.go" \
    && upx -9 -o "/usr/bin/${APP_NAME}" "${APP_NAME}"

# Start From Scratch For Running Environment
FROM scratch
# Start From Alpine For Debug Environment
# FROM alpine:latest

ENV APP_NAME backup

# Copy certificates for SSL/TLS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy permissions
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builder /usr/bin/${APP_NAME} /go/bin/${APP_NAME}

ENTRYPOINT ["/go/bin/backup"]
