package models

type DynamicVariableAssignment struct {
	Source          string `json:"source,omitempty"`
	DynamicVariable string `json:"dynamic_variable"`
	ValuePath       string `json:"value_path"`
}

type MCPToolAddApprovalRequest struct {
	ToolName        string                 `json:"tool_name"`
	ToolDescription string                 `json:"tool_description"`
	InputSchema     map[string]interface{} `json:"input_schema,omitempty"`
	ApprovalPolicy  string                 `json:"approval_policy,omitempty"`
}

type MCPToolConfigOverride struct {
	ToolName              string                      `json:"tool_name"`
	ForcePreToolSpeech    *bool                       `json:"force_pre_tool_speech,omitempty"`
	DisableInterruptions  *bool                       `json:"disable_interruptions,omitempty"`
	ToolCallSound         *string                     `json:"tool_call_sound,omitempty"`
	ToolCallSoundBehavior *string                     `json:"tool_call_sound_behavior,omitempty"`
	ExecutionMode         *string                     `json:"execution_mode,omitempty"`
	Assignments           []DynamicVariableAssignment `json:"assignments,omitempty"`
}

type MCPToolConfigOverrideCreateRequest struct {
	ToolName              string                      `json:"tool_name"`
	ForcePreToolSpeech    *bool                       `json:"force_pre_tool_speech,omitempty"`
	DisableInterruptions  *bool                       `json:"disable_interruptions,omitempty"`
	ToolCallSound         *string                     `json:"tool_call_sound,omitempty"`
	ToolCallSoundBehavior *string                     `json:"tool_call_sound_behavior,omitempty"`
	ExecutionMode         *string                     `json:"execution_mode,omitempty"`
	Assignments           []DynamicVariableAssignment `json:"assignments,omitempty"`
}

type MCPToolConfigOverrideUpdateRequest struct {
	ForcePreToolSpeech    *bool                       `json:"force_pre_tool_speech,omitempty"`
	DisableInterruptions  *bool                       `json:"disable_interruptions,omitempty"`
	ToolCallSound         *string                     `json:"tool_call_sound,omitempty"`
	ToolCallSoundBehavior *string                     `json:"tool_call_sound_behavior,omitempty"`
	ExecutionMode         *string                     `json:"execution_mode,omitempty"`
	Assignments           []DynamicVariableAssignment `json:"assignments,omitempty"`
}
