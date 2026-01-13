terraform {
  required_providers {
    elevenlabs = {
      source = "j4ng5y/elevenlabs"
    }
  }
}

provider "elevenlabs" {
  # api_key can be set here or via ELEVENLABS_API_KEY env var
}

data "elevenlabs_models" "all" {}

output "available_models" {
  value = data.elevenlabs_models.all.models
}

resource "elevenlabs_voice" "my_cloned_voice" {
  name        = "My Custom Voice"
  description = "A voice cloned from sample files"
  files       = ["samples/sample1.mp3", "samples/sample2.mp3"]
  labels = {
    "accent" = "american"
    "gender" = "male"
  }
}

resource "elevenlabs_project" "my_audiobook" {
  name             = "My Audiobook Project"
  default_model_id = "eleven_multilingual_v2"
}
