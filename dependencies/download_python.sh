#!/bin/bash

MAJOR_VERSION=3
MINOR_VERSION=8
PATCH_VERSION=1
FULL_VERSION=${MAJOR_VERSION}.${MINOR_VERSION}.${PATCH_VERSION}
URL="https://www.python.org/ftp/python/${FULL_VERSION}/Python-${FULL_VERSION}.tgz"

if [ ! -d python ]; then
    if [ ! -f python.gz ]; then
        curl -o python.tgz ${URL}
    fi

    mkdir -p python/python-${FULL_VERSION}
    tar xvf python.tgz --strip-components=1 -C python/python-${FULL_VERSION}
    rm python.tgz
fi

cd python/python-${FULL_VERSION}
cp ../../python-helper/Setup Modules/Setup
./configure LDFLAGS="-static -static-libgcc" CPPFLAGS="-static" --enable-optimizations --disable-shared --prefix="${PWD}/bin"
make LDFLAGS="-static" LINKFORSHARED=" " -j`nproc`
make install
