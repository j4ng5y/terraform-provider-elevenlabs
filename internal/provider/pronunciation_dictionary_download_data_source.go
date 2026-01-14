package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

var (
	_ datasource.DataSource              = &PronunciationDictionaryDownloadDataSource{}
	_ datasource.DataSourceWithConfigure = &PronunciationDictionaryDownloadDataSource{}
)

func NewPronunciationDictionaryDownloadDataSource() datasource.DataSource {
	return &PronunciationDictionaryDownloadDataSource{}
}

type PronunciationDictionaryDownloadDataSource struct {
	client *client.Client
}

type PronunciationDictionaryDownloadDataSourceModel struct {
	DictionaryID types.String `tfsdk:"dictionary_id"`
	VersionID    types.String `tfsdk:"version_id"`
	OutputPath   types.String `tfsdk:"output_path"`
	FileName     types.String `tfsdk:"file_name"`
	FileSize     types.Int64  `tfsdk:"file_size"`
	DownloadedAt types.String `tfsdk:"downloaded_at"`
}

func (d *PronunciationDictionaryDownloadDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pronunciation_dictionary_download"
}

func (d *PronunciationDictionaryDownloadDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for downloading ElevenLabs Pronunciation Dictionary PLS files.",
		Attributes: map[string]schema.Attribute{
			"dictionary_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the pronunciation dictionary to download.",
			},
			"version_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The version ID of the pronunciation dictionary. Defaults to latest version.",
			},
			"output_path": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The local path where the PLS file should be saved.",
			},
			"file_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the downloaded file.",
			},
			"file_size": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The size of the downloaded file in bytes.",
			},
			"downloaded_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp when the file was downloaded.",
			},
		},
	}
}

func (d *PronunciationDictionaryDownloadDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PronunciationDictionaryDownloadDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PronunciationDictionaryDownloadDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get version ID - use provided or get latest
	versionID := data.VersionID.ValueString()
	if versionID == "" {
		dict, err := d.client.GetPronunciationDictionary(data.DictionaryID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error getting pronunciation dictionary", err.Error())
			return
		}
		versionID = dict.LatestVersionID
	}

	// Download the PLS file
	plsData, err := d.client.DownloadPronunciationDictionary(data.DictionaryID.ValueString(), versionID)
	if err != nil {
		resp.Diagnostics.AddError("Error downloading pronunciation dictionary", err.Error())
		return
	}

	// Ensure the output directory exists
	outputDir := filepath.Dir(data.OutputPath.ValueString())
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		resp.Diagnostics.AddError("Error creating output directory", err.Error())
		return
	}

	// Write the file
	if err := os.WriteFile(data.OutputPath.ValueString(), plsData, 0644); err != nil {
		resp.Diagnostics.AddError("Error writing PLS file", err.Error())
		return
	}

	// Get file info for response
	fileInfo, err := os.Stat(data.OutputPath.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting file info", err.Error())
		return
	}

	// Set computed values
	data.FileName = types.StringValue(filepath.Base(data.OutputPath.ValueString()))
	data.FileSize = types.Int64Value(fileInfo.Size())
	data.VersionID = types.StringValue(versionID)
	data.DownloadedAt = types.StringValue(fileInfo.ModTime().Format("2006-01-02T15:04:05Z"))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
