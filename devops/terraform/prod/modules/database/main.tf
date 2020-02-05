resource "google_sql_database_instance" "main-db" {
    name = var.postgres_instance_name
    database_version = "POSTGRES_11"
    region = "us-central1"

    settings {
        tier = "db-f1-micro"
        availability_type = "ZONAL"

        backup_configuration {
            enabled = true
        }

        ip_configuration {
            ipv4_enabled = true
            require_ssl = true
        }

        location_preference {
            zone = "us-central1-c"
        }

        maintenance_window {
            day = 7
            hour = 3
            update_track = "stable"
        }
    }
}

resource "google_sql_database" "vault-database" {
    name     = "vault"
    instance = google_sql_database_instance.main-db.name
}

resource "google_sql_database" "grchive-database" {
    name     = "audit"
    instance = google_sql_database_instance.main-db.name
}

resource "google_sql_user" "web-user" {
    name     = var.postgres_user
    password = var.postgres_password
    instance = google_sql_database_instance.main-db.name
}

resource "google_sql_ssl_cert" "main-db-cert" {
    common_name = "main-db-cert"
    instance    = google_sql_database_instance.main-db.name
}