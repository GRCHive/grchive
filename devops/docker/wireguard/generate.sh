#!/bin/bash
WG_FOLDER="etc/wireguard"

mkdir -p ${WG_FOLDER}
envsubst < wg0.conf.tmpl > ${WG_FOLDER}/wg0.conf
cat ${WG_FOLDER}/wg0.conf

docker build --tag registry.gitlab.com/grchive/grchive/wireguard:latest .

rm -rf ${WG_FOLDER}
