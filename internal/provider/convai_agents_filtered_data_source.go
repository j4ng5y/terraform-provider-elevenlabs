package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

var (
	_ datasource.DataSource              = &ConvAIAgentsFilteredDataSource{}
	_ datasource.DataSourceWithConfigure = &ConvAIAgentsFilteredDataSource{}
)

func NewConvAIAgentsFilteredDataSource() datasource.DataSource {
	return &ConvAIAgentsFilteredDataSource{}
}

type ConvAIAgentsFilteredDataSource struct {
	client *client.Client
}

type ConvAIAgentsFilteredDataSourceModel struct {
	PageSize      types.Int64        `tfsdk:"page_size"`
	Search        types.String       `tfsdk:"search"`
	Archived      types.Bool         `tfsdk:"archived"`
	ShowOnlyOwned types.Bool         `tfsdk:"show_only_owned_agents"`
	Agents        []AgentDetailModel `tfsdk:"agents"`
}

type AgentDetailModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func (d *ConvAIAgentsFilteredDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_agents_filtered"
}

func (d *ConvAIAgentsFilteredDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for ElevenLabs Conversational AI agents with advanced filtering. Supports pagination, search, and status filtering.",
		Attributes: map[string]schema.Attribute{
			"page_size": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Maximum number of agents to return per page.",
			},
			"search": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Search term to filter agents by name.",
			},
			"archived": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to include archived agents.",
			},
			"show_only_owned_agents": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to only show agents owned by the current user.",
			},
			"agents": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *ConvAIAgentsFilteredDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *ConvAIAgentsFilteredDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConvAIAgentsFilteredDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	agents, err := d.client.GetConvAIAgentsFiltered(
		int(data.PageSize.ValueInt64()),
		data.Search.ValueString(),
		data.Archived.ValueBool(),
		data.ShowOnlyOwned.ValueBool(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Error getting filtered agents", err.Error())
		return
	}

	data.Agents = make([]AgentDetailModel, len(agents))
	for i, agent := range agents {
		data.Agents[i] = AgentDetailModel{
			ID:        types.StringValue(agent["id"].(string)),
			Name:      types.StringValue(agent["name"].(string)),
			CreatedAt: types.StringValue(agent["created_at"].(string)),
			UpdatedAt: types.StringValue(agent["updated_at"].(string)),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
