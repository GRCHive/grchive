#!/bin/bash
while getopts 'm' OPTION; do
    case "$OPTION" in
        m)
            MINIKUBE=1
            ;;
    esac
done

kubectl create secret generic gke-service-account --from-file=gcloud-service-account.json=devops/gcloud/gcloud-kubernetes-account.json -o yaml --dry-run --save-config | kubectl apply -f -

if [ -z $MINIKUBE ]; then
    kubectl create secret docker-registry regcred --docker-server=registry.gitlab.com --docker-username=${GKE_REGISTRY_USER} --docker-password=${GKE_REGISTRY_PASSWORD} -o yaml --dry-run --save-config | kubectl apply -f -
    DEV_PROD="prod"
else
    DEV_PROD="dev"
fi

cd devops/k8s

kubectl apply -f storage/${DEV_PROD}

if [ -z $MINIKUBE ]; then
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

    export RUNNER_IMAGE=registry.gitlab.com/grchive/grchive/database_query_runner:`git rev-parse HEAD`
    cd database_query_runner
    envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
    kubectl apply -f deployment.prod.yaml -f service.yaml
    cd ../

    export FETCHER_IMAGE=registry.gitlab.com/grchive/grchive/database_fetcher:`git rev-parse HEAD`
    cd database_fetcher
    envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
    kubectl apply -f deployment.prod.yaml
    cd ../

    export WEBSERVER_IMAGE=registry.gitlab.com/grchive/grchive/webserver:`git rev-parse HEAD`
    export NGINX_IMAGE=registry.gitlab.com/grchive/grchive/nginx:`git rev-parse HEAD`
    cd webserver
    envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
    kubectl apply -f deployment.prod.yaml -f service.prod.yaml -f ingress.${INGRESS_ENV}.yaml -f backendconfig.prod.yaml
    cd ../
else
    kubectl delete deployment vault-deployment  
    kubectl delete statefulset rabbitmq-set
    kubectl delete deployment preview-generator-deployment
    kubectl delete deployment database-fetcher-deployment
    kubectl delete deployment database-query-runner-deployment
    kubectl delete deployment webserver-deployment

    cd postgresql
    cp endpoint.yaml.tmpl endpoint.yaml
    kubectl apply -f .
    cd ../

    cd vault
    kubectl apply -f deployment.dev.yaml -f service-internal.yaml -f service-external.dev.yaml
    VAULT_PORT=$(kubectl get services -l app=vault | grep external | awk '{print $5}' | sed 's/.*:\([0-9]*\)\/TCP/\1/')
    sleep 10
    VAULT_DEPLOYMENT=$(kubectl get pods | grep vault-deployment | awk {'print $1'})
    kubectl exec -it ${VAULT_DEPLOYMENT} -- sh -c "VAULT_SKIP_VERIFY=1 vault operator unseal -address=\"https://localhost:8200\""
    cd ../

    cd rabbitmq
    kubectl apply -f service.yaml -f statefulset.dev.yaml
    cd ../

    cd preview_generator
    kubectl apply -f deployment.dev.yaml
    cd ../

    cd database_query_runner
    kubectl apply -f deployment.dev.yaml -f service.yaml
    cd ../

    cd database_fetcher
    kubectl apply -f deployment.dev.yaml
    cd ../

    cd webserver
    kubectl apply -f deployment.dev.yaml -f loadbalancer.yaml
    cd ../
fi
