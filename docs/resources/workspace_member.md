# workspace_member

Manages workspace member in ElevenLabs.

## Example Usage

```hcl
resource "workspace_member" "example" {
  id = "example-id"
}
```

## Argument Reference

- `id` (Required) - See provider schema for details.

## Attribute Reference

- `email` - Computed by the API.

## Import

```bash
terraform import workspace_member.example <resource_id>
```
