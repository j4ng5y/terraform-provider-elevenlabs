package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

var _ datasource.DataSource = &VoicesDataSource{}
var _ datasource.DataSourceWithConfigure = &VoicesDataSource{}

func NewVoicesDataSource() datasource.DataSource {
	return &VoicesDataSource{}
}

type VoicesDataSource struct {
	client *client.Client
}

type VoicesDataSourceModel struct {
	Voices []VoiceModel `tfsdk:"voices"`
}

type VoiceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Category types.String `tfsdk:"category"`
}

func (d *VoicesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_voices"
}

func (d *VoicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for ElevenLabs voices. Allows listing all available voices.",
		Attributes: map[string]schema.Attribute{
			"voices": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"category": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *VoicesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VoicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VoicesDataSourceModel

	voices, err := d.client.GetVoices()
	if err != nil {
		resp.Diagnostics.AddError("Error reading voices", err.Error())
		return
	}

	for _, v := range voices {
		data.Voices = append(data.Voices, VoiceModel{
			ID:       types.StringValue(v.VoiceID),
			Name:     types.StringValue(v.Name),
			Category: types.StringValue(v.Category),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
