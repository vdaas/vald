# syntax = docker/dockerfile:latest
# check=error=true
#
# Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

# skipcq: DOK-DL3007
FROM ubuntu:22.04
LABEL maintainer="vdaas.org vald team <vald@vdaas.org>"

ARG TARGETARCH
ARG TARGETOS

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Etc/UTC

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install required CI build tools including pigz for parallel (de)compression
# skipcq: DOK-DL3008
RUN --mount=type=cache,target=/var/lib/apt,sharing=locked,id=lib-ci-${TARGETARCH} \
    --mount=type=cache,target=/var/cache/apt,sharing=locked,id=cache-ci-${TARGETARCH} \
    set -eux \
    && echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache \
    && echo 'APT::Install-Recommends "false";' > /etc/apt/apt.conf.d/no-install-recommends \
    && apt-get update -y \
    && apt-get install -y --no-install-recommends --fix-missing \
        build-essential \
        ca-certificates \
        curl \
        docker.io \
        git \
        jq \
        make \
        pigz \
        tzdata \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
