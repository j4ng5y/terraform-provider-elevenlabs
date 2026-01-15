package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConvAIAgentResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/convai/agents":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"agents": [{"agent_id": "agent-123", "name": "Test Agent"}]}`))

		case r.Method == http.MethodGet && r.URL.Path == "/convai/agents/agent-123":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"agent_id": "agent-123",
				"name": "Test Agent",
				"config": {
					"prompt": "You are a helpful assistant",
					"first_message": "Hello!",
					"language": "en",
					"model_id": "model-1"
				}
			}`))

		case r.Method == http.MethodPost && r.URL.Path == "/convai/agents/create":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"agent_id": "agent-123",
				"name": "Test Agent",
				"config": {
					"prompt": "You are a helpful assistant"
				}
			}`))

		case r.Method == http.MethodPatch && r.URL.Path == "/convai/agents/agent-123":
			w.WriteHeader(http.StatusOK)

		case r.Method == http.MethodDelete && r.URL.Path == "/convai/agents/agent-123":
			w.WriteHeader(http.StatusOK)

		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	}))
	defer server.Close()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "elevenlabs" {
  api_key  = "test-key"
  base_url = "%s"
}

resource "elevenlabs_convai_agent" "test" {
  name          = "Test Agent"
  prompt        = "You are a helpful assistant"
  first_message = "Hello!"
  language      = "en"
  model_id      = "model-1"
}

data "elevenlabs_convai_agents" "all" {
  depends_on = [elevenlabs_convai_agent.test]
}
`, server.URL),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("elevenlabs_convai_agent.test", "name", "Test Agent"),
					resource.TestCheckResourceAttr("elevenlabs_convai_agent.test", "id", "agent-123"),
					resource.TestCheckResourceAttr("elevenlabs_convai_agent.test", "prompt", "You are a helpful assistant"),
					resource.TestCheckResourceAttr("data.elevenlabs_convai_agents.all", "agents.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_convai_agents.all", "agents.0.name", "Test Agent"),
				),
			},
		},
	})
}
