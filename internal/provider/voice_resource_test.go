package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVoiceResource(t *testing.T) {
	// Create a dummy file for upload
	tmpDir := t.TempDir()
	dummyFile := filepath.Join(tmpDir, "sample.mp3")
	if err := os.WriteFile(dummyFile, []byte("dummy audio content"), 0644); err != nil {
		t.Fatalf("Failed to create dummy file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/voices":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"voices": [{"voice_id": "voice-123", "name": "Test Voice", "category": "cloned"}]}`))
		
		case r.Method == http.MethodGet && r.URL.Path == "/voices/voice-123":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"voice_id": "voice-123",
				"name": "Test Voice",
				"description": "A test voice",
				"labels": {"gender": "female"},
				"settings": {
					"stability": 0.5,
					"similarity_boost": 0.75,
					"style": 0.0,
					"use_speaker_boost": true
				}
			}`))

		case r.Method == http.MethodPost && r.URL.Path == "/voices/add":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"voice_id": "voice-123"}`))

		case r.Method == http.MethodPost && r.URL.Path == "/voices/voice-123/edit":
			w.WriteHeader(http.StatusOK)

		case r.Method == http.MethodPost && r.URL.Path == "/voices/voice-123/settings/edit":
			w.WriteHeader(http.StatusOK)

		case r.Method == http.MethodDelete && r.URL.Path == "/voices/voice-123":
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

resource "elevenlabs_voice" "test" {
  name        = "Test Voice"
  description = "A test voice"
  files       = ["%s"]
  labels = {
    gender = "female"
  }
  settings = {
    stability         = 0.5
    similarity_boost  = 0.75
    style             = 0.0
    use_speaker_boost = true
  }
}

data "elevenlabs_voices" "all" {
  depends_on = [elevenlabs_voice.test]
}
`, server.URL, dummyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("elevenlabs_voice.test", "name", "Test Voice"),
					resource.TestCheckResourceAttr("elevenlabs_voice.test", "id", "voice-123"),
					resource.TestCheckResourceAttr("elevenlabs_voice.test", "description", "A test voice"),
					resource.TestCheckResourceAttr("elevenlabs_voice.test", "labels.gender", "female"),
					resource.TestCheckResourceAttr("elevenlabs_voice.test", "settings.stability", "0.5"),
					resource.TestCheckResourceAttr("data.elevenlabs_voices.all", "voices.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_voices.all", "voices.0.name", "Test Voice"),
				),
			},
		},
	})
}
