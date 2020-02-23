#!/bin/bash
set -xe
DIR=`dirname $0`

PREFIX=$1
ROOT_BASE=$2
DOMAIN=$3

CERT_DIR=$(mktemp -d)
echo "Generating ${PREFIX} Certificates into ${CERT_DIR}..."
${DIR}/../certs/generate_self_signed_certs.sh -c ${ROOT_BASE} -d ${CERT_DIR} -p ${PREFIX} -n ${DOMAIN}
kubectl create secret generic ${PREFIX}-certificate \
    --from-file=${CERT_DIR}/${PREFIX}.crt \
    --from-file=${CERT_DIR}/${PREFIX}.key \
    -o yaml --dry-run --save-config | kubectl apply -f -
rm -rf ${CERT_DIR}
