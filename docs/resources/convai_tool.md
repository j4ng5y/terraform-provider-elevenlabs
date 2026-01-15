# convai_tool

Manages convai tool in ElevenLabs.

## Example Usage

```hcl
resource "convai_tool" "example" {
  name = "example"
  description = "example"
}
```

## Argument Reference

- `name` (Required) - See provider schema for details.
- `description` (Required) - See provider schema for details.

## Attribute Reference

- `id` - Computed by the API.

## Import

You can find the ID in the ElevenLabs dashboard or retrieve it via the relevant data source in this provider.

```bash
terraform import convai_tool.example <resource_id>
```