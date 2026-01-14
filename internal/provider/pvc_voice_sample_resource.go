package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/models"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &PVCVoiceSampleResource{}
	_ resource.ResourceWithConfigure   = &PVCVoiceSampleResource{}
	_ resource.ResourceWithImportState = &PVCVoiceSampleResource{}
)

// PVCVoiceSampleResourceModel describes the resource data model.
type PVCVoiceSampleResourceModel struct {
	ID            types.String  `tfsdk:"id"`
	VoiceID       types.String  `tfsdk:"voice_id"`
	FilePath      types.String  `tfsdk:"file_path"`
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

func NewPVCVoiceSampleResource() resource.Resource {
	return &PVCVoiceSampleResource{}
}

// PVCVoiceSampleResource defines the resource implementation.
type PVCVoiceSampleResource struct {
	client *client.Client
}

func (r *PVCVoiceSampleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pvc_voice_sample"
}

func (r *PVCVoiceSampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "PVC Voice Sample resource for ElevenLabs. Allows managing training samples for Professional Voice Cloning voices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the PVC voice sample.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"voice_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the PVC voice this sample belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_path": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Local file path to the audio sample file.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the uploaded file.",
			},
			"mime_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The MIME type of the audio file.",
			},
			"size_bytes": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The size of the file in bytes.",
			},
			"hash": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The hash of the file content.",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The processing state of the sample.",
			},
			"transcription": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The transcription of the audio sample.",
			},
			"duration": schema.Float64Attribute{
				Computed:            true,
				MarkdownDescription: "The duration of the audio sample in seconds.",
			},
			"sample_rate": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The sample rate of the audio file.",
			},
			"channels": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The number of audio channels.",
			},
		},
	}
}

func (r *PVCVoiceSampleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *PVCVoiceSampleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PVCVoiceSampleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.AddPVCVoiceSampleRequest{
		FilePath: data.FilePath.ValueString(),
	}

	sample, err := r.client.AddPVCVoiceSample(data.VoiceID.ValueString(), addReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add PVC voice sample, got error: %s", err))
		return
	}

	data.ID = types.StringValue(sample.SampleID)
	data.FileName = types.StringValue(sample.FileName)
	data.MimeType = types.StringValue(sample.MimeType)
	data.SizeBytes = types.Int64Value(int64(sample.SizeBytes))
	data.Hash = types.StringValue(sample.Hash)
	data.State = types.StringValue(sample.State)
	data.Transcription = types.StringValue(sample.Transcription)
	data.Duration = types.Float64Value(sample.Duration)
	data.SampleRate = types.Int64Value(int64(sample.SampleRate))
	data.Channels = types.Int64Value(int64(sample.Channels))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PVCVoiceSampleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PVCVoiceSampleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get all samples for the voice and find the specific one
	samplesResp, err := r.client.ListPVCVoiceSamples(data.VoiceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list PVC voice samples, got error: %s", err))
		return
	}

	var foundSample *models.PVCVoiceSample
	for _, sample := range samplesResp.Samples {
		if sample.SampleID == data.ID.ValueString() {
			foundSample = &sample
			break
		}
	}

	if foundSample == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("PVC voice sample with ID %s not found", data.ID.ValueString()))
		return
	}

	data.FileName = types.StringValue(foundSample.FileName)
	data.MimeType = types.StringValue(foundSample.MimeType)
	data.SizeBytes = types.Int64Value(int64(foundSample.SizeBytes))
	data.Hash = types.StringValue(foundSample.Hash)
	data.State = types.StringValue(foundSample.State)
	data.Transcription = types.StringValue(foundSample.Transcription)
	data.Duration = types.Float64Value(foundSample.Duration)
	data.SampleRate = types.Int64Value(int64(foundSample.SampleRate))
	data.Channels = types.Int64Value(int64(foundSample.Channels))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PVCVoiceSampleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PVCVoiceSampleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &models.UpdatePVCVoiceSampleRequest{
		Transcription: data.Transcription.ValueString(),
	}

	err := r.client.UpdatePVCVoiceSample(data.VoiceID.ValueString(), data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update PVC voice sample, got error: %s", err))
		return
	}

	// Read the updated sample to get current state
	samplesResp, err := r.client.ListPVCVoiceSamples(data.VoiceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list PVC voice samples, got error: %s", err))
		return
	}

	var foundSample *models.PVCVoiceSample
	for _, sample := range samplesResp.Samples {
		if sample.SampleID == data.ID.ValueString() {
			foundSample = &sample
			break
		}
	}

	if foundSample == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("PVC voice sample with ID %s not found after update", data.ID.ValueString()))
		return
	}

	data.State = types.StringValue(foundSample.State)
	data.Duration = types.Float64Value(foundSample.Duration)
	data.SampleRate = types.Int64Value(int64(foundSample.SampleRate))
	data.Channels = types.Int64Value(int64(foundSample.Channels))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PVCVoiceSampleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PVCVoiceSampleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePVCVoiceSample(data.VoiceID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete PVC voice sample, got error: %s", err))
		return
	}
}

func (r *PVCVoiceSampleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: voice_id/sample_id
	idParts := strings.Split(req.ID, "/")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in the format: voice_id/sample_id",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("voice_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}
