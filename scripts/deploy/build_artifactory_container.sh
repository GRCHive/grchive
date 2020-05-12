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

bazel run ${BUILD_OPT} --platforms=//build:k8s //devops/docker/artifactory:artifactory

if [[ -z $MINIKUBE ]]; then
    FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/artifactory:`git rev-parse HEAD`
    docker tag bazel/devops/docker/artifactory:artifactory $FULL_IMAGE_URL
    docker push $FULL_IMAGE_URL
fi