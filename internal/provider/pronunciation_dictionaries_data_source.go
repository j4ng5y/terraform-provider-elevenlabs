package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

var _ datasource.DataSource = &PronunciationDictionariesDataSource{}
var _ datasource.DataSourceWithConfigure = &PronunciationDictionariesDataSource{}

func NewPronunciationDictionariesDataSource() datasource.DataSource {
	return &PronunciationDictionariesDataSource{}
}

type PronunciationDictionariesDataSource struct {
	client *client.Client
}

type PronunciationDictionariesDataSourceModel struct {
	Dictionaries []DictionaryModel `tfsdk:"dictionaries"`
}

type DictionaryModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *PronunciationDictionariesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pronunciation_dictionaries"
}

func (d *PronunciationDictionariesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for ElevenLabs pronunciation dictionaries. Allows listing all available dictionaries.",
		Attributes: map[string]schema.Attribute{
			"dictionaries": schema.ListNestedAttribute{
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

func (d *PronunciationDictionariesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PronunciationDictionariesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PronunciationDictionariesDataSourceModel

	dicts, err := d.client.GetPronunciationDictionaries()
	if err != nil {
		resp.Diagnostics.AddError("Error reading pronunciation dictionaries", err.Error())
		return
	}

	for _, dict := range dicts {
		data.Dictionaries = append(data.Dictionaries, DictionaryModel{
			ID:   types.StringValue(dict.ID),
			Name: types.StringValue(dict.Name),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
