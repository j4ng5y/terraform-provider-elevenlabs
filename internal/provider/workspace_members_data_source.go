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

var workspaceMemberAttrTypes = map[string]attr.Type{
	"user_id":              types.StringType,
	"email":                types.StringType,
	"workspace_permission": types.StringType,
	"is_invited":           types.BoolType,
}

var workspaceMemberObjectType = types.ObjectType{AttrTypes: workspaceMemberAttrTypes}

func NewWorkspaceMembersDataSource() datasource.DataSource {
	return &WorkspaceMembersDataSource{}
}

type WorkspaceMembersDataSource struct {
	client *client.Client
}

type workspaceMembersDataSourceModel struct {
	Members types.List `tfsdk:"members"`
}

func (d *WorkspaceMembersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_members"
}

func (d *WorkspaceMembersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches workspace members.",
		Attributes: map[string]schema.Attribute{
			"members": schema.ListAttribute{
				Computed:            true,
				ElementType:         workspaceMemberObjectType,
				MarkdownDescription: "List of workspace members.",
			},
		},
	}
}

func (d *WorkspaceMembersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkspaceMembersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data workspaceMembersDataSourceModel

	members, err := d.client.GetWorkspaceMembers()
	if err != nil {
		resp.Diagnostics.AddError("Error fetching workspace members", err.Error())
		return
	}

	membersList, diags := flattenWorkspaceMembers(ctx, members)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Members = membersList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenWorkspaceMembers(ctx context.Context, members []models.WorkspaceMember) (types.List, diag.Diagnostics) {
	if len(members) == 0 {
		return types.ListNull(workspaceMemberObjectType), nil
	}

	values := make([]attr.Value, 0, len(members))
	for _, member := range members {
		obj, objDiags := types.ObjectValue(workspaceMemberAttrTypes, map[string]attr.Value{
			"user_id":              types.StringValue(member.UserID),
			"email":                types.StringValue(member.Email),
			"workspace_permission": types.StringValue(member.WorkspacePermission),
			"is_invited":           types.BoolValue(member.IsInvited),
		})
		if objDiags.HasError() {
			return types.ListNull(workspaceMemberObjectType), objDiags
		}

		values = append(values, obj)
	}

	return types.ListValue(workspaceMemberObjectType, values)
}
