# voice_sample

Manages voice sample in ElevenLabs.

## Example Usage

```hcl
resource "voice_sample" "example" {
  voice_id = "example-id"
  file_path = "example"
}
```

## Argument Reference

- `voice_id` (Required) - See provider schema for details.
- `file_path` (Required) - See provider schema for details.

## Attribute Reference

- `id` - Computed by the API.
- `file_name` - Computed by the API.

## Import

```bash
terraform import voice_sample.example <resource_id>
```
