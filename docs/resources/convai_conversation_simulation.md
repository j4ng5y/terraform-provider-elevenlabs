# convai_conversation_simulation

Manages convai conversation simulation in ElevenLabs.

## Example Usage

```hcl
resource "convai_conversation_simulation" "example" {
  chat_history = "example"
  role = "example"
  content = "example"
}
```

## Argument Reference

- `chat_history` (Required) - See provider schema for details.
- `role` (Required) - See provider schema for details.
- `content` (Required) - See provider schema for details.

## Attribute Reference

- `simulated_conversation` - Computed by the API.
- `role` - Computed by the API.
- `content` - Computed by the API.
- `timestamp` - Computed by the API.
- `analysis` - Computed by the API.
