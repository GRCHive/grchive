#!/bin/bash
FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/rabbitmq:`git rev-parse HEAD`
bazel run -c opt --platforms=//build:k8s //devops/docker/rabbitmq:rabbitmq 
docker tag bazel/devops/docker/rabbitmq:rabbitmq $FULL_IMAGE_URL
docker push $FULL_IMAGE_URL
