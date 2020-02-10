#!/bin/bash
FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/vault:`git rev-parse HEAD`
bazel run -c opt //devops/docker/vault:vault 
docker tag bazel/devops/docker/vault:vault $FULL_IMAGE_URL
docker push $FULL_IMAGE_URL
