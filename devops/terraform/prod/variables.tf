variable "terraform-bucket-name" {
    type = string
    default = "grchive-tf-state-prod"
}

variable "terraform-credentials-file" {
    type = string
    default = "../../gcloud/gcloud-terraform-account.json"
}
