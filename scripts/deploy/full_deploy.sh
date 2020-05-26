#!/bin/bash

DIR=`dirname $0`

while getopts 'e:dbwk' OPTION; do
    case "$OPTION" in
        e)
            ENV=$OPTARG
            ;;
        d)
            NODEPLOY=1
            ;;
        b)
            NOBUILD=1
            ;;
        w)
            NOWIREGUARD=1
            ;;
        k)
            NOK8S=1
            ;;
    esac
done

set -xe

EXTRA_BUILD_OPTIONS=""

case "$ENV" in
    prod)
        . ${DIR}/pull_env_variables.sh
        echo $GCLOUD_WEBSERVER_ACCOUNT > devops/gcloud/gcloud-webserver-account.json
        echo $GCLOUD_TERRAFORM_ACCOUNT > devops/gcloud/gcloud-terraform-account.json
        echo $GCLOUD_KUBERNETES_ACCOUNT > devops/gcloud/gcloud-kubernetes-account.json
        GRCHIVE_PROJECT="grchive"
        export OKTA_CLIENT_ID="0oa1n0o8fPR0iSsIC4x6"
        export OKTA_CLIENT_SECRET=$PRODUCTION_OKTA_CLIENT_SECRET
        export GRCHIVE_PROJECT="grchive"
        export GRCHIVE_URI="https://www.grchive.com"
        export GRCHIVE_DOMAIN="grchive.com"
        export GRCHIVE_DOC_BUCKET="grchive-prod"
        export TERRAFORM_FOLDER="prod"
        export INGRESS_ENV="prod"
        export ARTIFACTORY_ENCRYPTED_PASSWORD=$PRODUCTION_ARTIFACTORY_ENCRYPTED_PASSWORD
        export DRONE_TOKEN=$PRODUCTION_DRONE_TOKEN
        export CONTAINER_REGISTRY_FOLDER="prod"

        USE_ENV_VARIABLES=1
        DO_TERRAFORM=1
        DEPLOY_GCLOUD=1
        ;;

    staging)
        . ${DIR}/pull_env_variables.sh
        echo $STAGING_GCLOUD_WEBSERVER_ACCOUNT > devops/gcloud/gcloud-webserver-account.json
        echo $STAGING_GCLOUD_TERRAFORM_ACCOUNT > devops/gcloud/gcloud-terraform-account.json
        echo $STAGING_GCLOUD_KUBERNETES_ACCOUNT > devops/gcloud/gcloud-kubernetes-account.json
        GRCHIVE_PROJECT="grchive-staging"
        export GRCHIVE_PROJECT=$STAGING_GRCHIVE_PROJECT
        export OKTA_CLIENT_ID="0oa25j979s1Txrkiz4x6"
        export OKTA_CLIENT_SECRET=$STAGING_OKTA_CLIENT_SECRET
        export GRCHIVE_PROJECT="grchive-staging"
        export GRCHIVE_URI="https://staging.grchive.com"
        export GRCHIVE_DOMAIN="staging.grchive.com"
        export GRCHIVE_DOC_BUCKET="grchive-staging"
        export TERRAFORM_FOLDER="staging"
        export INGRESS_ENV="staging"
        export ARTIFACTORY_ENCRYPTED_PASSWORD=$STAGING_ARTIFACTORY_ENCRYPTED_PASSWORD
        export DRONE_TOKEN=$STAGING_DRONE_TOKEN
        export CONTAINER_REGISTRY_FOLDER="staging"

        USE_ENV_VARIABLES=1
        DO_TERRAFORM=1
        DEPLOY_GCLOUD=1
        ;;

    minikube)
        HOST_STATUS=$(minikube status | grep host)
        if [[ "$HOST_STATUS" == *"Stopped"* ]]; then
            echo "Minikube is not running."
            exit 1
        fi
        NOWIREGUARD=1

        EXTRA_BUILD_OPTIONS="-m"
        eval $(minikube docker-env)
        ;;
esac

export DRONE_RUNNER_IMAGE="${CONTAINER_REGISTRY}/${CONTAINER_REGISTRY_FOLDER}/drone_runner_worker_image:latest"
export SCRIPT_RUNNER_IMAGE="${CONTAINER_REGISTRY}/${CONTAINER_REGISTRY_FOLDER}/script_runner_worker_image:latest"

if [[ ! -z "$USE_ENV_VARIABLES" ]]; then
    envsubst < build/variables.bzl.prod.tmpl > build/variables.bzl
fi

if [[ -z "$NOBUILD" ]]; then
    ${DIR}/build_gitea_container.sh ${EXTRA_BUILD_OPTIONS}
    ${DIR}/build_artifactory_container.sh ${EXTRA_BUILD_OPTIONS}
    ${DIR}/build_drone_container.sh ${EXTRA_BUILD_OPTIONS}
    ${DIR}/build_nginx_container.sh ${EXTRA_BUILD_OPTIONS}
    ${DIR}/build_rabbitmq_container.sh ${EXTRA_BUILD_OPTIONS}
    ${DIR}/build_vault_container.sh ${EXTRA_BUILD_OPTIONS}
    ${DIR}/build_preview_generator_container.sh ${EXTRA_BUILD_OPTIONS}
    ${DIR}/build_webserver_container.sh ${EXTRA_BUILD_OPTIONS}
    ${DIR}/build_database_refresh_worker.sh ${EXTRA_BUILD_OPTIONS}
    ${DIR}/build_database_runner_worker.sh ${EXTRA_BUILD_OPTIONS}
    ${DIR}/build_notification_hub.sh ${EXTRA_BUILD_OPTIONS}
    ${DIR}/build_script_runner_container.sh ${EXTRA_BUILD_OPTIONS}
    ${DIR}/build_task_manager_container.sh ${EXTRA_BUILD_OPTIONS}
fi

if [[ ! -z "$DO_TERRAFORM" ]]; then
    gcloud auth activate-service-account --key-file devops/gcloud/gcloud-terraform-account.json
    gcloud config set project ${GRCHIVE_PROJECT}
    gcloud config set compute/zone us-central1-c

    cloud_sql_proxy -instances=${GRCHIVE_PROJECT}:us-central1:${POSTGRES_INSTANCE_NAME}=tcp:5555 &
    PROXY_PID=$!

    ${DIR}/terraform.sh

    # Do this twice - sometimes if the kubernetes is cluster is destroyed, recreating it won't
    # also re-create the node pool.
    ${DIR}/terraform.sh

    kill -9 $PROXY_PID
fi

if [[ -z "$NODEPLOY" ]]; then
    if [[ ! -z "$DEPLOY_GCLOUD" ]]; then
        gcloud auth activate-service-account --key-file devops/gcloud/gcloud-kubernetes-account.json
        gcloud config set project ${GRCHIVE_PROJECT}
        gcloud config set compute/zone us-central1-c
        gcloud container clusters get-credentials webserver-gke-cluster
    fi

    if [[ -z "$NOWIREGUARD" ]]; then
        ${DIR}/deploy_wireguard.sh
        sleep 5
    fi

    if [[ -z "$NOK8S" ]]; then
        sudo wg-quick up wg0-client

        ${DIR}/deploy_k8s.sh ${EXTRA_BUILD_OPTIONS}

        sudo wg-quick down wg0-client
    fi
fi
