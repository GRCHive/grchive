#!/bin/bash
export DRONE_SERVER="${DRONE_SERVER_PROTO}://${DRONE_SERVER_HOST}:${DRONE_SERVER_PORT}"
export DRONE_TOKEN="${DRONE_TOKEN}"

vault login -address="${VAULT_HOST}:${VAULT_PORT}" ${VAULT_TOKEN}
TOKEN=$(vault kv get -address="${VAULT_HOST}:${VAULT_PORT}" -field token ${GITEA_TOKEN})

# Create Drone OAuth -- maybe move this to Drone setup...?
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

vault kv put -address="${VAULT_HOST}:${VAULT_PORT}" ${DRONE_GITEA_CLIENT_SECRET_VAULT} id="${CLIENT_ID}" secret="${CLIENT_SECRET}"

echo $CLIENT_ID
echo $CLIENT_SECRET

# Create Gitlab Registry Secret
DOCKER_SECRET="{ \"auths\": {\"registry.gitlab.com\": { \"auth\": \"$(echo ${GKE_REGISTRY_USER}:${GKE_REGISTRY_PASSWORD} | base64)\" } } }"
${CLI} orgsecret add ${GITEA_GLOBAL_ORG} dockerregistry "${DOCKER_SECRET}" || ${CLI} orgsecret update ${GITEA_GLOBAL_ORG} dockerregistry "${DOCKER_SECRET}"
