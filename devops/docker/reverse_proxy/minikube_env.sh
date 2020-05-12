#!/bin/bash
export MINIKUBE_IP=$(minikube ip)
export MINIKUBE_GITEA_PORT=$(kubectl get services -l app=gitea | grep external | awk '{split($5, a, ":");split(a[2],b,"/");print b[1];}')
