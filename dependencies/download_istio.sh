#!/bin/bash

MAJOR_VERSION=1
MINOR_VERSION=6
PATCH_VERSION=0

export ISTIO_VERSION=${MAJOR_VERSION}.${MINOR_VERSION}.${PATCH_VERSION}
curl -L https://istio.io/downloadIstio | sh -
mv istio-${ISTIO_VERSION} istio
