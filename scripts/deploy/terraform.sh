#!/bin/bash
cd devops/terraform/${TERRAFORM_FOLDER}
export TF_VAR_postgres_user=${POSTGRES_USER}
export TF_VAR_postgres_password=${POSTGRES_PASSWORD}
export TF_VAR_postgres_instance_name=${POSTGRES_INSTANCE_NAME}
terraform init
terraform apply
cd ../../../

cd devops/database/vault
envsubst < flyway/prod-flyway.conf.tmpl > flyway/prod-flyway.conf
flyway -configFiles=./flyway/prod-flyway.conf -url="jdbc:postgresql://localhost:5555/vault" migrate
cd ../../../

cd devops/database/webserver
envsubst < flyway/prod-flyway.conf.tmpl > flyway/prod-flyway.conf
flyway -configFiles=./flyway/prod-flyway.conf -url="jdbc:postgresql://localhost:5555/audit" migrate
cd ../../../
