# convai_phone_number

Manages convai phone number in ElevenLabs.

## Example Usage

```hcl
resource "convai_phone_number" "example" {
  phone_number = 1
  telephony_provider = "example-id"
}
```

## Argument Reference

- `phone_number` (Required) - See provider schema for details.
- `telephony_provider` (Required) - See provider schema for details.
- `label` (Optional) - See provider schema for details.

## Attribute Reference

- `id` - Computed by the API.

## Import

```bash
terraform import convai_phone_number.example <resource_id>
```
