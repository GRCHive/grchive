resource "google_kms_key_ring" "vault-keyring" {
    project         = "grchive"
    name            = "vault-keyring"
    location        = "us-central1"
}

resource "google_kms_crypto_key" "vault-key" {
    name            = "vault-key"
    key_ring        = google_kms_key_ring.vault-keyring.self_link
}
