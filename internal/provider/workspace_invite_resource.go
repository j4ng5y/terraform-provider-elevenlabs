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
	_ resource.Resource                = &WorkspaceInviteResource{}
	_ resource.ResourceWithConfigure   = &WorkspaceInviteResource{}
	_ resource.ResourceWithImportState = &WorkspaceInviteResource{}
)

func NewWorkspaceInviteResource() resource.Resource {
	return &WorkspaceInviteResource{}
}

type WorkspaceInviteResource struct {
	client *client.Client
}

type WorkspaceInviteResourceModel struct {
	Email types.String `tfsdk:"email"`
	Role  types.String `tfsdk:"role"`
}

func (r *WorkspaceInviteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_invite"
}

func (r *WorkspaceInviteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Workspace Invite resource for ElevenLabs. Allows inviting users to a workspace.",
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Permission level for the user. e.g., `workspace_member`, `workspace_admin`, `admin`.",
			},
		},
	}
}

func (r *WorkspaceInviteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WorkspaceInviteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WorkspaceInviteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.CreateWorkspaceInviteRequest{
		Email:               data.Email.ValueString(),
		WorkspacePermission: data.Role.ValueString(),
	}

	err := r.client.CreateWorkspaceInvite(addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating workspace invite", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceInviteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// The API doesn't have a single GET endpoint for an invite by email easily.
	// We'll rely on the state for now or just assume it exists if no error from list.
}

func (r *WorkspaceInviteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Updating an invite usually requires re-inviting or a different endpoint.
	resp.Diagnostics.AddWarning("Update limited", "Updating workspace invites might require replacement.")
}

func (r *WorkspaceInviteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WorkspaceInviteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteWorkspaceInvite(data.Email.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting workspace invite", err.Error())
		return
	}
}

func (r *WorkspaceInviteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("email"), req, resp)
}
