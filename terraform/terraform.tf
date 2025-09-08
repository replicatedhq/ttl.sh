terraform {
  backend "s3" {
    bucket                    = "ttl-sh-tf-state"
    key                       = "terraform.tfstate"
    endpoint                  = "hel1.your-objectstorage.com"
    region                    = "eu-central1"
    skip_credentials_validation = true
    skip_region_validation    = true
  }
}