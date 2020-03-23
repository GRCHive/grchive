#!/bin/bash

MAJOR_VERSION=1
MINOR_VERSION=3
PATCH_VERSION=70
FULL_VERSION=${MAJOR_VERSION}.${MINOR_VERSION}.${PATCH_VERSION}
URL="https://github.com/JetBrains/kotlin/releases/download/v${FULL_VERSION}/kotlin-compiler-${FULL_VERSION}.zip"

curl -Lo kotlinc.zip $URL
unzip kotlinc.zip
rm kotlinc.zip
