#!/bin/bash
FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/webserver:`git rev-parse HEAD`
bazel build -c opt --platforms=//build:k8s //devops/docker/webserver:docker_webserver.tar
docker load -i bazel-bin/devops/docker/webserver/docker_webserver.tar
docker tag bazel/devops/docker/webserver:docker_webserver $FULL_IMAGE_URL
docker push $FULL_IMAGE_URL
