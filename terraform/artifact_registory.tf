resource "google_artifact_registry_repository" "icc" {
  location      = "asia-northeast1"
  repository_id = "icc"
  description   = "instagram caption crawler"
  format        = "DOCKER"
}