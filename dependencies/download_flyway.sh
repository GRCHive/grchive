#!/bin/bash

MAJOR_VERSION=6
MINOR_VERSION=1
PATCH_VERSION=3
FLYWAY_VERSION_NAME=${MAJOR_VERSION}.${MINOR_VERSION}.${PATCH_VERSION}
OS=$OSTYPE

if [[ "$OSTYPE" == "darwin"* ]]; then
    OS="macosx"
else
    OS="linux"
fi
curl -o flyway.tar.gz https://repo1.maven.org/maven2/org/flywaydb/flyway-commandline/${FLYWAY_VERSION_NAME}/flyway-commandline-${FLYWAY_VERSION_NAME}-${OS}-x64.tar.gz
mkdir -p flyway
tar xvf flyway.tar.gz -C ./flyway --strip-components=1
rm flyway.tar.gz 
