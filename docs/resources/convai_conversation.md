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

You can find the ID in the ElevenLabs dashboard or retrieve it via the relevant data source in this provider.

```bash
terraform import convai_conversation.example <resource_id>
```