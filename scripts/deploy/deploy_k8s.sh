#!/bin/bash

gcloud container clusters get-credentials webserver-gke-cluster
kubectl create secret generic gke-service-account --from-file=gcloud-service-account.json=devops/gcloud/gcloud-kubernetes-account.json -o yaml --dry-run --save-config | kubectl apply -f -
kubectl create secret docker-registry regcred --docker-server=registry.gitlab.com --docker-username=${GKE_REGISTRY_USER} --docker-password=${GKE_REGISTRY_PASSWORD} -o yaml --dry-run --save-config | kubectl apply -f -

cd devops/k8s

kubectl apply -f storage/prod

kubectl apply -f ./cert-manager/letsencrypt-staging.yaml -f ./cert-manager/letsencrypt-prod.yaml

export VAULT_IMAGE=registry.gitlab.com/grchive/grchive/vault:`git rev-parse HEAD`
cd vault
envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
kubectl apply -f service-internal.yaml -f deployment.prod.yaml
cd ../

export RABBITMQ_IMAGE=registry.gitlab.com/grchive/grchive/rabbitmq:`git rev-parse HEAD`
cd rabbitmq
envsubst < statefulset.prod.yaml.tmpl > statefulset.prod.yaml
kubectl apply -f service.yaml -f statefulset.prod.yaml
cd ../

export PREVIEW_IMAGE=registry.gitlab.com/grchive/grchive/preview_generator:`git rev-parse HEAD`
cd preview_generator
envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
kubectl apply -f deployment.prod.yaml
cd ../

export WEBSERVER_IMAGE=registry.gitlab.com/grchive/grchive/webserver:`git rev-parse HEAD`
export NGINX_IMAGE=registry.gitlab.com/grchive/grchive/nginx:`git rev-parse HEAD`
cd webserver
envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
kubectl apply -f deployment.prod.yaml -f service.prod.yaml -f ingress.prod.yaml -f backendconfig.prod.yaml
cd ../
