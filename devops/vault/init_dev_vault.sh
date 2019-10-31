#!/bin/bash

echo "========= YOU MUST ALREADY BE LOGGED IN WITH A SUFFICIENTLY PRIVILEGED TOKEN FOR THIS TO WORK ========="
vault secrets enable -path="encryption-keys" $@ kv
vault kv enable-versioning $@ encryption-keys
