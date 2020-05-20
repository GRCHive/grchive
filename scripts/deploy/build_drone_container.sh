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

bazel run ${BUILD_OPT} --platforms=//build:k8s //devops/docker/drone:drone-build
bazel run ${BUILD_OPT} --platforms=//build:k8s //devops/docker/drone_runner:drone-runner-k8s
bazel run ${BUILD_OPT} --platforms=//build:k8s //devops/docker/drone_runner_worker_image:latest

if [[ -z $MINIKUBE ]]; then
    FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/${CONTAINER_REGISTRY_FOLDER}/drone:`git rev-parse HEAD`
    docker tag bazel/devops/docker/drone:drone $FULL_IMAGE_URL
    docker push $FULL_IMAGE_URL

    FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/${CONTAINER_REGISTRY_FOLDER}/drone_runner_k8s:`git rev-parse HEAD`
    docker tag bazel/devops/docker/drone_runner:drone-runner-k8s $FULL_IMAGE_URL
    docker push $FULL_IMAGE_URL

    FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/${CONTAINER_REGISTRY_FOLDER}/drone_runner_worker_image:latest
    docker tag bazel/devops/docker/drone_runner_worker_image:latest $FULL_IMAGE_URL
    docker push $FULL_IMAGE_URL
fi
