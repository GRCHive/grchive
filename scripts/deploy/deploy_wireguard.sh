#!/bin/bash
DIR=`dirname $0`
set -xe

envsubst < ${DIR}/../../devops/wireguard/wg0.conf.tmpl > wg0.conf

gcloud --quiet compute scp ${DIR}/prepare_wireguard_server.sh grchive-wireguard-central1-c:~/
gcloud --quiet compute scp wg0.conf grchive-wireguard-central1-c:~/

gssh="gcloud compute ssh grchive-wireguard-central1-c"
$gssh --command "bash ~/prepare_wireguard_server.sh"
$gssh --command "sudo cp ~/wg0.conf /etc/wireguard/wg0.conf; sudo wg-quick down wg0; sudo wg-quick up wg0"

rm wg0.conf
