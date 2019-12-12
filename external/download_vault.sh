#!/bin/bash

if [[ "$OSTYPE" == "darwin"* ]]; then
    curl -o vault.zip https://releases.hashicorp.com/vault/1.2.3/vault_1.2.3_darwin_amd64.zip
else
    curl -o vault.zip https://releases.hashicorp.com/vault/1.2.3/vault_1.2.3_linux_amd64.zip
fi

mkdir -p vault
unzip vault.zip -d vault
rm vault.zip
