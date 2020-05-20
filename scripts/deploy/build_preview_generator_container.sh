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

bazel build ${BUILD_OPT} --platforms=//build:k8s //devops/docker/preview_generator:docker_preview_generator.tar
docker load -i bazel-bin/devops/docker/preview_generator/docker_preview_generator.tar

if [[ -z $MINIKUBE ]]; then
    FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/${CONTAINER_REGISTRY_FOLDER}/preview_generator:`git rev-parse HEAD`
    docker tag bazel/devops/docker/preview_generator:docker_preview_generator $FULL_IMAGE_URL
    docker push $FULL_IMAGE_URL
fi
