# elevenlabs_convai_mcp_tool_config

Manages a ConvAI MCP tool configuration override for a specific MCP server/tool pair.

## Example Usage

```hcl
resource "elevenlabs_convai_mcp_tool_config" "lookup" {
  mcp_server_id = elevenlabs_convai_mcp_server.main.id
  tool_name     = "lookup"

  force_pre_tool_speech    = true
  disable_interruptions    = false
  tool_call_sound          = "beep"
  tool_call_sound_behavior = "auto"
  execution_mode           = "immediate"

  assignments = [
    {
      source           = "response"
      dynamic_variable = "user_id"
      value_path       = "data.id"
    }
  ]
}
```

## Argument Reference

- `mcp_server_id` (Required) - The MCP server ID that owns the tool.
- `tool_name` (Required) - The name of the MCP tool to configure.
- `force_pre_tool_speech` (Optional) - Override to require pre-tool speech for this tool.
- `disable_interruptions` (Optional) - Override to disable interruptions while this tool runs.
- `tool_call_sound` (Optional) - Override the sound played during tool execution.
- `tool_call_sound_behavior` (Optional) - Override when tool call sound plays (`auto`, `always`, etc.).
- `execution_mode` (Optional) - Override how the tool executes (`immediate`, `post_tool_speech`, `async`).
- `assignments` (Optional) - Dynamic variable assignments derived from tool responses.
  - `source` (Optional) - Value source (currently `response`).
  - `dynamic_variable` (Required) - Dynamic variable name to set.
  - `value_path` (Required) - Dot-notation path in the response payload.

## Attribute Reference

- `id` - Composite ID in the format `mcp_server_id:tool_name`.

## Import

Tool configs can be imported using the composite ID:

```bash
terraform import elevenlabs_convai_mcp_tool_config.lookup mcp-123:lookup
```
