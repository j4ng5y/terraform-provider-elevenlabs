package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConvAISecretResource(t *testing.T) {
	// Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Create Secret
		if r.Method == http.MethodPost && r.URL.Path == "/convai/secrets" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"secret_id": "secret-123",
				"name": "test-secret",
				"created_at_unix_secs": 1600000000
			}`))
			return
		}

		// Read Secret
		if r.Method == http.MethodGet && r.URL.Path == "/convai/secrets/secret-123" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"secret_id": "secret-123",
				"name": "test-secret",
				"created_at_unix_secs": 1600000000
			}`))
			return
		}

		// Update Secret
		if r.Method == http.MethodPatch && r.URL.Path == "/convai/secrets/secret-123" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Delete Secret
		if r.Method == http.MethodDelete && r.URL.Path == "/convai/secrets/secret-123" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Default: Not Found
		http.Error(w, "Not Found", http.StatusNotFound)
	}))
	defer server.Close()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: fmt.Sprintf(`
provider "elevenlabs" {
  api_key  = "test-key"
  base_url = "%s"
}

resource "elevenlabs_convai_secret" "test" {
  name  = "test-secret"
  value = "super-secret-value"
}
`, server.URL),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("elevenlabs_convai_secret.test", "name", "test-secret"),
					resource.TestCheckResourceAttr("elevenlabs_convai_secret.test", "id", "secret-123"),
				),
			},
			// Update (Note: In this mock, update doesn't change the ID or Name returned by Read, 
			// but we can verify the state updates if we changed the config. 
			// However, since Name is the only non-sensitive attribute we can check, let's keep it simple.)
			{
				Config: fmt.Sprintf(`
provider "elevenlabs" {
  api_key  = "test-key"
  base_url = "%s"
}

resource "elevenlabs_convai_secret" "test" {
  name  = "test-secret"
  value = "new-super-secret-value"
}
`, server.URL),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("elevenlabs_convai_secret.test", "name", "test-secret"),
				),
			},
		},
	})
}
