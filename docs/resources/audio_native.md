# audio_native

Manages audio native in ElevenLabs.

## Example Usage

```hcl
resource "audio_native" "example" {
  name = "example"
  file_path = "example"
}
```

## Argument Reference

- `name` (Required) - See provider schema for details.
- `file_path` (Required) - See provider schema for details.
- `voice_id` (Optional) - See provider schema for details.
- `model_id` (Optional) - See provider schema for details.
- `title` (Optional) - See provider schema for details.
- `author` (Optional) - See provider schema for details.
- `text_color` (Optional) - See provider schema for details.
- `background_color` (Optional) - See provider schema for details.
- `auto_convert` (Optional) - See provider schema for details.

## Attribute Reference

- `id` - Computed by the API.
- `html_snippet` - Computed by the API.
- `status` - Computed by the API.

## Import

```bash
terraform import audio_native.example <resource_id>
```
