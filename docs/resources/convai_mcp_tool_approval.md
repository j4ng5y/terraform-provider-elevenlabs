# elevenlabs_convai_mcp_tool_approval

Manages ConvAI MCP tool approvals for a specific MCP server/tool pair.

## Example Usage

```hcl
resource "elevenlabs_convai_mcp_tool_approval" "lookup" {
  mcp_server_id    = elevenlabs_convai_mcp_server.main.id
  tool_name        = "lookup"
  tool_description = "Lookup tool"
  input_schema     = jsonencode({ type = "object" })
  approval_policy  = "auto_approved"
}
```

## Argument Reference

- `mcp_server_id` (Required) - The MCP server ID that owns the tool.
- `tool_name` (Required) - The name of the MCP tool to approve.
- `tool_description` (Required) - Description of the MCP tool.
- `input_schema` (Optional) - JSON-encoded input schema for the tool.
- `approval_policy` (Optional) - Tool-level approval policy (`auto_approved` or `requires_approval`).

## Attribute Reference

- `id` - Composite ID in the format `mcp_server_id:tool_name`.

## Import

Use the MCP server ID from `elevenlabs_convai_mcp_servers` and the tool name from your MCP server configuration (tools list).

Tool approvals can be imported using the composite ID:

```bash
terraform import elevenlabs_convai_mcp_tool_approval.lookup mcp-123:lookup
```