# convai_llm_usage_calculator

Fetches convai llm usage calculator from ElevenLabs.

## Example Usage

```hcl
data "convai_llm_usage_calculator" "example" {
}
```

## Argument Reference

- No configurable arguments.

## Attribute Reference

- `llm_prices` - Computed by the API.
- `model_name` - Computed by the API.
- `input_cost_per_million_tokens` - Computed by the API.
- `output_cost_per_million_tokens` - Computed by the API.
- `estimated_input_tokens` - Computed by the API.
- `estimated_output_tokens` - Computed by the API.
- `estimated_cost` - Computed by the API.
