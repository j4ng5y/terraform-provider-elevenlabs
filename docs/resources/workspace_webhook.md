# workspace_webhook

Manages workspace webhook in ElevenLabs.

## Example Usage

```hcl
resource "workspace_webhook" "example" {
  url = "https://example.com"
}
```

## Argument Reference

- `url` (Required) - See provider schema for details.

## Attribute Reference

- `id` - Computed by the API.

## Import

```bash
terraform import workspace_webhook.example <resource_id>
```
