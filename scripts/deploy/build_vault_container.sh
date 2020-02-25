#!/bin/bash
while getopts 'm' OPTION; do
    case "$OPTION" in
        m)
            MINIKUBE=1
            ;;
    esac
done

bazel run -c opt --platforms=//build:k8s //devops/docker/vault:vault 

if [[ -z $MINIKUBE ]]; then
    FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/vault:`git rev-parse HEAD`
    docker tag bazel/devops/docker/vault:vault $FULL_IMAGE_URL
    docker push $FULL_IMAGE_URL
fi
