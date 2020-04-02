#!/bin/bash
FULL_SECRET_PATH=$1
STRIPPED_PATH="${FULL_SECRET_PATH#secret/}"

AUTH_RESP=$(curl \
    -X POST \
    -d "{\"password\": \"${VAULT_PASSWORD}\"}" \
    ${VAULT_HOST}:${VAULT_PORT}/v1/auth/userpass/login/${VAULT_USER})
TOKEN=$(echo $AUTH_RESP | jq -j '.auth.client_token')

SECRET_RESP=$(curl \
    -X GET \
    --header "X-Vault-Token: ${TOKEN}"\
    ${VAULT_HOST}:${VAULT_PORT}/v1/secret/data/${STRIPPED_PATH})
RESULT=$(echo $SECRET_RESP | jq '.data.data')
echo $RESULT
