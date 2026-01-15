package provider

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPronunciationAndAudioResources(t *testing.T) {
	audioFile := writeTempFile(t, "audio.txt", []byte("audio content"))
	updateFile := writeTempFile(t, "audio-update.txt", []byte("updated audio"))
	voiceSampleFile := writeTempFile(t, "sample.wav", []byte("sample content"))
	outputPath := filepath.Join(t.TempDir(), "dict.pls")

	server := newTestServer(t, []testRoute{
		{
			Method: httpMethodPost,
			Path:   "/pronunciation-dictionaries/add-from-rules",
			Body:   `{"id":"dict-123","name":"Test Dict","latest_version_id":"version-1"}`,
		},
		{
			Method: httpMethodGet,
			Path:   "/pronunciation-dictionaries/dict-123",
			Body:   `{"id":"dict-123","name":"Test Dict","latest_version_id":"version-1"}`,
		},
		{
			Method: httpMethodGet,
			Path:   "/pronunciation-dictionaries",
			Body:   `{"pronunciation_dictionaries":[{"id":"dict-123","name":"Test Dict"}]}`,
		},
		{
			Method: httpMethodPost,
			Path:   "/pronunciation-dictionaries/dict-123/add-rules",
		},
		{
			Method: httpMethodPatch,
			Path:   "/pronunciation-dictionaries/dict-123",
		},
		{
			Method: httpMethodGet,
			Path:   "/pronunciation-dictionaries/dict-123/version-1/download",
			Body:   `"dGVzdA=="`,
		},
		{
			Method: httpMethodPost,
			Path:   "/audio-native",
			Body:   `{"project_id":"audio-123","html_snippet":"<div>audio</div>"}`,
		},
		{
			Method: httpMethodGet,
			Path:   "/audio-native/audio-123/settings",
			Body:   `{"title":"Test Title","author":"Test Author","text_color":"#000000","background_color":"#ffffff","status":"ready"}`,
		},
		{
			Method: httpMethodPost,
			Path:   "/audio-native/audio-123/content",
		},
		{
			Method: httpMethodDelete,
			Path:   "/projects/audio-123",
		},
		{
			Method: httpMethodPost,
			Path:   "/voices/pvc/voice-123/samples",
			Body:   `{"sample_id":"sample-123","file_name":"sample.wav"}`,
		},
		{
			Method: httpMethodDelete,
			Path:   "/voices/voice-123/samples/sample-123",
		},
	})
	defer server.Close()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s

resource "elevenlabs_pronunciation_dictionary" "dict" {
  name        = "Test Dict"
  description = "Test description"
  rules = [
    {
      type              = "alias"
      string_to_replace = "NY"
      alias             = "New York"
    }
  ]
}

resource "elevenlabs_pronunciation_dictionary_rules" "dict_rules" {
  dictionary_id = elevenlabs_pronunciation_dictionary.dict.id
  action        = "add"
  rules = [
    {
      type              = "alias"
      string_to_replace = "CA"
      alias             = "California"
    }
  ]
}

data "elevenlabs_pronunciation_dictionaries" "all" {}

data "elevenlabs_pronunciation_dictionary_download" "pls" {
  dictionary_id = elevenlabs_pronunciation_dictionary.dict.id
  output_path   = "%s"
}

resource "elevenlabs_audio_native" "audio" {
  name             = "Test Audio"
  file_path        = "%s"
  title            = "Test Title"
  author           = "Test Author"
  text_color       = "#000000"
  background_color = "#ffffff"
  auto_convert     = true
}

resource "elevenlabs_audio_native_content_update" "audio_update" {
  project_id = elevenlabs_audio_native.audio.id
  file_path  = "%s"
}

resource "elevenlabs_voice_sample" "sample" {
  voice_id  = "voice-123"
  file_path = "%s"
}
`, testAccProviderConfig(server.URL), outputPath, audioFile, updateFile, voiceSampleFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("elevenlabs_pronunciation_dictionary.dict", "id", "dict-123"),
					resource.TestCheckResourceAttr("elevenlabs_pronunciation_dictionary.dict", "latest_version_id", "version-1"),
					resource.TestCheckResourceAttr("data.elevenlabs_pronunciation_dictionaries.all", "dictionaries.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_pronunciation_dictionary_download.pls", "file_name", "dict.pls"),
					resource.TestCheckResourceAttr("elevenlabs_audio_native.audio", "id", "audio-123"),
					resource.TestCheckResourceAttr("elevenlabs_audio_native.audio", "status", "ready"),
					resource.TestCheckResourceAttr("elevenlabs_voice_sample.sample", "id", "sample-123"),
				),
			},
		},
	})
}

const (
	httpMethodGet    = "GET"
	httpMethodPost   = "POST"
	httpMethodPatch  = "PATCH"
	httpMethodDelete = "DELETE"
)
