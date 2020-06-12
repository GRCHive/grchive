#!/bin/bash
while getopts 'm' OPTION; do
    case "$OPTION" in
        m)
            MINIKUBE=1
            ;;
    esac
done

kubectl create secret generic gke-service-account --from-file=gcloud-service-account.json=devops/gcloud/gcloud-kubernetes-account.json -o yaml --dry-run --save-config | kubectl apply -f -
kubectl create secret generic gke-service-account --from-file=gcloud-service-account.json=devops/gcloud/gcloud-kubernetes-account.json -o yaml --dry-run --save-config | kubectl apply --namespace backend -f -

cd devops/k8s

kubectl apply -f backendNamespace.yaml
kubectl label namespace default istio-injection=enabled
kubectl label namespace backend istio-injection=enabled

kubectl apply -f istio/mtls.yaml

if [ -z $MINIKUBE ]; then
    kubectl create secret docker-registry regcred --docker-server=registry.gitlab.com --docker-username=${GKE_REGISTRY_USER} --docker-password=${GKE_REGISTRY_PASSWORD} -o yaml --dry-run --save-config | kubectl apply -f -
    kubectl create secret docker-registry regcred --docker-server=registry.gitlab.com --docker-username=${GKE_REGISTRY_USER} --docker-password=${GKE_REGISTRY_PASSWORD} -o yaml --dry-run --save-config | kubectl apply --namespace backend -f -
    DEV_PROD="prod"
else
    DEV_PROD="dev"
fi

kubectl apply -f storage/${DEV_PROD}

if [ -z $MINIKUBE ]; then
    kubectl apply -f ./cert-manager/letsencrypt-staging.yaml -f ./cert-manager/letsencrypt-prod.yaml

    cd cloud_sql_proxy
    envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
    kubectl apply -f service.yaml -f deployment.prod.yaml
    cd ../

    CONTAINER_REGISTRY_URL=${CONTAINER_REGISTRY}/${CONTAINER_REGISTRY_FOLDER}

    if [[ ! -z $BASH_DISABLE_DASHBOARD ]]; then
        export VAULT_IMAGE=${CONTAINER_REGISTRY_URL}/vault:`git rev-parse HEAD`
        cd vault
        envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
        kubectl apply -f service-internal.yaml -f deployment.prod.yaml
        cd ../
    fi

    if [[ ! -z $BASH_DISABLE_DASHBOARD ]]; then
        export GITEA_IMAGE=${CONTAINER_REGISTRY_URL}/gitea:`git rev-parse HEAD`
        cd gitea
        envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
        kubectl apply -f service.yaml -f deployment.prod.yaml
        cd ../
    fi

    if [[ ! -z $BASH_DISABLE_DASHBOARD ]]; then
        export ARTIFACTORY_IMAGE=${CONTAINER_REGISTRY_URL}/artifactory:`git rev-parse HEAD`
        cd artifactory
        envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
        kubectl apply -f service.yaml -f deployment.prod.yaml
        cd ../
    fi

    if [[ ! -z $BASH_DISABLE_DASHBOARD ]]; then
        export DRONE_IMAGE=${CONTAINER_REGISTRY_URL}/drone:`git rev-parse HEAD`
        cd drone
        envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
        kubectl apply -f service.yaml -f deployment.prod.yaml
        cd ../
    fi

    if [[ ! -z $BASH_DISABLE_DASHBOARD ]]; then
        export DRONE_RUNNER_IMAGE=${CONTAINER_REGISTRY_URL}/drone_runner_k8s:`git rev-parse HEAD`
        cd drone_runner
        envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
        kubectl apply -f role.yaml -f deployment.prod.yaml
        cd ../
    fi

    export RABBITMQ_IMAGE=${CONTAINER_REGISTRY_URL}/rabbitmq:`git rev-parse HEAD`
    cd rabbitmq
    envsubst < statefulset.prod.yaml.tmpl > statefulset.prod.yaml
    kubectl apply -f service.yaml -f statefulset.prod.yaml
    cd ../

    if [[ ! -z $BASH_DISABLE_DASHBOARD ]]; then
        export PREVIEW_IMAGE=${CONTAINER_REGISTRY_URL}/preview_generator:`git rev-parse HEAD`
        cd preview_generator
        envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
        kubectl apply -f deployment.prod.yaml
        cd ../
    fi

    if [[ ! -z $BASH_DISABLE_DASHBOARD ]]; then
        export RUNNER_IMAGE=${CONTAINER_REGISTRY_URL}/database_query_runner:`git rev-parse HEAD`
        cd database_query_runner
        envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
        kubectl apply -f deployment.prod.yaml -f service.yaml
        cd ../
    fi

    if [[ ! -z $BASH_DISABLE_DASHBOARD ]]; then
        export FETCHER_IMAGE=${CONTAINER_REGISTRY_URL}/database_fetcher:`git rev-parse HEAD`
        cd database_fetcher
        envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
        kubectl apply -f deployment.prod.yaml
        cd ../
    fi

    export NOTIFICATION_HUB_IMAGE=${CONTAINER_REGISTRY_URL}/notification_hub:`git rev-parse HEAD`
    cd notification_hub
    envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
    kubectl apply -f deployment.prod.yaml
    cd ../

    if [[ ! -z $BASH_DISABLE_DASHBOARD ]]; then
        export SCRIPT_RUNNER_IMAGE=${CONTAINER_REGISTRY_URL}/script_runner:`git rev-parse HEAD`
        cd script_runner
        envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
        kubectl apply -f deployment.prod.yaml -f role.yaml
        cd ../
    fi

    if [[ ! -z $BASH_DISABLE_DASHBOARD ]]; then
        export TASK_MANAGER_IMAGE=${CONTAINER_REGISTRY_URL}/task_manager:`git rev-parse HEAD`
        cd task_manager
        envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
        kubectl apply -f deployment.prod.yaml
        cd ../
    fi

    if [[ ! -z $BASH_DISABLE_DASHBOARD ]]; then
        export SHELL_RUNNER_IMAGE=${CONTAINER_REGISTRY_URL}/shell_runner:`git rev-parse HEAD`
        cd shell_runner
        envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
        kubectl apply -f deployment.prod.yaml
        cd ../
    fi

    if [[ ! -z $BASH_DISABLE_DASHBOARD ]]; then
        export INTEGRATION_RUNNER_IMAGE=${CONTAINER_REGISTRY_URL}/integration_runner:`git rev-parse HEAD`
        cd integration_runner
        envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
        kubectl apply -f deployment.prod.yaml
        cd ../
    fi

    export WEBSERVER_IMAGE=${CONTAINER_REGISTRY_URL}/webserver:`git rev-parse HEAD`
    export NGINX_IMAGE=${CONTAINER_REGISTRY_URL}/nginx:`git rev-parse HEAD`
    cd webserver
    envsubst < deployment.prod.yaml.tmpl > deployment.prod.yaml
    kubectl apply -f deployment.prod.yaml -f service.prod.yaml -f ingress.${INGRESS_ENV}.yaml -f backendconfig.prod.yaml
    cd ../
