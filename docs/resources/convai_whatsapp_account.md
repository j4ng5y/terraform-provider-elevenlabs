# convai_whatsapp_account

Manages convai whatsapp account in ElevenLabs.

## Example Usage

```hcl
resource "convai_whatsapp_account" "example" {
  phone_number_id = "example-id"
  business_account_id = "example-id"
}
```

## Argument Reference

- `phone_number_id` (Required) - See provider schema for details.
- `business_account_id` (Required) - See provider schema for details.

## Attribute Reference

- No computed attributes.

## Import

You can find the ID in the ElevenLabs dashboard or retrieve it via the relevant data source in this provider.

```bash
terraform import convai_whatsapp_account.example <resource_id>
```