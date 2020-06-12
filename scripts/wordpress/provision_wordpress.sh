#!/bin/bash

DIR=`dirname $0`

. ${DIR}/../deploy/pull_env_variables.sh ${DIR}/../deploy
echo $GCLOUD_WEBSERVER_ACCOUNT > devops/gcloud/gcloud-webserver-account.json
echo $GCLOUD_TERRAFORM_ACCOUNT > devops/gcloud/gcloud-terraform-account.json
echo $GCLOUD_KUBERNETES_ACCOUNT > devops/gcloud/gcloud-kubernetes-account.json

cd devops/terraform/blog
export TF_VAR_wp_database_user=${WORDPRESS_DB_USER}
export TF_VAR_wp_database_password=${WORDPRESS_DB_PASSWORD}
export TF_VAR_wp_database_name=${WORDPRESS_DB_NAME}
export TF_VAR_wp_instance_name=${WORDPRESS_INSTANCE_NAME}
terraform init
terraform apply
cd ../../../
