#!/bin/bash
docker start psql
sleep 5

docker start vault
vault operator unseal -address="http://172.22.0.3:8200"

docker start rabbitmq

docker start nfssrv
sleep 5

docker start gitea
docker start artifactory
docker start drone drone-runner
