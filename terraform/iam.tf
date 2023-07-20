resource "google_service_account" "icc_crawler_account" {
  account_id   = "icc-crawler"
  display_name = "icc-crawler"
}

resource "google_project_iam_binding" "icc_crawler_account_publisher_binding" {
  project = local.project_id
  role    = "roles/pubsub.publisher"
  members = [
    "serviceAccount:${google_service_account.icc_crawler_account.email}"
  ]
}

resource "google_project_iam_binding" "icc_crawler_account_subscriber_binding" {
  project = local.project_id
  role    = "roles/pubsub.subscriber"
  members = [
    "serviceAccount:${google_service_account.icc_crawler_account.email}"
  ]
}

resource "google_project_iam_binding" "storage_reader_binding" {
  project = local.project_id
  role    = "roles/storage.objectViewer"
  members = [
    "serviceAccount:${google_service_account.icc_crawler_account.email}"
  ]
}
