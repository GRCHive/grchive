#!/bin/bash

MAJOR_VERSION=1
MINOR_VERSION=17
PATCH_VERSION=0
FULL_VERSION=${MAJOR_VERSION}.${MINOR_VERSION}.${PATCH_VERSION}

mkdir -p kubectl
curl -Lo kubectl/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.17.0/bin/linux/amd64/kubectl 
chmod +x ./kubectl/kubectl
