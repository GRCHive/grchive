#!/bin/bash

DIR=`dirname $0`

. ${DIR}/../deploy/pull_env_variables.sh ${DIR}/../deploy
echo $GCLOUD_WEBSERVER_ACCOUNT > devops/gcloud/gcloud-webserver-account.json
echo $GCLOUD_TERRAFORM_ACCOUNT > devops/gcloud/gcloud-terraform-account.json
echo $GCLOUD_KUBERNETES_ACCOUNT > devops/gcloud/gcloud-kubernetes-account.json

envsubst < build/variables.bzl.prod.tmpl > build/variables.bzl

bazel run ${BUILD_OPT} --platforms=//build:k8s //devops/docker/wordpress:latest
FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/wordpress:latest
docker tag bazel/devops/docker/wordpress:latest $FULL_IMAGE_URL
docker push $FULL_IMAGE_URL

cd devops/docker/wordpress/nginx
NGINX_URL=registry.gitlab.com/grchive/grchive/wordpress/nginx:latest
docker build --tag $NGINX_URL .
docker push $NGINX_URL
