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
)

func NewConvAIConversationResource() resource.Resource {
	return &ConvAIConversationResource{}
}

type ConvAIConversationResource struct {
	client *client.Client
}

type ConvAIConversationResourceModel struct {
	ConversationID types.String `tfsdk:"conversation_id"`
	AgentID        types.String `tfsdk:"agent_id"`
	Name           types.String `tfsdk:"name"`
	CreatedAt      types.String `tfsdk:"created_at"`
}

func (r *ConvAIConversationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_conversation"
}

func (r *ConvAIConversationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Conversational AI Conversation resource for ElevenLabs. Allows importing/managing conversations.",
		Attributes: map[string]schema.Attribute{
			"conversation_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Unique identifier of the conversation.",
			},
			"agent_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the agent this conversation belongs to.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Name of the conversation.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Creation timestamp of the conversation.",
			},
		},
	}
}

func (r *ConvAIConversationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *ConvAIConversationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConvAIConversationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	conversation, err := r.client.GetConvAIConversation(data.ConversationID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading conversation", err.Error())
		return
	}

	if conversation != nil {
		if (*conversation)["agent_id"] != data.AgentID.ValueString() {
			resp.Diagnostics.AddError("Agent ID mismatch", "The conversation belongs to a different agent")
			return
		}
		data.Name = types.StringValue((*conversation)["name"].(string))
		data.CreatedAt = types.StringValue((*conversation)["created_at"].(string))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIConversationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConvAIConversationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	conversation, err := r.client.GetConvAIConversation(data.ConversationID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading conversation", err.Error())
		return
	}

	if conversation == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.AgentID = types.StringValue((*conversation)["agent_id"].(string))
	data.Name = types.StringValue((*conversation)["name"].(string))
	data.CreatedAt = types.StringValue((*conversation)["created_at"].(string))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIConversationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *ConvAIConversationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConvAIConversationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConvAIConversation(data.ConversationID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting conversation", err.Error())
		return
	}
}

func (r *ConvAIConversationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("conversation_id"), req, resp)
}

