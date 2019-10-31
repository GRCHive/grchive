storage "postgresql" {
    connection_url = "postgres://DEVUSER@localhost:5432/audit?sslmode=disable&timezone=UTC"
    table = "vault_kv_store"
}

listener "tcp" {
    address     = "localhost:8081"
    tls_disable = 1
}

disable_mlock = true
