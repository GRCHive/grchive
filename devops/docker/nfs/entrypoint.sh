#!/bin/bash
mv /exports /etc/exports

rpcbind
service nfs-kernel-server start

function ctrl_c() {
    exit 0
}
trap ctrl_c SIGTERM SIGINT

while true; do sleep 1; done