else
    kubectl delete deployment vault-deployment  
    kubectl delete deployment gitea-deployment  
    kubectl delete deployment artifactory-deployment  
    kubectl delete deployment drone-deployment  
    kubectl delete deployment drone-runner-deployment  
    kubectl delete deployment task-manager-deployment
    kubectl delete deployment script-runner-deployment
    kubectl delete statefulset rabbitmq-set
    kubectl delete deployment preview-generator-deployment
    kubectl delete deployment database-fetcher-deployment
    kubectl delete deployment database-query-runner-deployment
    kubectl delete deployment webserver-deployment
    kubectl delete deployment notification-hub-deployment

    cd postgresql
    envsubst < deployment.dev.yaml.tmpl > deployment.dev.yaml
    kubectl apply -f .
    cd ../

    cd vault
    kubectl apply -f deployment.dev.yaml -f service-internal.yaml -f service-external.dev.yaml
    VAULT_PORT=$(kubectl get services -l app=vault | grep external | awk '{print $5}' | sed 's/.*:\([0-9]*\)\/TCP/\1/')
    sleep 10
    VAULT_DEPLOYMENT=$(kubectl get pods | grep vault-deployment | awk {'print $1'})
    kubectl exec -it ${VAULT_DEPLOYMENT} -- sh -c "VAULT_SKIP_VERIFY=1 vault operator unseal -address=\"https://localhost:8200\""
    cd ../

    cd gitea
    kubectl apply -f ./deployment.dev.yaml -f ./service-external.dev.yaml -f ./service.yaml
    cd ../

    cd artifactory
    kubectl apply -f ./deployment.dev.yaml -f ./service-external.dev.yaml -f ./service.yaml
    cd ../

    cd drone
    kubectl apply -f ./deployment.dev.yaml -f ./service-external.dev.yaml -f ./service.yaml
    cd ../

    cd drone_runner
    kubectl apply -f ./deployment.dev.yaml -f ./role.yaml
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

    cd notification_hub
    kubectl apply -f deployment.dev.yaml
    cd ../

    cd task_manager
    kubectl apply -f deployment.dev.yaml
    cd ../

    cd script_runner
    kubectl apply -f deployment.dev.yaml -f role.yaml
    cd ../

    cd webserver
    kubectl apply -f deployment.dev.yaml -f loadbalancer.yaml
    cd ../
fi
