package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/models"
)

var (
	_ resource.Resource                = &ConvAIMCPToolConfigResource{}
	_ resource.ResourceWithConfigure   = &ConvAIMCPToolConfigResource{}
	_ resource.ResourceWithImportState = &ConvAIMCPToolConfigResource{}
)

func NewConvAIMCPToolConfigResource() resource.Resource {
	return &ConvAIMCPToolConfigResource{}
}

type ConvAIMCPToolConfigResource struct {
	client *client.Client
}

type ConvAIMCPToolConfigAssignmentModel struct {
	Source          types.String `tfsdk:"source"`
	DynamicVariable types.String `tfsdk:"dynamic_variable"`
	ValuePath       types.String `tfsdk:"value_path"`
}

type ConvAIMCPToolConfigResourceModel struct {
	ID                    types.String                         `tfsdk:"id"`
	MCPServerID           types.String                         `tfsdk:"mcp_server_id"`
	ToolName              types.String                         `tfsdk:"tool_name"`
	ForcePreToolSpeech    types.Bool                           `tfsdk:"force_pre_tool_speech"`
	DisableInterruptions  types.Bool                           `tfsdk:"disable_interruptions"`
	ToolCallSound         types.String                         `tfsdk:"tool_call_sound"`
	ToolCallSoundBehavior types.String                         `tfsdk:"tool_call_sound_behavior"`
	ExecutionMode         types.String                         `tfsdk:"execution_mode"`
	Assignments           []ConvAIMCPToolConfigAssignmentModel `tfsdk:"assignments"`
}

func (r *ConvAIMCPToolConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_mcp_tool_config"
}

func (r *ConvAIMCPToolConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Conversational AI MCP tool configuration override resource for ElevenLabs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mcp_server_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tool_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"force_pre_tool_speech": schema.BoolAttribute{
				Optional: true,
			},
			"disable_interruptions": schema.BoolAttribute{
				Optional: true,
			},
			"tool_call_sound": schema.StringAttribute{
				Optional: true,
			},
			"tool_call_sound_behavior": schema.StringAttribute{
				Optional: true,
			},
			"execution_mode": schema.StringAttribute{
				Optional: true,
			},
			"assignments": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"source": schema.StringAttribute{
							Optional: true,
						},
						"dynamic_variable": schema.StringAttribute{
							Required: true,
						},
						"value_path": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}

