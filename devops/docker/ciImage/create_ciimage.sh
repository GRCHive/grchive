#!/bin/bash

IMAGE_NAME=$1
DEPENDENCIES="bootstrap_bazel.sh download_libffi.sh download_python.sh download_terraform.sh download_flyway.sh download_kubectl.sh python-helper"

# Copy needed files
for DEPEND in ${DEPENDENCIES};
do
    cp -r ../../../dependencies/${DEPEND} .
done

docker build --tag ${IMAGE_NAME} .
docker push ${IMAGE_NAME}

# Cleanup
for DEPEND in ${DEPENDENCIES};
do
    rm -rf ${DEPEND}
done
