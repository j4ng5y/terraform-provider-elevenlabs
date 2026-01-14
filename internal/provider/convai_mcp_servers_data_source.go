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

var (
	_ datasource.DataSource              = &ConvAIMCPServersDataSource{}
	_ datasource.DataSourceWithConfigure = &ConvAIMCPServersDataSource{}
)

var convaiMCPServerAttrTypes = map[string]attr.Type{
	"mcp_server_id":  types.StringType,
	"name":           types.StringType,
	"description":    types.StringType,
	"server_url":     types.StringType,
	"status":         types.StringType,
	"last_connected": types.StringType,
}

var convaiMCPServerObjectType = types.ObjectType{AttrTypes: convaiMCPServerAttrTypes}

func NewConvAIMCPServersDataSource() datasource.DataSource {
	return &ConvAIMCPServersDataSource{}
}

type ConvAIMCPServersDataSource struct {
	client *client.Client
}

type convaiMCPServersDataSourceModel struct {
	MCPServers types.List `tfsdk:"mcp_servers"`
}

func (d *ConvAIMCPServersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_mcp_servers"
}

func (d *ConvAIMCPServersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches ElevenLabs ConvAI MCP servers.",
		Attributes: map[string]schema.Attribute{
			"mcp_servers": schema.ListAttribute{
				Computed:            true,
				ElementType:         convaiMCPServerObjectType,
				MarkdownDescription: "List of ConvAI MCP servers configured in the workspace.",
			},
		},
	}
}

func (d *ConvAIMCPServersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *ConvAIMCPServersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data convaiMCPServersDataSourceModel

	servers, err := d.client.GetConvAIMCPServers()
	if err != nil {
		resp.Diagnostics.AddError("Error fetching ConvAI MCP servers", err.Error())
		return
	}

	serversList, diags := flattenConvAIMCPServers(ctx, servers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.MCPServers = serversList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenConvAIMCPServers(ctx context.Context, servers []models.ConvAIMCPServer) (types.List, diag.Diagnostics) {
	if len(servers) == 0 {
		return types.ListNull(convaiMCPServerObjectType), nil
	}

	values := make([]attr.Value, 0, len(servers))
	for _, server := range servers {
		obj, objDiags := types.ObjectValue(convaiMCPServerAttrTypes, map[string]attr.Value{
			"mcp_server_id":  types.StringValue(server.MCPServerID),
			"name":           types.StringValue(server.Name),
			"description":    types.StringNull(),
			"server_url":     types.StringValue(server.URL),
			"status":         types.StringNull(),
			"last_connected": types.StringNull(),
		})
		if objDiags.HasError() {
			return types.ListNull(convaiMCPServerObjectType), objDiags
		}

		values = append(values, obj)
	}

	return types.ListValue(convaiMCPServerObjectType, values)
}
