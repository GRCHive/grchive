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

## Build and Run

- `cd $SRC`
- `bazel build //src/webserver:webserver`
- `bazel run //src/webserver:webserver`

## Run Tests

- `cd $SRC`
- `bazel test --test_output=errors //...`
