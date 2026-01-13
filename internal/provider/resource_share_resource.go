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

var (
	_ resource.Resource                = &ResourceShareResource{}
	_ resource.ResourceWithConfigure   = &ResourceShareResource{}
	_ resource.ResourceWithImportState = &ResourceShareResource{}
)

func NewResourceShareResource() resource.Resource {
	return &ResourceShareResource{}
}

type ResourceShareResource struct {
	client *client.Client
}

type ResourceShareResourceModel struct {
	ID           types.String `tfsdk:"id"`
	ResourceID   types.String `tfsdk:"resource_id"`
	ResourceType types.String `tfsdk:"resource_type"`
	Email        types.String `tfsdk:"email"`
	Role         types.String `tfsdk:"role"`
}

func (r *ResourceShareResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_share"
}

func (r *ResourceShareResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource Share resource for ElevenLabs. Allows sharing voices, agents, etc. with users.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Type of resource: `voice`, `agent`, etc.",
			},
			"email": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *ResourceShareResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourceShareResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ResourceShareResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.ShareResource(data.ResourceID.ValueString(), data.ResourceType.ValueString(), data.Email.ValueString(), data.Role.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error sharing resource", err.Error())
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%s:%s:%s", data.ResourceType.ValueString(), data.ResourceID.ValueString(), data.Email.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceShareResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// API logic to check current sharing state
}

func (r *ResourceShareResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Role update might be possible via re-share or specific endpoint
	var data ResourceShareResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.ShareResource(data.ResourceID.ValueString(), data.ResourceType.ValueString(), data.Email.ValueString(), data.Role.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error updating resource share", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceShareResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ResourceShareResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UnshareResource(data.ResourceID.ValueString(), data.ResourceType.ValueString(), data.Email.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error unsharing resource", err.Error())
		return
	}
}

func (r *ResourceShareResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
