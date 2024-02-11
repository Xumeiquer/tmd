#!/bin/bash

# file: .github/scripts/prepare-tdlib.sh

set -x

export DEBIAN_FRONTEND=noninteractive

apt update
apt install -y make git zlib1g-dev libssl-dev gperf php-cli cmake g++

mkdir -p /usr/src/tdlib/

git clone https://github.com/tdlib/td.git /usr/src/tdlib/td
cd /usr/src/tdlib/td
git checkout ${{env.TDLIB_VERSION}}

mkdir -p /usr/src/tdlib/td/build-native

TMD_TARGETS=($(echo ${{env.TARGETS}} | tr ";" "\n"))
for qtarget in "${TMD_TARGETS[@]}"; do
    echo "$qtarget"
    target=$(echo "$qtarget" | sed -e 's/^"//' -e 's/"$//')
    echo "$target"
    mkdir -p /usr/src/tdlib/td/build-${target}
    chmod 777 /usr/src/tdlib/td/build-${target}
done

cd /usr/src/tdlib/td/build-native
cmake -DCMAKE_BUILD_TYPE=Release ..
cmake --build . --target prepare_cross_compiling
