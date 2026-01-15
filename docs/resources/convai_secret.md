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

You can find the ID in the ElevenLabs dashboard or retrieve it via the relevant data source in this provider.

```bash
terraform import convai_secret.example <resource_id>
```