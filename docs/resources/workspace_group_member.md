# workspace_group_member

Manages workspace group member in ElevenLabs.

## Example Usage

```hcl
resource "workspace_group_member" "example" {
  group_id = "example-id"
  email = "user@example.com"
}
```

## Argument Reference

- `group_id` (Required) - See provider schema for details.
- `email` (Required) - See provider schema for details.

## Attribute Reference

- No computed attributes.
