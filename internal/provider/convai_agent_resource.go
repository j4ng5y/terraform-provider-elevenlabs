package provider

import (
	"context"
	"fmt"

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
	_ resource.Resource                = &ConvAIAgentResource{}
	_ resource.ResourceWithConfigure   = &ConvAIAgentResource{}
	_ resource.ResourceWithImportState = &ConvAIAgentResource{}
)

func NewConvAIAgentResource() resource.Resource {
	return &ConvAIAgentResource{}
}

type ConvAIAgentResource struct {
	client *client.Client
}

type ConvAIAgentResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Prompt       types.String `tfsdk:"prompt"`
	FirstMessage types.String `tfsdk:"first_message"`
	Language     types.String `tfsdk:"language"`
	ModelID      types.String `tfsdk:"model_id"`
}

func (r *ConvAIAgentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_agent"
}

func (r *ConvAIAgentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Conversational AI Agent resource for ElevenLabs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"prompt": schema.StringAttribute{
				Required: true,
			},
			"first_message": schema.StringAttribute{
				Optional: true,
			},
			"language": schema.StringAttribute{
				Optional: true,
			},
			"model_id": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *ConvAIAgentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConvAIAgentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConvAIAgentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.CreateConvAIAgentRequest{
		Name: data.Name.ValueString(),
		Config: &models.ConvAIAgentConfig{
			Prompt:       data.Prompt.ValueString(),
			FirstMessage: data.FirstMessage.ValueString(),
			Language:     data.Language.ValueString(),
			ModelID:      data.ModelID.ValueString(),
		},
	}

	agent, err := r.client.CreateConvAIAgent(addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ConvAI agent", err.Error())
		return
	}

	data.ID = types.StringValue(agent.AgentID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIAgentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConvAIAgentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	agent, err := r.client.GetConvAIAgent(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ConvAI agent", err.Error())
		return
	}

	data.Name = types.StringValue(agent.Name)
	if agent.Config != nil {
		data.Prompt = types.StringValue(agent.Config.Prompt)
		data.FirstMessage = types.StringValue(agent.Config.FirstMessage)
		data.Language = types.StringValue(agent.Config.Language)
		data.ModelID = types.StringValue(agent.Config.ModelID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIAgentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConvAIAgentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &models.CreateConvAIAgentRequest{
		Name: data.Name.ValueString(),
		Config: &models.ConvAIAgentConfig{
			Prompt:       data.Prompt.ValueString(),
			FirstMessage: data.FirstMessage.ValueString(),
			Language:     data.Language.ValueString(),
			ModelID:      data.ModelID.ValueString(),
		},
	}

	err := r.client.UpdateConvAIAgent(data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating ConvAI agent", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIAgentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConvAIAgentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConvAIAgent(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting ConvAI agent", err.Error())
		return
	}
}

func (r *ConvAIAgentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
