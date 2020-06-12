#!/bin/bash

DIR=`dirname $0`

. ${DIR}/../deploy/pull_env_variables.sh ${DIR}/../deploy
echo $GCLOUD_WEBSERVER_ACCOUNT > devops/gcloud/gcloud-webserver-account.json
echo $GCLOUD_TERRAFORM_ACCOUNT > devops/gcloud/gcloud-terraform-account.json
echo $GCLOUD_KUBERNETES_ACCOUNT > devops/gcloud/gcloud-kubernetes-account.json

gcloud compute ssh grchive-wordpress-central1-c --zone=us-central1-c --command="\
docker login registry.gitlab.com --username ${GKE_REGISTRY_USER} --password ${GKE_REGISTRY_PASSWORD};
docker pull registry.gitlab.com/grchive/grchive/wordpress/nginx:latest;
docker pull registry.gitlab.com/grchive/grchive/wordpress:latest;
sudo mkdir -p /mnt/stateful_partition/wordpress/html;
sudo chown -R `whoami` /mnt/stateful_partition/wordpress;
"

gcloud compute scp devops/docker/wordpress/docker-compose.yml grchive-wordpress-central1-c:/mnt/stateful_partition/wordpress/docker-compose.yml --zone=us-central1-c


gcloud compute ssh grchive-wordpress-central1-c --zone=us-central1-c --command='\
cd /mnt/stateful_partition/wordpress;
docker run -d --name wp_compose --rm -v /var/run/docker.sock:/var/run/docker.sock -v $PWD:$PWD -w=$PWD docker/compose:1.24.0 up;
'
