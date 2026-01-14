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
	_ datasource.DataSource              = &ConvAIToolsDataSource{}
	_ datasource.DataSourceWithConfigure = &ConvAIToolsDataSource{}
)

var convAIToolAttrTypes = map[string]attr.Type{
	"tool_id":             types.StringType,
	"name":                types.StringType,
	"description":         types.StringType,
	"parameters_json":     types.StringType,
	"dependent_agent_ids": types.ListType{ElemType: types.StringType},
	"access_info_json":    types.StringType,
	"tool_config_json":    types.StringType,
	"usage_stats_json":    types.StringType,
}

var convAIToolObjectType = types.ObjectType{AttrTypes: convAIToolAttrTypes}

func NewConvAIToolsDataSource() datasource.DataSource {
	return &ConvAIToolsDataSource{}
}

type ConvAIToolsDataSource struct {
	client *client.Client
}

type convaiToolsDataSourceModel struct {
	Tools types.List `tfsdk:"tools"`
}

func (d *ConvAIToolsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_tools"
}

func (d *ConvAIToolsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches ElevenLabs ConvAI tools.",
		Attributes: map[string]schema.Attribute{
			"tools": schema.ListAttribute{
				Computed:            true,
				ElementType:         convAIToolObjectType,
				MarkdownDescription: "List of ConvAI tools configured in the workspace.",
			},
		},
	}
}

func (d *ConvAIToolsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ConvAIToolsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data convaiToolsDataSourceModel

	tools, err := d.client.GetConvAITools()
	if err != nil {
		resp.Diagnostics.AddError("Error fetching ConvAI tools", err.Error())
		return
	}

	toolsList, diags := flattenConvAITools(ctx, tools)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Tools = toolsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenConvAITools(ctx context.Context, tools []models.ConvAITool) (types.List, diag.Diagnostics) {
	if len(tools) == 0 {
		return types.ListNull(convAIToolObjectType), nil
	}

	values := make([]attr.Value, 0, len(tools))
	for _, tool := range tools {
		agentIds, diags := stringsToListValue(ctx, tool.DependentAgentIDs)
		if diags.HasError() {
			return types.ListNull(convAIToolObjectType), diags
		}

		obj, objDiags := types.ObjectValue(convAIToolAttrTypes, map[string]attr.Value{
			"tool_id":             types.StringValue(tool.ToolID),
			"name":                types.StringValue(tool.Name),
			"description":         types.StringValue(tool.Description),
			"parameters_json":     jsonStringValue(tool.Parameters),
			"dependent_agent_ids": agentIds,
			"access_info_json":    jsonStringValue(tool.AccessInfo),
			"tool_config_json":    jsonStringValue(tool.ToolConfig),
			"usage_stats_json":    jsonStringValue(tool.UsageStats),
		})
		if objDiags.HasError() {
			return types.ListNull(convAIToolObjectType), objDiags
		}

		values = append(values, obj)
	}

	return types.ListValue(convAIToolObjectType, values)
}
