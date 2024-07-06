terraform {
  required_providers {
    bhyve = {
      source = "github.com/gillcaleb/bhyve"
    }
  }
}

provider "bhyve" {}

data "bhyve_zone" "zone" {}
