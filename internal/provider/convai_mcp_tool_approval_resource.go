package provider

import (
	"context"
	"encoding/json"
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

const defaultMCPToolApprovalPolicy = "requires_approval"

var (
	_ resource.Resource                = &ConvAIMCPToolApprovalResource{}
	_ resource.ResourceWithConfigure   = &ConvAIMCPToolApprovalResource{}
	_ resource.ResourceWithImportState = &ConvAIMCPToolApprovalResource{}
)

func NewConvAIMCPToolApprovalResource() resource.Resource {
	return &ConvAIMCPToolApprovalResource{}
}

type ConvAIMCPToolApprovalResource struct {
	client *client.Client
}

type ConvAIMCPToolApprovalResourceModel struct {
	ID              types.String `tfsdk:"id"`
	MCPServerID     types.String `tfsdk:"mcp_server_id"`
	ToolName        types.String `tfsdk:"tool_name"`
	ToolDescription types.String `tfsdk:"tool_description"`
	InputSchema     types.String `tfsdk:"input_schema"`
	ApprovalPolicy  types.String `tfsdk:"approval_policy"`
}

func (r *ConvAIMCPToolApprovalResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_mcp_tool_approval"
}

func (r *ConvAIMCPToolApprovalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Conversational AI MCP tool approval resource for ElevenLabs.",
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
			"tool_description": schema.StringAttribute{
				Required: true,
			},
			"input_schema": schema.StringAttribute{
				Optional: true,
			},
			"approval_policy": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *ConvAIMCPToolApprovalResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConvAIMCPToolApprovalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConvAIMCPToolApprovalResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq, err := buildMCPToolApprovalRequest(&data)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing MCP tool approval input schema", err.Error())
		return
	}

	err = r.client.CreateConvAIMCPToolApproval(data.MCPServerID.ValueString(), addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ConvAI MCP tool approval", err.Error())
		return
	}

	data.ID = types.StringValue(formatMCPToolConfigID(data.MCPServerID.ValueString(), data.ToolName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIMCPToolApprovalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// No read endpoint available; keep state as-is.
}

func (r *ConvAIMCPToolApprovalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConvAIMCPToolApprovalResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConvAIMCPToolApproval(data.MCPServerID.ValueString(), data.ToolName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting ConvAI MCP tool approval", err.Error())
		return
	}

	addReq, err := buildMCPToolApprovalRequest(&data)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing MCP tool approval input schema", err.Error())
		return
	}

	err = r.client.CreateConvAIMCPToolApproval(data.MCPServerID.ValueString(), addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ConvAI MCP tool approval", err.Error())
		return
	}

	data.ID = types.StringValue(formatMCPToolConfigID(data.MCPServerID.ValueString(), data.ToolName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIMCPToolApprovalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConvAIMCPToolApprovalResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConvAIMCPToolApproval(data.MCPServerID.ValueString(), data.ToolName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting ConvAI MCP tool approval", err.Error())
		return
	}
}

func (r *ConvAIMCPToolApprovalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid import ID", "expected ID in format mcp_server_id:tool_name")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("mcp_server_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tool_name"), parts[1])...)
}

func buildMCPToolApprovalRequest(data *ConvAIMCPToolApprovalResourceModel) (*models.MCPToolAddApprovalRequest, error) {
	req := &models.MCPToolAddApprovalRequest{
		ToolName:        data.ToolName.ValueString(),
		ToolDescription: data.ToolDescription.ValueString(),
	}

	if !data.ApprovalPolicy.IsNull() && !data.ApprovalPolicy.IsUnknown() {
		req.ApprovalPolicy = data.ApprovalPolicy.ValueString()
	}
	if req.ApprovalPolicy == "" {
		req.ApprovalPolicy = defaultMCPToolApprovalPolicy
	}

	if data.InputSchema.IsNull() || data.InputSchema.IsUnknown() || data.InputSchema.ValueString() == "" {
		return req, nil
	}

	var schemaValue map[string]interface{}
	if err := json.Unmarshal([]byte(data.InputSchema.ValueString()), &schemaValue); err != nil {
		return nil, err
	}
	req.InputSchema = schemaValue

	return req, nil
}
