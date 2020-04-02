#!/bin/bash

# Need to make things not symlinks for Docker to be happy.
rm -rf devops/docker/drone-real
mkdir -p devops/docker/drone-real

cp scripts/vault/get_vault_secret_http.sh devops/docker/drone-real

cd devops/docker/drone-real
cp -RL ../drone/* .

docker build -t bazel/devops/docker/drone:drone .

cd ../
rm -rf drone-real
