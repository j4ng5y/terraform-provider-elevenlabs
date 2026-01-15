# convai_agent_test

Manages convai agent test in ElevenLabs.

## Example Usage

```hcl
resource "convai_agent_test" "example" {
  name = "example"
  success_condition = "example"
}
```

## Argument Reference

- `name` (Required) - See provider schema for details.
- `success_condition` (Required) - See provider schema for details.

## Attribute Reference

- `id` - Computed by the API.

## Import

You can find the ID in the ElevenLabs dashboard or retrieve it via the relevant data source in this provider.

```bash
terraform import convai_agent_test.example <resource_id>
```