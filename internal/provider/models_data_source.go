package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

var _ datasource.DataSource = &ModelsDataSource{}
var _ datasource.DataSourceWithConfigure = &ModelsDataSource{}

func NewModelsDataSource() datasource.DataSource {
	return &ModelsDataSource{}
}

type ModelsDataSource struct {
	client *client.Client
}

type ModelsDataSourceModel struct {
	Models []ModelModel `tfsdk:"models"`
}

type ModelModel struct {
	ModelID              types.String `tfsdk:"model_id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	CanDoTextToSpeech    types.Bool   `tfsdk:"can_do_text_to_speech"`
	CanDoVoiceConversion types.Bool   `tfsdk:"can_do_voice_conversion"`
}

func (d *ModelsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_models"
}

func (d *ModelsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for ElevenLabs models. Allows retrieving information about available models.",
		Attributes: map[string]schema.Attribute{
			"models": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"model_id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"can_do_text_to_speech": schema.BoolAttribute{
							Computed: true,
						},
						"can_do_voice_conversion": schema.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *ModelsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *ModelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ModelsDataSourceModel

	models, err := d.client.GetModels()
	if err != nil {
		resp.Diagnostics.AddError("Error reading models", err.Error())
		return
	}

	for _, m := range models {
		data.Models = append(data.Models, ModelModel{
			ModelID:              types.StringValue(m.ModelID),
			Name:                 types.StringValue(m.Name),
			Description:          types.StringValue(m.Description),
			CanDoTextToSpeech:    types.BoolValue(m.CanDoTextToSpeech),
			CanDoVoiceConversion: types.BoolValue(m.CanDoVoiceConversion),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
