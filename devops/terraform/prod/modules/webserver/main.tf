resource "google_container_cluster" "webserver-gke" {
    name     = "webserver-gke-cluster"
    location = "us-central1-c"

    remove_default_node_pool = true
    initial_node_count       = 1

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
        machine_type = "f1-micro"
    }
}
