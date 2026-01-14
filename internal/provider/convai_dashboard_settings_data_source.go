package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

func NewConvAIDashboardSettingsDataSource() datasource.DataSource {
	return &ConvAIDashboardSettingsDataSource{}
}

type ConvAIDashboardSettingsDataSource struct {
	client *client.Client
}

type convaiDashboardSettingsDataSourceModel struct {
	AnalyticsEnabled       types.Bool `tfsdk:"analytics_enabled"`
	RecordingEnabled       types.Bool `tfsdk:"recording_enabled"`
	TranscriptionEnabled   types.Bool `tfsdk:"transcription_enabled"`
	LLMOptimizationEnabled types.Bool `tfsdk:"llm_optimization_enabled"`
}

func (d *ConvAIDashboardSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_dashboard_settings"
}

func (d *ConvAIDashboardSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches ElevenLabs ConvAI dashboard settings.",
		Attributes: map[string]schema.Attribute{
			"analytics_enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether analytics are enabled.",
			},
			"recording_enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether call recording is enabled.",
			},
			"transcription_enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether transcription is enabled.",
			},
			"llm_optimization_enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether LLM optimization is enabled.",
			},
		},
	}
}

func (d *ConvAIDashboardSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ConvAIDashboardSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data convaiDashboardSettingsDataSourceModel

	settings, err := d.client.GetConvAIDashboardSettings()
	if err != nil {
		resp.Diagnostics.AddError("Error fetching ConvAI dashboard settings", err.Error())
		return
	}

	data.AnalyticsEnabled = types.BoolValue(false)
	data.RecordingEnabled = types.BoolValue(false)
	data.TranscriptionEnabled = types.BoolValue(false)
	data.LLMOptimizationEnabled = types.BoolValue(false)

	if analytics, ok := settings["analytics_enabled"]; ok {
		if b, ok := analytics.(bool); ok {
			data.AnalyticsEnabled = types.BoolValue(b)
		}
	}
	if recording, ok := settings["recording_enabled"]; ok {
		if b, ok := recording.(bool); ok {
			data.RecordingEnabled = types.BoolValue(b)
		}
	}
	if transcription, ok := settings["transcription_enabled"]; ok {
		if b, ok := transcription.(bool); ok {
			data.TranscriptionEnabled = types.BoolValue(b)
		}
	}
	if llmOpt, ok := settings["LLm_optimization_enabled"]; ok {
		if b, ok := llmOpt.(bool); ok {
			data.LLMOptimizationEnabled = types.BoolValue(b)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
