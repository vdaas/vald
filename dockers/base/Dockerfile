#
# Copyright (C) 2019-2019 kpango (Yusuke Kato)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

FROM ubuntu:latest AS builder

ENV NGT_VERSION 1.7.9
ENV GO_VERSION 1.12.9

ENV DEBIAN_FRONTEND noninteractive
ENV INITRD No
ENV LANG en_US.UTF-8
ENV GOROOT /opt/go
ENV GOPATH /go

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    build-essential \
    cmake \
    upx \
    curl \
    unzip \
    git \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /tmp
RUN curl -LO https://github.com/yahoojapan/NGT/archive/v${NGT_VERSION}.tar.gz
	tar zxf v${NGT_VERSION}.tar.gz -C /tmp
	cd /tmp/NGT-${NGT_VERSION}; cmake .
	make -j -C /tmp/NGT-${NGT_VERSION}
	make install -C /tmp/NGT-${NGT_VERSION}

WORKDIR /opt
RUN curl -sSL -O https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar zxf go${GO_VERSION}.linux-amd64.tar.gz && rm go${GO_VERSION}.linux-amd64.tar.gz && \
    ln -s /opt/go/bin/go /usr/bin/ && \
    mkdir $GOPATH

