terraform {
    backend "gcs" {
        credentials = "../gcloud/gcloud-terraform-account.json"
        bucket = "grchive-tf-state-prod"
        prefix = "terraform/state"
    }
}
