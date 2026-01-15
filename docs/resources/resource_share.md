# resource_share

Manages resource share in ElevenLabs.

## Example Usage

```hcl
resource "resource_share" "example" {
  resource_id = "example-id"
  resource_type = "example"
  email = "user@example.com"
  role = "example"
}
```

## Argument Reference

- `resource_id` (Required) - See provider schema for details.
- `resource_type` (Required) - See provider schema for details.
- `email` (Required) - See provider schema for details.
- `role` (Required) - See provider schema for details.

## Attribute Reference

- `id` - Computed by the API.

## Import

```bash
terraform import resource_share.example <resource_id>
```
