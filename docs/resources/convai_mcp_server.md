# convai_mcp_server

Manages convai mcp server in ElevenLabs.

## Example Usage

```hcl
resource "convai_mcp_server" "example" {
  name = "example"
  url = "https://example.com"
}
```

## Argument Reference

- `name` (Required) - See provider schema for details.
- `url` (Required) - See provider schema for details.

## Attribute Reference

- `id` - Computed by the API.

## Import

You can find the ID in the ElevenLabs dashboard or retrieve it via the relevant data source in this provider.

```bash
terraform import convai_mcp_server.example <resource_id>
```