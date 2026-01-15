# voice

Manages voice in ElevenLabs.

## Example Usage

```hcl
resource "voice" "example" {
}
```

## Argument Reference

- `settings` (Optional) - See provider schema for details.
- `stability` (Optional) - See provider schema for details.
- `similarity_boost` (Optional) - See provider schema for details.
- `style` (Optional) - See provider schema for details.
- `use_speaker_boost` (Optional) - See provider schema for details.

## Attribute Reference

- `stability` - Computed by the API.
- `similarity_boost` - Computed by the API.
- `style` - Computed by the API.
- `use_speaker_boost` - Computed by the API.

## Import

```bash
terraform import voice.example <resource_id>
```
