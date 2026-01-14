package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &PVCVoiceSamplesDataSource{}
)

// PVCVoiceSamplesDataSourceModel describes the data source data model.
type PVCVoiceSamplesDataSourceModel struct {
	VoiceID types.String                    `tfsdk:"voice_id"`
	Samples []PVCVoiceSampleDataSourceModel `tfsdk:"samples"`
}

func NewPVCVoiceSamplesDataSource() datasource.DataSource {
	return &PVCVoiceSamplesDataSource{}
}

// PVCVoiceSamplesDataSource defines the data source implementation.
type PVCVoiceSamplesDataSource struct {
	client *client.Client
}

func (d *PVCVoiceSamplesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pvc_voice_samples"
}

func (d *PVCVoiceSamplesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "PVC Voice Samples data source for ElevenLabs. Allows listing training samples for a Professional Voice Cloning voice.",
		Attributes: map[string]schema.Attribute{
			"voice_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the PVC voice to list samples for.",
			},
			"samples": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sample_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the sample.",
						},
						"file_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the sample file.",
						},
						"mime_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The MIME type of the sample.",
						},
						"size_bytes": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The size of the sample in bytes.",
						},
						"hash": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The hash of the sample content.",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The processing state of the sample.",
						},
						"transcription": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The transcription of the sample.",
						},
						"duration": schema.Float64Attribute{
							Computed:            true,
							MarkdownDescription: "The duration of the sample in seconds.",
						},
						"sample_rate": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The sample rate of the audio.",
						},
						"channels": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of audio channels.",
						},
					},
				},
			},
		},
	}
}

func (d *PVCVoiceSamplesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PVCVoiceSamplesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PVCVoiceSamplesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	samplesResp, err := d.client.ListPVCVoiceSamples(data.VoiceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list PVC voice samples, got error: %s", err))
		return
	}

	for _, sample := range samplesResp.Samples {
		sampleModel := PVCVoiceSampleDataSourceModel{
			SampleID:      types.StringValue(sample.SampleID),
			FileName:      types.StringValue(sample.FileName),
			MimeType:      types.StringValue(sample.MimeType),
			SizeBytes:     types.Int64Value(int64(sample.SizeBytes)),
			Hash:          types.StringValue(sample.Hash),
			State:         types.StringValue(sample.State),
			Transcription: types.StringValue(sample.Transcription),
			Duration:      types.Float64Value(sample.Duration),
			SampleRate:    types.Int64Value(int64(sample.SampleRate)),
			Channels:      types.Int64Value(int64(sample.Channels)),
		}

		data.Samples = append(data.Samples, sampleModel)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
