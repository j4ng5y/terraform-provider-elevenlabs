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
	_ resource.Resource                = &ConvAIPhoneNumberResource{}
	_ resource.ResourceWithConfigure   = &ConvAIPhoneNumberResource{}
	_ resource.ResourceWithImportState = &ConvAIPhoneNumberResource{}
)

func NewConvAIPhoneNumberResource() resource.Resource {
	return &ConvAIPhoneNumberResource{}
}

type ConvAIPhoneNumberResource struct {
	client *client.Client
}

type ConvAIPhoneNumberResourceModel struct {
	ID          types.String `tfsdk:"id"`
	PhoneNumber types.String `tfsdk:"phone_number"`
	Provider    types.String `tfsdk:"telephony_provider"`
	Label       types.String `tfsdk:"label"`
}

func (r *ConvAIPhoneNumberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_phone_number"
}

func (r *ConvAIPhoneNumberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Conversational AI Phone Number resource for ElevenLabs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"phone_number": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"telephony_provider": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The provider of the phone number (e.g., `twilio`).",
			},
			"label": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *ConvAIPhoneNumberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConvAIPhoneNumberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConvAIPhoneNumberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.ImportPhoneNumberRequest{
		PhoneNumber: data.PhoneNumber.ValueString(),
		Provider:    data.Provider.ValueString(),
		Label:       data.Label.ValueString(),
	}

	phone, err := r.client.ImportConvAIPhoneNumber(addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error importing ConvAI phone number", err.Error())
		return
	}

	data.ID = types.StringValue(phone.PhoneNumberID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIPhoneNumberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// API Read logic
}

func (r *ConvAIPhoneNumberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConvAIPhoneNumberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &models.ImportPhoneNumberRequest{
		PhoneNumber: data.PhoneNumber.ValueString(),
		Provider:    data.Provider.ValueString(),
		Label:       data.Label.ValueString(),
	}

	err := r.client.UpdateConvAIPhoneNumber(data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating ConvAI phone number", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIPhoneNumberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConvAIPhoneNumberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConvAIPhoneNumber(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting ConvAI phone number", err.Error())
		return
	}
}

func (r *ConvAIPhoneNumberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
