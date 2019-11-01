# Audit Stuff

This document will assume that the git checkout directory is set in an environment variable `$SRC`.

## Prerequisites

- G++ v?
- JDK8
- zip
- unzip
- go 1.13+ 
- PostgreSQL
- Node v12.10

## Setup

- Install Blaze.

    ```
    cd $SRC/external
    ./bootstrap_blaze.sh
    ```
- Add the Blaze output directory to your `$PATH` (`$SRC/external/bazel/output`).
- Download Flyway.

    ```
    cd $SRC/external
    ./download_flyway.sh
    ```
- Add the Flyway directory to your `$PATH` (`$SRC/external/flyway`).
- Download Vault.

    ```
    cd $SRC/external
    ./download_vault.sh
    ```
- Add the Vault directory to your `$PATH` (`$SRC/external/vault`).


## Setup Dev Environment

- Setup PostgreSQL
    ```
    cd $SRC/devops/database
    ./init_dev_db.sh
    ```
- Create the PostgreSQL schema
    ```
    flyway -configFiles=./flyway/dev-flyway.conf migrate
    ```
## Setup and Unseal Vault

- `cd $SRC`
- `vault server -config=devops/vault/config/dev.hcl &`
- `vault operator init -address="http://localhost:8081" -n 1 -t 1`
- Store the unseal key and the root token somehwere.
- `vault operator unseal -address="http://localhost:8081"`
- Past the previously copied unseal key to unseal the vault.
- `vault login -address="http://localhost:8081" $ROOT_TOKEN_HERE`

Note that the Vault server must be started before attempting to run the webserver.
Every time the Vault server is restarted, it must be unsealed.

## Build and Run

- `cd $SRC`
- `bazel build //src/webserver:webserver`
- `bazel run //src/webserver:webserver`

## Run Tests

- `cd $SRC`
- `bazel test --test_output=errors //...`