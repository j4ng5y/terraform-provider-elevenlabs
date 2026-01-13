package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

var _ datasource.DataSource = &ConvAIAgentsDataSource{}
var _ datasource.DataSourceWithConfigure = &ConvAIAgentsDataSource{}

func NewConvAIAgentsDataSource() datasource.DataSource {
	return &ConvAIAgentsDataSource{}
}

type ConvAIAgentsDataSource struct {
	client *client.Client
}

type ConvAIAgentsDataSourceModel struct {
	Agents []AgentModel `tfsdk:"agents"`
}

type AgentModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *ConvAIAgentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_agents"
}

func (d *ConvAIAgentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for ElevenLabs Conversational AI agents. Allows listing all available agents.",
		Attributes: map[string]schema.Attribute{
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
					},
				},
			},
		},
	}
}

func (d *ConvAIAgentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ConvAIAgentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConvAIAgentsDataSourceModel

	agents, err := d.client.GetConvAIAgents()
	if err != nil {
		resp.Diagnostics.AddError("Error reading ConvAI agents", err.Error())
		return
	}

	for _, agent := range agents {
		data.Agents = append(data.Agents, AgentModel{
			ID:   types.StringValue(agent.AgentID),
			Name: types.StringValue(agent.Name),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
