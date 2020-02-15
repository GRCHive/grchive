#!/bin/bash
FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/database_query_runner:`git rev-parse HEAD`
bazel build -c opt //devops/docker/database_query_runner:docker_database_query_runner.tar
docker load -i bazel-bin/devops/docker/database_query_runner/docker_database_query_runner.tar
docker tag bazel/devops/docker/database_query_runner:docker_database_query_runner $FULL_IMAGE_URL
docker push $FULL_IMAGE_URL
