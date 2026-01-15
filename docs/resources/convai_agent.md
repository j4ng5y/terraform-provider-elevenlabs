# convai_agent

Manages convai agent in ElevenLabs.

## Example Usage

```hcl
resource "convai_agent" "example" {
  name = "example"
  prompt = "example"
}
```

## Argument Reference

- `name` (Required) - See provider schema for details.
- `prompt` (Required) - See provider schema for details.
- `first_message` (Optional) - See provider schema for details.
- `language` (Optional) - See provider schema for details.
- `model_id` (Optional) - See provider schema for details.

## Attribute Reference

- `id` - Computed by the API.

## Import

You can find the ID in the ElevenLabs dashboard or retrieve it via the relevant data source in this provider.

```bash
terraform import convai_agent.example <resource_id>
```