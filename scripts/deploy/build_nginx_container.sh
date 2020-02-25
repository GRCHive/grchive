#!/bin/bash
while getopts 'm' OPTION; do
    case "$OPTION" in
        m)
            MINIKUBE=1
            ;;
    esac
done

bazel run -c opt --platforms=//build:k8s //devops/docker/nginx:nginx 

if [[ -z $MINIKUBE ]]; then
    FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/nginx:`git rev-parse HEAD`
    docker tag bazel/devops/docker/nginx:nginx $FULL_IMAGE_URL
    docker push $FULL_IMAGE_URL
fi
