#!/bin/bash
USERNAME="grchive-gitea-admin"
PASSWORD="$(openssl rand -hex 32)A!"

while getopts 'p:mh:' OPTION; do
    case "$OPTION" in
        p)
            POD=$OPTARG
            ;;
    esac
done

if [[ -z "$POD" ]]; then
    echo "!!! DOCKER MODE"
    # Create admin user OR change password if that failed (e.g. if we executed this already)
    docker exec -it gitea sh -c "gitea admin create-user --username ${USERNAME} --password ${PASSWORD} --email support@grchive.com --must-change-password=false --admin || gitea admin change-password --username ${USERNAME} --password ${PASSWORD}"
else
    echo "!!! KUBE MODE -- $POD"
    kubectl exec -it $POD -- sh -c "gitea admin create-user --username ${USERNAME} --password ${PASSWORD} --email support@grchive.com --must-change-password=false --admin || gitea admin change-password --username ${USERNAME} --password ${PASSWORD}"
fi

# Create token for admin user
echo "CREATE ADMIN TOKEN"
USER_AUTH_URL="${GITEA_PROTOCOL}://${USERNAME}:${PASSWORD}@${GITEA_HOST}:${GITEA_PORT}"
RESULT=$(curl \
    -H "Content-Type: application/json" \
    -H "accept: application/json" \
    -X POST \
    -d '{"name": "grchive-global-token"}' \
    "${USER_AUTH_URL}/api/v1/users/${USERNAME}/tokens")
echo $RESULT
TOKEN=$(echo $RESULT | jq -j '.sha1')

# Create global organization
echo "CREATE ORG"
ORG_RESULT=$(curl \
    -H "Content-Type: application/json" \
    -H "accept: application/json" \
    -H "Authorization: token ${TOKEN}" \
    -X POST \
    -d "{\"username\": \"${GITEA_GLOBAL_ORG}\" }" \
    "${GITEA_PROTOCOL}://${GITEA_HOST}:${GITEA_PORT}/api/v1/admin/users/${USERNAME}/orgs")
echo $ORG_RESULT

echo "STORE CREDENTIALS"
export VAULT_SKIP_VERIFY=1
vault login -address="${VAULT_HOST}:${VAULT_PORT}" ${VAULT_TOKEN}
vault kv put -address="${VAULT_HOST}:${VAULT_PORT}" ${GITEA_TOKEN} token="${TOKEN}" password="${PASSWORD}"

echo $TOKEN
