# service_account_key

Manages service account key in ElevenLabs.

## Example Usage

```hcl
resource "service_account_key" "example" {
  user_id = "example-id"
  name = "example"
}
```

## Argument Reference

- `user_id` (Required) - See provider schema for details.
- `name` (Required) - See provider schema for details.
- `character_limit` (Optional) - See provider schema for details.

## Attribute Reference

- `id` - Computed by the API.

## Import

```bash
terraform import service_account_key.example <resource_id>
```
