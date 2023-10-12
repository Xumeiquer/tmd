# Build TDLIB stage
FROM alpine:latest AS tdlib

WORKDIR /

RUN apk update && apk upgrade && apk add --no-cache \
        openssl-libs-static \
        ca-certificates \
        linux-headers \
        alpine-sdk \
        openssl-dev \
        zlib-static\ 
        zlib-dev \
        cmake \
        gperf \
        git \
        php

RUN git clone https://github.com/tdlib/td.git && \
    cd td && \
    rm -rf build && \
    mkdir build && \
    cd build && \
    cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=/usr/local .. && \
    cmake --build . --target prepare_cross_compiling && \
    cd .. && \
    php SplitSource.php && \
    cd build && \
    cmake --build . --target install && \
    cd .. && \
    php SplitSource.php --undo && \
    cd .. && \
    ls -l /usr/local
