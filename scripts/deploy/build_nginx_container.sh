#!/bin/bash
FULL_IMAGE_URL=registry.gitlab.com/grchive/grchive/nginx:`git rev-parse HEAD`
bazel run -c opt --platforms=//build:k8s //devops/docker/nginx:nginx 
docker tag bazel/devops/docker/nginx:nginx $FULL_IMAGE_URL
docker push $FULL_IMAGE_URL
