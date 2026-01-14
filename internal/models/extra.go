package models

type WorkspaceMember struct {
	UserID              string `json:"user_id"`
	Email               string `json:"email"`
	WorkspacePermission string `json:"workspace_permission"`
	IsInvited           bool   `json:"is_invited"`
}

type UpdateWorkspaceMemberRequest struct {
	WorkspacePermission string `json:"workspace_permission"`
}

type WorkspaceGroup struct {
	GroupID string `json:"group_id"`
	Name    string `json:"name"`
}

type CreateWorkspaceGroupRequest struct {
	Name string `json:"name"`
}

type ConvAIMCPServer struct {
	MCPServerID string `json:"mcp_server_id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
}

type CreateConvAIMCPServerRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ConvAIPhoneNumberAgent struct {
	AgentID   string `json:"agent_id"`
	AgentName string `json:"agent_name"`
}

type ConvAIPhoneNumber struct {
	PhoneNumberID    string                  `json:"phone_number_id"`
	PhoneNumber      string                  `json:"phone_number"`
	Provider         string                  `json:"provider"`
	Label            string                  `json:"label,omitempty"`
	SupportsInbound  bool                    `json:"supports_inbound,omitempty"`
	SupportsOutbound bool                    `json:"supports_outbound,omitempty"`
	AssignedAgent    *ConvAIPhoneNumberAgent `json:"assigned_agent,omitempty"`
	ProviderConfig   map[string]interface{}  `json:"provider_config,omitempty"`
	OutboundTrunk    map[string]interface{}  `json:"outbound_trunk,omitempty"`
	InboundTrunk     map[string]interface{}  `json:"inbound_trunk,omitempty"`
	LivekitStack     map[string]interface{}  `json:"livekit_stack,omitempty"`
}

type ImportPhoneNumberRequest struct {
	PhoneNumber string `json:"phone_number"`
	Provider    string `json:"provider"`
	Label       string `json:"label,omitempty"`
}

type StudioChapter struct {
	ChapterID string `json:"chapter_id"`
	Name      string `json:"name"`
}

type CreateStudioChapterRequest struct {
	Name string `json:"name"`
}
