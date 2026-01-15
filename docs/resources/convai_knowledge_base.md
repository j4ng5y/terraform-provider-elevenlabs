# convai_knowledge_base

Manages convai knowledge base in ElevenLabs.

## Example Usage

```hcl
resource "convai_knowledge_base" "example" {
  name = "example"
}
```

## Argument Reference

- `name` (Required) - See provider schema for details.
- `url` (Optional) - See provider schema for details.
- `content` (Optional) - See provider schema for details.
- `file_path` (Optional) - See provider schema for details.

## Attribute Reference

- `id` - Computed by the API.
- `type` - Computed by the API.
- `status` - Computed by the API.

## Import

You can find the ID in the ElevenLabs dashboard or retrieve it via the relevant data source in this provider.

```bash
terraform import convai_knowledge_base.example <resource_id>
```