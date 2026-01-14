package models

type ConvAIAgent struct {
	AgentID string             `json:"agent_id"`
	Name    string             `json:"name"`
	Config  *ConvAIAgentConfig `json:"config"`
}

type ConvAIAgentConfig struct {
	Prompt       string `json:"prompt"`
	FirstMessage string `json:"first_message,omitempty"`
	Language     string `json:"language,omitempty"`
	ModelID      string `json:"model_id,omitempty"`
}

type CreateConvAIAgentRequest struct {
	Name   string             `json:"name"`
	Config *ConvAIAgentConfig `json:"config"`
}

type ConvAIKnowledgeBase struct {
	DocumentationID string                                 `json:"documentation_id"`
	Name            string                                 `json:"name"`
	Type            string                                 `json:"type"`
	Status          string                                 `json:"status"`
	Metadata        map[string]interface{}                 `json:"metadata,omitempty"`
	SupportedUsages []string                               `json:"supported_usages,omitempty"`
	FolderParentID  string                                 `json:"folder_parent_id,omitempty"`
	FolderPath      []ConvAIKnowledgeBaseFolderPathSegment `json:"folder_path,omitempty"`
	AccessInfo      map[string]interface{}                 `json:"access_info,omitempty"`
}

type CreateConvAIKnowledgeBaseRequest struct {
	Name     string `json:"name"`
	URL      string `json:"url,omitempty"`
	Content  string `json:"content,omitempty"`
	FilePath string `json:"-"`
}

type ListConvAIKnowledgeBaseDocumentsParams struct {
	PageSize      *int
	Search        string
	Cursor        string
	ShowOnlyOwned *bool
	Types         []string
}

type ConvAIKnowledgeBaseListResponse struct {
	Documents  []ConvAIKnowledgeBase `json:"documents"`
	HasMore    bool                  `json:"has_more"`
	NextCursor string                `json:"next_cursor"`
}

type ConvAIKnowledgeBaseFolderPathSegment struct {
	FolderID   string `json:"folder_id"`
	FolderName string `json:"folder_name"`
}

type ConvAITool struct {
	ToolID            string                 `json:"tool_id"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Parameters        map[string]interface{} `json:"parameters"`
	DependentAgentIDs []string               `json:"dependent_agent_ids,omitempty"`
	AccessInfo        map[string]interface{} `json:"access_info,omitempty"`
	ToolConfig        map[string]interface{} `json:"tool_config,omitempty"`
	UsageStats        map[string]interface{} `json:"usage_stats,omitempty"`
}

type CreateConvAIToolRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type ConvAISecret struct {
	SecretID string `json:"secret_id"`
	Name     string `json:"name"`
}

type CreateConvAISecretRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ConvAIAgentTest struct {
	TestID           string `json:"test_id"`
	Name             string `json:"name"`
	SuccessCondition string `json:"success_condition"`
}

type CreateConvAIAgentTestRequest struct {
	Name             string        `json:"name"`
	ChatHistory      []interface{} `json:"chat_history"`
	SuccessCondition string        `json:"success_condition"`
	SuccessExamples  []string      `json:"success_examples"`
	FailureExamples  []string      `json:"failure_examples"`
}

type ConvAIToolsResponse struct {
	Tools []ConvAITool `json:"tools"`
}

type ConvAIWhatsAppAccount struct {
	BusinessAccountID   string `json:"business_account_id"`
	BusinessAccountName string `json:"business_account_name"`
	PhoneNumberID       string `json:"phone_number_id"`
	PhoneNumberName     string `json:"phone_number_name"`
	PhoneNumber         string `json:"phone_number"`
	AssignedAgentID     string `json:"assigned_agent_id,omitempty"`
	AssignedAgentName   string `json:"assigned_agent_name,omitempty"`
}

type ConvAIWhatsAppAccountListResponse struct {
	Items []ConvAIWhatsAppAccount `json:"items"`
}

type ImportWhatsAppAccountRequest struct {
	BusinessAccountID string `json:"business_account_id"`
	PhoneNumberID     string `json:"phone_number_id"`
	TokenCode         string `json:"token_code"`
}

type UpdateWhatsAppAccountRequest struct {
	AssignedAgentID *string `json:"assigned_agent_id,omitempty"`
}
