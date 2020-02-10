#!/bin/bash

DIR=`dirname $0`
. ${DIR}/pull_env_variables.sh

while getopts 'e:' OPTION; do
    case "$OPTION" in
        e)
            ENV=$OPTARG
            ;;
    esac
done

set -xe

case "$ENV" in
    prod)
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
        ;;

    staging)
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
        ;;
esac

envsubst < build/variables.bzl.prod.tmpl > build/variables.bzl

${DIR}/build_nginx_container.sh
${DIR}/build_rabbitmq_container.sh
${DIR}/build_vault_container.sh
${DIR}/build_preview_generator_container.sh
${DIR}/build_webserver_container.sh

gcloud auth activate-service-account --key-file devops/gcloud/gcloud-terraform-account.json
gcloud config set project ${GRCHIVE_PROJECT}
gcloud config set compute/zone us-central1-c

cloud_sql_proxy -instances=${GRCHIVE_PROJECT}:us-central1:${POSTGRES_INSTANCE_NAME}=tcp:5555 &
PROXY_PID=$!

${DIR}/terraform.sh

gcloud auth activate-service-account --key-file devops/gcloud/gcloud-kubernetes-account.json
gcloud config set project ${GRCHIVE_PROJECT}
gcloud config set compute/zone us-central1-c
${DIR}/deploy_k8s.sh

kill -9 $PROXY_PID
