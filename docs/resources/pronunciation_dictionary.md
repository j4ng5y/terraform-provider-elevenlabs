# pronunciation_dictionary

Manages pronunciation dictionary in ElevenLabs.

## Example Usage

```hcl
resource "pronunciation_dictionary" "example" {
  name = "example"
  string_to_replace = "example"
}
```

## Argument Reference

- `name` (Required) - See provider schema for details.
- `string_to_replace` (Required) - See provider schema for details.
- `description` (Optional) - See provider schema for details.
- `rules` (Optional) - See provider schema for details.
- `alias` (Optional) - See provider schema for details.
- `phoneme` (Optional) - See provider schema for details.
- `alphabet` (Optional) - See provider schema for details.

## Attribute Reference

- `id` - Computed by the API.
- `latest_version_id` - Computed by the API.

## Import

You can find the ID in the ElevenLabs dashboard or retrieve it via the relevant data source in this provider.

```bash
terraform import pronunciation_dictionary.example <resource_id>
```