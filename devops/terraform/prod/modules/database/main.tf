resource "google_sql_database_instance" "main-db" {
    name = "main-db-instance"
    database_version = "POSTGRES_11"
    region = "us-central1"

    settings = {
        tier = "db-f1-micro"
        availability_type = "REGIONAL"

        location_preference = {
            zone = "us-central1-c"
        }

        maintenance_window = {
            day = 7
            hour = 3
            update_track = "stable"
        }
    }
}
