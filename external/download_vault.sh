#!/bin/bash

MAJOR_VERSION=1
MINOR_VERSION=3
PATCH_VERSION=0
VAULT_VERSION_NAME=${MAJOR_VERSION}.${MINOR_VERSION}.${PATCH_VERSION}
OS=$OSTYPE

if [[ "$OSTYPE" == "darwin"* ]]; then
    OS="darwin"
else
    OS="linux"
fi
curl -o vault.zip https://releases.hashicorp.com/vault/${VAULT_VERSION_NAME}/vault_${VAULT_VERSION_NAME}_${OS}_amd64.zip

mkdir -p vault
unzip vault.zip -d vault
rm vault.zip
