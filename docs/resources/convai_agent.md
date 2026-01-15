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

```bash
terraform import convai_agent.example <resource_id>
```
