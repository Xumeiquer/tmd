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
        zlib-static \ 
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
        -DOPENSSL_INCLUDE_DIR=/usr/include/openssl \
        -DOPENSSL_SSL_LIBRARY=/usr/lib/libssl.a \
        -DOPENSSL_CRYPTO_LIBRARY=/usr/lib/libcrypto.a \
        -DZLIB_USE_STATIC_LIBS=ON \
        -DZLIB_INCLUDE_DIR=/usr/include \
        -DCMAKE_FIND_LIBRARY_SUFFIXES=.a \
        -DBUILD_SHARED_LIBS=OFF \
        -DCMAKE_EXE_LINKER_FLAGS=-static .. && \
    cmake --build . --target install && \
    cd ../.. && \
    ls -l /usr/local

# Build Telegram Media Downloader
FROM golang:alpine AS golang

RUN mkdir -p /usr/local/lib/pkgconfig/
COPY --from=tdlib /usr/local/include/td /usr/local/include/td
COPY --from=tdlib /usr/local/lib/libtd*.a /usr/local/lib/
COPY --from=tdlib /usr/local/lib/pkgconfig/libtd*.pc /usr/local/lib/pkgconfig/

RUN apk update && apk upgrade && apk add --no-cache \
    openssl-libs-static \
    build-base \
    zlib-static \
    pkgconfig \
    musl-dev \
    bash \
    git 

RUN mkdir -p /usr/src/tmd
WORKDIR /usr/src/tmd

#RUN git clone https://github.com/Xumeiquer/tmd.git .
COPY . /usr/src/tmd
RUN go mod tidy

ENV PKG_CONFIG_PATH="$PKG_CONFIG_PATH:/usr/lib/pkgconfig:/usr/local/lib/pkgconfig"
ENV CGO_CFLAGS="-I/usr/local/include/td/telegram/ -I/usr/local/include/td/tl"
ENV CGO_LDFLAGS="-L/lib -L/usr/lib -L/usr/local/lib -ltdjson -ldl -lm  -lstdc++ -lz"
RUN bash build.sh

# Final image
FROM gcr.io/distroless/base:latest

COPY --from=golang /usr/src/tmd/build/tmd /usr/bin/tmd

WORKDIR /

ENTRYPOINT [ "/usr/bin/tmd"]
