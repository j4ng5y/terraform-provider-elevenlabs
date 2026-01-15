package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/projects":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"projects": [{"project_id": "project-123", "name": "Test Project", "state": "default"}]}`))

		case r.Method == http.MethodGet && r.URL.Path == "/projects/project-123":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"project_id": "project-123",
				"name": "Test Project",
				"default_model_id": "model-1",
				"default_paragraph_voice_id": "voice-1",
				"default_title_voice_id": "voice-2",
				"state": "default"
			}`))

		case r.Method == http.MethodPost && r.URL.Path == "/projects":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"project_id": "project-123",
				"name": "Test Project",
				"default_model_id": "model-1",
				"default_paragraph_voice_id": "voice-1",
				"default_title_voice_id": "voice-2",
				"state": "default"
			}`))

		case r.Method == http.MethodDelete && r.URL.Path == "/projects/project-123":
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

resource "elevenlabs_project" "test" {
  name                        = "Test Project"
  default_model_id            = "model-1"
  default_paragraph_voice_id  = "voice-1"
  default_title_voice_id      = "voice-2"
}

data "elevenlabs_projects" "all" {
  depends_on = [elevenlabs_project.test]
}
`, server.URL),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("elevenlabs_project.test", "name", "Test Project"),
					resource.TestCheckResourceAttr("elevenlabs_project.test", "id", "project-123"),
					resource.TestCheckResourceAttr("elevenlabs_project.test", "default_model_id", "model-1"),
					resource.TestCheckResourceAttr("data.elevenlabs_projects.all", "projects.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_projects.all", "projects.0.name", "Test Project"),
				),
			},
		},
	})
}
