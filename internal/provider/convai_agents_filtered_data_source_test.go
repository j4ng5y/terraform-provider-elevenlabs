package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConvAIAgentsFilteredDataSource(t *testing.T) {
	server := newTestServer(t, []testRoute{
		{
			Method: http.MethodGet,
			Path:   "/convai/agents",
			Body:   `{"agents":[{"id":"agent-123","name":"Filtered","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-02T00:00:00Z"}]}`,
		},
	})
	defer server.Close()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s

data "elevenlabs_convai_agents_filtered" "all" {
  page_size             = 10
  search                = "Filtered"
  archived              = false
  show_only_owned_agents = true
}
`, testAccProviderConfig(server.URL)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.elevenlabs_convai_agents_filtered.all", "agents.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_convai_agents_filtered.all", "agents.0.id", "agent-123"),
				),
			},
		},
	})
}
