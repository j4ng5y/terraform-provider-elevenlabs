package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPronunciationDictionaryUpdateResource(t *testing.T) {
	readCount := 0

	server := newTestServer(t, []testRoute{
		{
			Method: http.MethodGet,
			Path:   "/pronunciation-dictionaries/dict-123",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if readCount == 0 {
					readCount++
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{"id":"dict-123","name":"Test Dict","latest_version_id":"version-1"}`))
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"id":"dict-123","name":"Updated Dict","latest_version_id":"version-2"}`))
			},
		},
		{
			Method: http.MethodPatch,
			Path:   "/pronunciation-dictionaries/dict-123",
		},
	})
	defer server.Close()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s

resource "elevenlabs_pronunciation_dictionary_update" "dict" {
  id       = "dict-123"
  name     = "Initial Dict"
  archived = true
}
`, testAccProviderConfig(server.URL)),
				ResourceName:  "elevenlabs_pronunciation_dictionary_update.dict",
				ImportState:   true,
				ImportStateId: "dict-123",
				ImportStatePersist: true,
			},
			{
				Config: fmt.Sprintf(`
%s

resource "elevenlabs_pronunciation_dictionary_update" "dict" {
  id       = "dict-123"
  name     = "Updated Dict"
  archived = true
}
`, testAccProviderConfig(server.URL)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("elevenlabs_pronunciation_dictionary_update.dict", "version_id", "version-2"),
				),
			},
		},
	})
}
