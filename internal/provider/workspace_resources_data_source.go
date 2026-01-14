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

var workspaceResourceAttrTypes = map[string]attr.Type{
	"resource_id":   types.StringType,
	"resource_type": types.StringType,
	"name":          types.StringType,
	"owner_id":      types.StringType,
	"shared":        types.BoolType,
}

var workspaceResourceObjectType = types.ObjectType{AttrTypes: workspaceResourceAttrTypes}

func NewWorkspaceResourcesDataSource() datasource.DataSource {
	return &WorkspaceResourcesDataSource{}
}

type WorkspaceResourcesDataSource struct {
	client *client.Client
}

type workspaceResourcesDataSourceModel struct {
	Resources types.List `tfsdk:"resources"`
}

func (d *WorkspaceResourcesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_resources"
}

func (d *WorkspaceResourcesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches workspace resources.",
		Attributes: map[string]schema.Attribute{
			"resources": schema.ListAttribute{
				Computed:            true,
				ElementType:         workspaceResourceObjectType,
				MarkdownDescription: "List of workspace resources.",
			},
		},
	}
}

func (d *WorkspaceResourcesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkspaceResourcesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data workspaceResourcesDataSourceModel

	resources, err := d.client.GetWorkspaceResources()
	if err != nil {
		resp.Diagnostics.AddError("Error fetching workspace resources", err.Error())
		return
	}

	resourcesList, diags := flattenWorkspaceResources(ctx, resources)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Resources = resourcesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenWorkspaceResources(ctx context.Context, resources []models.WorkspaceResource) (types.List, diag.Diagnostics) {
	if len(resources) == 0 {
		return types.ListNull(workspaceResourceObjectType), nil
	}

	values := make([]attr.Value, 0, len(resources))
	for _, resource := range resources {
		obj, objDiags := types.ObjectValue(workspaceResourceAttrTypes, map[string]attr.Value{
			"resource_id":   types.StringValue(resource.ResourceID),
			"resource_type": types.StringValue(resource.ResourceType),
			"name":          types.StringValue(resource.Name),
			"owner_id":      types.StringValue(resource.OwnerID),
			"shared":        types.BoolValue(resource.Shared),
		})
		if objDiags.HasError() {
			return types.ListNull(workspaceResourceObjectType), objDiags
		}

		values = append(values, obj)
	}

	return types.ListValue(workspaceResourceObjectType, values)
}
