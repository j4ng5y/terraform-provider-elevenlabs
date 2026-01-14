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
	_ resource.Resource                = &ServiceAccountKeyResource{}
	_ resource.ResourceWithConfigure   = &ServiceAccountKeyResource{}
	_ resource.ResourceWithImportState = &ServiceAccountKeyResource{}
)

func NewServiceAccountKeyResource() resource.Resource {
	return &ServiceAccountKeyResource{}
}

type ServiceAccountKeyResource struct {
	client *client.Client
}

type ServiceAccountKeyResourceModel struct {
	ID             types.String `tfsdk:"id"`
	UserID         types.String `tfsdk:"user_id"`
	Name           types.String `tfsdk:"name"`
	XiApiKey       types.String `tfsdk:"api_key"`
	Permissions    types.List   `tfsdk:"permissions"`
	CharacterLimit types.Int64  `tfsdk:"character_limit"`
}

func (r *ServiceAccountKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_account_key"
}

func (r *ServiceAccountKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Service Account Key resource for ElevenLabs. Allows managing API keys for service accounts.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The ID of the service account user.",
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"api_key": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"permissions": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"character_limit": schema.Int64Attribute{
				Optional: true,
			},
		},
	}
}

func (r *ServiceAccountKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ServiceAccountKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServiceAccountKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var permissions []string
	resp.Diagnostics.Append(data.Permissions.ElementsAs(ctx, &permissions, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.CreateServiceAccountKeyRequest{
		Name: data.Name.ValueString(),
	}

	key, err := r.client.CreateServiceAccountKey(data.UserID.ValueString(), addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating service account key", err.Error())
		return
	}

	data.ID = types.StringValue(key.KeyID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceAccountKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// API supports listing keys. We could find the specific key by ID in the list.
}

func (r *ServiceAccountKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning("Update limited", "Updating service account keys might require replacement.")
}

func (r *ServiceAccountKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ServiceAccountKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteServiceAccountKey(data.UserID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting service account key", err.Error())
		return
	}
}

func (r *ServiceAccountKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import requires both userID and keyID, usually handled via a slash-separated string
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
