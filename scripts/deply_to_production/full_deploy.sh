#!/bin/bash

DIR=`dirname $0`
. ${DIR}/pull_env_variables.sh

set -xe

envsubst < build/variables.bzl.prod.tmpl > build/variables.bzl
echo $GCLOUD_WEBSERVER_ACCOUNT > devops/gcloud/gcloud-webserver-account.json
echo $GCLOUD_TERRAFORM_ACCOUNT > devops/gcloud/gcloud-terraform-account.json
echo $GCLOUD_KUBERNETES_ACCOUNT > devops/gcloud/gcloud-kubernetes-account.json

gcloud auth activate-service-account --key-file devops/gcloud/gcloud-kubernetes-account.json
gcloud config set project grchive
gcloud config set compute/zone us-central1-c

${DIR}/build_nginx_container.sh
${DIR}/build_rabbitmq_container.sh
${DIR}/build_vault_container.sh
${DIR}/build_preview_generator_container.sh
${DIR}/build_webserver_container.sh

${DIR}/terraform.sh
${DIR}/deploy_k8s.sh
