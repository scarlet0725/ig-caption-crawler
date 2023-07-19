terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "4.74.0"
    }
  }
  backend "gcs" {
    bucket = "prism-prod-terraform-state"
    prefix = "ig-caption-crawler"
  }
}