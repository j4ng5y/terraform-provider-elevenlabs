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
	_ resource.Resource                = &ConvAISecretResource{}
	_ resource.ResourceWithConfigure   = &ConvAISecretResource{}
	_ resource.ResourceWithImportState = &ConvAISecretResource{}
)

func NewConvAISecretResource() resource.Resource {
	return &ConvAISecretResource{}
}

type ConvAISecretResource struct {
	client *client.Client
}

type ConvAISecretResourceModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func (r *ConvAISecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_secret"
}

func (r *ConvAISecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Conversational AI Secret resource for ElevenLabs.",
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
			"value": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func (r *ConvAISecretResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConvAISecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConvAISecretResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.CreateConvAISecretRequest{
		Name:  data.Name.ValueString(),
		Value: data.Value.ValueString(),
	}

	secret, err := r.client.CreateConvAISecret(addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ConvAI secret", err.Error())
		return
	}

	data.ID = types.StringValue(secret.SecretID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAISecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConvAISecretResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := r.client.GetConvAISecret(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ConvAI secret", err.Error())
		return
	}

	data.Name = types.StringValue(secret.Name)
	// Value is not returned by the API for security reasons

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAISecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConvAISecretResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &models.CreateConvAISecretRequest{
		Name:  data.Name.ValueString(),
		Value: data.Value.ValueString(),
	}

	err := r.client.UpdateConvAISecret(data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating ConvAI secret", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAISecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConvAISecretResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConvAISecret(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting ConvAI secret", err.Error())
		return
	}
}

func (r *ConvAISecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
