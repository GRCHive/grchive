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

bazel build ${BUILD_OPT} --platforms=//build:k8s //devops/docker/script_runner:latest.tar
docker load -i bazel-bin/devops/docker/script_runner/latest.tar

bazel build ${BUILD_OPT} --platforms=//build:k8s //devops/docker/script_runner_worker_image:latest.tar
docker load -i bazel-bin/devops/docker/script_runner_worker_image/latest.tar

if [[ -z $MINIKUBE ]]; then
    FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/${CONTAINER_REGISTRY_FOLDER}/script_runner:`git rev-parse HEAD`
    docker tag bazel/devops/docker/script_runner:latest $FULL_IMAGE_URL
    docker push $FULL_IMAGE_URL

    FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/${CONTAINER_REGISTRY_FOLDER}/script_runner_worker_image:latest
    docker tag bazel/devops/docker/script_runner_worker_image:latest $FULL_IMAGE_URL
    docker push $FULL_IMAGE_URL
fi
