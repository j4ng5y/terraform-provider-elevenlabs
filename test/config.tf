terraform {
  required_version = ">=1"
  required_providers {
    elevenlabs = {
      source = "j4ng5y/elevenlabs"
    }
  }
}

variable "api_key" {
  type      = string
  sensitive = true
}

variable "base_url" {
  type    = string
  default = ""
}

provider "elevenlabs" {
  api_key  = var.api_key
  base_url = var.base_url
}
