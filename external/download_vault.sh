#!/bin/bash

curl -o vault.zip https://releases.hashicorp.com/vault/1.2.3/vault_1.2.3_linux_amd64.zip
mkdir -p vault
unzip vault.zip -d vault
rm vault.zip
