# Audit Stuff

This document will assume that the git checkout directory is an set in an environment variable `$SRC`.

## Prerequisites

- G++ v?
- JDK8
- zip
- unzip
- go 1.13+ 

## Setup

- Install Blaze.

    ```
    cd $SRC/external
    ./bootstrap_blaze.sh
    ```
- Add the Blaze binary to your path (`external/bazel/output`).

## Build and Run

- `cd $SRC`
- `bazel build //src/webserver:webserver`
- `bazel run //src/webserver:webserver`
