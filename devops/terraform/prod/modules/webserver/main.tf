data "google_container_engine_versions" "gke-version-central1c" {
    location       = "us-central1-c"
    version_prefix = "1.15."
}

resource "google_container_cluster" "webserver-gke" {
    name     = "webserver-gke-cluster"
    location = "us-central1-c"

    remove_default_node_pool = true
    initial_node_count       = 1
    min_master_version       = data.google_container_engine_versions.gke-version-central1c.latest_master_version

    master_auth {
        username = ""
        password = ""

        client_certificate_config {
            issue_client_certificate = false
        }
    }
}

resource "google_container_node_pool" "webserver-node-pool" {
    name        = "webserver-node-pool"
    location    = "us-central1-c"
    cluster     = google_container_cluster.webserver-gke.name
    node_count  = 1

    node_config {
        disk_size_gb        = 10
        disk_type           = "pd-standard"
        image_type          = "COS"
        labels = {
            app = "webserver"
        }
        machine_type = "n1-standard-2"
    }
}

resource "google_compute_global_address" "webserver-static-ip" {
    name = "webserver-static-ip"
}

resource "google_storage_bucket" "webserver-control-doc-store" {
    name     = "grchive-prod"
    location = "US-CENTRAL1"
    bucket_policy_only = true

    versioning {
        enabled = true
    }
}
