package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPVCResources(t *testing.T) {
	sampleFile := writeTempFile(t, "pvc-sample.wav", []byte("sample content"))
	description := "Initial"

	server := newTestServer(t, []testRoute{
		{
			Method: http.MethodPost,
			Path:   "/voices/pvc",
			Body:   `{"voice_id":"voice-123","name":"PVC Voice","language":"en","description":"Initial","state":"ready","verification":"verified"}`,
		},
		{
			Method: http.MethodGet,
			Path:   "/voices/pvc/voice-123",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, `{"voice_id":"voice-123","name":"PVC Voice","language":"en","description":"%s","state":"ready","verification":"verified"}`, description)
			},
		},
		{
			Method: http.MethodPatch,
			Path:   "/voices/pvc/voice-123",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				description = "Updated"
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			Method: http.MethodDelete,
			Path:   "/voices/pvc/voice-123",
		},
		{
			Method: http.MethodGet,
			Path:   "/voices/pvc",
			Body:   `{"voices":[{"voice_id":"voice-123","name":"PVC Voice","language":"en"}]}`,
		},
		{
			Method: http.MethodPost,
			Path:   "/voices/pvc/voice-123/samples",
			Body:   `{"sample_id":"sample-123","file_name":"pvc-sample.wav","mime_type":"audio/wav","size_bytes":42,"hash":"hash","state":"ready","transcription":"hello","duration":1.23,"sample_rate":44100,"channels":1}`,
		},
		{
			Method: http.MethodGet,
			Path:   "/voices/pvc/voice-123/samples",
			Body:   `{"samples":[{"sample_id":"sample-123","file_name":"pvc-sample.wav","mime_type":"audio/wav","size_bytes":42,"hash":"hash","state":"ready","transcription":"hello","duration":1.23,"sample_rate":44100,"channels":1}]}`,
		},
		{
			Method: http.MethodDelete,
			Path:   "/voices/pvc/voice-123/samples/sample-123",
		},
	})
	defer server.Close()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s

resource "elevenlabs_pvc_voice" "voice" {
  name        = "PVC Voice"
  language    = "en"
  description = "Initial"
  labels = {
    env = "test"
  }
}

resource "elevenlabs_pvc_voice_sample" "sample" {
  voice_id  = elevenlabs_pvc_voice.voice.id
  file_path = "%s"
}

data "elevenlabs_pvc_voices" "all" {}

data "elevenlabs_pvc_voice_samples" "all" {
  voice_id = elevenlabs_pvc_voice.voice.id
}
`, testAccProviderConfig(server.URL), sampleFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("elevenlabs_pvc_voice.voice", "id", "voice-123"),
					resource.TestCheckResourceAttr("elevenlabs_pvc_voice_sample.sample", "id", "sample-123"),
					resource.TestCheckResourceAttr("data.elevenlabs_pvc_voices.all", "voices.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_pvc_voice_samples.all", "samples.#", "1"),
				),
			},
			{
				Config: fmt.Sprintf(`
%s

resource "elevenlabs_pvc_voice" "voice" {
  name        = "PVC Voice"
  language    = "en"
  description = "Updated"
}

resource "elevenlabs_pvc_voice_sample" "sample" {
  voice_id  = elevenlabs_pvc_voice.voice.id
  file_path = "%s"
}
`, testAccProviderConfig(server.URL), sampleFile),
			},
		},
	})
}
