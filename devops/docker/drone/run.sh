#!/bin/bash

# Wait for access to vault to come up..
sleep 30;

# Pull secrets from Vault.
VAULT_RESP=$(/app/get_vault_secret_http.sh ${DRONE_GITEA_CLIENT_SECRET_VAULT})


export DRONE_GITEA_CLIENT_ID=$(echo $VAULT_RESP | jq -j '.id')
export DRONE_GITEA_CLIENT_SECRET=$(echo $VAULT_RESP | jq -j '.secret')

echo "Starting Drone Server..."
exec /app/drone-server
