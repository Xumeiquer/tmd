# Build TDLIB stage

FROM alpine:latest AS tdlib

RUN mkdir -p /usr/src/tdlib
WORKDIR /usr/src/tdlib

RUN apk update && apk upgrade && apk add --no-cache \
        openssl-libs-static \
        ca-certificates \
        linux-headers \
        alpine-sdk \
        openssl-dev \
        zlib-static\ 
        musl-dev \
        zlib-dev \
        cmake \
        gperf \
        git

RUN git clone https://github.com/tdlib/td.git /usr/src/tdlib/td && ls -la /usr/src/tdlib && \
    cd /usr/src/tdlib/td && \
    rm -rf build && \
    git checkout $(git describe --tags "$(git rev-list --tags --max-count=1)") && \
    mkdir build && \
    cd build && \
    cmake -DCMAKE_BUILD_TYPE=Release \
        -DCMAKE_INSTALL_PREFIX:PATH=/usr/local \
        -DCMAKE_FIND_LIBRARY_SUFFIXES=.a \
        -DBUILD_SHARED_LIBS=OFF \
        -DCMAKE_EXE_LINKER_FLAGS=-static .. && \
    cmake --build . --target install && \
    cd ../.. && \
    ls -l /usr/local

# Build Telegram Media Downloader
FROM golang:alpine AS golang

COPY --from=tdlib /usr/local/include/td /usr/local/include/td
COPY --from=tdlib /usr/local/lib/libtd* /usr/local/lib/

RUN apk update && apk upgrade && apk add --no-cache \
    openssl-libs-static \
    build-base \
    zlib-static\ 
    pkgconfig \
    bash \
    git 

RUN mkdir -p /usr/src/tmd
WORKDIR /usr/src/tmd

#RUN git clone https://github.com/Xumeiquer/tmd.git .
COPY . /usr/src/tmd
RUN go mod tidy

RUN bash build.sh

# Final image
FROM gcr.io/distroless/base:latest

COPY --from=golang /usr/src/tmd/build/tmd /usr/bin/tmd

WORKDIR /

ENTRYPOINT [ "/usr/bin/tmd"]
