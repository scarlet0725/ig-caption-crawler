resource "google_storage_bucket" "icc" {
  name          = "ig-caption-crawler"
  location      = "ASIA-NORTHEAST1"
  force_destroy = true
}
