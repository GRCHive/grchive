#!/bin/bash
while getopts 'm' OPTION; do
    case "$OPTION" in
        m)
            MINIKUBE=1
            ;;
    esac
done

bazel build -c opt --platforms=//build:k8s //devops/docker/database_fetcher:docker_database_fetcher.tar
docker load -i bazel-bin/devops/docker/database_fetcher/docker_database_fetcher.tar

if [[ -z $MINIKUBE ]]; then
    FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/database_fetcher:`git rev-parse HEAD`
    docker tag bazel/devops/docker/database_fetcher:docker_database_fetcher $FULL_IMAGE_URL
    docker push $FULL_IMAGE_URL
fi
