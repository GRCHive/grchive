#!/bin/bash
BUILD_OPT="-c opt"
while getopts 'm' OPTION; do
    case "$OPTION" in
        m)
            MINIKUBE=1
            BUILD_OPT=""
            ;;
    esac
done

bazel build ${BUILD_OPT} --platforms=//build:k8s //devops/docker/notification_hub:docker_notification_hub.tar
docker load -i bazel-bin/devops/docker/notification_hub/docker_notification_hub.tar

if [[ -z $MINIKUBE ]]; then
    FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/${CONTAINER_REGISTRY_FOLDER}/notification_hub:`git rev-parse HEAD`
    docker tag bazel/devops/docker/notification_hub:docker_notification_hub $FULL_IMAGE_URL
    docker push $FULL_IMAGE_URL
fi
