package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

var (
	_ resource.Resource              = &ConvAIAgentDuplicatorResource{}
	_ resource.ResourceWithConfigure = &ConvAIAgentDuplicatorResource{}
)

func NewConvAIAgentDuplicatorResource() resource.Resource {
	return &ConvAIAgentDuplicatorResource{}
}

type ConvAIAgentDuplicatorResource struct {
	client *client.Client
}

type ConvAIAgentDuplicatorResourceModel struct {
	SourceAgentID   types.String `tfsdk:"source_agent_id"`
	NewAgentName    types.String `tfsdk:"new_agent_name"`
	NewAgentID      types.String `tfsdk:"new_agent_id"`
	NewAgentNameGet types.String `tfsdk:"new_agent_name_get"`
	CreatedAt       types.String `tfsdk:"created_at"`
}

func (r *ConvAIAgentDuplicatorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_agent_duplicate"
}

func (r *ConvAIAgentDuplicatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for duplicating ElevenLabs ConvAI agents. Creates a copy of an existing agent with optional name override.",
		Attributes: map[string]schema.Attribute{
			"source_agent_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the source agent to duplicate.",
			},
			"new_agent_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The name for the new duplicated agent.",
			},
			"new_agent_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the newly created agent.",
			},
			"new_agent_name_get": schema.StringAttribute{
				Computed: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *ConvAIAgentDuplicatorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConvAIAgentDuplicatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConvAIAgentDuplicatorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	agent, err := r.client.DuplicateConvAIAgent(
		plan.SourceAgentID.ValueString(),
		plan.NewAgentName.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Error duplicating agent", err.Error())
		return
	}

	plan.NewAgentID = types.StringValue(agent.AgentID)
	plan.NewAgentNameGet = types.StringValue(agent.Name)
	plan.CreatedAt = types.StringNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ConvAIAgentDuplicatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ConvAIAgentDuplicatorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	agent, err := r.client.GetConvAIAgent(state.NewAgentID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading duplicated agent", err.Error())
		return
	}

	state.NewAgentNameGet = types.StringValue(agent.Name)
	state.CreatedAt = types.StringNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ConvAIAgentDuplicatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ConvAIAgentDuplicatorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Duplicate again for update (creates another copy)
	agent, err := r.client.DuplicateConvAIAgent(
		plan.SourceAgentID.ValueString(),
		plan.NewAgentName.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Error duplicating agent", err.Error())
		return
	}

	plan.NewAgentID = types.StringValue(agent.AgentID)
	plan.CreatedAt = types.StringNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ConvAIAgentDuplicatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ConvAIAgentDuplicatorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConvAIAgent(state.NewAgentID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting duplicated agent", err.Error())
		return
	}
}
