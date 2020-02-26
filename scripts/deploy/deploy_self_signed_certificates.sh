#!/bin/bash
set -xe
DIR=`dirname $0`

ROOT_CERT_DIR=$(mktemp -d)
openssl genrsa -out $ROOT_CERT_DIR/rootca.key 4096
openssl req -x509 \
    -new \
    -nodes \
    -key $ROOT_CERT_DIR/rootca.key \
    -sha512 \
    -days 3650 \
    -out $ROOT_CERT_DIR/rootca.crt \
    -subj "/C=US/ST=NJ/L=Livingston/O=GRCHive/OU=GRCHive/CN=grchive-rootca"

kubectl create secret generic rootca \
    --from-file=${ROOT_CERT_DIR}/rootca.crt -o yaml --dry-run --save-config | kubectl apply -f -

$DIR/deploy_single_signed_cert.sh rabbitmq-service $ROOT_CERT_DIR/rootca
$DIR/deploy_single_signed_cert.sh internal-vault-service $ROOT_CERT_DIR/rootca
$DIR/deploy_single_signed_cert.sh query-runner-service $ROOT_CERT_DIR/rootca

rm -rf $ROOT_CERT_DIR
