ARG GO_VERSION=latest
ARG BASE_TAG=latest
ARG DISTROLESS_IMAGE=gcr.io/distroless/static
ARG DISTROLESS_IMAGE_TAG=nonroot
ARG UPX_OPTIONS=-9

FROM golang:${GO_VERSION} AS golang

FROM vdaas/vald-base:${BASE_TAG} AS builder
ARG UPX_OPTIONS

ENV ORG vdaas
ENV REPO vald
ENV PKG rebalancer/storage/job
ENV APP_NAME rebalancer-storage-job

COPY --from=golang /usr/local/go $GOROOT

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

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/Makefile.d
COPY Makefile.d .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
COPY Makefile .

COPY .git .

RUN make REPO=${ORG} NAME=${REPO} cmd/${PKG}/${APP_NAME} \
    && upx ${UPX_OPTIONS} -o "/usr/bin/${APP_NAME}" "cmd/${PKG}/${APP_NAME}"

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd/${PKG}
RUN cp sample.yaml /tmp/config.yaml

FROM ${DISTROLESS_IMAGE}:${DISTROLESS_IMAGE_TAG}
LABEL maintainer "Vald team <vald@vdaas.org>"

ENV APP_NAME rebalancer-storage-job

COPY --from=builder /usr/bin/${APP_NAME} /go/bin/${APP_NAME}
COPY --from=builder /tmp/config.yaml /etc/server/config.yaml

USER nonroot:nonroot

ENTRYPOINT ["/go/bin/rebalancer-storage-job"]