package models

type WorkspaceWebhook struct {
	WebhookID string   `json:"webhook_id"`
	URL       string   `json:"url"`
	Events    []string `json:"events"`
	Secret    string   `json:"secret,omitempty"`
}

type CreateWorkspaceWebhookRequest struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

type WorkspaceInvite struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type CreateWorkspaceInviteRequest struct {
	Email               string `json:"email"`
	WorkspacePermission string `json:"workspace_permission"`
}

type ServiceAccountKey struct {
	KeyID       string   `json:"key_id"`
	XiApiKey    string   `json:"xi-api-key,omitempty"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type CreateServiceAccountKeyRequest struct {
	Name           string   `json:"name"`
	Permissions    []string `json:"permissions"`
	CharacterLimit int      `json:"character_limit,omitempty"`
}
