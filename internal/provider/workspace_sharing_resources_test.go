package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkspaceAndSharingResources(t *testing.T) {
	readWebhookCount := 0

	server := newTestServer(t, []testRoute{
		{
			Method: http.MethodPost,
			Path:   "/workspace/webhooks",
			Body:   `{"webhook_id":"webhook-123","url":"https://example.com","events":["voice_created"],"secret":"secret-123"}`,
		},
		{
			Method: http.MethodGet,
			Path:   "/workspace/webhooks/webhook-123",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if readWebhookCount == 0 {
					readWebhookCount++
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{"webhook_id":"webhook-123","url":"https://example.com","events":["voice_created"],"secret":"secret-123"}`))
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"webhook_id":"webhook-123","url":"https://example.com/updated","events":["voice_created","voice_deleted"],"secret":"secret-123"}`))
			},
		},
		{
			Method: http.MethodPatch,
			Path:   "/workspace/webhooks/webhook-123",
		},
		{
			Method: http.MethodDelete,
			Path:   "/workspace/webhooks/webhook-123",
		},
		{
			Method: http.MethodGet,
			Path:   "/workspace/webhooks",
			Body:   `[{"webhook_id":"webhook-123","url":"https://example.com","events":["voice_created"],"secret":"secret-123"}]`,
		},
		{
			Method: http.MethodPost,
			Path:   "/workspace/invites/add",
		},
		{
			Method: http.MethodDelete,
			Path:   "/workspace/invites",
		},
		{
			Method: http.MethodGet,
			Path:   "/workspace/invites",
			Body:   `[{"email":"user@example.com","role":"workspace_member"}]`,
		},
		{
			Method: http.MethodPost,
			Path:   "/workspace/members/user-123",
		},
		{
			Method: http.MethodGet,
			Path:   "/workspace/members",
			Body:   `[{"user_id":"user-123","email":"member@example.com","workspace_permission":"workspace_admin","is_invited":false}]`,
		},
		{
			Method: http.MethodGet,
			Path:   "/workspace/groups/search",
			Body:   `[{"group_id":"group-123","name":"Test Group"}]`,
		},
		{
			Method: http.MethodPost,
			Path:   "/workspace/groups/group-123/members",
		},
		{
			Method: http.MethodPost,
			Path:   "/workspace/groups/group-123/members/remove",
		},
		{
			Method: http.MethodGet,
			Path:   "/service-accounts",
			Body:   `[{"service_account_user_id":"sa-123","name":"Service Account"}]`,
		},
		{
			Method: http.MethodPost,
			Path:   "/service-accounts/sa-123/api-keys",
			Body:   `{"key_id":"key-123","xi-api-key":"api-key","name":"Key","permissions":["voice.read"]}`,
		},
		{
			Method: http.MethodDelete,
			Path:   "/service-accounts/sa-123/api-keys/key-123",
		},
		{
			Method: http.MethodGet,
			Path:   "/workspace/resources",
			Body:   `[{"resource_id":"resource-123","resource_type":"voice","name":"Shared Voice","owner_id":"owner-1","shared":true}]`,
		},
		{
			Method: http.MethodPost,
			Path:   "/workspace/resources/resource-123/share",
		},
		{
			Method: http.MethodPost,
			Path:   "/workspace/resources/resource-123/unshare",
		},
		{
			Method: http.MethodPost,
			Path:   "/voices/add/public-user/voice-456",
			Body:   `{"voice_id":"voice-789"}`,
		},
		{
			Method: http.MethodGet,
			Path:   "/voices/voice-789",
			Body:   `{"voice_id":"voice-789","name":"Shared Name"}`,
		},
		{
			Method: http.MethodDelete,
			Path:   "/voices/voice-789",
		},
	})
	defer server.Close()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s

resource "elevenlabs_workspace_webhook" "test" {
  url    = "https://example.com"
  events = ["voice_created"]
}

resource "elevenlabs_workspace_invite" "invite" {
  email = "user@example.com"
  role  = "workspace_member"
}

resource "elevenlabs_workspace_member" "member" {
  id                   = "user-123"
  workspace_permission = "workspace_admin"
}

resource "elevenlabs_workspace_group_member" "member" {
  group_id = "group-123"
  email    = "group@example.com"
}

resource "elevenlabs_service_account_key" "key" {
  user_id         = "sa-123"
  name            = "Key"
  permissions     = ["voice.read"]
  character_limit = 1000
}

resource "elevenlabs_resource_share" "share" {
  resource_id   = "resource-123"
  resource_type = "voice"
  email         = "share@example.com"
  role          = "viewer"
}

resource "elevenlabs_shared_voice" "shared" {
  public_user_id = "public-user"
  voice_id       = "voice-456"
  name           = "Shared Name"
}

data "elevenlabs_workspace_webhooks" "all" {}

data "elevenlabs_workspace_invites" "all" {}

data "elevenlabs_workspace_members" "all" {}

data "elevenlabs_workspace_groups" "all" {}

data "elevenlabs_workspace_service_accounts" "all" {}

data "elevenlabs_workspace_resources" "all" {}
`, testAccProviderConfig(server.URL)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("elevenlabs_workspace_webhook.test", "id", "webhook-123"),
					resource.TestCheckResourceAttr("elevenlabs_workspace_invite.invite", "email", "user@example.com"),
					resource.TestCheckResourceAttr("elevenlabs_workspace_member.member", "id", "user-123"),
					resource.TestCheckResourceAttr("elevenlabs_workspace_group_member.member", "group_id", "group-123"),
					resource.TestCheckResourceAttr("elevenlabs_service_account_key.key", "id", "key-123"),
					resource.TestCheckResourceAttr("elevenlabs_resource_share.share", "id", "voice:resource-123:share@example.com"),
					resource.TestCheckResourceAttr("elevenlabs_shared_voice.shared", "id", "voice-789"),
					resource.TestCheckResourceAttr("data.elevenlabs_workspace_webhooks.all", "webhooks.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_workspace_invites.all", "invites.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_workspace_members.all", "members.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_workspace_groups.all", "groups.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_workspace_service_accounts.all", "service_accounts.#", "1"),
					resource.TestCheckResourceAttr("data.elevenlabs_workspace_resources.all", "resources.#", "1"),
				),
			},
			{
				Config: fmt.Sprintf(`
%s

resource "elevenlabs_workspace_webhook" "test" {
  url    = "https://example.com/updated"
  events = ["voice_created", "voice_deleted"]
}
`, testAccProviderConfig(server.URL)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("elevenlabs_workspace_webhook.test", "url", "https://example.com/updated"),
				),
			},
		},
	})
}
