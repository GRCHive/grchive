#!/bin/bash

MAJOR_VERSION=3
MINOR_VERSION=3
FULL_VERSION=${MAJOR_VERSION}.${MINOR_VERSION}
URL="https://github.com/libffi/libffi/releases/download/v${FULL_VERSION}/libffi-${FULL_VERSION}.tar.gz"

if [ ! -d libffi ]; then
    if [ ! -f libffi.tar.gz ]; then
        curl -Lo libffi.tar.gz ${URL}
    fi

    mkdir -p libffi
    tar xvf libffi.tar.gz -C libffi
    rm libffi.tar.gz
fi

cd libffi/libffi-${FULL_VERSION}
./configure --prefix="${PWD}/bin"
make -j`nproc`
make install
