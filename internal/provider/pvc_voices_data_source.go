package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &PVCVoicesDataSource{}
)

// PVCVoicesDataSourceModel describes the data source data model.
type PVCVoicesDataSourceModel struct {
	Voices []PVCVoiceDataSourceModel `tfsdk:"voices"`
}

type PVCVoiceDataSourceModel struct {
	ID           types.String   `tfsdk:"id"`
	Name         types.String   `tfsdk:"name"`
	Language     types.String   `tfsdk:"language"`
	Description  types.String   `tfsdk:"description"`
	Labels       types.Map      `tfsdk:"labels"`
	State        types.String   `tfsdk:"state"`
	Verification types.String   `tfsdk:"verification"`
	Samples      types.List     `tfsdk:"samples"`
	Settings     *VoiceSettings `tfsdk:"settings"`
	CreatedAt    types.String   `tfsdk:"created_at"`
	UpdatedAt    types.String   `tfsdk:"updated_at"`
}

type PVCVoiceSampleDataSourceModel struct {
	SampleID      types.String  `tfsdk:"sample_id"`
	FileName      types.String  `tfsdk:"file_name"`
	MimeType      types.String  `tfsdk:"mime_type"`
	SizeBytes     types.Int64   `tfsdk:"size_bytes"`
	Hash          types.String  `tfsdk:"hash"`
	State         types.String  `tfsdk:"state"`
	Transcription types.String  `tfsdk:"transcription"`
	Duration      types.Float64 `tfsdk:"duration"`
	SampleRate    types.Int64   `tfsdk:"sample_rate"`
	Channels      types.Int64   `tfsdk:"channels"`
}

func NewPVCVoicesDataSource() datasource.DataSource {
	return &PVCVoicesDataSource{}
}

// PVCVoicesDataSource defines the data source implementation.
type PVCVoicesDataSource struct {
	client *client.Client
}

func (d *PVCVoicesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pvc_voices"
}

func (d *PVCVoicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "PVC Voices data source for ElevenLabs. Allows listing Professional Voice Cloning voices.",
		Attributes: map[string]schema.Attribute{
			"voices": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the PVC voice.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the PVC voice.",
						},
						"language": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The language code for the PVC voice.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A description of the PVC voice.",
						},
						"labels": schema.MapAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "Labels associated with the PVC voice.",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The current training state of the PVC voice.",
						},
						"verification": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The verification status of the PVC voice.",
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
						"settings": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"stability": schema.Float64Attribute{
									Computed: true,
								},
								"similarity_boost": schema.Float64Attribute{
									Computed: true,
								},
								"style": schema.Float64Attribute{
									Computed: true,
								},
								"use_speaker_boost": schema.BoolAttribute{
									Computed: true,
								},
							},
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The creation timestamp.",
						},
						"updated_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The last update timestamp.",
						},
					},
				},
			},
		},
	}
}

func (d *PVCVoicesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PVCVoicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PVCVoicesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	voicesResp, err := d.client.ListPVCVoices()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list PVC voices, got error: %s", err))
		return
	}

	for _, voice := range voicesResp.Voices {
		voiceModel := PVCVoiceDataSourceModel{
			ID:           types.StringValue(voice.VoiceID),
			Name:         types.StringValue(voice.Name),
			Language:     types.StringValue(voice.Language),
			Description:  types.StringValue(voice.Description),
			State:        types.StringValue(voice.State),
			Verification: types.StringValue(voice.Verification),
			CreatedAt:    types.StringValue(voice.CreatedAt),
			UpdatedAt:    types.StringValue(voice.UpdatedAt),
		}

		if voice.Labels != nil {
			labels, diags := types.MapValueFrom(ctx, types.StringType, voice.Labels)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			voiceModel.Labels = labels
		}

		if voice.Samples != nil {
			samples := make([]PVCVoiceSampleDataSourceModel, 0, len(voice.Samples))
			for _, sample := range voice.Samples {
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
				samples = append(samples, sampleModel)
			}

			samplesList, diags := types.ListValueFrom(ctx, types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"sample_id":     types.StringType,
					"file_name":     types.StringType,
					"mime_type":     types.StringType,
					"size_bytes":    types.Int64Type,
					"hash":          types.StringType,
					"state":         types.StringType,
					"transcription": types.StringType,
					"duration":      types.Float64Type,
					"sample_rate":   types.Int64Type,
					"channels":      types.Int64Type,
				},
			}, samples)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			voiceModel.Samples = samplesList
		}

		if voice.Settings != nil {
			voiceModel.Settings = &VoiceSettings{
				Stability:       types.Float64Value(voice.Settings.Stability),
				SimilarityBoost: types.Float64Value(voice.Settings.SimilarityBoost),
				Style:           types.Float64Value(voice.Settings.Style),
				UseSpeakerBoost: types.BoolValue(voice.Settings.UseSpeakerBoost),
			}
		}

		data.Voices = append(data.Voices, voiceModel)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
