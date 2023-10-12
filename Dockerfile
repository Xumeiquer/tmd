# Build TDLIB stage
FROM ghcr.io/xumeiquer/tdlib:latest AS tdlib

# Build Telegram Media Downloader
FROM golang:alpine AS golang

COPY --from=tdlib /usr/local/include/td /usr/local/include/td
COPY --from=tdlib /usr/local/lib/libtd* /usr/local/lib/
COPY --from=tdlib /usr/lib/libssl.a /usr/local/lib/libssl.a
COPY --from=tdlib /usr/lib/libcrypto.a /usr/local/lib/libcrypto.a
COPY --from=tdlib /lib/libz.a /usr/local/lib/libz.a

RUN apk add build-base bash

WORKDIR /tmd

COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY .git .

RUN bash build.sh

# Final image
FROM gcr.io/distroless/base:latest

COPY --from=golang /tmd/build/tmd /tmd

ENTRYPOINT [ "/tmd" ]
