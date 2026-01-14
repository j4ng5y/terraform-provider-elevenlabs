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

var workspaceGroupAttrTypes = map[string]attr.Type{
	"group_id": types.StringType,
	"name":     types.StringType,
}

var workspaceGroupObjectType = types.ObjectType{AttrTypes: workspaceGroupAttrTypes}

func NewWorkspaceGroupsDataSource() datasource.DataSource {
	return &WorkspaceGroupsDataSource{}
}

type WorkspaceGroupsDataSource struct {
	client *client.Client
}

type workspaceGroupsDataSourceModel struct {
	Groups types.List `tfsdk:"groups"`
}

func (d *WorkspaceGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_groups"
}

func (d *WorkspaceGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches workspace user groups.",
		Attributes: map[string]schema.Attribute{
			"groups": schema.ListAttribute{
				Computed:            true,
				ElementType:         workspaceGroupObjectType,
				MarkdownDescription: "List of workspace user groups.",
			},
		},
	}
}

func (d *WorkspaceGroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkspaceGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data workspaceGroupsDataSourceModel

	groups, err := d.client.SearchWorkspaceGroups("")
	if err != nil {
		resp.Diagnostics.AddError("Error fetching workspace groups", err.Error())
		return
	}

	groupsList, diags := flattenWorkspaceGroups(ctx, groups)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Groups = groupsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenWorkspaceGroups(ctx context.Context, groups []models.WorkspaceGroup) (types.List, diag.Diagnostics) {
	if len(groups) == 0 {
		return types.ListNull(workspaceGroupObjectType), nil
	}

	values := make([]attr.Value, 0, len(groups))
	for _, group := range groups {
		obj, objDiags := types.ObjectValue(workspaceGroupAttrTypes, map[string]attr.Value{
			"group_id": types.StringValue(group.GroupID),
			"name":     types.StringValue(group.Name),
		})
		if objDiags.HasError() {
			return types.ListNull(workspaceGroupObjectType), objDiags
		}

		values = append(values, obj)
	}

	return types.ListValue(workspaceGroupObjectType, values)
}
