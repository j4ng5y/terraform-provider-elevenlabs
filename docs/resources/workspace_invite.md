# workspace_invite

Manages workspace invite in ElevenLabs.

## Example Usage

```hcl
resource "workspace_invite" "example" {
  email = "user@example.com"
}
```

## Argument Reference

- `email` (Required) - See provider schema for details.

## Attribute Reference

- No computed attributes.

## Import

You can find the ID in the ElevenLabs dashboard or retrieve it via the relevant data source in this provider.

```bash
terraform import workspace_invite.example <resource_id>
```