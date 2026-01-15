# convai_conversation

Manages convai conversation in ElevenLabs.

## Example Usage

```hcl
resource "convai_conversation" "example" {
  conversation_id = "example-id"
}
```

## Argument Reference

- `conversation_id` (Required) - See provider schema for details.

## Attribute Reference

- No computed attributes.

## Import

```bash
terraform import convai_conversation.example <resource_id>
```
