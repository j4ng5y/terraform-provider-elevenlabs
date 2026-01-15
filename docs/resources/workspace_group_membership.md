# workspace_group_membership

Manages workspace group membership in ElevenLabs.

## Example Usage

```hcl
resource "workspace_group_membership" "example" {
  group_id = "example-id"
  email = "user@example.com"
}
```

## Argument Reference

- `group_id` (Required) - See provider schema for details.
- `email` (Required) - See provider schema for details.

## Attribute Reference

- `id` - Computed by the API.

## Import

```bash
terraform import workspace_group_membership.example <resource_id>
```
