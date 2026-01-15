package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccModelsDataSource(t *testing.T) {
	server := newTestServer(t, []testRoute{
		{
			Method: http.MethodGet,
			Path:   "/models",
			Body:   `[{"model_id":"model-123","name":"Model","description":"Desc","can_do_text_to_speech":true,"can_do_voice_conversion":false}]`,
		},
	})
	defer server.Close()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s

data "elevenlabs_models" "all" {}
`, testAccProviderConfig(server.URL)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.elevenlabs_models.all", "models.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_models.all", "models.0.model_id", "model-123"),
				),
			},
		},
	})
}
