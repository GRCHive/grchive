#!/bin/bash

MAJOR_VERSION=279
MINOR_VERSION=0
PATCH_VERSION=0
FULL_VERSION=${MAJOR_VERSION}.${MINOR_VERSION}.${PATCH_VERSION}
URL="https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${FULL_VERSION}-linux-x86_64.tar.gz"

curl -o gcloud.tar.gz $URL
mkdir -p gcloud
tar xvf gcloud.tar.gz -C gcloud 
cd gcloud
./google-cloud-sdk/install.sh
cd ../
rm gcloud.tar.gz

curl -o gcloud/cloud_sql_proxy https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64
chmod +x ./gcloud/cloud_sql_proxy
