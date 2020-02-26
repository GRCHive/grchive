data "google_container_engine_versions" "gke-version-central1c" {
    location       = "us-central1-c"
    version_prefix = "1.15.9-"
}

resource "google_compute_network" "gke-outbound-network" {
    name                    = "gke-outbound-network"
    auto_create_subnetworks = false
    routing_mode            = "REGIONAL"
}

resource "google_compute_subnetwork" "gke-outbound-network-us-central1" {
    name                    = "gke-outbound-network-us-central1"
    region                  = "us-central1"
    network                 = google_compute_network.gke-outbound-network.self_link
    ip_cidr_range           = "192.168.1.0/24"
}

resource "google_container_cluster" "webserver-gke" {
    name     = "webserver-gke-cluster"
    location = "us-central1-c"

    remove_default_node_pool = true
    initial_node_count       = 1
    min_master_version       = data.google_container_engine_versions.gke-version-central1c.latest_master_version

    network    = google_compute_network.gke-outbound-network.self_link
    subnetwork = google_compute_subnetwork.gke-outbound-network-us-central1.self_link

    master_auth {
        username = ""
        password = ""

        client_certificate_config {
            issue_client_certificate = false
        }
    }

    private cluster_config {
        enable_private_nodes = true
        enable_private_endpoint = true
        master_ipv4_cidr_block = "172.16.0.0/28"
    }

    ip_allocation_policy {
        cluster_ipv4_cidr_block = ""
        services_ipv4_cidr_block = ""
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
    name     = var.doc_storage_bucket
    location = "US-CENTRAL1"
    bucket_policy_only = true

    versioning {
        enabled = true
    }
}

resource "google_compute_address" "gke-outbound-network-us-central1-nat-ip" {
    count   = 1
    name    = "gke-outbound-network-us-central1-nat-ip-${count.index}"
    region  = google_compute_subnetwork.gke-outbound-network-us-central1.region
}

resource "google_compute_router" "gke-outbound-network-us-central1-router" {
    name        = "gke-outbound-network-us-central1-router"
    region      = google_compute_subnetwork.gke-outbound-network-us-central1.region
    network     = google_compute_network.gke-outbound-network.self_link
}

resource "google_compute_router_nat" "gke-outbound-network-us-central1-nat" {
    name    = "gke-outbound-network-us-central1-nat"
    router  = google_compute_router.gke-outbound-network-us-central1-router.name
    region  = google_compute_router.gke-outbound-network-us-central1-router.region

    nat_ip_allocate_option = "MANUAL_ONLY"
    nat_ips                = google_compute_address.gke-outbound-network-us-central1-nat-ip.*.self_link

    source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
    subnetwork {
        name                    = google_compute_subnetwork.gke-outbound-network-us-central1.self_link
        source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
    }

    log_config {
        enable = true
        filter = "ALL"
    }
}
