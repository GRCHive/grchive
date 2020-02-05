provider "google" {
    credentials = file(var.terraform-credentials-file)
    project     = "grchive"
    region      = "us-central1"
    zone        = "us-central1-c"
}

resource "google_storage_bucket" "gc-terraform-store" {
    name = var.terraform-bucket-name
    location = "US-CENTRAL1"
    versioning = {
        enabled = true
    }
}

terraform {
    backend "gcs" {
        credentials = var.terraform-credentials-file
        bucket = var.terraform-bucket-name
        prefix = "terraform/state"
    }
}

module "webserver" {
    source = "./modules/webserver"
}
