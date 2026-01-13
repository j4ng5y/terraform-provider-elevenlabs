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
	DocumentationID string `json:"documentation_id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	Status          string `json:"status"`
}

type CreateConvAIKnowledgeBaseRequest struct {
	Name     string `json:"name"`
	URL      string `json:"url,omitempty"`
	Content  string `json:"content,omitempty"`
	FilePath string `json:"-"`
}

type ConvAITool struct {
	ToolID      string                 `json:"tool_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
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
