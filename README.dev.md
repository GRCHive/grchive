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
- RabbitMQ
- Docker

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

## Setup

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
- Download Python.
    
    ```
    cd $SRC/dependencies
    ./download_python.sh
    ```

## Setup Dev Environment

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
- `docker run -p ${RABBITMQ_PORT}:${RABBITMQ_PORT} bazel/devops/docker/rabbitmq:rabbitmq`

## Build and Run

Note that this steps relies on having an unsealed Vault, a running RabbitMQ server (Docker), and a running PostgreSQL database.

- `cd $SRC`
- `bazel build //src/webserver:webserver`
- `bazel run //src/webserver:webserver`

If you wish to run it in production mode:

- `bazel run -c opt //src/webserver:webserver`

## Run Tests

- `cd $SRC`
- `bazel test --test_output=errors //test/...`
