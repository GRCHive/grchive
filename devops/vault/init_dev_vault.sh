#!/bin/sh
vault secrets enable -address="${VAULT_HOST}:${VAULT_PORT}" transit
vault write -address="${VAULT_HOST}:${VAULT_PORT}" -f transit/keys/passwords
