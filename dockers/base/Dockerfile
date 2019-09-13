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

