# GRCHive

This document will assume that the git checkout directory is set in an environment variable `$SRC`.

## Prerequisites

- G++ v?
- JDK8
- zip
- unzip
- go 1.13+ 
- PostgreSQL
- Node v12.10
- Docker
- Libreoffice

## Build Variables

The build depends on the presence on a number of properly configured variables to properly generate configuration files.
These variables are stored in `$SRC/build/variables.bzl`.
To generate this file, copy `$SRC/build/variables.bzl.tmpl` to `$SRC/build/variables.bzl` and replace variables with values valid for your configuration.

### RabbitMQ

- `RABBITMQ_USER`: The username to use to connect to the RabbitMQ server in the Docker container.
- `RABBITMQ_PASSWORD`: The password to use to connect to the RabbitMQ server in the Docker container.
- `RABBITMQ_HOST`: The hostname/IP address of the RabbitMQ server to connect to.
- `RABBITMQ_PORT`: The port the RabbitMQ server is listening on.

### PostgreSQL

- `POSTGRES_USER` : The username to use to connect to the PostgreSQL server.
- `POSTGRES_PASSWORD` : The password to use to connect to the PostgreSQL server.
- `POSTGRES_HOST` : The hostname/IP address of the PostgreSQL server to connect to.
- `POSTGRES_PORT` : The port the PostgreSQL server is listening on.

### Vault

- `VAULT_HOST` : The hostname/IP address of the Vault server to connect to.
- `VAULT_PORT` : The port the Vault server is listening on.
- `VAULT_USER` : The username to authenticate with the Vault server.
- `VAULT_PASSWORD` : The passowrd to authenticate with the Vault server.

### Okta

- `OKTA_URL`: The Org URL for Okta.
- `OKTA_KEY`: The API Key to use for the given org (API > Tokens).
- `OKTA_CLIENT_ID` : The client application Client ID (Applications > Client Credentials > Client ID).
- `OKTA_CLIENT_SECRET` : The client application Client secret (Applications > Client Credentials > Client secret).

## Setup Dependencies

- Install Blaze.

    ```
    cd $SRC/dependencies
    ./bootstrap_blaze.sh
    ```
- Add the Blaze output directory to your `$PATH` (`$SRC/dependencies/bazel/output`).
- Download Flyway.

    ```
    cd $SRC/dependencies
    ./download_flyway.sh
    ```
- Add the Flyway directory to your `$PATH` (`$SRC/dependencies/flyway`).
- Download Vault.

    ```
    cd $SRC/dependencies
    ./download_vault.sh
    ```

- Add the Vault directory to your `$PATH` (`$SRC/dependencies/vault`).
- Download libffi.

    ```
    cd $SRC/dependencies
    ./download_libffi.sh
    ```
- Download Python.
    
    ```
    cd $SRC/dependencies
    ./download_python.sh
    ```
- Download Kubectl.

    ```
    cd $SRC/dependencies
    ./download_kubectl.sh
    ```
- Add the Kubectl directory to your `$PATH` (`$SRC/dependencies/kubectl`).
- Download Minikube.

    ```
    cd $SRC/dependencies
    ./download_minikube.sh
    ```
- Add the Minikube directory to your `$PATH` (`$SRC/dependencies/minikube`).

## Setup Google Cloud

Obtain a service account JSON key with the permissions
- `storage.objects.delete`
- `storage.objects.get`
- `storage.objects.create`
and place this file in `$SRC/devops/gcloud` and name it `gcloud-webserver-account.json`.

## Setup PostgreSQL

- Setup PostgreSQL
    ```
    cd $SRC/devops/database
    ./init_dev_db.sh
    ```
- Create the PostgreSQL schema
    ```
    cd $SRC/devops/database/vault
    flyway -configFiles=./flyway/dev-flyway.conf migrate
    cd $SRC/devops/database/webserver
    flyway -configFiles=./flyway/dev-flyway.conf migrate
    ```
## Setup and Unseal Vault (Docker)

Replace `${VAULT_HOST}` and `${VAULT_PORT}` with the corresponding values in the `variables.bzl` file.

- `cd $SRC`
- `bazel run //devops/docker/vault:vault`
- `docker run --network=host bazel/devops/docker/vault:vault`
- `vault operator init -address="${VAULT_HOST}:${VAULT_PORT}" -n 1 -t 1`
- Store the unseal key and the root token somehwere.
- `vault operator unseal -address="${VAULT_HOST}:${VAULT_PORT}"`
- Paste the previously copied unseal key to unseal the vault.
- `bazel run --action_env VAULT_TOKEN="$YOUR_ROOT_TOKEN" //devops/vault:dev_init_vault`

Note that the Vault server must be started before attempting to run the webserver.
Every time the Vault server is restarted, it must be unsealed.

## Run RabbitMQ (Docker)

Replace `${RABBITMQ_PORT}` with the corresponding value in the `variables.bzl` file.

