package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

var (
	_ resource.Resource              = &ConvAIConversationSimulatorResource{}
	_ resource.ResourceWithConfigure = &ConvAIConversationSimulatorResource{}
)

func NewConvAIConversationSimulatorResource() resource.Resource {
	return &ConvAIConversationSimulatorResource{}
}

type ConvAIConversationSimulatorResource struct {
	client *client.Client
}

type ConvAIConversationSimulatorResourceModel struct {
	AgentID       types.String            `tfsdk:"agent_id"`
	ChatHistory   []ChatMessageModel      `tfsdk:"chat_history"`
	AgentConfig   types.Map               `tfsdk:"agent_configuration"`
	SimulatedChat types.List              `tfsdk:"simulated_conversation"`
	Analysis      types.Object            `tfsdk:"analysis"`
}

type ChatMessageModel struct {
	Role    types.String `tfsdk:"role"`
	Content types.String `tfsdk:"content"`
}

type SimulatedMessageModel struct {
	Role      types.String `tfsdk:"role"`
	Content   types.String `tfsdk:"content"`
	Timestamp types.String `tfsdk:"timestamp"`
}

var simulatedMessageObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"role":      types.StringType,
	"content":   types.StringType,
	"timestamp": types.StringType,
}}

type AnalysisModel struct {
	GoalAchieved  types.Bool   `tfsdk:"goal_achieved"`
	TotalDuration types.Int64  `tfsdk:"total_duration_ms"`
	NumberOfTurns types.Int64  `tfsdk:"number_of_turns"`
	Summary       types.String `tfsdk:"summary"`
}

func (r *ConvAIConversationSimulatorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_conversation_simulation"
}

func (r *ConvAIConversationSimulatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for simulating conversations with ElevenLabs ConvAI agents. Tests agent behavior with custom chat history and configuration.",
		Attributes: map[string]schema.Attribute{
			"agent_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the agent to simulate conversation with.",
			},
			"chat_history": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"role": schema.StringAttribute{
							Required: true,
						},
						"content": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"agent_configuration": schema.MapAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Optional agent configuration to use during simulation.",
			},
			"simulated_conversation": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"role": schema.StringAttribute{
							Computed: true,
						},
						"content": schema.StringAttribute{
							Computed: true,
						},
						"timestamp": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"analysis": schema.ObjectAttribute{
				Computed: true,
				AttributeTypes: map[string]attr.Type{
					"goal_achieved":     types.BoolType,
					"total_duration_ms": types.Int64Type,
					"number_of_turns":   types.Int64Type,
					"summary":           types.StringType,
				},
			},
		},
	}
}

func (r *ConvAIConversationSimulatorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *ConvAIConversationSimulatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConvAIConversationSimulatorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.SimulatedChat = types.ListNull(simulatedMessageObjectType)

	// Convert chat history
	chatHistory := make([]map[string]interface{}, len(plan.ChatHistory))
	for i, msg := range plan.ChatHistory {
		chatHistory[i] = map[string]interface{}{
			"role":    msg.Role.ValueString(),
			"content": msg.Content.ValueString(),
		}
	}

	// Convert agent configuration
	var agentConfig map[string]interface{}
	if !plan.AgentConfig.IsNull() {
		agentConfigElements := plan.AgentConfig.Elements()
		agentConfig = make(map[string]interface{})
		for key, value := range agentConfigElements {
			agentConfig[key] = value
		}
	}

	result, err := r.client.SimulateConversation(plan.AgentID.ValueString(), chatHistory, agentConfig)
	if err != nil {
		resp.Diagnostics.AddError("Error simulating conversation", err.Error())
		return
	}

	// Process simulated conversation
	if simulated, ok := result["simulated_conversation"].([]interface{}); ok {
		simulatedMessages := make([]SimulatedMessageModel, len(simulated))
		for i, msg := range simulated {
			m := msg.(map[string]interface{})
			simulatedMessages[i] = SimulatedMessageModel{
				Role:      types.StringValue(m["role"].(string)),
				Content:   types.StringValue(m["content"].(string)),
				Timestamp: types.StringValue(m["timestamp"].(string)),
			}
		}
		listValue, diags := types.ListValueFrom(ctx, simulatedMessageObjectType, simulatedMessages)
		resp.Diagnostics.Append(diags...)
		plan.SimulatedChat = listValue
	}

	// Process analysis
	if analysis, ok := result["analysis"].(map[string]interface{}); ok {
		analysisValue, diags := types.ObjectValue(map[string]attr.Type{
			"goal_achieved":     types.BoolType,
			"total_duration_ms": types.Int64Type,
			"number_of_turns":   types.Int64Type,
			"summary":           types.StringType,
		}, map[string]attr.Value{
			"goal_achieved":     types.BoolValue(analysis["goal_achieved"].(bool)),
			"total_duration_ms": types.Int64Value(int64(analysis["total_duration_ms"].(float64))),
			"number_of_turns":   types.Int64Value(int64(analysis["number_of_turns"].(float64))),
			"summary":           types.StringValue(analysis["summary"].(string)),
		})
		resp.Diagnostics.Append(diags...)
		plan.Analysis = analysisValue
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ConvAIConversationSimulatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// This is an action resource, no persistent state needed
}

func (r *ConvAIConversationSimulatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ConvAIConversationSimulatorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.SimulatedChat = types.ListNull(simulatedMessageObjectType)

	// Convert chat history
	chatHistory := make([]map[string]interface{}, len(plan.ChatHistory))
	for i, msg := range plan.ChatHistory {
		chatHistory[i] = map[string]interface{}{
			"role":    msg.Role.ValueString(),
			"content": msg.Content.ValueString(),
		}
	}

	// Convert agent configuration
	var agentConfig map[string]interface{}
	if !plan.AgentConfig.IsNull() {
		agentConfigElements := plan.AgentConfig.Elements()
		agentConfig = make(map[string]interface{})
		for key, value := range agentConfigElements {
			agentConfig[key] = value
		}
	}

	result, err := r.client.SimulateConversation(plan.AgentID.ValueString(), chatHistory, agentConfig)
	if err != nil {
		resp.Diagnostics.AddError("Error simulating conversation", err.Error())
		return
	}

	// Process simulated conversation
	if simulated, ok := result["simulated_conversation"].([]interface{}); ok {
		simulatedMessages := make([]SimulatedMessageModel, len(simulated))
		for i, msg := range simulated {
			m := msg.(map[string]interface{})
			simulatedMessages[i] = SimulatedMessageModel{
				Role:      types.StringValue(m["role"].(string)),
				Content:   types.StringValue(m["content"].(string)),
				Timestamp: types.StringValue(m["timestamp"].(string)),
			}
		}
		listValue, diags := types.ListValueFrom(ctx, simulatedMessageObjectType, simulatedMessages)
		resp.Diagnostics.Append(diags...)
		plan.SimulatedChat = listValue
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ConvAIConversationSimulatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This is an action resource, no deletion needed
}
