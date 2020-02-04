provider "google" {
    credentials = file("../gcloud/gcloud-terraform-account.json")
    project     = "grchive"
    region      = "us-central1"
    zone        = "us-central1-c"
}
