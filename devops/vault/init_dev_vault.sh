#!/bin/bash

echo "========= YOU MUST ALREADY BE LOGGED IN WITH A SUFFICIENTLY PRIVILEGED TOKEN FOR THIS TO WORK ========="
vault secrets enable $@ transit
vault write $@ -f transit/keys/passwords