func (r *ConvAIMCPToolConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConvAIMCPToolConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConvAIMCPToolConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.MCPToolConfigOverrideCreateRequest{
		ToolName:              data.ToolName.ValueString(),
		ForcePreToolSpeech:    boolPointerFromValue(data.ForcePreToolSpeech),
		DisableInterruptions:  boolPointerFromValue(data.DisableInterruptions),
		ToolCallSound:         stringPointerFromValue(data.ToolCallSound),
		ToolCallSoundBehavior: stringPointerFromValue(data.ToolCallSoundBehavior),
		ExecutionMode:         stringPointerFromValue(data.ExecutionMode),
		Assignments:           expandMCPAssignments(data.Assignments),
	}

	err := r.client.CreateConvAIMCPToolConfigOverride(data.MCPServerID.ValueString(), addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ConvAI MCP tool config override", err.Error())
		return
	}

	config, err := r.client.GetConvAIMCPToolConfigOverride(data.MCPServerID.ValueString(), data.ToolName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ConvAI MCP tool config override", err.Error())
		return
	}

	applyMCPToolConfigToState(&data, config)
	data.ID = types.StringValue(formatMCPToolConfigID(data.MCPServerID.ValueString(), data.ToolName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIMCPToolConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConvAIMCPToolConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.client.GetConvAIMCPToolConfigOverride(data.MCPServerID.ValueString(), data.ToolName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ConvAI MCP tool config override", err.Error())
		return
	}

	applyMCPToolConfigToState(&data, config)
	data.ID = types.StringValue(formatMCPToolConfigID(data.MCPServerID.ValueString(), data.ToolName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIMCPToolConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConvAIMCPToolConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &models.MCPToolConfigOverrideUpdateRequest{
		ForcePreToolSpeech:    boolPointerFromValue(data.ForcePreToolSpeech),
		DisableInterruptions:  boolPointerFromValue(data.DisableInterruptions),
		ToolCallSound:         stringPointerFromValue(data.ToolCallSound),
		ToolCallSoundBehavior: stringPointerFromValue(data.ToolCallSoundBehavior),
		ExecutionMode:         stringPointerFromValue(data.ExecutionMode),
		Assignments:           expandMCPAssignments(data.Assignments),
	}

	err := r.client.UpdateConvAIMCPToolConfigOverride(data.MCPServerID.ValueString(), data.ToolName.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating ConvAI MCP tool config override", err.Error())
		return
	}

	config, err := r.client.GetConvAIMCPToolConfigOverride(data.MCPServerID.ValueString(), data.ToolName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ConvAI MCP tool config override", err.Error())
		return
	}

	applyMCPToolConfigToState(&data, config)
	data.ID = types.StringValue(formatMCPToolConfigID(data.MCPServerID.ValueString(), data.ToolName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIMCPToolConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConvAIMCPToolConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConvAIMCPToolConfigOverride(data.MCPServerID.ValueString(), data.ToolName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting ConvAI MCP tool config override", err.Error())
		return
	}
}

func (r *ConvAIMCPToolConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	mcpServerID, toolName, err := parseMCPToolConfigID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("mcp_server_id"), mcpServerID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tool_name"), toolName)...)
}

func formatMCPToolConfigID(mcpServerID, toolName string) string {
	return fmt.Sprintf("%s:%s", mcpServerID, toolName)
}

func parseMCPToolConfigID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("expected ID in format mcp_server_id:tool_name")
	}
	return parts[0], parts[1], nil
}

func applyMCPToolConfigToState(data *ConvAIMCPToolConfigResourceModel, config *models.MCPToolConfigOverride) {
	data.ToolName = types.StringValue(config.ToolName)
	data.ForcePreToolSpeech = boolValueOrNull(config.ForcePreToolSpeech)
	data.DisableInterruptions = boolValueOrNull(config.DisableInterruptions)
	data.ToolCallSound = stringValueOrNull(config.ToolCallSound)
	data.ToolCallSoundBehavior = stringValueOrNull(config.ToolCallSoundBehavior)
	data.ExecutionMode = stringValueOrNull(config.ExecutionMode)
	data.Assignments = flattenMCPAssignments(config.Assignments)
}

func expandMCPAssignments(assignments []ConvAIMCPToolConfigAssignmentModel) []models.DynamicVariableAssignment {
	if len(assignments) == 0 {
		return nil
	}

	result := make([]models.DynamicVariableAssignment, 0, len(assignments))
	for _, assignment := range assignments {
		item := models.DynamicVariableAssignment{
			Source:          assignment.Source.ValueString(),
			DynamicVariable: assignment.DynamicVariable.ValueString(),
			ValuePath:       assignment.ValuePath.ValueString(),
		}
		if assignment.Source.IsNull() || assignment.Source.IsUnknown() {
			item.Source = ""
		}
		result = append(result, item)
	}

	return result
}

func flattenMCPAssignments(assignments []models.DynamicVariableAssignment) []ConvAIMCPToolConfigAssignmentModel {
	if len(assignments) == 0 {
		return nil
	}

	result := make([]ConvAIMCPToolConfigAssignmentModel, 0, len(assignments))
	for _, assignment := range assignments {
		item := ConvAIMCPToolConfigAssignmentModel{
			DynamicVariable: types.StringValue(assignment.DynamicVariable),
			ValuePath:       types.StringValue(assignment.ValuePath),
		}
		if assignment.Source == "" {
			item.Source = types.StringNull()
		} else {
			item.Source = types.StringValue(assignment.Source)
		}
		result = append(result, item)
	}

	return result
}

func boolPointerFromValue(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	v := value.ValueBool()
	return &v
}

func stringPointerFromValue(value types.String) *string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	v := value.ValueString()
	return &v
}

func boolValueOrNull(value *bool) types.Bool {
	if value == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*value)
}

func stringValueOrNull(value *string) types.String {
	if value == nil {
		return types.StringNull()
	}
	return types.StringValue(*value)
}
