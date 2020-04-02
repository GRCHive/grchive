#!/bin/bash
USERNAME="grchive-gitea-admin"
PASSWORD="$(openssl rand -hex 32)A!"

# Create admin user OR change password if that failed (e.g. if we executed this already)
docker exec -it gitea sh -c "gitea admin create-user --username ${USERNAME} --password ${PASSWORD} --email support@grchive.com --must-change-password=false --admin || gitea admin change-password --username ${USERNAME} --password ${PASSWORD}"

USER_AUTH_URL="${GITEA_PROTOCOL}://${USERNAME}:${PASSWORD}@${GITEA_HOST}:${GITEA_PORT}"
RESULT=$(curl \
    -H "Content-Type: application/json" \
    -H "accept: application/json" \
    -X POST \
    -d '{"name": "grchive-global-token"}' \
    "${USER_AUTH_URL}/api/v1/users/${USERNAME}/tokens")
echo $RESULT
TOKEN=$(echo $RESULT | jq -j '.sha1')

OAUTH_RESULT=$(curl \
    -H "Content-Type: application/json" \
    -H "accept: application/json" \
    -H "Authorization: token ${TOKEN}" \
    -X POST \
    -d "{\"name\": \"grchive-drone-ci\", \"redirect_uris\": [ \"${DRONE_SERVER_PROTO}://${DRONE_SERVER_HOST}:${DRONE_SERVER_PORT}/login\" ]}" \
    "${GITEA_PROTOCOL}://${GITEA_HOST}:${GITEA_PORT}/api/v1/user/applications/oauth2")
echo $OAUTH_RESULT
CLIENT_ID=$(echo $OAUTH_RESULT | jq -j '.client_id')
CLIENT_SECRET=$(echo $OAUTH_RESULT | jq -j '.client_secret')

vault login -address="${VAULT_HOST}:${VAULT_PORT}" ${VAULT_TOKEN}
vault kv put -address="${VAULT_HOST}:${VAULT_PORT}" ${GITEA_TOKEN} token="${TOKEN}" password="${PASSWORD}"
vault kv put -address="${VAULT_HOST}:${VAULT_PORT}" ${DRONE_GITEA_CLIENT_SECRET_VAULT} id="${CLIENT_ID}" secret="${CLIENT_SECRET}"

echo $TOKEN
echo $CLIENT_ID
echo $CLIENT_SECRET
