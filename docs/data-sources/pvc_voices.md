# elevenlabs_pvc_voices

Lists Professional Voice Cloning (PVC) voices in ElevenLabs.

## Example Usage

```hcl
data "elevenlabs_pvc_voices" "all" {}
```

## Attribute Reference

- `voices` - A list of PVC voices, each containing:
  - `id` - The unique identifier for the PVC voice.
  - `name` - The name of the PVC voice.
  - `language` - The language code for the PVC voice.
  - `description` - A description of the PVC voice.
  - `labels` - Labels associated with the PVC voice.
  - `state` - The current training state of the PVC voice.
  - `verification` - The verification status of the PVC voice.
  - `samples` - A list of training samples for the voice.
  - `settings` - Voice settings including stability, similarity_boost, style, and use_speaker_boost.
  - `created_at` - The creation timestamp.
  - `updated_at` - The last update timestamp.