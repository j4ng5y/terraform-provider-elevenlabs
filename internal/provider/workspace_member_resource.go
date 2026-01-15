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
	_ resource.Resource                = &WorkspaceMemberResource{}
	_ resource.ResourceWithConfigure   = &WorkspaceMemberResource{}
	_ resource.ResourceWithImportState = &WorkspaceMemberResource{}
)

func NewWorkspaceMemberResource() resource.Resource {
	return &WorkspaceMemberResource{}
}

type WorkspaceMemberResource struct {
	client *client.Client
}

type WorkspaceMemberResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Email               types.String `tfsdk:"email"`
	WorkspacePermission types.String `tfsdk:"workspace_permission"`
}

func (r *WorkspaceMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_member"
}

func (r *WorkspaceMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Workspace Member resource for ElevenLabs. Allows managing roles of existing workspace members.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The User ID of the member.",
			},
			"email": schema.StringAttribute{
				Computed: true,
			},
			"workspace_permission": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Permission level for the user. e.g., `workspace_member`, `workspace_admin`, `admin`.",
			},
		},
	}
}

func (r *WorkspaceMemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WorkspaceMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Create for a member usually means updating an existing user's role in the workspace.
	var data WorkspaceMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &models.UpdateWorkspaceMemberRequest{
		WorkspacePermission: data.WorkspacePermission.ValueString(),
	}

	err := r.client.UpdateWorkspaceMember(data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating workspace member", err.Error())
		return
	}

	data.Email = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// API Read logic
}

func (r *WorkspaceMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WorkspaceMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &models.UpdateWorkspaceMemberRequest{
		WorkspacePermission: data.WorkspacePermission.ValueString(),
	}

	err := r.client.UpdateWorkspaceMember(data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating workspace member", err.Error())
		return
	}

	data.Email = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Deleting a member from workspace might require a different endpoint or just setting permission to none.
	resp.Diagnostics.AddWarning("Delete not fully implemented", "Removing members from workspace is currently limited.")
}

func (r *WorkspaceMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
