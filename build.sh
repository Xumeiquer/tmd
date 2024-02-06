#!/usr/bin/env bash

# STEP 1: Determinate the required values

PACKAGE="github.com/Xumeiquer/tmd"
VERSION=$(git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')
COMMIT_HASH=$(git rev-parse --short HEAD 2> /dev/null)
BUILD_TIMESTAMP=$(date '+%Y-%m-%dT%H:%M:%S')

SOURCE="main.go"
OUTPUT="tmd"

# STEP 2: Build the ldflags

CFLAGS=$(pkg-config --cflags openssl tdclient tdcore tdactor tddb tdsqlite tdnet tdutils tdjson_static tdjson_private)
LIBS=$(pkg-config --libs openssl tdclient tdcore tdactor tddb tdsqlite tdnet tdutils tdjson_static tdjson_private)

LDFLAGS=(
  "-s -w"
  "-X '${PACKAGE}/cmd.Version=${VERSION:-v0.0.0-dev}'"
  "-X '${PACKAGE}/cmd.CommitHash=${COMMIT_HASH:-000000}'"
  "-X '${PACKAGE}/cmd.BuildTime=${BUILD_TIMESTAMP}'"
  "-extldflags '-static $LIBS $CFLAGS -ldl -lm  -lstdc++ -lz'"
)

# STEP 3: Actual Go build process

echo "go build -ldflags=\"${LDFLAGS[*]}\" -trimpath -o build/$OUTPUT $SOURCE"
go build -ldflags="${LDFLAGS[*]}" -trimpath -o build/$OUTPUT $SOURCE
