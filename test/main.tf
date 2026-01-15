data "elevenlabs_models" "all" {}

data "elevenlabs_voices" "all" {}

output "models" {
  value = data.elevenlabs_models.all
}

output "voices" {
  value = data.elevenlabs_voices.all
}
