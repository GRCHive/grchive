#!/bin/bash
USERNAME="grchive-gitea-admin"
PASSWORD="$(openssl rand -hex 32)A!"

# Create admin user OR change password if that failed (e.g. if we executed this already)
docker exec -it gitea sh -c "gitea admin create-user --username ${USERNAME} --password ${PASSWORD} --email support@grchive.com --must-change-password=false || gitea admin change-password --username ${USERNAME} --password ${PASSWORD}"

TOKEN=$(curl \
    -H "Content-Type: application/json" \
    -H "accept: application/json" \
    -X POST \
    -d '{"name": "grchive-global-token"}' \
    "http://${USERNAME}:${PASSWORD}@localhost:3000/api/v1/users/${USERNAME}/tokens" | jq -j '.sha1')

echo $TOKEN
vault login -address="${VAULT_HOST}:${VAULT_PORT}" ${VAULT_TOKEN}
vault kv put -address="${VAULT_HOST}:${VAULT_PORT}" ${GITEA_TOKEN} token="${TOKEN}"
