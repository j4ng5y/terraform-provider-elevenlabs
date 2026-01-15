# convai_secret

Manages convai secret in ElevenLabs.

## Example Usage

```hcl
resource "convai_secret" "example" {
  name = "example"
}
```

## Argument Reference

- `name` (Required) - See provider schema for details.

## Attribute Reference

- `id` - Computed by the API.

## Import

```bash
terraform import convai_secret.example <resource_id>
```
