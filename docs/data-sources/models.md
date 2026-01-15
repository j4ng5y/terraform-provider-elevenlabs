# models

Fetches models from ElevenLabs.

## Example Usage

```hcl
data "models" "example" {
}
```

## Argument Reference

- No configurable arguments.

## Attribute Reference

- `models` - Computed by the API.
- `model_id` - Computed by the API.
- `name` - Computed by the API.
- `description` - Computed by the API.
- `can_do_text_to_speech` - Computed by the API.
- `can_do_voice_conversion` - Computed by the API.
