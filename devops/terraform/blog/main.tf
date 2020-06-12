terraform {
    backend "gcs" {
        credentials = "../../gcloud/gcloud-terraform-account.json"
        bucket = "grchive-tf-state-prod"
        prefix = "terraform/wp-state"
    }
}

provider "google" {
    credentials = file("../../gcloud/gcloud-terraform-account.json")
    project     = "grchive"
    region      = "us-central1"
    zone        = "us-central1-c"
    version     =  "~> 3.7"
    scopes      = [
        "https://www.googleapis.com/auth/compute",
        "https://www.googleapis.com/auth/cloud-platform",
        "https://www.googleapis.com/auth/ndev.clouddns.readwrite",
        "https://www.googleapis.com/auth/devstorage.full_control",
        "https://www.googleapis.com/auth/userinfo.email",
        "https://www.googleapis.com/auth/cloud-platform",
        "https://www.googleapis.com/auth/sqlservice.admin",
    ]
}

resource "google_compute_address" "wordpress-static-ip" {
    name    = "wordpress-static-ip"
    region  = "us-central1"
}

data "google_compute_image" "wordpress-image" {
    name = "cos-stable-81-12871-119-0"
    project = "cos-cloud"
}

resource "google_compute_network" "wordpress-network" {
    name                    = "wordpress-network"
    auto_create_subnetworks = false
    routing_mode            = "REGIONAL"
}

resource "google_compute_subnetwork" "wordpress-network-us-central1" {
    name                    = "wordpress-network-us-central1"
    region                  = "us-central1"
    network                 = google_compute_network.wordpress-network.self_link
    ip_cidr_range           = "192.168.1.0/24"
}

resource "google_compute_firewall" "wordpress-network-firewall-http-ingress" {
    name                    = "wordpress-network-firewall-http-ingress"
    network                 = google_compute_network.wordpress-network.name
    direction               = "INGRESS"

    allow {
        protocol = "tcp"
        ports = ["80", "443", "8080", "22"]
    }
}

resource "google_compute_instance" "wordpress" {
    name            = "grchive-wordpress-central1-c"
    machine_type    = "f1-micro"
    zone            = "us-central1-c"

    boot_disk {
        initialize_params {
            size    = 20
            type    = "pd-standard"
            image   = data.google_compute_image.wordpress-image.self_link
        }
    }

    network_interface {
        network    = google_compute_network.wordpress-network.self_link
        subnetwork = google_compute_subnetwork.wordpress-network-us-central1.self_link

        access_config {
            nat_ip = google_compute_address.wordpress-static-ip.address
        }
    }
}

resource "google_sql_database_instance" "wordpress-db" {
    name = var.wp_instance_name
    database_version = "MYSQL_5_7"
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

            authorized_networks {
                name = "blog-${google_compute_instance.wordpress.name}"
                value = google_compute_address.wordpress-static-ip.address
            }
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

resource "google_sql_database" "wordpress-database" {
    name     = var.wp_database_name
    instance = google_sql_database_instance.wordpress-db.name
}

resource "google_sql_user" "wordpress-user" {
    name     = var.wp_database_user
    password = var.wp_database_password
    instance = google_sql_database_instance.wordpress-db.name
}

resource "google_sql_ssl_cert" "wordpress-db-cert" {
    common_name = "wordpress-db-cert"
    instance    = google_sql_database_instance.wordpress-db.name
}
