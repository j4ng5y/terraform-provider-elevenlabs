package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

var (
	_ resource.Resource              = &AudioNativeContentUpdateResource{}
	_ resource.ResourceWithConfigure = &AudioNativeContentUpdateResource{}
)

func NewAudioNativeContentUpdateResource() resource.Resource {
	return &AudioNativeContentUpdateResource{}
}

type AudioNativeContentUpdateResource struct {
	client *client.Client
}

type AudioNativeContentUpdateResourceModel struct {
	ProjectID types.String `tfsdk:"project_id"`
	FilePath  types.String `tfsdk:"file_path"`
	VoiceID   types.String `tfsdk:"voice_id"`
	ModelID   types.String `tfsdk:"model_id"`
}

func (r *AudioNativeContentUpdateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_audio_native_content_update"
}

func (r *AudioNativeContentUpdateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for updating Audio Native project content in ElevenLabs. Allows changing the audio file and voice/model settings.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Audio Native project to update.",
			},
			"file_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Path to the new audio file.",
			},
			"voice_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The voice ID to use for TTS.",
			},
			"model_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The model ID to use for TTS.",
			},
		},
	}
}

func (r *AudioNativeContentUpdateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *AudioNativeContentUpdateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AudioNativeContentUpdateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateAudioNativeContent(
		data.ProjectID.ValueString(),
		data.FilePath.ValueString(),
		data.VoiceID.ValueString(),
		data.ModelID.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Error updating Audio Native content", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AudioNativeContentUpdateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// This is an action resource, no persistent state needed
}

func (r *AudioNativeContentUpdateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AudioNativeContentUpdateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateAudioNativeContent(
		data.ProjectID.ValueString(),
		data.FilePath.ValueString(),
		data.VoiceID.ValueString(),
		data.ModelID.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Error updating Audio Native content", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AudioNativeContentUpdateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This is an action resource, no deletion needed
}
