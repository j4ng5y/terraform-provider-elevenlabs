# shared_voice

Manages shared voice in ElevenLabs.

## Example Usage

```hcl
resource "shared_voice" "example" {
  public_user_id = "example-id"
  voice_id = "example-id"
  name = "example"
}
```

## Argument Reference

- `public_user_id` (Required) - See provider schema for details.
- `voice_id` (Required) - See provider schema for details.
- `name` (Required) - See provider schema for details.

## Attribute Reference

- `id` - Computed by the API.

## Import

You can find the ID in the ElevenLabs dashboard or retrieve it via the relevant data source in this provider.

```bash
terraform import shared_voice.example <resource_id>
```