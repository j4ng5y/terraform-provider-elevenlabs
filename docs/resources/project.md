# project

Manages project in ElevenLabs.

## Example Usage

```hcl
resource "project" "example" {
  name = "example"
}
```

## Argument Reference

- `name` (Required) - See provider schema for details.
- `default_model_id` (Optional) - See provider schema for details.
- `default_paragraph_voice_id` (Optional) - See provider schema for details.
- `default_title_voice_id` (Optional) - See provider schema for details.

## Attribute Reference

- `id` - Computed by the API.
- `default_model_id` - Computed by the API.
- `default_paragraph_voice_id` - Computed by the API.
- `default_title_voice_id` - Computed by the API.
- `state` - Computed by the API.

## Import

You can find the ID in the ElevenLabs dashboard or retrieve it via the relevant data source in this provider.

```bash
terraform import project.example <resource_id>
```