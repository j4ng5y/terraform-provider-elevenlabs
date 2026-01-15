# elevenlabs_pvc_voice

Manages a Professional Voice Cloning (PVC) voice in ElevenLabs.

## Example Usage

```hcl
resource "elevenlabs_pvc_voice" "enterprise_voice" {
  name        = "Enterprise Narrator"
  language    = "en"
  description = "High-quality voice for corporate training videos."
  labels = {
    "department" = "training"
    "quality"    = "enterprise"
  }
}
```

## Argument Reference

- `name` (Required) - The name of the PVC voice.
- `language` (Required) - The language code for the PVC voice (e.g., 'en', 'es', 'fr').
- `description` (Optional) - A description of the PVC voice.
- `labels` (Optional) - Labels associated with the PVC voice.

## Attribute Reference

- `id` - The unique identifier for the PVC voice.
- `state` - The current training state of the PVC voice.
- `verification` - The verification status of the PVC voice.
- `settings` - Voice settings including stability, similarity_boost, style, and use_speaker_boost.

## Import

You can find the ID in the ElevenLabs dashboard or retrieve it via the relevant data source in this provider.

PVC voices can be imported using the voice ID:

```bash
terraform import elevenlabs_pvc_voice.enterprise_voice voice_123456789
```