terraform {
    backend "gcs" {
        credentials = "../../gcloud/gcloud-terraform-account.json"
        bucket = "grchive-tf-state-prod"
        prefix = "terraform/state"
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

module "database" {
    source = "../modules/database"

    postgres_user = var.postgres_user
    postgres_password = var.postgres_password
    postgres_instance_name = var.postgres_instance_name
}

module "vault" {
    source = "../modules/vault"
}

module "webserver" {
    source = "../modules/webserver"

    doc_storage_bucket = "grchive-prod"
    kotlin_lib_bucket = "grchive-kotlin-lib-prod"
}
