package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/models"
)

var workspaceInviteAttrTypes = map[string]attr.Type{
	"email":      types.StringType,
	"role":       types.StringType,
	"status":     types.StringType,
	"invited_at": types.StringType,
}

var workspaceInviteObjectType = types.ObjectType{AttrTypes: workspaceInviteAttrTypes}

func NewWorkspaceInvitesDataSource() datasource.DataSource {
	return &WorkspaceInvitesDataSource{}
}

type WorkspaceInvitesDataSource struct {
	client *client.Client
}

type workspaceInvitesDataSourceModel struct {
	Invites types.List `tfsdk:"invites"`
}

func (d *WorkspaceInvitesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_invites"
}

func (d *WorkspaceInvitesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches workspace invites.",
		Attributes: map[string]schema.Attribute{
			"invites": schema.ListAttribute{
				Computed:            true,
				ElementType:         workspaceInviteObjectType,
				MarkdownDescription: "List of workspace invites.",
			},
		},
	}
}

func (d *WorkspaceInvitesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = c
}

func (d *WorkspaceInvitesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data workspaceInvitesDataSourceModel

	invites, err := d.client.GetWorkspaceInvites()
	if err != nil {
		resp.Diagnostics.AddError("Error fetching workspace invites", err.Error())
		return
	}

	invitesList, diags := flattenWorkspaceInvites(ctx, invites)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Invites = invitesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenWorkspaceInvites(ctx context.Context, invites []models.WorkspaceInvite) (types.List, diag.Diagnostics) {
	if len(invites) == 0 {
		return types.ListNull(workspaceInviteObjectType), nil
	}

	values := make([]attr.Value, 0, len(invites))
	for _, invite := range invites {
		obj, objDiags := types.ObjectValue(workspaceInviteAttrTypes, map[string]attr.Value{
			"email":      types.StringValue(invite.Email),
			"role":       types.StringValue(invite.Role),
			"status":     types.StringNull(),
			"invited_at": types.StringNull(),
		})
		if objDiags.HasError() {
			return types.ListNull(workspaceInviteObjectType), objDiags
		}

		values = append(values, obj)
	}

	return types.ListValue(workspaceInviteObjectType, values)
}
