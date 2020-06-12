output "wp_address" {
    value = "${google_compute_address.wordpress-static-ip.address}"
}

output "wp_client_key" {
    value = "${google_sql_ssl_cert.wordpress-db-cert.private_key}"
}

output "wp_server_cert" {
    value = "${google_sql_ssl_cert.wordpress-db-cert.server_ca_cert}"
}

output "wp_client_cert" {
    value = "${google_sql_ssl_cert.wordpress-db-cert.cert}"
}
