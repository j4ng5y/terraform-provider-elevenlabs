package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConvAIResourcesAndDataSources(t *testing.T) {
	server := newTestServer(t, []testRoute{
		{
			Method: http.MethodPost,
			Path:   "/convai/knowledge-base",
			Body:   `{"documentation_id":"kb-123","name":"KB","type":"text","status":"ready"}`,
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/knowledge-base/kb-123",
			Body:   `{"documentation_id":"kb-123","name":"KB","type":"text","status":"ready"}`,
		},
		{
			Method: http.MethodDelete,
			Path:   "/convai/knowledge-base/kb-123",
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/knowledge-base",
			Body:   `{"documents":[{"documentation_id":"kb-123","name":"KB","type":"text","status":"ready"}],"has_more":false,"next_cursor":""}`,
		},
		{
			Method: http.MethodPost,
			Path:   "/convai/tools",
			Body:   `{"tool_id":"tool-123","name":"Tool","description":"Desc"}`,
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/tools/tool-123",
			Body:   `{"tool_id":"tool-123","name":"Tool","description":"Desc"}`,
		},
		{
			Method: http.MethodPatch,
			Path:   "/convai/tools/tool-123",
		},
		{
			Method: http.MethodDelete,
			Path:   "/convai/tools/tool-123",
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/tools",
			Body:   `{"tools":[{"tool_id":"tool-123","name":"Tool","description":"Desc"}]}`,
		},
		{
			Method: http.MethodPost,
			Path:   "/convai/mcp-servers",
			Body:   `{"mcp_server_id":"mcp-123","name":"MCP","url":"https://mcp.example.com"}`,
		},
		{
			Method: http.MethodPatch,
			Path:   "/convai/mcp-servers/mcp-123",
		},
		{
			Method: http.MethodDelete,
			Path:   "/convai/mcp-servers/mcp-123",
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/mcp-servers",
			Body:   `[{"mcp_server_id":"mcp-123","name":"MCP","url":"https://mcp.example.com"}]`,
		},
		{
			Method: http.MethodPost,
			Path:   "/convai/phone-numbers",
			Body:   `{"phone_number_id":"phone-123","phone_number":"+123","provider":"twilio","label":"Main","supports_inbound":true,"supports_outbound":true,"assigned_agent":{"agent_id":"agent-123","agent_name":"Agent"}}`,
		},
		{
			Method: http.MethodPatch,
			Path:   "/convai/phone-numbers/phone-123",
		},
		{
			Method: http.MethodDelete,
			Path:   "/convai/phone-numbers/phone-123",
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/phone-numbers",
			Body:   `[{"phone_number_id":"phone-123","phone_number":"+123","provider":"twilio","label":"Main","supports_inbound":true,"supports_outbound":true,"assigned_agent":{"agent_id":"agent-123","agent_name":"Agent"}}]`,
		},
		{
			Method: http.MethodPost,
			Path:   "/convai/whatsapp-accounts",
			Body:   `{"business_account_id":"ba-123","business_account_name":"Biz","phone_number_id":"wa-123","phone_number_name":"WA","phone_number":"+456","assigned_agent_id":"agent-123","assigned_agent_name":"Agent"}`,
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/whatsapp-accounts/wa-123",
			Body:   `{"business_account_id":"ba-123","business_account_name":"Biz","phone_number_id":"wa-123","phone_number_name":"WA","phone_number":"+456","assigned_agent_id":"agent-123","assigned_agent_name":"Agent"}`,
		},
		{
			Method: http.MethodPatch,
			Path:   "/convai/whatsapp-accounts/wa-123",
			Body:   `{"business_account_id":"ba-123","business_account_name":"Biz","phone_number_id":"wa-123","phone_number_name":"WA","phone_number":"+456","assigned_agent_id":"agent-123","assigned_agent_name":"Agent"}`,
		},
		{
			Method: http.MethodDelete,
			Path:   "/convai/whatsapp-accounts/wa-123",
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/whatsapp-accounts",
			Body:   `{"items":[{"business_account_id":"ba-123","business_account_name":"Biz","phone_number_id":"wa-123","phone_number_name":"WA","phone_number":"+456","assigned_agent_id":"agent-123","assigned_agent_name":"Agent"}]}`,
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/settings",
			Body:   `{"recording_enabled":true}`,
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/conversations/conv-123",
			Body:   `{"conversation_id":"conv-123","agent_id":"agent-123","name":"Conversation","created_at":"2024-01-01T00:00:00Z"}`,
		},
		{
			Method: http.MethodDelete,
			Path:   "/convai/conversations/conv-123",
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/conversations",
			Body:   `[{"conversation_id":"conv-123","name":"Conversation","agent_id":"agent-123","agent_name":"Agent","created_at":"2024-01-01T00:00:00Z"}]`,
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/conversation/get-signed-url",
			Body:   `{"conversation_signature":"sig-123","conversation_id":"conv-123"}`,
		},
		{
			Method: http.MethodPost,
			Path:   "/convai/agents/agent-123/simulate-conversation",
			Body:   `{"simulated_conversation":[{"role":"assistant","content":"Hi","timestamp":"2024-01-01T00:00:00Z"}],"analysis":{"goal_achieved":true,"total_duration_ms":12,"number_of_turns":1,"summary":"ok"}}`,
		},
		{
			Method: http.MethodPost,
			Path:   "/convai/agent-testing/create",
			Body:   `{"test_id":"test-123","name":"Test","success_condition":"ok"}`,
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/agent-testing/test-123",
			Body:   `{"test_id":"test-123","name":"Test","success_condition":"ok"}`,
		},
		{
			Method: http.MethodDelete,
			Path:   "/convai/agent-testing/test-123",
		},
		{
			Method: http.MethodPost,
			Path:   "/convai/agents/agent-123/run-tests",
			Body:   `{"test_invocation_id":"inv-123","status":"completed","results":[{"test_id":"test-123","status":"passed","passed":true,"message":"ok","duration_ms":5}],"started_at":"2024-01-01T00:00:00Z","completed_at":"2024-01-01T00:00:05Z"}`,
		},
		{
			Method: http.MethodPost,
			Path:   "/convai/agents/agent-123/duplicate",
			Body:   `{"agent_id":"agent-dup","name":"Agent Copy","config":{"prompt":"hi"}}`,
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/agents/agent-dup",
			Body:   `{"agent_id":"agent-dup","name":"Agent Copy","config":{"prompt":"hi"}}`,
		},
		{
			Method: http.MethodDelete,
			Path:   "/convai/agents/agent-dup",
		},
		{
			Method: http.MethodPost,
			Path:   "/convai/agent/agent-123/llm-usage/calculate",
			Body:   `{"llm_prices":[{"model_name":"gpt","input_cost_per_million_tokens":1.0,"output_cost_per_million_tokens":2.0,"estimated_input_tokens":100,"estimated_output_tokens":200,"estimated_cost":0.3}]}`,
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/dashboard/settings",
			Body:   `{"analytics_enabled":true,"recording_enabled":false,"transcription_enabled":true,"LLm_optimization_enabled":true}`,
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/secrets",
			Body:   `[{"secret_id":"secret-123","name":"Secret"}]`,
		},
		{
			Method: http.MethodGet,
			Path:   "/convai/batch-calling/workspace",
			Body:   `[{"batch_id":"batch-123","agent_id":"agent-123","phone_number_id":"phone-123","status":"running","total_calls":5,"completed_calls":2,"failed_calls":0,"created_at":"2024-01-01T00:00:00Z"}]`,
		},
	})
	defer server.Close()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s

resource "elevenlabs_convai_knowledge_base" "kb" {
  name    = "KB"
  content = "Hello"
}

resource "elevenlabs_convai_tool" "tool" {
  name        = "Tool"
  description = "Desc"
}

resource "elevenlabs_convai_mcp_server" "mcp" {
  name = "MCP"
  url  = "https://mcp.example.com"
}

resource "elevenlabs_convai_phone_number" "phone" {
  phone_number       = "+123"
  telephony_provider = "twilio"
  label              = "Main"
}

resource "elevenlabs_convai_whatsapp_account" "wa" {
  phone_number_id     = "wa-123"
  business_account_id = "ba-123"
  token_code          = "token"
  assigned_agent_id   = "agent-123"
}

resource "elevenlabs_convai_settings" "settings" {}

resource "elevenlabs_convai_conversation" "conv" {
  conversation_id = "conv-123"
  agent_id        = "agent-123"
}

resource "elevenlabs_convai_conversation_simulation" "sim" {
  agent_id = "agent-123"
  chat_history = [
    {
      role    = "user"
      content = "Hi"
    }
  ]
}

resource "elevenlabs_convai_agent_test" "test" {
  name              = "Test"
  success_condition = "ok"
}

resource "elevenlabs_convai_agent_test_runner" "runner" {
  agent_id = "agent-123"
  test_ids = [elevenlabs_convai_agent_test.test.id]
}

resource "elevenlabs_convai_agent_duplicate" "dup" {
  source_agent_id = "agent-123"
  new_agent_name  = "Agent Copy"
}

data "elevenlabs_convai_knowledge_bases" "all" {}

data "elevenlabs_convai_tools" "all" {}

data "elevenlabs_convai_mcp_servers" "all" {}

data "elevenlabs_convai_phone_numbers" "all" {}

data "elevenlabs_convai_whatsapp_accounts" "all" {}

data "elevenlabs_convai_conversations" "all" {}

data "elevenlabs_convai_dashboard_settings" "all" {}

data "elevenlabs_convai_secrets" "all" {}

data "elevenlabs_convai_batch_calling" "all" {}

data "elevenlabs_convai_signed_url" "signed" {
  agent_id                = "agent-123"
  include_conversation_id = true
}

data "elevenlabs_convai_llm_usage_calculator" "usage" {
  agent_id        = "agent-123"
  prompt_length   = 10
  number_of_pages = 2
  rag_enabled     = true
}
`, testAccProviderConfig(server.URL)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("elevenlabs_convai_knowledge_base.kb", "id", "kb-123"),
					resource.TestCheckResourceAttr("elevenlabs_convai_tool.tool", "id", "tool-123"),
					resource.TestCheckResourceAttr("elevenlabs_convai_mcp_server.mcp", "id", "mcp-123"),
					resource.TestCheckResourceAttr("elevenlabs_convai_phone_number.phone", "id", "phone-123"),
					resource.TestCheckResourceAttr("elevenlabs_convai_whatsapp_account.wa", "phone_number", "+456"),
					resource.TestCheckResourceAttr("elevenlabs_convai_conversation.conv", "name", "Conversation"),
					resource.TestCheckResourceAttr("elevenlabs_convai_agent_test.test", "id", "test-123"),
					resource.TestCheckResourceAttr("elevenlabs_convai_agent_test_runner.runner", "test_invocation_id", "inv-123"),
					resource.TestCheckResourceAttr("elevenlabs_convai_agent_duplicate.dup", "new_agent_id", "agent-dup"),
					resource.TestCheckResourceAttr("data.elevenlabs_convai_knowledge_bases.all", "documents.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_convai_tools.all", "tools.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_convai_mcp_servers.all", "mcp_servers.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_convai_phone_numbers.all", "phone_numbers.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_convai_whatsapp_accounts.all", "whatsapp_accounts.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_convai_conversations.all", "conversations.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_convai_dashboard_settings.all", "analytics_enabled", "true"),
					resource.TestCheckResourceAttr("data.elevenlabs_convai_secrets.all", "secrets.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_convai_batch_calling.all", "batch_calls.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_convai_signed_url.signed", "conversation_signature", "sig-123"),
					resource.TestCheckResourceAttr("data.elevenlabs_convai_llm_usage_calculator.usage", "llm_prices.#", "1"),
				),
			},
		},
	})
}
