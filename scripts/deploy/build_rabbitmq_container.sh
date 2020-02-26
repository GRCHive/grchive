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

bazel run ${BUILD_OPT} --platforms=//build:k8s //devops/docker/rabbitmq:rabbitmq 

if [[ -z $MINIKUBE ]]; then
    FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/rabbitmq:`git rev-parse HEAD`
    docker tag bazel/devops/docker/rabbitmq:rabbitmq $FULL_IMAGE_URL
    docker push $FULL_IMAGE_URL
fi
