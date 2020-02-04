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
}

module "webserver" {
    source = "./modules/webserver"
}
