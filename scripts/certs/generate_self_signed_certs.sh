#!/bin/bash

while getopts 'd:c:p:' OPTION; do
    case "$OPTION" in
        d)
            OUTDIR=$OPTARG
            ;;
        c)
            CA_CERT_BASE=$OPTARG
            ;;
        p)
            PREFIX=$OPTARG
            ;;
    esac
done

set -xe
mkdir -p ${OUTDIR}

openssl genrsa -out $OUTDIR/$PREFIX.key 4096
openssl req -new \
    -sha512 \
    -key $OUTDIR/$PREFIX.key \
    -subj "/C=US/ST=NJ/L=Livingston/O=GRCHive/OU=GRCHive/CN=${PREFIX}" \
    -out $OUTDIR/$PREFIX.csr
openssl x509 -req -in $OUTDIR/$PREFIX.csr \
    -sha512 \
    -days 3650 \
    -CA $CA_CERT_BASE.crt \
    -CAkey $CA_CERT_BASE.key \
    -CAcreateserial \
    -out $OUTDIR/$PREFIX.crt