- `cd $SRC`
- `bazel run //devops/docker/rabbitmq:rabbitmq`
- `docker run --hostname rabbitmq --mount source=rabbitmqmnt,target=/var/lib/rabbitmq -p ${RABBITMQ_PORT}:${RABBITMQ_PORT} bazel/devops/docker/rabbitmq:rabbitmq`

## Build and Run Webserver

Note that this step relies on having an unsealed Vault (Docker), a running RabbitMQ server (Docker), and a running PostgreSQL database.

- `cd $SRC`
- `bazel build //src/webserver:webserver`
- `bazel run //src/webserver:webserver`

If you wish to run it in production mode:

- `bazel run -c opt //src/webserver:webserver`

If you wish to run the Docker container:

- `bazel run //devops/docker/webserver:docker_webserver`

## Build and Run Preview Generator

The preview generator worker listens on RabbitMQ for an files that clients upload and generates a PDF preview.
Thus this step also relies on having an unsealed Vault (Docker), a running RabbitMQ server (Docker), and a running PostgreSQL database.

- `cd $SRC`
- `bazel build //src/preview_generator:frontend`
- `bazel run //src/preview_generator:frontend`

If you wish to run the Docker container:

- `bazel run //devops/docker/preview_generator:docker_preview_generator`

## Run Tests

- `cd $SRC`
- `bazel test --test_output=errors //test/...`

## Running on Kubernetes (Minikube)

### Setup Minikube

- Ensure that you have a virtualization driver installed. See <https://kubernetes.io/docs/setup/learning-environment/minikube/#specifying-the-vm-driver>. kvm2 is what this guide will use.
- Run `minikube start --vm-driver=kvm2`. 

If you encounter errors when running the command:
- Ensure both libvirt and QEMU are installed.
- Ensure that `lsmod | grep kvm` returns something like `kvm                   790528  1 kvm_amd`.
- Ensure that your current user is added to the `libvirt` group: `sudo usermod --append --groups libvirt $(whoami)`
- Ensure that you have the `ebtables`, `iptables`, and `dnsmasq` packages installed.
- Ensure that `virt-host-validate` passes all tests.

After ensuring all these things, you may need to

- Reboot
- Restart the libvirt service (e.g. `sudo systemctl restart libvirtd.service`)
- Run `minikube delete`

Finally, run `eval $(minikube docker-env)` to ensure docker containers are available to minikube.
At this point you will need to rebuild all Docker containers.

### Storage

- `cd $SRC/devops/k8s/storage/dev`
- `kubectl apply -f .`

### PostgreSQL

- Modify your `postgresql.conf` file (e.g. `/var/lib/postgres/data/postgresql.conf`) to have

    ```
    listen_addresses = '*'
    ```
- Find the IPv4 address of your current ethernet link using `ip addr` (e.g. `192.168.1.160`).
- Modify your `pg_hba.conf` file (e.g. `/var/lib/postgres/data/pg_hba.conf`) to have

    ```
    host    all             all             192.168.1.160/16        trust
    ``` 

  where `192.168.1.160` is the IP address found by `ip addr`.
- Set the `POSTGRES_HOST` build variable to be `postgresql-dev-service`.
- `cd $SRC/devops/k8s/postgresql`
- `cp endpoint.yaml.tmpl endpoint.yaml`
- Modify `endpoint.yaml` to have the correct IP address (as found by `ip addr` earlier).
- `kubectl apply -f .`

### Vault

- `cd $SRC/devops/k8s/vault`
- `kubectl apply -f ./deployment.dev.yaml`
- `kubectl apply -f ./service-internal.yaml`
- `kubectl apply -f ./service-external.dev.yaml`
- Set the `VAULT_HOST` build variable to be `internal-vault-service.grchive.com` (this will require you to rebuild some Docker containers).
- Run `kubectl get services -l app=vault` and get the external port of the service. For example

    ```
    external-vault-service   NodePort    10.96.203.250   <none>        8200:30296/TCP   75s
    ```

  would indicate that the external port is `30296`.
- Check that you can reach the Vault server: `vault status -address="http://$(minikube ip):30296`
- Proceed to init/unseal the Vault server using the steps from the Docker Vault but using the new address.


### RabbitMQ

- `cd $SRC/devops/k8s/rabbitmq`
- `kubectl apply -f ./service.yaml`
- `kubectl apply -f ./statefulset.dev.yaml`
- Set the `RABBITMQ_HOST` build variable to be `rabbitmq-service` (this will require you to rebuild some Docker containers).

### Preview Generator

- `cd $SRC/devops/k8s/preview_generator`
- `kubectl apply -f ./deployment.dev.yaml`

### Webserver

- `cd $SRC`
- `bazel run //devops/docker/nginx:nginx`
- `cd $SRC/devops/k8s/webserver`
- `kubectl apply -f ./deployment.dev.yaml`
- `kubectl apply -f ./loadbalancer.yaml`
- Run `minikube service external-webserver-service --url` to obtain the URL to put into your web browser to access the webserver.
