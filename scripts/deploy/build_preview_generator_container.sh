#!/bin/bash
FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/preview_generator:`git rev-parse HEAD`
bazel build -c opt --platforms=//build:k8s //devops/docker/preview_generator:docker_preview_generator.tar
docker load -i bazel-bin/devops/docker/preview_generator/docker_preview_generator.tar
docker tag bazel/devops/docker/preview_generator:docker_preview_generator $FULL_IMAGE_URL
docker push $FULL_IMAGE_URL
