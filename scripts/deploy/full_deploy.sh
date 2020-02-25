#!/bin/bash

DIR=`dirname $0`

while getopts 'e:' OPTION; do
    case "$OPTION" in
        e)
            ENV=$OPTARG
            ;;
    esac
done

set -xe

DO_TERRAFORM=0
USE_ENV_VARIABLES=0
EXTRA_BUILD_OPTIONS=""
DEPLOY_GCLOUD=0

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

        USE_ENV_VARIABLES=1
        DO_TERRAFORM=1
        DEPLOY_GCLOUD=1
        ;;

    minikube)
        USE_ENV_VARIABLES=0
        DO_TERRAFORM=0
        DEPLOY_GCLOUD=0

        HOST_STATUS=$(minikube status | grep host)
        if [[ "$HOST_STATUS" == *"Stopped"* ]];
            echo "Minikube is not running."
            return 1
        fi

        EXTRA_BUILD_OPTIONS="-m"
        ;;
esac

if [[ ! -z $USE_ENV_VARIABLES ]]; then
    envsubst < build/variables.bzl.prod.tmpl > build/variables.bzl
fi

${DIR}/build_nginx_container.sh ${EXTRA_BUILD_OPTIONS}
${DIR}/build_rabbitmq_container.sh ${EXTRA_BUILD_OPTIONS}
${DIR}/build_vault_container.sh ${EXTRA_BUILD_OPTIONS}
${DIR}/build_preview_generator_container.sh ${EXTRA_BUILD_OPTIONS}
${DIR}/build_webserver_container.sh ${EXTRA_BUILD_OPTIONS}
${DIR}/build_database_refresh_worker.sh ${EXTRA_BUILD_OPTIONS}
${DIR}/build_database_runner_worker.sh ${EXTRA_BUILD_OPTIONS}

if [[ ! -z $DO_TERRAFORM ]]; then
    gcloud auth activate-service-account --key-file devops/gcloud/gcloud-terraform-account.json
    gcloud config set project ${GRCHIVE_PROJECT}
    gcloud config set compute/zone us-central1-c

    cloud_sql_proxy -instances=${GRCHIVE_PROJECT}:us-central1:${POSTGRES_INSTANCE_NAME}=tcp:5555 &
    PROXY_PID=$!

    ${DIR}/terraform.sh

    kill -9 $PROXY_PID
fi

if [[ ! -z $DEPLOY_GCLOUD ]]; then
    gcloud auth activate-service-account --key-file devops/gcloud/gcloud-kubernetes-account.json
    gcloud config set project ${GRCHIVE_PROJECT}
    gcloud config set compute/zone us-central1-c
    gcloud container clusters get-credentials webserver-gke-cluster
fi

${DIR}/deploy_self_signed_certificates.sh
${DIR}/deploy_k8s.sh
