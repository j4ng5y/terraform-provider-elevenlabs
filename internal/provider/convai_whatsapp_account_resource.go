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
	_ resource.Resource                = &ConvAIWhatsAppAccountResource{}
	_ resource.ResourceWithConfigure   = &ConvAIWhatsAppAccountResource{}
	_ resource.ResourceWithImportState = &ConvAIWhatsAppAccountResource{}
)

func NewConvAIWhatsAppAccountResource() resource.Resource {
	return &ConvAIWhatsAppAccountResource{}
}

type ConvAIWhatsAppAccountResource struct {
	client *client.Client
}

type ConvAIWhatsAppAccountResourceModel struct {
	PhoneNumberID     types.String `tfsdk:"phone_number_id"`
	BusinessAccountID types.String `tfsdk:"business_account_id"`
	TokenCode         types.String `tfsdk:"token_code"`
	AssignedAgentID   types.String `tfsdk:"assigned_agent_id"`

	// Computed
	BusinessAccountName types.String `tfsdk:"business_account_name"`
	PhoneNumberName     types.String `tfsdk:"phone_number_name"`
	PhoneNumber         types.String `tfsdk:"phone_number"`
	AssignedAgentName   types.String `tfsdk:"assigned_agent_name"`
}

func (r *ConvAIWhatsAppAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_whatsapp_account"
}

func (r *ConvAIWhatsAppAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "ConvAI WhatsApp Account resource for ElevenLabs.",
		Attributes: map[string]schema.Attribute{
			"phone_number_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Unique identifier of the phone number to associate with this WhatsApp account.",
			},
			"business_account_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Business account ID from your WhatsApp Business API.",
			},
			"token_code": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Token code from WhatsApp Business API setup flow.",
			},
			"assigned_agent_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "ID of the ConvAI agent to assign to this WhatsApp account.",
			},
			"business_account_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the WhatsApp business account.",
			},
			"phone_number_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Display name of the phone number.",
			},
			"phone_number": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The actual phone number.",
			},
			"assigned_agent_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the assigned ConvAI agent.",
			},
		},
	}
}

func (r *ConvAIWhatsAppAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConvAIWhatsAppAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConvAIWhatsAppAccountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.ImportWhatsAppAccountRequest{
		BusinessAccountID: data.BusinessAccountID.ValueString(),
		PhoneNumberID:     data.PhoneNumberID.ValueString(),
		TokenCode:         data.TokenCode.ValueString(),
	}

	account, err := r.client.ImportConvAIWhatsAppAccount(addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error importing ConvAI WhatsApp account", err.Error())
		return
	}

	data.PhoneNumberID = types.StringValue(account.PhoneNumberID)
	data.BusinessAccountID = types.StringValue(account.BusinessAccountID)
	data.BusinessAccountName = types.StringValue(account.BusinessAccountName)
	data.PhoneNumberName = types.StringValue(account.PhoneNumberName)
	data.PhoneNumber = types.StringValue(account.PhoneNumber)
	data.AssignedAgentID = optionalStringValue(account.AssignedAgentID)
	data.AssignedAgentName = optionalStringValue(account.AssignedAgentName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIWhatsAppAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConvAIWhatsAppAccountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account, err := r.client.GetConvAIWhatsAppAccount(data.PhoneNumberID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ConvAI WhatsApp account", err.Error())
		return
	}

	data.PhoneNumberID = types.StringValue(account.PhoneNumberID)
	data.BusinessAccountID = types.StringValue(account.BusinessAccountID)
	data.BusinessAccountName = types.StringValue(account.BusinessAccountName)
	data.PhoneNumberName = types.StringValue(account.PhoneNumberName)
	data.PhoneNumber = types.StringValue(account.PhoneNumber)
	data.AssignedAgentID = optionalStringValue(account.AssignedAgentID)
	data.AssignedAgentName = optionalStringValue(account.AssignedAgentName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIWhatsAppAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConvAIWhatsAppAccountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &models.UpdateWhatsAppAccountRequest{}
	if !data.AssignedAgentID.IsNull() && !data.AssignedAgentID.IsUnknown() {
		agentID := data.AssignedAgentID.ValueString()
		updateReq.AssignedAgentID = &agentID
	} else {
		updateReq.AssignedAgentID = nil
	}

	account, err := r.client.UpdateConvAIWhatsAppAccount(data.PhoneNumberID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating ConvAI WhatsApp account", err.Error())
		return
	}

	data.PhoneNumberID = types.StringValue(account.PhoneNumberID)
	data.BusinessAccountID = types.StringValue(account.BusinessAccountID)
	data.BusinessAccountName = types.StringValue(account.BusinessAccountName)
	data.PhoneNumberName = types.StringValue(account.PhoneNumberName)
	data.PhoneNumber = types.StringValue(account.PhoneNumber)
	data.AssignedAgentID = optionalStringValue(account.AssignedAgentID)
	data.AssignedAgentName = optionalStringValue(account.AssignedAgentName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIWhatsAppAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConvAIWhatsAppAccountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConvAIWhatsAppAccount(data.PhoneNumberID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting ConvAI WhatsApp account", err.Error())
		return
	}
}

func (r *ConvAIWhatsAppAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("phone_number_id"), req, resp)
}
