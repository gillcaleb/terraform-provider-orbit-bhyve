terraform {
  required_providers {
    bhyve = {
      source = "github.com/gillcaleb/bhyve"
    }
  }
}

provider "bhyve" {}

resource "bhyve_zone" "zone" {
  id      = 5
  minutes = 5
}